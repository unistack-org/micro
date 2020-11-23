package server

import (
	"fmt"
	"sort"
	"sync"
	"time"

	//	cjson "github.com/unistack-org/micro-codec-json"
	//	cjsonrpc "github.com/unistack-org/micro-codec-jsonrpc"
	//	cproto "github.com/unistack-org/micro-codec-proto"
	//	cprotorpc "github.com/unistack-org/micro-codec-protorpc"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/registry"
)

var (
	DefaultCodecs = map[string]codec.Codec{
		//"application/json":         cjson.NewCodec,
		//"application/json-rpc":     cjsonrpc.NewCodec,
		//"application/protobuf":     cproto.NewCodec,
		//"application/proto-rpc":    cprotorpc.NewCodec,
		"application/octet-stream": codec.NewCodec(),
	}
)

const (
	defaultContentType = "application/json"
)

type noopServer struct {
	h           Handler
	opts        Options
	rsvc        *registry.Service
	handlers    map[string]Handler
	subscribers map[*subscriber][]broker.Subscriber
	registered  bool
	started     bool
	exit        chan chan error
	wg          *sync.WaitGroup
	sync.RWMutex
}

// NewServer returns new noop server
func NewServer(opts ...Option) Server {
	return &noopServer{opts: NewOptions(opts...)}
}

func (n *noopServer) newCodec(contentType string) (codec.Codec, error) {
	if cf, ok := n.opts.Codecs[contentType]; ok {
		return cf, nil
	}
	if cf, ok := DefaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, codec.ErrUnknownContentType
}

func (n *noopServer) Handle(handler Handler) error {
	n.h = handler
	return nil
}

func (n *noopServer) Subscribe(sb Subscriber) error {
	sub, ok := sb.(*subscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *subscriber")
	}
	if len(sub.handlers) == 0 {
		return fmt.Errorf("invalid subscriber: no handler functions")
	}

	if err := ValidateSubscriber(sb); err != nil {
		return err
	}

	n.Lock()
	if _, ok = n.subscribers[sub]; ok {
		n.Unlock()
		return fmt.Errorf("subscriber %v already exists", sub)
	}

	n.subscribers[sub] = nil
	n.Unlock()
	return nil
}

func (n *noopServer) NewHandler(h interface{}, opts ...HandlerOption) Handler {
	return newRpcHandler(h, opts...)
}

func (n *noopServer) NewSubscriber(topic string, sb interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, sb, opts...)
}

func (n *noopServer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	if n.handlers == nil {
		n.handlers = make(map[string]Handler)
	}
	if n.subscribers == nil {
		n.subscribers = make(map[*subscriber][]broker.Subscriber)
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

func (n *noopServer) Register() error {
	n.RLock()
	rsvc := n.rsvc
	config := n.opts
	n.RUnlock()

	// if service already filled, reuse it and return early
	if rsvc != nil {
		if err := DefaultRegisterFunc(rsvc, config); err != nil {
			return err
		}
		return nil
	}

	var err error
	var service *registry.Service
	var cacheService bool

	service, err = NewRegistryService(n)
	if err != nil {
		return err
	}

	n.RLock()
	// Maps are ordered randomly, sort the keys for consistency
	var handlerList []string
	for n, e := range n.handlers {
		// Only advertise non internal handlers
		if !e.Options().Internal {
			handlerList = append(handlerList, n)
		}
	}

	sort.Strings(handlerList)

	var subscriberList []*subscriber
	for e := range n.subscribers {
		// Only advertise non internal subscribers
		if !e.Options().Internal {
			subscriberList = append(subscriberList, e)
		}
	}
	sort.Slice(subscriberList, func(i, j int) bool {
		return subscriberList[i].topic > subscriberList[j].topic
	})

	endpoints := make([]*registry.Endpoint, 0, len(handlerList)+len(subscriberList))
	for _, h := range handlerList {
		endpoints = append(endpoints, n.handlers[h].Endpoints()...)
	}
	for _, e := range subscriberList {
		endpoints = append(endpoints, e.Endpoints()...)
	}
	n.RUnlock()

	service.Nodes[0].Metadata["protocol"] = "noop"
	service.Nodes[0].Metadata["transport"] = "noop"
	service.Endpoints = endpoints

	n.RLock()
	registered := n.registered
	n.RUnlock()

	if !registered {
		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof("registry [%s] Registering node: %s", config.Registry.String(), service.Nodes[0].Id)
		}
	}

	// register the service
	if err := DefaultRegisterFunc(service, config); err != nil {
		return err
	}

	// already registered? don't need to register subscribers
	if registered {
		return nil
	}

	n.Lock()
	defer n.Unlock()

	cx := config.Context

	for sb := range n.subscribers {
		handler := n.createSubHandler(sb, config)
		var opts []broker.SubscribeOption
		if queue := sb.Options().Queue; len(queue) > 0 {
			opts = append(opts, broker.SubscribeGroup(queue))
		}

		if sb.Options().Context != nil {
			cx = sb.Options().Context
		}

		opts = append(opts, broker.SubscribeContext(cx), broker.SubscribeAutoAck(sb.Options().AutoAck))

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof("subscribing to topic: %s", sb.Topic())
		}
		sub, err := config.Broker.Subscribe(cx, sb.Topic(), handler, opts...)
		if err != nil {
			return err
		}
		n.subscribers[sb] = []broker.Subscriber{sub}
	}

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

	service, err := NewRegistryService(n)
	if err != nil {
		return err
	}

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Infof("deregistering node: %s", service.Nodes[0].Id)
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

	cx := config.Context

	wg := sync.WaitGroup{}
	for sb, subs := range n.subscribers {
		for _, sub := range subs {
			if sb.Options().Context != nil {
				cx = sb.Options().Context
			}

			wg.Add(1)
			go func(s broker.Subscriber) {
				defer wg.Done()
				if config.Logger.V(logger.InfoLevel) {
					config.Logger.Infof("unsubscribing from topic: %s", s.Topic())
				}
				if err := s.Unsubscribe(cx); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf("unsubscribing from topic: %s err: %v", s.Topic(), err)
					}
				}
			}(sub)
		}
		n.subscribers[sb] = nil
	}
	wg.Wait()

	n.Unlock()
	return nil
}

func (n *noopServer) Start() error {
	n.RLock()
	if n.started {
		n.RUnlock()
		return nil
	}
	config := n.Options()
	n.RUnlock()

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Infof("server [noop] Listening on %s", config.Address)
	}
	n.Lock()
	if len(config.Advertise) == 0 {
		config.Advertise = config.Address
	}
	n.Unlock()

	// only connect if we're subscribed
	if len(n.subscribers) > 0 {
		// connect to the broker
		if err := config.Broker.Connect(config.Context); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Errorf("broker [%s] connect error: %v", config.Broker.String(), err)
			}
			return err
		}

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof("broker [%s] Connected to %s", config.Broker.String(), config.Broker.Address())
		}
	}

	// use RegisterCheck func before register
	if err := config.RegisterCheck(config.Context); err != nil {
		if config.Logger.V(logger.ErrorLevel) {
			config.Logger.Errorf("server %s-%s register check error: %s", config.Name, config.Id, err)
		}
	} else {
		// announce self to the world
		if err := n.Register(); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Errorf("server register error: %v", err)
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
				if rerr != nil && registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf("server %s-%s register check error: %s, deregister it", config.Name, config.Id, rerr)
					}
					// deregister self in case of error
					if err := n.Deregister(); err != nil {
						if config.Logger.V(logger.ErrorLevel) {
							config.Logger.Errorf("server %s-%s deregister error: %s", config.Name, config.Id, err)
						}
					}
				} else if rerr != nil && !registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf("server %s-%s register check error: %s", config.Name, config.Id, rerr)
					}
					continue
				}
				if err := n.Register(); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf("server %s-%s register error: %s", config.Name, config.Id, err)
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
				config.Logger.Errorf("server deregister error: ", err)
			}
		}

		// wait for waitgroup
		if n.wg != nil {
			n.wg.Wait()
		}

		// close transport
		ch <- nil

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof("broker [%s] Disconnected from %s", config.Broker.String(), config.Broker.Address())
		}
		// disconnect broker
		if err := config.Broker.Disconnect(config.Context); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Errorf("broker [%s] disconnect error: %v", config.Broker.String(), err)
			}
		}
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
