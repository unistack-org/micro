package server

import (
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/options"
	"go.unistack.org/micro/v3/register"
	maddr "go.unistack.org/micro/v3/util/addr"
	mnet "go.unistack.org/micro/v3/util/net"
	"go.unistack.org/micro/v3/util/rand"
)

// DefaultCodecs will be used to encode/decode
var DefaultCodecs = map[string]codec.Codec{
	"application/octet-stream": codec.NewCodec(),
}

const (
	defaultContentType = "application/json"
)

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

type rpcMessage struct {
	payload     interface{}
	codec       codec.Codec
	header      metadata.Metadata
	topic       string
	contentType string
}

func (r *rpcMessage) ContentType() string {
	return r.contentType
}

func (r *rpcMessage) Topic() string {
	return r.topic
}

func (r *rpcMessage) Body() interface{} {
	return r.payload
}

func (r *rpcMessage) Header() metadata.Metadata {
	return r.header
}

func (r *rpcMessage) Codec() codec.Codec {
	return r.codec
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
	registered := n.registered
	n.RUnlock()

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
		config.Logger.Info(n.opts.Context, fmt.Sprintf("deregistering node: %s", service.Nodes[0].ID))
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
					config.Logger.Info(n.opts.Context, "unsubscribing from topic: "+s.Topic())
				}
				if err := s.Unsubscribe(ncx); err != nil {
					if config.Logger.V(logger.ErrorLevel) {
						config.Logger.Error(n.opts.Context, "unsubscribing from topic: "+s.Topic(), err)
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
		config.Logger.Info(n.opts.Context, "server [noop] Listening on "+config.Address)
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
				config.Logger.Error(n.opts.Context, fmt.Sprintf("broker [%s] connect error", config.Broker.String()), err)
			}
			return err
		}

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Info(n.opts.Context, fmt.Sprintf("broker [%s] Connected to %s", config.Broker.String(), config.Broker.Address()))
		}
	}

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

	if err := n.subscribe(); err != nil {
		return err
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
	n.Lock()
	n.started = true
	n.Unlock()

	return nil
}

func (n *noopServer) subscribe() error {
	config := n.Options()

	subCtx := config.Context

	for sb := range n.subscribers {

		if cx := sb.Options().Context; cx != nil {
			subCtx = cx
		}

		opts := []broker.SubscribeOption{
			broker.SubscribeContext(subCtx),
			broker.SubscribeAutoAck(sb.Options().AutoAck),
			broker.SubscribeBodyOnly(sb.Options().BodyOnly),
		}

		if queue := sb.Options().Queue; len(queue) > 0 {
			opts = append(opts, broker.SubscribeGroup(queue))
		}

		if config.Logger.V(logger.InfoLevel) {
			config.Logger.Info(n.opts.Context, "subscribing to topic: "+sb.Topic())
		}

		sub, err := config.Broker.Subscribe(subCtx, sb.Topic(), n.createSubHandler(sb, config), opts...)
		if err != nil {
			return err
		}

		n.subscribers[sb] = []broker.Subscriber{sub}
	}

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

func newSubscriber(topic string, sub interface{}, opts ...SubscriberOption) Subscriber {
	var handlers []*handler

	options := NewSubscriberOptions(opts...)

	if typ := reflect.TypeOf(sub); typ.Kind() == reflect.Func {
		h := &handler{
			method: reflect.ValueOf(sub),
		}

		switch typ.NumIn() {
		case 1:
			h.reqType = typ.In(0)
		case 2:
			h.ctxType = typ.In(0)
			h.reqType = typ.In(1)
		}

		handlers = append(handlers, h)
	} else {
		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			h := &handler{
				method: method.Func,
			}

			switch method.Type.NumIn() {
			case 2:
				h.reqType = method.Type.In(1)
			case 3:
				h.ctxType = method.Type.In(1)
				h.reqType = method.Type.In(2)
			}

			handlers = append(handlers, h)
		}
	}

	return &subscriber{
		rcvr:       reflect.ValueOf(sub),
		typ:        reflect.TypeOf(sub),
		topic:      topic,
		subscriber: sub,
		handlers:   handlers,
		opts:       options,
	}
}

//nolint:gocyclo
func (n *noopServer) createSubHandler(sb *subscriber, opts Options) broker.Handler {
	return func(p broker.Event) (err error) {
		defer func() {
			if r := recover(); r != nil {
				n.RLock()
				config := n.opts
				n.RUnlock()
				if config.Logger.V(logger.ErrorLevel) {
					config.Logger.Error(n.opts.Context, "panic recovered: ", r)
					config.Logger.Error(n.opts.Context, string(debug.Stack()))
				}
				err = errors.InternalServerError(n.opts.Name+".subscriber", "panic recovered: %v", r)
			}
		}()

		msg := p.Message()
		// if we don't have headers, create empty map
		if msg.Header == nil {
			msg.Header = metadata.New(2)
		}

		ct := msg.Header["Content-Type"]
		if len(ct) == 0 {
			msg.Header.Set(metadata.HeaderContentType, defaultContentType)
			ct = defaultContentType
		}
		cf, err := n.newCodec(ct)
		if err != nil {
			return err
		}

		hdr := metadata.New(len(msg.Header))
		for k, v := range msg.Header {
			hdr.Set(k, v)
		}

		ctx := metadata.NewIncomingContext(sb.opts.Context, hdr)

		results := make(chan error, len(sb.handlers))

		for i := 0; i < len(sb.handlers); i++ {
			handler := sb.handlers[i]

			var isVal bool
			var req reflect.Value

			if handler.reqType.Kind() == reflect.Ptr {
				req = reflect.New(handler.reqType.Elem())
			} else {
				req = reflect.New(handler.reqType)
				isVal = true
			}
			if isVal {
				req = req.Elem()
			}

			if err = cf.Unmarshal(msg.Body, req.Interface()); err != nil {
				return err
			}

			fn := func(ctx context.Context, msg Message) error {
				var vals []reflect.Value
				if sb.typ.Kind() != reflect.Func {
					vals = append(vals, sb.rcvr)
				}
				if handler.ctxType != nil {
					vals = append(vals, reflect.ValueOf(ctx))
				}

				vals = append(vals, reflect.ValueOf(msg.Body()))

				returnValues := handler.method.Call(vals)
				if rerr := returnValues[0].Interface(); rerr != nil {
					return rerr.(error)
				}
				return nil
			}

			opts.Hooks.EachPrev(func(hook options.Hook) {
				if h, ok := hook.(HookSubHandler); ok {
					fn = h(fn)
				}
			})

			if n.wg != nil {
				n.wg.Add(1)
			}
			go func() {
				if n.wg != nil {
					defer n.wg.Done()
				}
				cerr := fn(ctx, &rpcMessage{
					topic:       sb.topic,
					contentType: ct,
					payload:     req.Interface(),
					header:      msg.Header,
				})
				results <- cerr
			}()
		}
		var errors []string
		for i := 0; i < len(sb.handlers); i++ {
			if rerr := <-results; rerr != nil {
				errors = append(errors, rerr.Error())
			}
		}
		if len(errors) > 0 {
			err = fmt.Errorf("subscriber error: %s", strings.Join(errors, "\n"))
		}
		return err
	}
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Subscriber() interface{} {
	return s.subscriber
}

func (s *subscriber) Options() SubscriberOptions {
	return s.opts
}

type subscriber struct {
	topic string

	typ        reflect.Type
	subscriber interface{}

	handlers []*handler

	rcvr reflect.Value
	opts SubscriberOptions
}

type handler struct {
	reqType reflect.Type
	ctxType reflect.Type
	method  reflect.Value
}
