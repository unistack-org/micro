package broker

import (
	"context"
	"sync"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/options"
	maddr "go.unistack.org/micro/v3/util/addr"
	"go.unistack.org/micro/v3/util/id"
	mnet "go.unistack.org/micro/v3/util/net"
	"go.unistack.org/micro/v3/util/rand"
)

type memoryBroker struct {
	funcPublish        broker.FuncPublish
	funcBatchPublish   broker.FuncBatchPublish
	funcSubscribe      broker.FuncSubscribe
	funcBatchSubscribe broker.FuncBatchSubscribe
	subscribers        map[string][]*memorySubscriber
	addr               string
	opts               broker.Options
	sync.RWMutex
	connected bool
}

type memoryEvent struct {
	err     error
	message interface{}
	topic   string
	opts    broker.Options
}

type memorySubscriber struct {
	ctx          context.Context
	exit         chan bool
	handler      broker.Handler
	batchhandler broker.BatchHandler
	id           string
	topic        string
	opts         broker.SubscribeOptions
}

func (m *memoryBroker) Options() broker.Options {
	return m.opts
}

func (m *memoryBroker) Address() string {
	return m.addr
}

func (m *memoryBroker) Connect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.Lock()
	defer m.Unlock()

	if m.connected {
		return nil
	}

	// use 127.0.0.1 to avoid scan of all network interfaces
	addr, err := maddr.Extract("127.0.0.1")
	if err != nil {
		return err
	}
	var rng rand.Rand
	i := rng.Intn(20000)
	// set addr with port
	addr = mnet.HostPort(addr, 10000+i)

	m.addr = addr
	m.connected = true

	return nil
}

func (m *memoryBroker) Disconnect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.Lock()
	defer m.Unlock()

	if !m.connected {
		return nil
	}

	m.connected = false
	return nil
}

func (m *memoryBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&m.opts)
	}

	m.funcPublish = m.fnPublish
	m.funcBatchPublish = m.fnBatchPublish
	m.funcSubscribe = m.fnSubscribe
	m.funcBatchSubscribe = m.fnBatchSubscribe

	m.opts.Hooks.EachNext(func(hook options.Hook) {
		switch h := hook.(type) {
		case broker.HookPublish:
			m.funcPublish = h(m.funcPublish)
		case broker.HookBatchPublish:
			m.funcBatchPublish = h(m.funcBatchPublish)
		case broker.HookSubscribe:
			m.funcSubscribe = h(m.funcSubscribe)
		case broker.HookBatchSubscribe:
			m.funcBatchSubscribe = h(m.funcBatchSubscribe)
		}
	})

	return nil
}

func (m *memoryBroker) Publish(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	return m.funcPublish(ctx, topic, msg, opts...)
}

func (m *memoryBroker) fnPublish(ctx context.Context, topic string, msg *broker.Message, opts ...broker.PublishOption) error {
	msg.Header.Set(metadata.HeaderTopic, topic)
	return m.publish(ctx, []*broker.Message{msg}, opts...)
}

func (m *memoryBroker) BatchPublish(ctx context.Context, msgs []*broker.Message, opts ...broker.PublishOption) error {
	return m.funcBatchPublish(ctx, msgs, opts...)
}

func (m *memoryBroker) fnBatchPublish(ctx context.Context, msgs []*broker.Message, opts ...broker.PublishOption) error {
	return m.publish(ctx, msgs, opts...)
}

func (m *memoryBroker) publish(ctx context.Context, msgs []*broker.Message, opts ...broker.PublishOption) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return broker.ErrNotConnected
	}
	m.RUnlock()

	var err error

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		options := broker.NewPublishOptions(opts...)

		msgTopicMap := make(map[string]broker.Events)
		for _, v := range msgs {
			p := &memoryEvent{opts: m.opts}

			if m.opts.Codec == nil || options.BodyOnly {
				p.topic, _ = v.Header.Get(metadata.HeaderTopic)
				p.message = v.Body
			} else {
				p.topic, _ = v.Header.Get(metadata.HeaderTopic)
				p.message, err = m.opts.Codec.Marshal(v)
				if err != nil {
					return err
				}
			}
			msgTopicMap[p.topic] = append(msgTopicMap[p.topic], p)
		}

		beh := m.opts.BatchErrorHandler
		eh := m.opts.ErrorHandler

		for t, ms := range msgTopicMap {
			m.RLock()
			subs, ok := m.subscribers[t]
			m.RUnlock()
			if !ok {
				continue
			}

			for _, sub := range subs {
				if sub.opts.BatchErrorHandler != nil {
					beh = sub.opts.BatchErrorHandler
				}
				if sub.opts.ErrorHandler != nil {
					eh = sub.opts.ErrorHandler
				}

				switch {
				// batch processing
				case sub.batchhandler != nil:
					if err = sub.batchhandler(ms); err != nil {
						ms.SetError(err)
						if beh != nil {
							_ = beh(ms)
						} else if m.opts.Logger.V(logger.ErrorLevel) {
							m.opts.Logger.Error(m.opts.Context, err.Error())
						}
					} else if sub.opts.AutoAck {
						if err = ms.Ack(); err != nil {
							m.opts.Logger.Error(m.opts.Context, "broker ack error", err)
						}
					}
					// single processing
				case sub.handler != nil:
					for _, p := range ms {
						if err = sub.handler(p); err != nil {
							p.SetError(err)
							if eh != nil {
								_ = eh(p)
							} else if m.opts.Logger.V(logger.ErrorLevel) {
								m.opts.Logger.Error(m.opts.Context, "broker handler error", err)
							}
						} else if sub.opts.AutoAck {
							if err = p.Ack(); err != nil {
								m.opts.Logger.Error(m.opts.Context, "broker ack error", err)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (m *memoryBroker) BatchSubscribe(ctx context.Context, topic string, handler broker.BatchHandler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return m.funcBatchSubscribe(ctx, topic, handler, opts...)
}

func (m *memoryBroker) fnBatchSubscribe(ctx context.Context, topic string, handler broker.BatchHandler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, broker.ErrNotConnected
	}
	m.RUnlock()

	sid, err := id.New()
	if err != nil {
		return nil, err
	}

	options := broker.NewSubscribeOptions(opts...)

	sub := &memorySubscriber{
		exit:         make(chan bool, 1),
		id:           sid,
		topic:        topic,
		batchhandler: handler,
		opts:         options,
		ctx:          ctx,
	}

	m.Lock()
	m.subscribers[topic] = append(m.subscribers[topic], sub)
	m.Unlock()

	go func() {
		<-sub.exit
		m.Lock()
		newSubscribers := make([]*memorySubscriber, 0, len(m.subscribers)-1)
		for _, sb := range m.subscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		m.subscribers[topic] = newSubscribers
		m.Unlock()
	}()

	return sub, nil
}

func (m *memoryBroker) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return m.funcSubscribe(ctx, topic, handler, opts...)
}

func (m *memoryBroker) fnSubscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, broker.ErrNotConnected
	}
	m.RUnlock()

	sid, err := id.New()
	if err != nil {
		return nil, err
	}

	options := broker.NewSubscribeOptions(opts...)

	sub := &memorySubscriber{
		exit:    make(chan bool, 1),
		id:      sid,
		topic:   topic,
		handler: handler,
		opts:    options,
		ctx:     ctx,
	}

	m.Lock()
	m.subscribers[topic] = append(m.subscribers[topic], sub)
	m.Unlock()

	go func() {
		<-sub.exit
		m.Lock()
		newSubscribers := make([]*memorySubscriber, 0, len(m.subscribers)-1)
		for _, sb := range m.subscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		m.subscribers[topic] = newSubscribers
		m.Unlock()
	}()

	return sub, nil
}

func (m *memoryBroker) String() string {
	return "memory"
}

func (m *memoryBroker) Name() string {
	return m.opts.Name
}

func (m *memoryBroker) Live() bool {
	return true
}

func (m *memoryBroker) Ready() bool {
	return true
}

func (m *memoryBroker) Health() bool {
	return true
}

func (m *memoryEvent) Topic() string {
	return m.topic
}

func (m *memoryEvent) Message() *broker.Message {
	switch v := m.message.(type) {
	case *broker.Message:
		return v
	case []byte:
		msg := &broker.Message{}
		if err := m.opts.Codec.Unmarshal(v, msg); err != nil {
			if m.opts.Logger.V(logger.ErrorLevel) {
				m.opts.Logger.Error(m.opts.Context, "[memory]: failed to unmarshal: %v", err)
			}
			return nil
		}
		return msg
	}

	return nil
}

func (m *memoryEvent) Ack() error {
	return nil
}

func (m *memoryEvent) Error() error {
	return m.err
}

func (m *memoryEvent) SetError(err error) {
	m.err = err
}

func (m *memoryEvent) Context() context.Context {
	return m.opts.Context
}

func (m *memorySubscriber) Options() broker.SubscribeOptions {
	return m.opts
}

func (m *memorySubscriber) Topic() string {
	return m.topic
}

func (m *memorySubscriber) Unsubscribe(ctx context.Context) error {
	m.exit <- true
	return nil
}

// NewBroker return new memory broker
func NewBroker(opts ...broker.Option) broker.Broker {
	return &memoryBroker{
		opts:        broker.NewOptions(opts...),
		subscribers: make(map[string][]*memorySubscriber),
	}
}
