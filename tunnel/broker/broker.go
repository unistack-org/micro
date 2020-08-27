// Package broker is a tunnel broker
package broker

import (
	"context"
	"fmt"

	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/transport"
	"github.com/unistack-org/micro/v3/tunnel"
)

type tunBroker struct {
	opts   broker.Options
	tunnel tunnel.Tunnel
}

type tunSubscriber struct {
	topic   string
	handler broker.Handler
	opts    broker.SubscribeOptions

	closed   chan bool
	listener tunnel.Listener
}

type tunEvent struct {
	topic   string
	message *broker.Message
}

// used to access tunnel from options context
type tunnelKey struct{}
type tunnelAddr struct{}

func (t *tunBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *tunBroker) Options() broker.Options {
	return t.opts
}

func (t *tunBroker) Address() string {
	return t.tunnel.Address()
}

func (t *tunBroker) Connect() error {
	return t.tunnel.Connect()
}

func (t *tunBroker) Disconnect() error {
	return t.tunnel.Close()
}

func (t *tunBroker) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	// TODO: this is probably inefficient, we might want to just maintain an open connection
	// it may be easier to add broadcast to the tunnel
	c, err := t.tunnel.Dial(topic, tunnel.DialMode(tunnel.Multicast))
	if err != nil {
		return err
	}
	defer c.Close()

	return c.Send(&transport.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (t *tunBroker) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	l, err := t.tunnel.Listen(topic, tunnel.ListenMode(tunnel.Multicast))
	if err != nil {
		return nil, err
	}

	var options broker.SubscribeOptions
	for _, o := range opts {
		o(&options)
	}

	tunSub := &tunSubscriber{
		topic:    topic,
		handler:  h,
		opts:     options,
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
			logger.Error(err)
			if err = c.Close(); err != nil {
				logger.Error(err)
			}
			continue
		}

		// close the connection
		c.Close()

		// handle the message
		go t.handler(&tunEvent{
			topic: t.topic,
			message: &broker.Message{
				Header: m.Header,
				Body:   m.Body,
			},
		})
	}
}

func (t *tunSubscriber) Options() broker.SubscribeOptions {
	return t.opts
}

func (t *tunSubscriber) Topic() string {
	return t.topic
}

func (t *tunSubscriber) Unsubscribe() error {
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
	return nil
}

func NewBroker(opts ...broker.Option) (broker.Broker, error) {
	options := broker.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}

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
