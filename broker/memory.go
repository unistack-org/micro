//go:build ignore

package broker

import (
	"context"
	"sync"
	"time"

	"go.unistack.org/micro/v4/logger"
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

func (m *memoryBroker) NewMessage(endpoint string, req interface{}, opts ...MessageOption) Message {
	return &memoryMessage{}
}

func (m *memoryBroker) Publish(ctx context.Context, message interface{}, opts ...PublishOption) error {
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
		var msgs []*memoryMessage
		switch v := message.(type) {
		case *memoryMessage:
			msgs = []*memoryMessage{v}
		case []*memoryMessage:
			msgs = v
		default:
			return ErrInvalidMessage
		}
		msgTopicMap := make(map[string][]*memoryMessage)
		for _, msg := range msgs {
			p := &memoryMessage{opts: options}
			/*
				if mb, ok := msg.Body().(*codec.Frame); ok {
					p.message = v.Body
				} else {
					p.topic, _ = v.Header.Get(metadata.HeaderTopic)
					p.message, err = m.opts.Codec.Marshal(v)
					if err != nil {
						return err
					}
				}
			*/
			msgTopicMap[msg.Topic()] = append(msgTopicMap[p.topic], p)
		}

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
				if sub.opts.ErrorHandler != nil {
					eh = sub.opts.ErrorHandler
				}

				for _, p := range ms {
					if err = sub.handler(p); err != nil {
						m.opts.Meter.Counter(SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
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
			te := time.Since(ts)
			m.opts.Meter.Summary(PublishMessageLatencyMicroseconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Histogram(PublishMessageDurationSeconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Summary(SubscribeMessageLatencyMicroseconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Histogram(SubscribeMessageDurationSeconds, "endpoint", t).Update(te.Seconds())
		}

	}

	return nil
}

func (m *memoryBroker) Subscribe(ctx context.Context, topic string, handler interface{}, opts ...SubscribeOption) (Subscriber, error) {
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

type memoryMessage struct {
	err   error
	body  interface{}
	topic string
	opts  PublishOptions
	ctx   context.Context
}

func (m *memoryMessage) Topic() string {
	return m.topic
}

func (m *memoryMessage) Body() interface{} {
	return m.body
}

func (m *memoryMessage) Ack() error {
	return nil
}

func (m *memoryMessage) Error() error {
	return m.err
}

func (m *memoryMessage) Context() context.Context {
	return m.ctx
}

type memorySubscriber struct {
	ctx     context.Context
	exit    chan bool
	handler interface{}
	id      string
	topic   string
	opts    SubscribeOptions
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
func NewBroker(opts ...Option) *memoryBroker {
	return &memoryBroker{
		opts:        NewOptions(opts...),
		subscribers: make(map[string][]*memorySubscriber),
	}
}
