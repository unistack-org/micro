package broker

type noopBroker struct {
	opts Options
}

type noopSubscriber struct {
	topic string
	opts  SubscribeOptions
}

func (n *noopBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	return nil
}

func (n *noopBroker) Options() Options {
	return n.opts
}

func (n *noopBroker) Address() string {
	return ""
}

func (n *noopBroker) Connect() error {
	return nil
}

func (n *noopBroker) Disconnect() error {
	return nil
}

func (n *noopBroker) Publish(topic string, m *Message, opts ...PublishOption) error {
	return nil
}

func (n *noopBroker) Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error) {
	options := NewSubscribeOptions()

	for _, o := range opts {
		o(&options)
	}

	return &noopSubscriber{topic: topic, opts: options}, nil
}

func (n *noopBroker) String() string {
	return "noop"
}

func (n *noopSubscriber) Options() SubscribeOptions {
	return n.opts
}

func (n *noopSubscriber) Topic() string {
	return n.topic
}

func (n *noopSubscriber) Unsubscribe() error {
	return nil
}

// newBroker returns a new noop broker
func newBroker(opts ...Option) Broker {
	options := NewOptions()

	for _, o := range opts {
		o(&options)
	}
	return &noopBroker{opts: options}
}
