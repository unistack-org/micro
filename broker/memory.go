package broker

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
	maddr "github.com/unistack-org/micro/v3/util/addr"
	mnet "github.com/unistack-org/micro/v3/util/net"
	"github.com/unistack-org/micro/v3/util/rand"
)

type memoryBroker struct {
	subscribers      map[string][]*memorySubscriber
	batchsubscribers map[string][]*memoryBatchSubscriber
	addr             string
	opts             Options
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
	ctx     context.Context
	exit    chan bool
	handler Handler
	id      string
	topic   string
	opts    SubscribeOptions
}

type memoryBatchSubscriber struct {
	ctx     context.Context
	exit    chan bool
	handler BatchHandler
	id      string
	topic   string
	opts    SubscribeOptions
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

func (m *memoryBroker) BatchPublish(ctx context.Context, msgs []*Message, opts ...PublishOption) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return ErrNotConnected
	}
	m.RUnlock()

	type msgWrapper struct {
		topic string
		body  interface{}
	}

	vs := make([]msgWrapper, 0, len(msgs))
	if m.opts.Codec == nil {
		m.RLock()
		for _, msg := range msgs {
			topic, _ := msg.Header.Get(metadata.HeaderTopic)
			vs = append(vs, msgWrapper{topic: topic, body: m})
		}
		m.RUnlock()
	} else {
		m.RLock()
		for _, msg := range msgs {
			topic, _ := msg.Header.Get(metadata.HeaderTopic)
			buf, err := m.opts.Codec.Marshal(msg)
			if err != nil {
				m.RUnlock()
				return err
			}
			vs = append(vs, msgWrapper{topic: topic, body: buf})
		}
		m.RUnlock()
	}

	if len(m.batchsubscribers) > 0 {
		eh := m.opts.BatchErrorHandler

		msgTopicMap := make(map[string]Events)
		for _, v := range vs {
			p := &memoryEvent{
				topic:   v.topic,
				message: v.body,
				opts:    m.opts,
			}
			msgTopicMap[p.topic] = append(msgTopicMap[p.topic], p)
		}

		for t, ms := range msgTopicMap {
			m.RLock()
			subs, ok := m.batchsubscribers[t]
			m.RUnlock()
			if !ok {
				continue
			}
			for _, sub := range subs {
				if err := sub.handler(ms); err != nil {
					ms.SetError(err)
					if sub.opts.BatchErrorHandler != nil {
						eh = sub.opts.BatchErrorHandler
					}
					if eh != nil {
						eh(ms)
					} else if m.opts.Logger.V(logger.ErrorLevel) {
						m.opts.Logger.Error(m.opts.Context, err.Error())
					}
				} else if sub.opts.AutoAck {
					if err := ms.Ack(); err != nil {
						m.opts.Logger.Errorf(m.opts.Context, "ack failed: %v", err)
					}
				}
			}
		}

	}

	eh := m.opts.ErrorHandler

	for _, v := range vs {
		p := &memoryEvent{
			topic:   v.topic,
			message: v.body,
			opts:    m.opts,
		}

		m.RLock()
		subs, ok := m.subscribers[p.topic]
		m.RUnlock()
		if !ok {
			continue
		}
		for _, sub := range subs {
			if err := sub.handler(p); err != nil {
				p.SetError(err)
				if sub.opts.ErrorHandler != nil {
					eh = sub.opts.ErrorHandler
				}
				if eh != nil {
					eh(p)
				} else if m.opts.Logger.V(logger.ErrorLevel) {
					m.opts.Logger.Error(m.opts.Context, err.Error())
				}
			} else if sub.opts.AutoAck {
				if err := p.Ack(); err != nil {
					m.opts.Logger.Errorf(m.opts.Context, "ack failed: %v", err)
				}
			}
		}
	}

	return nil
}

func (m *memoryBroker) Publish(ctx context.Context, topic string, msg *Message, opts ...PublishOption) error {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return ErrNotConnected
	}

	subs, ok := m.subscribers[topic]
	m.RUnlock()
	if !ok {
		return nil
	}

	var v interface{}
	if m.opts.Codec != nil {
		buf, err := m.opts.Codec.Marshal(msg)
		if err != nil {
			return err
		}
		v = buf
	} else {
		v = msg
	}

	p := &memoryEvent{
		topic:   topic,
		message: v,
		opts:    m.opts,
	}

	eh := m.opts.ErrorHandler

	for _, sub := range subs {
		if err := sub.handler(p); err != nil {
			p.err = err
			if sub.opts.ErrorHandler != nil {
				eh = sub.opts.ErrorHandler
			}
			if eh != nil {
				eh(p)
			} else if m.opts.Logger.V(logger.ErrorLevel) {
				m.opts.Logger.Error(m.opts.Context, err.Error())
			}
			continue
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

	options := NewSubscribeOptions(opts...)

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sub := &memoryBatchSubscriber{
		exit:    make(chan bool, 1),
		id:      id.String(),
		topic:   topic,
		handler: handler,
		opts:    options,
		ctx:     ctx,
	}

	m.Lock()
	m.batchsubscribers[topic] = append(m.batchsubscribers[topic], sub)
	m.Unlock()

	go func() {
		<-sub.exit
		m.Lock()
		var newSubscribers []*memoryBatchSubscriber
		for _, sb := range m.batchsubscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		m.batchsubscribers[topic] = newSubscribers
		m.Unlock()
	}()

	return sub, nil

	return nil, nil
}

func (m *memoryBroker) Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	m.RLock()
	if !m.connected {
		m.RUnlock()
		return nil, ErrNotConnected
	}
	m.RUnlock()

	options := NewSubscribeOptions(opts...)

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	sub := &memorySubscriber{
		exit:    make(chan bool, 1),
		id:      id.String(),
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
		var newSubscribers []*memorySubscriber
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

func (m *memoryBatchSubscriber) Options() SubscribeOptions {
	return m.opts
}

func (m *memoryBatchSubscriber) Topic() string {
	return m.topic
}

func (m *memoryBatchSubscriber) Unsubscribe(ctx context.Context) error {
	m.exit <- true
	return nil
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
		opts:             NewOptions(opts...),
		subscribers:      make(map[string][]*memorySubscriber),
		batchsubscribers: make(map[string][]*memoryBatchSubscriber),
	}
}
