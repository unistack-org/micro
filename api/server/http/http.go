// Package http provides a http server with features; acme, cors, etc
package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"

	"github.com/unistack-org/micro/v3/api/server"
	"github.com/unistack-org/micro/v3/logger"
)

type httpServer struct {
	mux  *http.ServeMux
	opts server.Options

	sync.RWMutex
	address string
	exit    chan chan error
}

func NewServer(address string, opts ...server.Option) server.Server {
	return &httpServer{
		opts:    server.NewOptions(opts...),
		mux:     http.NewServeMux(),
		address: address,
		exit:    make(chan chan error),
	}
}

func (s *httpServer) Address() string {
	s.RLock()
	defer s.RUnlock()
	return s.address
}

func (s *httpServer) Init(opts ...server.Option) error {
	for _, o := range opts {
		o(&s.opts)
	}
	return nil
}

func (s *httpServer) Handle(path string, handler http.Handler) {
	// TODO: move this stuff out to one place with ServeHTTP

	// apply the wrappers, e.g. auth
	for _, wrapper := range s.opts.Wrappers {
		handler = wrapper(handler)
	}

	s.mux.Handle(path, handler)
}

func (s *httpServer) Start() error {
	var l net.Listener
	var err error

	s.RLock()
	config := s.opts
	s.RUnlock()
	if s.opts.EnableACME && s.opts.ACMEProvider != nil {
		// should we check the address to make sure its using :443?
		l, err = s.opts.ACMEProvider.Listen(s.opts.ACMEHosts...)
	} else if s.opts.EnableTLS && s.opts.TLSConfig != nil {
		l, err = tls.Listen("tcp", s.address, s.opts.TLSConfig)
	} else {
		// otherwise plain listen
		l, err = net.Listen("tcp", s.address)
	}
	if err != nil {
		return err
	}

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Infof("HTTP API Listening on %s", l.Addr().String())
	}

	s.Lock()
	s.address = l.Addr().String()
	s.Unlock()

	go func() {
		if err := http.Serve(l, s.mux); err != nil {
			// temporary fix
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Errorf("serve err: %v", err)
			}
			s.Stop()
		}
	}()

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	return nil
}

func (s *httpServer) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func (s *httpServer) String() string {
	return "http"
}
