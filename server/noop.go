package server

import "github.com/unistack-org/micro/v3/registry"

type noopServer struct {
	h    Handler
	opts Options
}

type noopHandler struct {
	opts HandlerOptions
	h    interface{}
}

type noopSubscriber struct {
	topic string
	opts  SubscriberOptions
	h     interface{}
}

func (n *noopSubscriber) Topic() string {
	return n.topic
}

func (n *noopSubscriber) Subscriber() interface{} {
	return n.h
}

func (n *noopSubscriber) Endpoints() []*registry.Endpoint {
	return nil
}

func (n *noopSubscriber) Options() SubscriberOptions {
	return n.opts
}

func (n *noopHandler) Endpoints() []*registry.Endpoint {
	return nil
}

func (n *noopHandler) Handler() interface{} {
	return nil
}

func (n *noopHandler) Options() HandlerOptions {
	return n.opts
}

func (n *noopHandler) Name() string {
	return "noop"
}

func (n *noopServer) Handle(handler Handler) error {
	n.h = handler
	return nil
}

func (n *noopServer) Subscribe(subscriber Subscriber) error {
	//	n.s = handler
	return nil
}

func (n *noopServer) NewHandler(h interface{}, opts ...HandlerOption) Handler {
	options := NewHandlerOptions()
	for _, o := range opts {
		o(&options)
	}
	return &noopHandler{opts: options, h: h}
}

func (n *noopServer) NewSubscriber(topic string, h interface{}, opts ...SubscriberOption) Subscriber {
	options := NewSubscriberOptions()
	for _, o := range opts {
		o(&options)
	}
	return &noopSubscriber{topic: topic, opts: options, h: h}
}

func (n *noopServer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return nil
}

func (n *noopServer) Start() error {
	return nil
}

func (n *noopServer) Stop() error {
	return nil
}

func (n *noopServer) Options() Options {
	return n.opts
}

func (n *noopServer) String() string {
	return "noop"
}

func newServer(opts ...Option) Server {
	options := NewOptions()
	for _, o := range opts {
		o(&options)
	}
	return &noopServer{opts: options}
}
