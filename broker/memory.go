package broker

import (
	"context"
	"sync"
	"time"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	maddr "go.unistack.org/micro/v4/util/addr"
	"go.unistack.org/micro/v4/util/id"
	mnet "go.unistack.org/micro/v4/util/net"
	"go.unistack.org/micro/v4/util/rand"
)

type memoryBroker struct {
	subscribers map[string][]*memorySubscriber
	addr        string
	opts        Options
	sync.RWMutex
	connected bool
}

type memoryEvent struct {
	err     error
	message interface{}
	topic   string
	opts    Options
}

type memorySubscriber struct {
	ctx          context.Context
	exit         chan bool
	handler      Handler
	batchhandler BatchHandler
	id           string
	topic        string
	opts         SubscribeOptions
}

func (m *memoryBroker) Options() Options {
	return m.opts
}

func (m *memoryBroker) Address() string {
	return m.addr
}

func (m *memoryBroker) Connect(ctx context.Context) error {
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
	m.Lock()
	defer m.Unlock()

	if !m.connected {
		return nil
	}

	m.connected = false
	return nil
}

func (m *memoryBroker) Init(opts ...Option) error {
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

func (m *memoryBroker) Publish(ctx context.Context, topic string, msg *Message, opts ...PublishOption) error {
	msg.Header.Set(metadata.HeaderTopic, topic)
	return m.publish(ctx, []*Message{msg}, opts...)
}

func (m *memoryBroker) BatchPublish(ctx context.Context, msgs []*Message, opts ...PublishOption) error {
	return m.publish(ctx, msgs, opts...)
}

func (m *memoryBroker) publish(ctx context.Context, msgs []*Message, opts ...PublishOption) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return ErrNotConnected
	}
	m.RUnlock()

	var err error

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		options := NewPublishOptions(opts...)

		msgTopicMap := make(map[string]Events)
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
			ts := time.Now()

			m.opts.Meter.Counter(PublishMessageInflight, "endpoint", t).Add(len(ms))
			m.opts.Meter.Counter(SubscribeMessageInflight, "endpoint", t).Add(len(ms))

			m.RLock()
			subs, ok := m.subscribers[t]
			m.RUnlock()
			if !ok {
				m.opts.Meter.Counter(PublishMessageTotal, "endpoint", t, "status", "failure").Add(len(ms))
				m.opts.Meter.Counter(PublishMessageInflight, "endpoint", t).Add(-len(ms))
				m.opts.Meter.Counter(SubscribeMessageInflight, "endpoint", t).Add(-len(ms))
				continue
			}

			m.opts.Meter.Counter(PublishMessageTotal, "endpoint", t, "status", "success").Add(len(ms))
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
						m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
						ms.SetError(err)
						if beh != nil {
							_ = beh(ms)
						} else if m.opts.Logger.V(logger.ErrorLevel) {
							m.opts.Logger.Error(m.opts.Context, err.Error())
						}
					} else {
						if sub.opts.AutoAck {
							if err = ms.Ack(); err != nil {
								m.opts.Logger.Errorf(m.opts.Context, "ack failed: %v", err)
								m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
							} else {
								m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "success").Inc()
							}
						} else {
							m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "success").Inc()
						}
					}
					m.opts.Meter.Counter(PublishMessageInflight, "endpoint", t).Add(-len(ms))
					m.opts.Meter.Counter(SubscribeMessageInflight, "endpoint", t).Add(-len(ms))
					// single processing
				case sub.handler != nil:
					for _, p := range ms {
						if err = sub.handler(p); err != nil {
							m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
							p.SetError(err)
							if eh != nil {
								_ = eh(p)
							} else if m.opts.Logger.V(logger.ErrorLevel) {
								m.opts.Logger.Error(m.opts.Context, err.Error())
							}
						} else {
							if sub.opts.AutoAck {
								if err = p.Ack(); err != nil {
									m.opts.Logger.Errorf(m.opts.Context, "ack failed: %v", err)
									m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
								} else {
									m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "success").Inc()
								}
							} else {
								m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "success").Inc()
							}
						}
						m.opts.Meter.Counter(PublishMessageInflight, "endpoint", t).Add(-1)
						m.opts.Meter.Counter(SubscribeMessageInflight, "endpoint", t).Add(-1)
					}
				}

			}
			te := time.Since(ts)
			m.opts.Meter.Summary(PublishMessageLatencyMicroseconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Histogram(PublishMessageDurationSeconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Summary(SubscribeMessageLatencyMicroseconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Histogram(SubscribeMessageDurationSeconds, "endpoint", t).Update(te.Seconds())
		}

	}

	return nil
}

func (m *memoryBroker) BatchSubscribe(ctx context.Context, topic string, handler BatchHandler, opts ...SubscribeOption) (Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, ErrNotConnected
	}
	m.RUnlock()

	sid, err := id.New()
	if err != nil {
		return nil, err
	}

	options := NewSubscribeOptions(opts...)

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

func (m *memoryBroker) Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, ErrNotConnected
	}
	m.RUnlock()

	sid, err := id.New()
	if err != nil {
		return nil, err
	}

	options := NewSubscribeOptions(opts...)

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

func (m *memoryEvent) Topic() string {
	return m.topic
}

func (m *memoryEvent) Message() *Message {
	switch v := m.message.(type) {
	case *Message:
		return v
	case []byte:
		msg := &Message{}
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

func (m *memorySubscriber) Options() SubscribeOptions {
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
func NewBroker(opts ...Option) Broker {
	return &memoryBroker{
		opts:        NewOptions(opts...),
		subscribers: make(map[string][]*memorySubscriber),
	}
}
