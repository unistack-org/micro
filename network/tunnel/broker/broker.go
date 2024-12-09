// Package broker is a tunnel broker
package broker

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/network/transport"
	"go.unistack.org/micro/v3/network/tunnel"
)

type tunBroker struct {
	tunnel tunnel.Tunnel
	opts   broker.Options
}

type tunSubscriber struct {
	listener tunnel.Listener
	handler  broker.Handler
	closed   chan bool
	topic    string
	opts     broker.SubscribeOptions
}

type tunBatchSubscriber struct {
	listener tunnel.Listener
	handler  broker.BatchHandler
	closed   chan bool
	topic    string
	opts     broker.SubscribeOptions
}

type tunEvent struct {
	err     error
	message *broker.Message
	topic   string
}

// used to access tunnel from options context
type (
	tunnelKey  struct{}
	tunnelAddr struct{}
)

func (t *tunBroker) Live() bool {
	return true
}

func (t *tunBroker) Ready() bool {
	return true
}

func (t *tunBroker) Health() bool {
	return true
}

func (t *tunBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tunBroker) Name() string {
	return t.opts.Name
}

func (t *tunBroker) Options() broker.Options {
	return t.opts
}

func (t *tunBroker) Address() string {
	return t.tunnel.Address()
}

func (t *tunBroker) Connect(ctx context.Context) error {
	return t.tunnel.Connect(ctx)
}

func (t *tunBroker) Disconnect(ctx context.Context) error {
	return t.tunnel.Close(ctx)
}

func (t *tunBroker) BatchPublish(ctx context.Context, msgs []*broker.Message, _ ...broker.PublishOption) error {
	// TODO: this is probably inefficient, we might want to just maintain an open connection
	// it may be easier to add broadcast to the tunnel
	topicMap := make(map[string]tunnel.Session)

	var err error
	for _, msg := range msgs {
		topic, _ := msg.Header.Get(metadata.HeaderTopic)
		c, ok := topicMap[topic]
		if !ok {
			c, err = t.tunnel.Dial(ctx, topic, tunnel.DialMode(tunnel.Multicast))
			if err != nil {
				return err
			}
			defer c.Close()
			topicMap[topic] = c
		}

		if err = c.Send(&transport.Message{
			Header: msg.Header,
			Body:   msg.Body,
		}); err != nil {
			//	msg.SetError(err)
			return err
		}
	}

	return nil
}

func (t *tunBroker) Publish(ctx context.Context, topic string, m *broker.Message, _ ...broker.PublishOption) error {
	// TODO: this is probably inefficient, we might want to just maintain an open connection
	// it may be easier to add broadcast to the tunnel
	c, err := t.tunnel.Dial(ctx, topic, tunnel.DialMode(tunnel.Multicast))
	if err != nil {
		return err
	}
	defer c.Close()

	return c.Send(&transport.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (t *tunBroker) BatchSubscribe(ctx context.Context, topic string, h broker.BatchHandler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	l, err := t.tunnel.Listen(ctx, topic, tunnel.ListenMode(tunnel.Multicast))
	if err != nil {
		return nil, err
	}

	tunSub := &tunBatchSubscriber{
		topic:    topic,
		handler:  h,
		opts:     broker.NewSubscribeOptions(opts...),
		closed:   make(chan bool),
		listener: l,
	}

	// start processing
	go tunSub.run()

	return tunSub, nil
}

func (t *tunBroker) Subscribe(ctx context.Context, topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	l, err := t.tunnel.Listen(ctx, topic, tunnel.ListenMode(tunnel.Multicast))
	if err != nil {
		return nil, err
	}

	tunSub := &tunSubscriber{
		topic:    topic,
		handler:  h,
		opts:     broker.NewSubscribeOptions(opts...),
		closed:   make(chan bool),
		listener: l,
	}

	// start processing
	go tunSub.run()

	return tunSub, nil
}

func (t *tunBroker) String() string {
	return "tunnel"
}

func (t *tunBatchSubscriber) run() {
	for {
		// accept a new connection
		c, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.closed:
				return
			default:
				continue
			}
		}

		// receive message
		m := new(transport.Message)
		if err := c.Recv(m); err != nil {
			if logger.DefaultLogger.V(logger.ErrorLevel) {
				logger.DefaultLogger.Error(t.opts.Context, err.Error(), err)
			}
			if err = c.Close(); err != nil {
				if logger.DefaultLogger.V(logger.ErrorLevel) {
					logger.DefaultLogger.Error(t.opts.Context, err.Error(), err)
				}
			}
			continue
		}

		// close the connection
		c.Close()

		evts := broker.Events{&tunEvent{
			topic: t.topic,
			message: &broker.Message{
				Header: m.Header,
				Body:   m.Body,
			},
		}}
		// handle the message
		go func() {
			_ = t.handler(evts)
		}()

	}
}

func (t *tunSubscriber) run() {
	for {
		// accept a new connection
		c, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.closed:
				return
			default:
				continue
			}
		}

		// receive message
		m := new(transport.Message)
		if err := c.Recv(m); err != nil {
			if logger.DefaultLogger.V(logger.ErrorLevel) {
				logger.DefaultLogger.Error(t.opts.Context, err.Error(), err)
			}
			if err = c.Close(); err != nil {
				if logger.DefaultLogger.V(logger.ErrorLevel) {
					logger.DefaultLogger.Error(t.opts.Context, err.Error(), err)
				}
			}
			continue
		}

		// close the connection
		c.Close()

		// handle the message
		go func() {
			_ = t.handler(&tunEvent{
				topic: t.topic,
				message: &broker.Message{
					Header: m.Header,
					Body:   m.Body,
				},
			})
		}()
	}
}

func (t *tunBatchSubscriber) Options() broker.SubscribeOptions {
	return t.opts
}

func (t *tunBatchSubscriber) Topic() string {
	return t.topic
}

func (t *tunBatchSubscriber) Unsubscribe(ctx context.Context) error {
	select {
	case <-t.closed:
		return nil
	default:
		close(t.closed)
		return t.listener.Close()
	}
}

func (t *tunSubscriber) Options() broker.SubscribeOptions {
	return t.opts
}

func (t *tunSubscriber) Topic() string {
	return t.topic
}

func (t *tunSubscriber) Unsubscribe(ctx context.Context) error {
	select {
	case <-t.closed:
		return nil
	default:
		close(t.closed)
		return t.listener.Close()
	}
}

func (t *tunEvent) Topic() string {
	return t.topic
}

func (t *tunEvent) Message() *broker.Message {
	return t.message
}

func (t *tunEvent) Ack() error {
	return nil
}

func (t *tunEvent) Error() error {
	return t.err
}

func (t *tunEvent) SetError(err error) {
	t.err = err
}

func (t *tunEvent) Context() context.Context {
	return context.TODO()
}

// NewBroker returns new tunnel broker
func NewBroker(opts ...broker.Option) (broker.Broker, error) {
	options := broker.NewOptions(opts...)

	t, ok := options.Context.Value(tunnelKey{}).(tunnel.Tunnel)
	if !ok {
		return nil, fmt.Errorf("tunnel not set")
	}

	a, ok := options.Context.Value(tunnelAddr{}).(string)
	if ok {
		// initialise address
		if err := t.Init(tunnel.Address(a)); err != nil {
			return nil, err
		}
	}

	if len(options.Addrs) > 0 {
		// initialise nodes
		if err := t.Init(tunnel.Nodes(options.Addrs...)); err != nil {
			return nil, err
		}
	}

	return &tunBroker{
		opts:   options,
		tunnel: t,
	}, nil
}

// WithAddress sets the tunnel address
func WithAddress(a string) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tunnelAddr{}, a)
	}
}

// WithTunnel sets the internal tunnel
func WithTunnel(t tunnel.Tunnel) broker.Option {
	return func(o *broker.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tunnelKey{}, t)
	}
}
