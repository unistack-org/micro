package broker

import "context"

type NoopBroker struct {
	opts Options
}

type noopSubscriber struct {
	topic string
	opts  SubscribeOptions
}

func (n *NoopBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	return nil
}

func (n *NoopBroker) Options() Options {
	return n.opts
}

func (n *NoopBroker) Address() string {
	return ""
}

func (n *NoopBroker) Connect(ctx context.Context) error {
	return nil
}

func (n *NoopBroker) Disconnect(ctx context.Context) error {
	return nil
}

func (n *NoopBroker) Publish(ctx context.Context, topic string, m *Message, opts ...PublishOption) error {
	return nil
}

func (n *NoopBroker) Subscribe(ctx context.Context, topic string, h Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := NewSubscribeOptions()

	for _, o := range opts {
		o(&options)
	}

	return &noopSubscriber{topic: topic, opts: options}, nil
}

func (n *NoopBroker) String() string {
	return "noop"
}

func (n *noopSubscriber) Options() SubscribeOptions {
	return n.opts
}

func (n *noopSubscriber) Topic() string {
	return n.topic
}

func (n *noopSubscriber) Unsubscribe(ctx context.Context) error {
	return nil
}
