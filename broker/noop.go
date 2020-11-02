package broker

import "context"

type noopBroker struct {
	opts Options
}

type noopSubscriber struct {
	topic string
	opts  SubscribeOptions
}

// NewBroker returns new noop broker
func NewBroker(opts ...Option) Broker {
	return &noopBroker{opts: NewOptions(opts...)}
}

// Init initialize broker
func (n *noopBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	return nil
}

// Options returns broker Options
func (n *noopBroker) Options() Options {
	return n.opts
}

// Address returns broker address
func (n *noopBroker) Address() string {
	return ""
}

// Connect connects to broker
func (n *noopBroker) Connect(ctx context.Context) error {
	return nil
}

// Disconnect disconnects from broker
func (n *noopBroker) Disconnect(ctx context.Context) error {
	return nil
}

// Publish publishes message to broker
func (n *noopBroker) Publish(ctx context.Context, topic string, m *Message, opts ...PublishOption) error {
	return nil
}

// Subscribe subscribes to broker topic
func (n *noopBroker) Subscribe(ctx context.Context, topic string, h Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := NewSubscribeOptions(opts...)
	return &noopSubscriber{topic: topic, opts: options}, nil
}

// String return broker string representation
func (n *noopBroker) String() string {
	return "noop"
}

// Options returns subscriber options
func (n *noopSubscriber) Options() SubscribeOptions {
	return n.opts
}

// TOpic returns subscriber topic
func (n *noopSubscriber) Topic() string {
	return n.topic
}

// Unsubscribe unsbscribes from broker topic
func (n *noopSubscriber) Unsubscribe(ctx context.Context) error {
	return nil
}
