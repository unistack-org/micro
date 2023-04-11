package server

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"go.unistack.org/micro/v4/broker"
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

const (
	defaultContentType = "application/json"
)

type noopServer struct {
	h           Handler
	wg          *sync.WaitGroup
	rsvc        *register.Service
	handlers    map[string]Handler
	subscribers map[*subscriber][]broker.Subscriber
	exit        chan chan error
	opts        Options
	sync.RWMutex
	registered bool
	started    bool
}

// NewServer returns new noop server
func NewServer(opts ...Option) Server {
	n := &noopServer{opts: NewOptions(opts...)}
	if n.handlers == nil {
		n.handlers = make(map[string]Handler)
	}
	if n.subscribers == nil {
		n.subscribers = make(map[*subscriber][]broker.Subscriber)
	}
	if n.exit == nil {
		n.exit = make(chan chan error)
	}
	return n
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

func (n *noopServer) Name() string {
	return n.opts.Name
}

func (n *noopServer) Subscribe(sb Subscriber) error {
	sub, ok := sb.(*subscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *subscriber")
	} else if len(sub.handlers) == 0 {
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
	return newRPCHandler(h, opts...)
}

func (n *noopServer) NewSubscriber(topic string, sb interface{}, opts ...SubscriberOption) Subscriber {
	return newSubscriber(topic, sb, opts...)
}

func (n *noopServer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	if n.handlers == nil {
		n.handlers = make(map[string]Handler, 1)
	}
	if n.subscribers == nil {
		n.subscribers = make(map[*subscriber][]broker.Subscriber, 1)
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

	subscriberList := make([]*subscriber, 0, len(n.subscribers))
	for e := range n.subscribers {
		subscriberList = append(subscriberList, e)
	}
	sort.Slice(subscriberList, func(i, j int) bool {
		return subscriberList[i].topic > subscriberList[j].topic
	})

	endpoints := make([]*register.Endpoint, 0, len(handlerList)+len(subscriberList))
	for _, h := range handlerList {
		endpoints = append(endpoints, n.handlers[h].Endpoints()...)
	}
	for _, e := range subscriberList {
		endpoints = append(endpoints, e.Endpoints()...)
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
			config.Logger.Infof(n.opts.Context, "register [%s] Registering node: %s", config.Register.String(), service.Nodes[0].ID)
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

	cx := config.Context

	var sub broker.Subscriber

	for sb := range n.subscribers {
		if sb.Options().Context != nil {
			cx = sb.Options().Context
		}

		opts := []broker.SubscribeOption{broker.SubscribeContext(cx), broker.SubscribeAutoAck(sb.Options().AutoAck)}
		if queue := sb.Options().Queue; len(queue) > 0 {
			opts = append(opts, broker.SubscribeGroup(queue))
		}

		if sb.Options().Batch {
			// batch processing handler
			sub, err = config.Broker.BatchSubscribe(cx, sb.Topic(), n.newBatchSubHandler(sb, config), opts...)
		} else {
			// single processing handler
			sub, err = config.Broker.Subscribe(cx, sb.Topic(), n.newSubHandler(sb, config), opts...)
		}

		if err != nil {
			return err
		}

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof(n.opts.Context, "subscribing to topic: %s", sb.Topic())
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

	service, err := NewRegisterService(n)
	if err != nil {
		return err
	}

	if config.Logger.V(logger.InfoLevel) {
		config.Logger.Infof(n.opts.Context, "deregistering node: %s", service.Nodes[0].ID)
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
		for idx := range subs {
			if sb.Options().Context != nil {
				cx = sb.Options().Context
			}

			ncx := cx
			wg.Add(1)
			go func(s broker.Subscriber) {
				defer wg.Done()
				if config.Logger.V(logger.InfoLevel) {
					config.Logger.Infof(n.opts.Context, "unsubscribing from topic: %s", s.Topic())
				}
				if err := s.Unsubscribe(ncx); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf(n.opts.Context, "unsubscribing from topic: %s err: %v", s.Topic(), err)
					}
				}
			}(subs[idx])
		}
		n.subscribers[sb] = nil
	}
	wg.Wait()

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
		config.Logger.Infof(n.opts.Context, "server [noop] Listening on %s", config.Address)
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
				config.Logger.Errorf(n.opts.Context, "broker [%s] connect error: %v", config.Broker.String(), err)
			}
			return err
		}

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof(n.opts.Context, "broker [%s] Connected to %s", config.Broker.String(), config.Broker.Address())
		}
	}

	// use RegisterCheck func before register
	// nolint: nestif
	if err := config.RegisterCheck(config.Context); err != nil {
		if config.Logger.V(logger.ErrorLevel) {
			config.Logger.Errorf(n.opts.Context, "server %s-%s register check error: %s", config.Name, config.ID, err)
		}
	} else {
		// announce self to the world
		if err := n.Register(); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Errorf(n.opts.Context, "server register error: %v", err)
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
						config.Logger.Errorf(n.opts.Context, "server %s-%s register check error: %s, deregister it", config.Name, config.ID, rerr)
					}
					// deregister self in case of error
					if err := n.Deregister(); err != nil {
						if config.Logger.V(logger.ErrorLevel) {
							config.Logger.Errorf(n.opts.Context, "server %s-%s deregister error: %s", config.Name, config.ID, err)
						}
					}
				} else if rerr != nil && !registered {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf(n.opts.Context, "server %s-%s register check error: %s", config.Name, config.ID, rerr)
					}
					continue
				}
				if err := n.Register(); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Errorf(n.opts.Context, "server %s-%s register error: %s", config.Name, config.ID, err)
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
				config.Logger.Errorf(n.opts.Context, "server deregister error: ", err)
			}
		}

		// wait for waitgroup
		if n.wg != nil {
			n.wg.Wait()
		}

		// close transport
		ch <- nil

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Infof(n.opts.Context, "broker [%s] Disconnected from %s", config.Broker.String(), config.Broker.Address())
		}
		// disconnect broker
		if err := config.Broker.Disconnect(config.Context); err != nil {
			if config.Logger.V(logger.ErrorLevel) {
				config.Logger.Errorf(n.opts.Context, "broker [%s] disconnect error: %v", config.Broker.String(), err)
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
