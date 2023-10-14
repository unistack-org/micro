package broker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/semconv"
	maddr "go.unistack.org/micro/v4/util/addr"
	"go.unistack.org/micro/v4/util/id"
	mnet "go.unistack.org/micro/v4/util/net"
	"go.unistack.org/micro/v4/util/rand"
)

type MemoryBroker struct {
	subscribers map[string][]*memorySubscriber
	addr        string
	opts        Options
	sync.RWMutex
	connected bool
}

func (m *MemoryBroker) Options() Options {
	return m.opts
}

func (m *MemoryBroker) Address() string {
	return m.addr
}

func (m *MemoryBroker) Connect(ctx context.Context) error {
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

func (m *MemoryBroker) Disconnect(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if m.connected {
			m.connected = false
		}
	}

	return nil
}

func (m *MemoryBroker) Init(opts ...options.Option) error {
	var err error
	for _, o := range opts {
		if err = o(&m.opts); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryBroker) Publish(ctx context.Context, message interface{}, opts ...options.Option) error {
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
		var msgs []Message
		switch v := message.(type) {
		case []Message:
			msgs = v
		case Message:
			msgs = append(msgs, v)
		default:
			return ErrInvalidMessage
		}
		msgTopicMap := make(map[string][]*memoryMessage)
		for _, msg := range msgs {
			p := &memoryMessage{opts: options}
			p.topic, _ = msg.Header().Get(metadata.HeaderTopic)
			if v, ok := msg.Body().(*codec.Frame); ok {
				p.body = msg.Body()
			} else if len(m.opts.Codecs) == 0 {
				p.body = msg.Body()
			} else {
				cf, ok := m.opts.Codecs[options.ContentType]
				if !ok {
					return fmt.Errorf("%s: %s", codec.ErrUnknownContentType, options.ContentType)
				}
				p.body, err = cf.Marshal(v)
				if err != nil {
					return err
				}
			}
			msgTopicMap[p.topic] = append(msgTopicMap[p.topic], p)
		}

		eh := m.opts.ErrorHandler

		for t, ms := range msgTopicMap {
			ts := time.Now()

			m.opts.Meter.Counter(semconv.PublishMessageInflight, "endpoint", t).Add(len(ms))
			m.opts.Meter.Counter(semconv.SubscribeMessageInflight, "endpoint", t).Add(len(ms))

			m.RLock()
			subs, ok := m.subscribers[t]
			m.RUnlock()
			if !ok {
				m.opts.Meter.Counter(semconv.PublishMessageTotal, "endpoint", t, "status", "failure").Add(len(ms))
				m.opts.Meter.Counter(semconv.PublishMessageInflight, "endpoint", t).Add(-len(ms))
				m.opts.Meter.Counter(semconv.SubscribeMessageInflight, "endpoint", t).Add(-len(ms))
				continue
			}

			m.opts.Meter.Counter(semconv.PublishMessageTotal, "endpoint", t, "status", "success").Add(len(ms))
			for _, sub := range subs {
				if sub.opts.ErrorHandler != nil {
					eh = sub.opts.ErrorHandler
				}

				switch mh := sub.handler.(type) {
				case MessagesHandler:
					mhs := make([]Message, 0, len(ms))
					for _, m := range ms {
						mhs = append(mhs, m)
					}
					if err = mh(mhs); err != nil {
						m.opts.Meter.Counter(semconv.SubscribeMessageTotal, "endpoint", t, "status", "failure").Add(len(ms))
						if eh != nil {
							switch meh := eh.(type) {
							case MessagesHandler:
								_ = meh(mhs)
							case MessageHandler:
								for _, me := range mhs {
									_ = meh(me)
								}
							}
						} else if m.opts.Logger.V(logger.ErrorLevel) {
							m.opts.Logger.Error(m.opts.Context, err.Error())
						}
					}
				case MessageHandler:
					for _, p := range ms {
						if err = mh(p); err != nil {
							m.opts.Meter.Counter(semconv.SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
							if eh != nil {
								switch meh := eh.(type) {
								case MessageHandler:
									_ = meh(p)
								case MessagesHandler:
									_ = meh([]Message{p})
								}
							} else if m.opts.Logger.V(logger.ErrorLevel) {
								m.opts.Logger.Error(m.opts.Context, err.Error())
							}
						} else {
							if sub.opts.AutoAck {
								if err = p.Ack(); err != nil {
									m.opts.Logger.Error(m.opts.Context, "ack failed: "+err.Error())
									m.opts.Meter.Counter(semconv.SubscribeMessageTotal, "endpoint", t, "status", "failure").Inc()
								} else {
									m.opts.Meter.Counter(semconv.SubscribeMessageTotal, "endpoint", t, "status", "success").Inc()
								}
							} else {
								m.opts.Meter.Counter(semconv.SubscribeMessageTotal, "endpoint", t, "status", "success").Inc()
							}
						}
						m.opts.Meter.Counter(semconv.PublishMessageInflight, "endpoint", t).Add(-1)
						m.opts.Meter.Counter(semconv.SubscribeMessageInflight, "endpoint", t).Add(-1)
					}
				}
			}

			te := time.Since(ts)
			m.opts.Meter.Summary(semconv.PublishMessageLatencyMicroseconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Histogram(semconv.PublishMessageDurationSeconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Summary(semconv.SubscribeMessageLatencyMicroseconds, "endpoint", t).Update(te.Seconds())
			m.opts.Meter.Histogram(semconv.SubscribeMessageDurationSeconds, "endpoint", t).Update(te.Seconds())
		}

	}

	return nil
}

func (m *MemoryBroker) Subscribe(ctx context.Context, topic string, handler interface{}, opts ...options.Option) (Subscriber, error) {
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

func (m *MemoryBroker) String() string {
	return "memory"
}

func (m *MemoryBroker) Name() string {
	return m.opts.Name
}

type memoryMessage struct {
	err    error
	body   interface{}
	topic  string
	header metadata.Metadata
	opts   PublishOptions
	ctx    context.Context
}

func (m *memoryMessage) Topic() string {
	return m.topic
}

func (m *memoryMessage) Header() metadata.Metadata {
	return m.header
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
func NewBroker(opts ...options.Option) *MemoryBroker {
	return &MemoryBroker{
		opts:        NewOptions(opts...),
		subscribers: make(map[string][]*memorySubscriber),
	}
}
