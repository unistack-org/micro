package server

import (
	"reflect"
	"sort"
	"sync"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/register"
	maddr "go.unistack.org/micro/v4/util/addr"
	mnet "go.unistack.org/micro/v4/util/net"
	"go.unistack.org/micro/v4/util/rand"
)

// DefaultCodecs will be used to encode/decode
var DefaultCodecs = map[string]codec.Codec{
	"application/octet-stream": codec.NewCodec(),
}

type noopServer struct {
	h        *rpcHandler
	wg       *sync.WaitGroup
	rsvc     *register.Service
	handlers map[string]*rpcHandler
	exit     chan chan error
	opts     Options
	sync.RWMutex
	registered bool
	started    bool
}

// NewServer returns new noop server
func NewServer(opts ...options.Option) Server {
	n := &noopServer{opts: NewOptions(opts...)}
	if n.handlers == nil {
		n.handlers = make(map[string]*rpcHandler)
	}
	if n.exit == nil {
		n.exit = make(chan chan error)
	}
	return n
}

func (n *noopServer) Handle(h interface{}, opts ...options.Option) error {
	n.h = newRPCHandler(h, opts...)
	return nil
}

func (n *noopServer) Name() string {
	return n.opts.Name
}

func (n *noopServer) Init(opts ...options.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	if n.handlers == nil {
		n.handlers = make(map[string]*rpcHandler, 1)
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
	n.RLock()
	rsvc := n.rsvc
	config := n.opts
	n.RUnlock()

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

	n.RLock()
	handlerList := make([]string, 0, len(n.handlers))
	for n := range n.handlers {
		handlerList = append(handlerList, n)
	}

	sort.Strings(handlerList)

	endpoints := make([]*register.Endpoint, 0, len(handlerList))
	for _, h := range handlerList {
		endpoints = append(endpoints, n.handlers[h].Endpoints()...)
	}
	n.RUnlock()

	service.Nodes[0].Metadata["protocol"] = "noop"
	service.Nodes[0].Metadata["transport"] = service.Nodes[0].Metadata["protocol"]
	service.Endpoints = endpoints

	n.RLock()
	registered := n.registered
	n.RUnlock()

	if !registered {
		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Info(n.opts.Context, "register ["+config.Register.String()+"] Registering node: "+service.Nodes[0].ID)
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

	n.Lock()
	defer n.Unlock()

	n.registered = true
	if cacheService {
		n.rsvc = service
	}

	return nil
}

func (n *noopServer) Deregister() error {
	var err error

	n.RLock()
	config := n.opts
	n.RUnlock()

	service, err := NewRegisterService(n)
	if err != nil {
		return err
	}

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Info(n.opts.Context, "deregistering node: "+service.Nodes[0].ID)
	}

	if err := DefaultDeregisterFunc(service, config); err != nil {
		return err
	}

	n.Lock()
	n.rsvc = nil

	if !n.registered {
		n.Unlock()
		return nil
	}

	n.registered = false

	n.Unlock()
	return nil
}

//nolint:gocyclo
func (n *noopServer) Start() error {
	n.RLock()
	if n.started {
		n.RUnlock()
		return nil
	}
	config := n.Options()
	n.RUnlock()

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

	n.Lock()
	if len(config.Advertise) == 0 {
		config.Advertise = config.Address
	}
	n.Unlock()

	// use RegisterCheck func before register
	// nolint: nestif
	if err := config.RegisterCheck(config.Context); err != nil {
		if config.Logger.V(logger.ErrorLevel) {
			config.Logger.Error(n.opts.Context, "server "+config.Name+"-"+config.ID+" register check error: "+err.Error())
		}
	} else {
		// announce self to the world
		if err := n.Register(); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Error(n.opts.Context, "server register error: "+err.Error())
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
				n.RLock()
				registered := n.registered
				n.RUnlock()
				rerr := config.RegisterCheck(config.Context)
				// nolint: nestif
				if rerr != nil && registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, "server "+config.Name+"-"+config.ID+" register check error: ", rerr.Error())
					}
					// deregister self in case of error
					if err := n.Deregister(); err != nil {
						if config.Logger.V(logger.ErrorLevel) {
							config.Logger.Error(n.opts.Context, "server "+config.Name+"-"+config.ID+" deregister error: ", err.Error())
						}
					}
				} else if rerr != nil && !registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, "server "+config.Name+"-"+config.ID+" register check error: ", rerr.Error())
					}
					continue
				}
				if err := n.Register(); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, "server "+config.Name+"-"+config.ID+" register error: ", err.Error())
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
				config.Logger.Error(n.opts.Context, "server deregister error: "+err.Error())
			}
		}

		// wait for waitgroup
		if n.wg != nil {
			n.wg.Wait()
		}

		// close transport
		ch <- nil
	}()

	// mark the server as started
	n.Lock()
	n.started = true
	n.Unlock()

	return nil
}

func (n *noopServer) Stop() error {
	n.RLock()
	if !n.started {
		n.RUnlock()
		return nil
	}
	n.RUnlock()

	ch := make(chan error)
	n.exit <- ch

	err := <-ch
	n.Lock()
	n.rsvc = nil
	n.started = false
	n.Unlock()

	return err
}

type rpcHandler struct {
	opts      HandleOptions
	handler   interface{}
	name      string
	endpoints []*register.Endpoint
}

func newRPCHandler(handler interface{}, opts ...options.Option) *rpcHandler {
	options := NewHandleOptions(opts...)

	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	var endpoints []*register.Endpoint

	for m := 0; m < typ.NumMethod(); m++ {
		if e := register.ExtractEndpoint(typ.Method(m)); e != nil {
			e.Name = name + "." + e.Name

			for k, v := range options.Metadata[e.Name] {
				e.Metadata[k] = v
			}

			endpoints = append(endpoints, e)
		}
	}

	return &rpcHandler{
		name:      name,
		handler:   handler,
		endpoints: endpoints,
		opts:      options,
	}
}

func (r *rpcHandler) Name() string {
	return r.name
}

func (r *rpcHandler) Handler() interface{} {
	return r.handler
}

func (r *rpcHandler) Endpoints() []*register.Endpoint {
	return r.endpoints
}

func (r *rpcHandler) Options() HandleOptions {
	return r.opts
}
