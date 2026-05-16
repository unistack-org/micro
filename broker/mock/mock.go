package mock

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/metadata"
)

var _ broker.Broker = (*MockBroker)(nil)

// expectation is the common interface for all expected broker calls.
type expectation interface {
	fulfilled() bool
	Lock()
	Unlock()
	String() string
}

type commonExpectation struct {
	sync.Mutex
	triggered bool
	err       error
}

func (e *commonExpectation) fulfilled() bool { return e.triggered }

// ExpectedConnect holds an expectation for broker.Connect.
type ExpectedConnect struct {
	commonExpectation
}

// WillReturnError configures the error to return from Connect.
func (e *ExpectedConnect) WillReturnError(err error) *ExpectedConnect {
	e.err = err
	return e
}

func (e *ExpectedConnect) String() string {
	msg := "ExpectedConnect => expecting broker.Connect"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedDisconnect holds an expectation for broker.Disconnect.
type ExpectedDisconnect struct {
	commonExpectation
}

// WillReturnError configures the error to return from Disconnect.
func (e *ExpectedDisconnect) WillReturnError(err error) *ExpectedDisconnect {
	e.err = err
	return e
}

func (e *ExpectedDisconnect) String() string {
	msg := "ExpectedDisconnect => expecting broker.Disconnect"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedPublish holds an expectation for broker.Publish.
type ExpectedPublish struct {
	commonExpectation
	topic string
	delay time.Duration
}

// WillReturnError configures the error to return from Publish.
func (e *ExpectedPublish) WillReturnError(err error) *ExpectedPublish {
	e.err = err
	return e
}

// WillDelayFor configures a delay before Publish returns.
func (e *ExpectedPublish) WillDelayFor(d time.Duration) *ExpectedPublish {
	e.delay = d
	return e
}

func (e *ExpectedPublish) String() string {
	msg := fmt.Sprintf("ExpectedPublish => expecting broker.Publish to topic %q", e.topic)
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedSubscribe holds an expectation for broker.Subscribe.
type ExpectedSubscribe struct {
	commonExpectation
	topic string
}

// WillReturnError configures the error to return from Subscribe.
func (e *ExpectedSubscribe) WillReturnError(err error) *ExpectedSubscribe {
	e.err = err
	return e
}

func (e *ExpectedSubscribe) String() string {
	msg := fmt.Sprintf("ExpectedSubscribe => expecting broker.Subscribe to topic %q", e.topic)
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

// ExpectedUnsubscribe holds an expectation for Subscriber.Unsubscribe.
type ExpectedUnsubscribe struct {
	commonExpectation
	topic string
}

// WillReturnError configures the error to return from Unsubscribe.
func (e *ExpectedUnsubscribe) WillReturnError(err error) *ExpectedUnsubscribe {
	e.err = err
	return e
}

func (e *ExpectedUnsubscribe) String() string {
	msg := fmt.Sprintf("ExpectedUnsubscribe => expecting Subscriber.Unsubscribe from topic %q", e.topic)
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	return msg
}

type handlerEntry struct {
	id      uint64
	handler any
}

// MockBroker is a mock implementation of broker.Broker for use in tests.
//
// Declare expected calls with ExpectConnect, ExpectPublish, etc. before exercising
// the code under test. After the test, call ExpectationsWereMet to verify all
// declared expectations were fulfilled.
//
// To test subscriber handlers, call InjectMessage after Subscribe.
type MockBroker struct {
	opts      broker.Options
	mu        sync.Mutex
	expected  []expectation
	handlers  map[string][]*handlerEntry
	connected bool
	nextID    uint64
}

// NewMockBroker creates a new MockBroker.
func NewMockBroker(opts ...broker.Option) *MockBroker {
	return &MockBroker{
		opts:     broker.NewOptions(opts...),
		handlers: make(map[string][]*handlerEntry),
	}
}

// ExpectConnect registers an expectation that broker.Connect will be called.
func (m *MockBroker) ExpectConnect() *ExpectedConnect {
	e := &ExpectedConnect{}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectDisconnect registers an expectation that broker.Disconnect will be called.
func (m *MockBroker) ExpectDisconnect() *ExpectedDisconnect {
	e := &ExpectedDisconnect{}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectPublish registers an expectation that broker.Publish will be called for topic.
func (m *MockBroker) ExpectPublish(topic string) *ExpectedPublish {
	e := &ExpectedPublish{topic: topic}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectSubscribe registers an expectation that broker.Subscribe will be called for topic.
func (m *MockBroker) ExpectSubscribe(topic string) *ExpectedSubscribe {
	e := &ExpectedSubscribe{topic: topic}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectUnsubscribe registers an expectation that Subscriber.Unsubscribe will be called for topic.
func (m *MockBroker) ExpectUnsubscribe(topic string) *ExpectedUnsubscribe {
	e := &ExpectedUnsubscribe{topic: topic}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectationsWereMet returns an error if any declared expectation was not fulfilled.
func (m *MockBroker) ExpectationsWereMet() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		e.Lock()
		fulfilled := e.fulfilled()
		e.Unlock()
		if !fulfilled {
			return fmt.Errorf("there is a remaining expectation which was not matched: %s", e)
		}
	}
	return nil
}

// InjectMessage delivers msg to all handlers subscribed to topic.
// Supports both func(broker.Message) error and func([]broker.Message) error.
// Returns the first handler error encountered.
func (m *MockBroker) InjectMessage(ctx context.Context, topic string, msg broker.Message) error {
	m.mu.Lock()
	entries := make([]*handlerEntry, len(m.handlers[topic]))
	copy(entries, m.handlers[topic])
	m.mu.Unlock()

	for _, entry := range entries {
		switch fn := entry.handler.(type) {
		case func(broker.Message) error:
			if err := fn(msg); err != nil {
				return err
			}
		case func([]broker.Message) error:
			if err := fn([]broker.Message{msg}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MockBroker) newCodec(ct string) (codec.Codec, error) {
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}
	m.mu.Lock()
	c, ok := m.opts.Codecs[ct]
	m.mu.Unlock()
	if ok {
		return c, nil
	}
	return nil, codec.ErrUnknownContentType
}

// Connect implements broker.Broker.
func (m *MockBroker) Connect(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, e := range m.expected {
		ec, ok := e.(*ExpectedConnect)
		if !ok {
			continue
		}
		ec.Lock()
		if ec.triggered {
			ec.Unlock()
			continue
		}
		ec.triggered = true
		err := ec.err
		ec.Unlock()
		if err == nil {
			m.connected = true
		}
		return err
	}
	return fmt.Errorf("unexpected call to broker.Connect")
}

// Disconnect implements broker.Broker.
func (m *MockBroker) Disconnect(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, e := range m.expected {
		ec, ok := e.(*ExpectedDisconnect)
		if !ok {
			continue
		}
		ec.Lock()
		if ec.triggered {
			ec.Unlock()
			continue
		}
		ec.triggered = true
		err := ec.err
		ec.Unlock()
		if err == nil {
			m.connected = false
		}
		return err
	}
	return fmt.Errorf("unexpected call to broker.Disconnect")
}

// Publish implements broker.Broker. Records the call without delivering to handlers.
func (m *MockBroker) Publish(_ context.Context, topic string, _ ...broker.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return broker.ErrNotConnected
	}

	for _, e := range m.expected {
		ep, ok := e.(*ExpectedPublish)
		if !ok {
			continue
		}
		ep.Lock()
		if ep.triggered || ep.topic != topic {
			ep.Unlock()
			continue
		}
		ep.triggered = true
		delay := ep.delay
		err := ep.err
		ep.Unlock()
		if delay > 0 {
			time.Sleep(delay)
		}
		return err
	}
	return fmt.Errorf("unexpected call to broker.Publish for topic %q", topic)
}

// Subscribe implements broker.Broker. Registers handler for later use with InjectMessage.
func (m *MockBroker) Subscribe(_ context.Context, topic string, handler any, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return nil, broker.ErrNotConnected
	}

	for _, e := range m.expected {
		es, ok := e.(*ExpectedSubscribe)
		if !ok {
			continue
		}
		es.Lock()
		if es.triggered || es.topic != topic {
			es.Unlock()
			continue
		}
		es.triggered = true
		err := es.err
		es.Unlock()
		if err != nil {
			return nil, err
		}
		m.nextID++
		id := m.nextID
		m.handlers[topic] = append(m.handlers[topic], &handlerEntry{id: id, handler: handler})
		return &MockSubscriber{
			topic:     topic,
			opts:      broker.NewSubscribeOptions(opts...),
			mock:      m,
			handlerID: id,
		}, nil
	}
	return nil, fmt.Errorf("unexpected call to broker.Subscribe for topic %q", topic)
}

// NewMessage implements broker.Broker. Marshals body using a codec configured via broker.Codec option.
func (m *MockBroker) NewMessage(ctx context.Context, hdr metadata.Metadata, body any, opts ...broker.MessageOption) (broker.Message, error) {
	options := broker.NewMessageOptions(opts...)
	if options.ContentType == "" {
		options.ContentType = m.opts.ContentType
	}
	c, err := m.newCodec(options.ContentType)
	if err != nil {
		return nil, err
	}
	b, err := c.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &mockMessage{ctx: ctx, hdr: hdr, body: b, c: c}, nil
}

// Init implements broker.Broker.
func (m *MockBroker) Init(opts ...broker.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

// Options implements broker.Broker.
func (m *MockBroker) Options() broker.Options { return m.opts }

// Name implements broker.Broker.
func (m *MockBroker) Name() string { return m.opts.Name }

// String implements broker.Broker.
func (m *MockBroker) String() string { return "mock" }

// Address implements broker.Broker.
func (m *MockBroker) Address() string { return "" }

// Live implements broker.Broker.
func (m *MockBroker) Live() bool { return true }

// Ready implements broker.Broker.
func (m *MockBroker) Ready() bool { return true }

// Health implements broker.Broker.
func (m *MockBroker) Health() bool { return true }

// MockSubscriber is returned by MockBroker.Subscribe.
type MockSubscriber struct {
	topic     string
	opts      broker.SubscribeOptions
	mock      *MockBroker
	handlerID uint64
}

// Options implements broker.Subscriber.
func (s *MockSubscriber) Options() broker.SubscribeOptions { return s.opts }

// Topic implements broker.Subscriber.
func (s *MockSubscriber) Topic() string { return s.topic }

// Unsubscribe implements broker.Subscriber.
func (s *MockSubscriber) Unsubscribe(_ context.Context) error {
	s.mock.mu.Lock()
	defer s.mock.mu.Unlock()

	for _, e := range s.mock.expected {
		eu, ok := e.(*ExpectedUnsubscribe)
		if !ok {
			continue
		}
		eu.Lock()
		if eu.triggered || eu.topic != s.topic {
			eu.Unlock()
			continue
		}
		eu.triggered = true
		err := eu.err
		eu.Unlock()
		if err != nil {
			return err
		}
		handlers := s.mock.handlers[s.topic]
		newHandlers := make([]*handlerEntry, 0, len(handlers))
		for _, h := range handlers {
			if h.id != s.handlerID {
				newHandlers = append(newHandlers, h)
			}
		}
		s.mock.handlers[s.topic] = newHandlers
		return nil
	}
	return fmt.Errorf("unexpected call to Subscriber.Unsubscribe for topic %q", s.topic)
}
