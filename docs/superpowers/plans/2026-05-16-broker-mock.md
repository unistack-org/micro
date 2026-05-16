# Broker Mock Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a strict expectation-based mock broker at `go.unistack.org/micro/v5/broker/mock` for unit testing application code that uses the broker without real message transport.

**Architecture:** Follows the MockClient pattern in `broker/mock/`. All broker calls (`Connect`, `Disconnect`, `Publish`, `Subscribe`, `Unsubscribe`) must have a pre-declared expectation; unexpected calls return an error. `Publish` records the call without delivering to handlers. Handlers are called only via explicit `InjectMessage`. `ExpectationsWereMet` verifies all declared expectations were fulfilled.

**Tech Stack:** Go 1.26, `go.unistack.org/micro/v5/broker`, `sync`, `fmt`, `time`.

---

## File Layout

| File | Responsibility |
|---|---|
| `broker/mock/mock_message.go` | `mockMessage` (implements `broker.Message`), `NewMockMessage` factory |
| `broker/mock/mock.go` | `MockBroker`, `MockSubscriber`, all `Expected*` types, `expectation` interface |
| `broker/mock/mock_test.go` | All tests (package `mock_test`) |
| `broker/mock/example_test.go` | Godoc `Example_*` functions (package `mock_test`) |

---

### Task 1: mockMessage

**Files:**
- Create: `broker/mock/mock_message.go`
- Create: `broker/mock/mock_test.go`

- [ ] **Step 1.1: Write the failing test**

Create `broker/mock/mock_test.go`:

```go
package mock_test

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/broker/mock"
	"go.unistack.org/micro/v5/metadata"
)

func TestNewMockMessage(t *testing.T) {
	ctx := context.Background()
	hdr := metadata.Metadata{"key": "val"}
	body := []byte(`{"id":1}`)

	msg := mock.NewMockMessage(ctx, "orders", hdr, body)

	if msg.Topic() != "orders" {
		t.Fatalf("want topic %q, got %q", "orders", msg.Topic())
	}
	if string(msg.Body()) != string(body) {
		t.Fatalf("want body %q, got %q", body, msg.Body())
	}
	if msg.Context() != ctx {
		t.Fatal("context mismatch")
	}
	if v, ok := msg.Header()["key"]; !ok || v != "val" {
		t.Fatal("header mismatch")
	}
	if err := msg.Ack(); err != nil {
		t.Fatalf("Ack: %v", err)
	}
	if err := msg.Error(); err != nil {
		t.Fatalf("Error: %v", err)
	}
}
```

- [ ] **Step 1.2: Run test to verify it fails**

```bash
cd broker/mock && go test ./... -run TestNewMockMessage -v
```

Expected: compile error — package does not exist yet.

- [ ] **Step 1.3: Create mock_message.go**

```go
// Package mock provides a mock implementation of broker.Broker for testing.
package mock

import (
	"context"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/metadata"
)

type mockMessage struct {
	ctx   context.Context
	topic string
	hdr   metadata.Metadata
	body  []byte
	c     codec.Codec
	err   error
}

func (m *mockMessage) Context() context.Context  { return m.ctx }
func (m *mockMessage) Topic() string             { return m.topic }
func (m *mockMessage) Header() metadata.Metadata { return m.hdr }
func (m *mockMessage) Body() []byte              { return m.body }
func (m *mockMessage) Ack() error                { return nil }
func (m *mockMessage) Error() error              { return m.err }

func (m *mockMessage) Unmarshal(dst any, opts ...codec.Option) error {
	if m.c == nil {
		return codec.ErrUnknownContentType
	}
	return m.c.Unmarshal(m.body, dst, opts...)
}

// NewMockMessage creates a broker.Message with pre-marshaled body for use with InjectMessage.
// Optionally pass a codec to enable Unmarshal in the handler under test.
func NewMockMessage(ctx context.Context, topic string, hdr metadata.Metadata, body []byte, c ...codec.Codec) broker.Message {
	msg := &mockMessage{ctx: ctx, topic: topic, hdr: hdr, body: body}
	if len(c) > 0 {
		msg.c = c[0]
	}
	return msg
}
```

- [ ] **Step 1.4: Run test to verify it passes**

```bash
cd broker/mock && go test ./... -run TestNewMockMessage -v
```

Expected: PASS

- [ ] **Step 1.5: Commit**

```bash
git add broker/mock/mock_message.go broker/mock/mock_test.go
git commit -m "feat(broker/mock): add mockMessage and NewMockMessage"
```

---

### Task 2: Full MockBroker implementation

**Files:**
- Create: `broker/mock/mock.go`

Write the complete implementation first (interface requires all methods), tests follow in Tasks 3–7.

- [ ] **Step 2.1: Create mock.go**

```go
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
```

- [ ] **Step 2.2: Verify it compiles**

```bash
cd broker/mock && go build ./...
```

Expected: no errors.

- [ ] **Step 2.3: Run existing test**

```bash
cd broker/mock && go test ./... -run TestNewMockMessage -v
```

Expected: PASS

- [ ] **Step 2.4: Commit**

```bash
git add broker/mock/mock.go
git commit -m "feat(broker/mock): add MockBroker, MockSubscriber, all Expected types"
```

---

### Task 3: Tests for Connect and Disconnect

**Files:**
- Modify: `broker/mock/mock_test.go`

- [ ] **Step 3.1: Add tests**

Append to `broker/mock/mock_test.go`:

```go
func TestConnect(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestConnect_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	if err := b.Connect(ctx); err == nil {
		t.Fatal("expected error for unexpected Connect, got nil")
	}
}

func TestConnect_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("conn failed")
	b.ExpectConnect().WillReturnError(want)

	if err := b.Connect(ctx); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestDisconnect(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectDisconnect()

	_ = b.Connect(ctx)
	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestDisconnect_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	if err := b.Disconnect(ctx); err == nil {
		t.Fatal("expected error for unexpected Disconnect, got nil")
	}
}

func TestDisconnect_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	want := fmt.Errorf("disconnect failed")
	b.ExpectDisconnect().WillReturnError(want)

	_ = b.Connect(ctx)
	if err := b.Disconnect(ctx); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestExpectationsWereMet_Unfulfilled(t *testing.T) {
	b := mock.NewMockBroker()
	b.ExpectConnect()

	if err := b.ExpectationsWereMet(); err == nil {
		t.Fatal("expected error for unfulfilled expectation, got nil")
	}
}

func TestExpectationsWereMet_Empty(t *testing.T) {
	b := mock.NewMockBroker()
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}
```

- [ ] **Step 3.2: Run tests**

```bash
cd broker/mock && go test ./... -run "TestConnect|TestDisconnect|TestExpectationsWereMet" -v
```

Expected: all PASS

- [ ] **Step 3.3: Commit**

```bash
git add broker/mock/mock_test.go
git commit -m "test(broker/mock): add Connect, Disconnect, ExpectationsWereMet tests"
```

---

### Task 4: Tests for Publish

**Files:**
- Modify: `broker/mock/mock_test.go`

- [ ] **Step 4.1: Add tests**

Append to `broker/mock/mock_test.go`:

```go
func TestPublish(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectPublish("orders")

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestPublish_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err == nil {
		t.Fatal("expected error for unexpected Publish, got nil")
	}
}

func TestPublish_WrongTopic(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectPublish("invoices")

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err == nil {
		t.Fatal("expected error for wrong topic, got nil")
	}
}

func TestPublish_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("publish failed")
	b.ExpectConnect()
	b.ExpectPublish("orders").WillReturnError(want)

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestPublish_NotConnected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectPublish("orders")

	if err := b.Publish(ctx, "orders"); err != broker.ErrNotConnected {
		t.Fatalf("want ErrNotConnected, got %v", err)
	}
}

func TestPublish_Delay(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectPublish("orders").WillDelayFor(50 * time.Millisecond)

	_ = b.Connect(ctx)
	start := time.Now()
	_ = b.Publish(ctx, "orders")
	if time.Since(start) < 50*time.Millisecond {
		t.Fatal("expected delay of at least 50ms")
	}
}
```

- [ ] **Step 4.2: Run tests**

```bash
cd broker/mock && go test ./... -run TestPublish -v
```

Expected: all PASS

- [ ] **Step 4.3: Commit**

```bash
git add broker/mock/mock_test.go
git commit -m "test(broker/mock): add Publish tests"
```

---

### Task 5: Tests for Subscribe and InjectMessage

**Files:**
- Modify: `broker/mock/mock_test.go`

- [ ] **Step 5.1: Add tests**

Append to `broker/mock/mock_test.go`:

```go
func TestSubscribe(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)
	sub, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if sub == nil {
		t.Fatal("expected non-nil subscriber")
	}
	if sub.Topic() != "orders" {
		t.Fatalf("want topic %q, got %q", "orders", sub.Topic())
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestSubscribe_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()

	_ = b.Connect(ctx)
	_, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err == nil {
		t.Fatal("expected error for unexpected Subscribe, got nil")
	}
}

func TestSubscribe_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("subscribe failed")
	b.ExpectConnect()
	b.ExpectSubscribe("orders").WillReturnError(want)

	_ = b.Connect(ctx)
	_, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestSubscribe_NotConnected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectSubscribe("orders")

	_, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err != broker.ErrNotConnected {
		t.Fatalf("want ErrNotConnected, got %v", err)
	}
}

func TestInjectMessage_SingleHandler(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)

	received := make([]broker.Message, 0)
	_, _ = b.Subscribe(ctx, "orders", func(m broker.Message) error {
		received = append(received, m)
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{"x": "y"}, []byte(`{}`))
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("want 1 message, got %d", len(received))
	}
	if string(received[0].Body()) != `{}` {
		t.Fatalf("unexpected body: %s", received[0].Body())
	}
}

func TestInjectMessage_BatchHandler(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)

	var batchReceived []broker.Message
	_, _ = b.Subscribe(ctx, "orders", func(msgs []broker.Message) error {
		batchReceived = append(batchReceived, msgs...)
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":2}`))
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage: %v", err)
	}
	if len(batchReceived) != 1 {
		t.Fatalf("want 1 message in batch, got %d", len(batchReceived))
	}
}

func TestInjectMessage_NoHandlers(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	msg := mock.NewMockMessage(ctx, "orders", nil, []byte(`{}`))
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage with no handlers should return nil, got: %v", err)
	}
}

func TestInjectMessage_HandlerError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)
	want := fmt.Errorf("handler error")
	_, _ = b.Subscribe(ctx, "orders", func(broker.Message) error { return want })

	msg := mock.NewMockMessage(ctx, "orders", nil, nil)
	if err := b.InjectMessage(ctx, "orders", msg); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}
```

- [ ] **Step 5.2: Run tests**

```bash
cd broker/mock && go test ./... -run "TestSubscribe|TestInjectMessage" -v
```

Expected: all PASS

- [ ] **Step 5.3: Commit**

```bash
git add broker/mock/mock_test.go
git commit -m "test(broker/mock): add Subscribe and InjectMessage tests"
```

---

### Task 6: Tests for Unsubscribe

**Files:**
- Modify: `broker/mock/mock_test.go`

- [ ] **Step 6.1: Add tests**

Append to `broker/mock/mock_test.go`:

```go
func TestUnsubscribe(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders")

	_ = b.Connect(ctx)
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unsubscribe: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestUnsubscribe_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })

	if err := sub.Unsubscribe(ctx); err == nil {
		t.Fatal("expected error for unexpected Unsubscribe, got nil")
	}
}

func TestUnsubscribe_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("unsubscribe failed")
	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders").WillReturnError(want)

	_ = b.Connect(ctx)
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })

	if err := sub.Unsubscribe(ctx); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestUnsubscribe_RemovesHandler(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders")

	_ = b.Connect(ctx)

	called := false
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error {
		called = true
		return nil
	})
	_ = sub.Unsubscribe(ctx)

	msg := mock.NewMockMessage(ctx, "orders", nil, nil)
	_ = b.InjectMessage(ctx, "orders", msg)

	if called {
		t.Fatal("handler should not be called after Unsubscribe")
	}
}
```

- [ ] **Step 6.2: Run tests**

```bash
cd broker/mock && go test ./... -run TestUnsubscribe -v
```

Expected: all PASS

- [ ] **Step 6.3: Commit**

```bash
git add broker/mock/mock_test.go
git commit -m "test(broker/mock): add Unsubscribe tests"
```

---

### Task 7: Full lifecycle integration test

**Files:**
- Modify: `broker/mock/mock_test.go`

- [ ] **Step 7.1: Add test**

Append to `broker/mock/mock_test.go`:

```go
func TestFullLifecycle(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectPublish("results")
	b.ExpectUnsubscribe("orders")
	b.ExpectDisconnect()

	// connect
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// subscribe
	processed := make([]string, 0)
	sub, err := b.Subscribe(ctx, "orders", func(m broker.Message) error {
		processed = append(processed, string(m.Body()))
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// inject a message — handler receives it
	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":42}`))
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage: %v", err)
	}

	// publish (recorded, not delivered)
	if err := b.Publish(ctx, "results"); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	// unsubscribe
	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unsubscribe: %v", err)
	}

	// disconnect
	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect: %v", err)
	}

	// verify all expectations met
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}

	if len(processed) != 1 || processed[0] != `{"id":42}` {
		t.Fatalf("unexpected processed messages: %v", processed)
	}
}
```

- [ ] **Step 7.2: Run all tests**

```bash
cd broker/mock && go test ./... -v
```

Expected: all PASS

- [ ] **Step 7.3: Commit**

```bash
git add broker/mock/mock_test.go
git commit -m "test(broker/mock): add full lifecycle integration test"
```

---

### Task 8: Godoc examples

**Files:**
- Create: `broker/mock/example_test.go`

- [ ] **Step 8.1: Create example_test.go**

```go
package mock_test

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/broker/mock"
	"go.unistack.org/micro/v5/metadata"
)

// Example demonstrates the full lifecycle: connect, subscribe, inject a message,
// publish, unsubscribe, disconnect, and verify all expectations were met.
func Example() {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectPublish("results")
	b.ExpectUnsubscribe("orders")
	b.ExpectDisconnect()

	_ = b.Connect(ctx)

	sub, _ := b.Subscribe(ctx, "orders", func(m broker.Message) error {
		fmt.Printf("received: %s\n", m.Body())
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":1}`))
	_ = b.InjectMessage(ctx, "orders", msg)

	_ = b.Publish(ctx, "results")
	_ = sub.Unsubscribe(ctx)
	_ = b.Disconnect(ctx)

	if err := b.ExpectationsWereMet(); err != nil {
		fmt.Println("unmet expectations:", err)
	}
	// Output:
	// received: {"id":1}
}

// ExampleMockBroker_ExpectPublish demonstrates testing that application code
// publishes to the correct topic.
func ExampleMockBroker_ExpectPublish() {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectPublish("orders")
	b.ExpectDisconnect()

	_ = b.Connect(ctx)

	// application code under test
	if err := b.Publish(ctx, "orders"); err != nil {
		fmt.Println("publish error:", err)
	}

	_ = b.Disconnect(ctx)

	if err := b.ExpectationsWereMet(); err != nil {
		fmt.Println("unmet:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

// ExampleMockBroker_InjectMessage demonstrates testing a subscriber handler
// by injecting a pre-built message.
func ExampleMockBroker_InjectMessage() {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders")
	b.ExpectDisconnect()

	_ = b.Connect(ctx)

	sub, _ := b.Subscribe(ctx, "orders", func(m broker.Message) error {
		fmt.Printf("handler got: %s\n", m.Body())
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":7}`))
	_ = b.InjectMessage(ctx, "orders", msg)

	_ = sub.Unsubscribe(ctx)
	_ = b.Disconnect(ctx)
	// Output:
	// handler got: {"id":7}
}
```

- [ ] **Step 8.2: Run example tests**

```bash
cd broker/mock && go test ./... -run Example -v
```

Expected: all PASS

- [ ] **Step 8.3: Commit**

```bash
git add broker/mock/example_test.go
git commit -m "docs(broker/mock): add godoc Example functions"
```

---

## Self-Review Checklist

- [x] **Spec coverage:** Connect/Disconnect/Publish/Subscribe/Unsubscribe expectations ✓, `ExpectationsWereMet` ✓, `InjectMessage` ✓, `NewMockMessage` ✓, both handler types ✓, godoc examples ✓
- [x] **Placeholders:** None — all steps contain complete code
- [x] **Type consistency:** `handlerEntry`, `MockSubscriber.handlerID`, `MockBroker.nextID` used consistently across Tasks 2–6; `NewMockMessage` signature matches Task 1 and usage in Tasks 5–8
