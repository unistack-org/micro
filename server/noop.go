package server

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/register"
	maddr "go.unistack.org/micro/v4/util/addr"
	mnet "go.unistack.org/micro/v4/util/net"
	"go.unistack.org/micro/v4/util/rand"
)

// DefaultCodecs will be used to encode/decode
var DefaultCodecs = map[string]codec.Codec{
	"application/octet-stream": codec.NewCodec(),
}

type rpcHandler struct {
	opts    HandlerOptions
	handler interface{}
	name    string
}

func newRPCHandler(handler interface{}, opts ...HandlerOption) Handler {
	options := NewHandlerOptions(opts...)

	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	return &rpcHandler{
		name:    name,
		handler: handler,
		opts:    options,
	}
}

func (r *rpcHandler) Name() string {
	return r.name
}

func (r *rpcHandler) Handler() interface{} {
	return r.handler
}

func (r *rpcHandler) Options() HandlerOptions {
	return r.opts
}

type noopServer struct {
	h          Handler
	wg         *sync.WaitGroup
	rsvc       *register.Service
	handlers   map[string]Handler
	exit       chan chan error
	opts       Options
	mu         sync.RWMutex
	registered bool
	started    bool
}

// NewServer returns new noop server
func NewServer(opts ...Option) Server {
	n := &noopServer{opts: NewOptions(opts...)}
	if n.handlers == nil {
		n.handlers = make(map[string]Handler)
	}
	if n.exit == nil {
		n.exit = make(chan chan error)
	}
	return n
}

func (n *noopServer) Live() bool {
	return true
}

func (n *noopServer) Ready() bool {
	return true
}

func (n *noopServer) Health() bool {
	return true
}

func (n *noopServer) Handle(handler Handler) error {
	n.h = handler
	return nil
}

func (n *noopServer) Name() string {
	return n.opts.Name
}

func (n *noopServer) NewHandler(h interface{}, opts ...HandlerOption) Handler {
	return newRPCHandler(h, opts...)
}

func (n *noopServer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	if n.handlers == nil {
		n.handlers = make(map[string]Handler, 1)
	}

	if n.exit == nil {
		n.exit = make(chan chan error)
	}

	return nil
}

func (n *noopServer) Options() Options {
	return n.opts
}

func (n *noopServer) String() string {
	return "noop"
}

//nolint:gocyclo
func (n *noopServer) Register() error {
	n.mu.RLock()
	rsvc := n.rsvc
	config := n.opts
	n.mu.RUnlock()

	// if service already filled, reuse it and return early
	if rsvc != nil {
		return DefaultRegisterFunc(rsvc, config)
	}

	var err error
	var service *register.Service
	var cacheService bool

	service, err = NewRegisterService(n)
	if err != nil {
		return err
	}

	n.mu.RLock()
	registered := n.registered
	n.mu.RUnlock()

	if !registered {
		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Info(n.opts.Context, fmt.Sprintf("register [%s] Registering node: %s", config.Register.String(), service.Nodes[0].ID))
		}
	}

	// register the service
	if err = DefaultRegisterFunc(service, config); err != nil {
		return err
	}

	// already registered? don't need to register subscribers
	if registered {
		return nil
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.registered = true
	if cacheService {
		n.rsvc = service
	}

	return nil
}

func (n *noopServer) Deregister() error {
	var err error

	n.mu.RLock()
	config := n.opts
	n.mu.RUnlock()

	service, err := NewRegisterService(n)
	if err != nil {
		return err
	}

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Info(n.opts.Context, fmt.Sprintf("deregistering node: %s", service.Nodes[0].ID))
	}

	if err := DefaultDeregisterFunc(service, config); err != nil {
		return err
	}

	n.mu.Lock()
	n.rsvc = nil

	if !n.registered {
		n.mu.Unlock()
		return nil
	}

	n.registered = false

	n.mu.Unlock()
	return nil
}

//nolint:gocyclo
func (n *noopServer) Start() error {
	n.mu.RLock()
	if n.started {
		n.mu.RUnlock()
		return nil
	}
	config := n.Options()
	n.mu.RUnlock()

	// use 127.0.0.1 to avoid scan of all network interfaces
	addr, err := maddr.Extract("127.0.0.1")
	if err != nil {
		return err
	}
	var rng rand.Rand
	i := rng.Intn(20000)
	// set addr with port
	addr = mnet.HostPort(addr, 10000+i)

	config.Address = addr

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Info(n.opts.Context, "server [noop] Listening on "+config.Address)
	}

	n.mu.Lock()
	if len(config.Advertise) == 0 {
		config.Advertise = config.Address
	}
	n.mu.Unlock()

	// use RegisterCheck func before register
	// nolint: nestif
	if err := config.RegisterCheck(config.Context); err != nil {
		if config.Logger.V(logger.ErrorLevel) {
			config.Logger.Error(n.opts.Context, fmt.Sprintf("server %s-%s register check error", config.Name, config.ID), err)
		}
	} else {
		// announce self to the world
		if err := n.Register(); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Error(n.opts.Context, "server register error", err)
			}
		}
	}

	go func() {
		t := new(time.Ticker)

		// only process if it exists
		if config.RegisterInterval > time.Duration(0) {
			// new ticker
			t = time.NewTicker(config.RegisterInterval)
		}

		// return error chan
		var ch chan error

	Loop:
		for {
			select {
			// register self on interval
			case <-t.C:
				n.mu.RLock()
				registered := n.registered
				n.mu.RUnlock()
				rerr := config.RegisterCheck(config.Context)
				// nolint: nestif
				if rerr != nil && registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, fmt.Sprintf("server %s-%s register check error, deregister it", config.Name, config.ID), rerr)
					}
					// deregister self in case of error
					if err := n.Deregister(); err != nil {
						if config.Logger.V(logger.ErrorLevel) {
							config.Logger.Error(n.opts.Context, fmt.Sprintf("server %s-%s deregister error", config.Name, config.ID), err)
						}
					}
				} else if rerr != nil && !registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, fmt.Sprintf("server %s-%s register check error", config.Name, config.ID), rerr)
					}
					continue
				}
				if err := n.Register(); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, fmt.Sprintf("server %s-%s register error", config.Name, config.ID), err)
					}
				}
			// wait for exit
			case ch = <-n.exit:
				break Loop
			}
		}

		// deregister self
		if err := n.Deregister(); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Error(n.opts.Context, "server deregister error", err)
			}
		}

		// wait for waitgroup
		if n.wg != nil {
			n.wg.Wait()
		}

		// close transport
		ch <- nil

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Info(n.opts.Context, fmt.Sprintf("broker [%s] Disconnected from %s", config.Broker.String(), config.Broker.Address()))
		}
		// disconnect broker
		if err := config.Broker.Disconnect(config.Context); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Error(n.opts.Context, fmt.Sprintf("broker [%s] disconnect error", config.Broker.String()), err)
			}
		}
	}()

	// mark the server as started
	n.mu.Lock()
	n.started = true
	n.mu.Unlock()

	return nil
}

func (n *noopServer) Stop() error {
	n.mu.RLock()
	if !n.started {
		n.mu.RUnlock()
		return nil
	}
	n.mu.RUnlock()

	ch := make(chan error)
	n.exit <- ch

	err := <-ch
	n.mu.Lock()
	n.rsvc = nil
	n.started = false
	n.mu.Unlock()

	return err
}
