# Broker Mock Design

**Date:** 2026-05-16
**Location:** `go.unistack.org/micro/v5/broker/mock`

## Goal

Provide a mock implementation of `broker.Broker` for unit testing application code that uses the broker without real message transport. The mock follows the same expectation-based pattern as `go.unistack.org/micro/v5/client/mock`.

## Approach

Strict expectation-based mock: every `Connect`, `Disconnect`, `Publish`, `Subscribe`, and `Unsubscribe` call must have a pre-declared expectation. Unexpected calls return an error. After the test, `ExpectationsWereMet()` verifies all declared expectations were fulfilled.

`Publish` records the call but does **not** deliver messages to subscribed handlers. Messages are delivered to handlers explicitly via `InjectMessage`.

## File Layout

```
micro/broker/mock/
├── mock.go          # MockBroker, MockSubscriber, expectation types, commonExpectation
└── mock_message.go  # MockMessage — implements broker.Message for InjectMessage, exported for assertions
```

## Types

### Internal

```go
// expectation is the common interface for all expected calls.
type expectation interface {
    fulfilled() bool
    Lock()
    Unlock()
    String() string
}

// commonExpectation is embedded in every Expect* type.
type commonExpectation struct {
    sync.Mutex
    triggered bool
    err       error
}

func (e *commonExpectation) fulfilled() bool { return e.triggered }
```

### Expected call types

```go
type ExpectedConnect struct {
    commonExpectation
}
func (e *ExpectedConnect) WillReturnError(err error) *ExpectedConnect

type ExpectedDisconnect struct {
    commonExpectation
}
func (e *ExpectedDisconnect) WillReturnError(err error) *ExpectedDisconnect

type ExpectedPublish struct {
    commonExpectation
    topic string
    delay time.Duration
}
func (e *ExpectedPublish) WillReturnError(err error) *ExpectedPublish
func (e *ExpectedPublish) WillDelayFor(d time.Duration) *ExpectedPublish

type ExpectedSubscribe struct {
    commonExpectation
    topic string
}
func (e *ExpectedSubscribe) WillReturnError(err error) *ExpectedSubscribe

type ExpectedUnsubscribe struct {
    commonExpectation
    topic string
}
func (e *ExpectedUnsubscribe) WillReturnError(err error) *ExpectedUnsubscribe
```

### MockBroker

```go
type MockBroker struct {
    opts     broker.Options
    mu       sync.Mutex
    expected []expectation        // ordered queue of declared expectations
    handlers map[string][]any     // topic → registered handlers (for InjectMessage)
}

var _ broker.Broker = (*MockBroker)(nil)

func NewMockBroker(opts ...broker.Option) *MockBroker

// Expectation setup
func (m *MockBroker) ExpectConnect() *ExpectedConnect
func (m *MockBroker) ExpectDisconnect() *ExpectedDisconnect
func (m *MockBroker) ExpectPublish(topic string) *ExpectedPublish
func (m *MockBroker) ExpectSubscribe(topic string) *ExpectedSubscribe
func (m *MockBroker) ExpectUnsubscribe(topic string) *ExpectedUnsubscribe

// Verification
func (m *MockBroker) ExpectationsWereMet() error

// Message injection — delivers msg to all handlers registered for topic
func (m *MockBroker) InjectMessage(ctx context.Context, topic string, msg broker.Message) error
```

### MockSubscriber

```go
type MockSubscriber struct {
    topic string
    opts  broker.SubscribeOptions
    mock  *MockBroker
}

func (s *MockSubscriber) Options() broker.SubscribeOptions
func (s *MockSubscriber) Topic() string
func (s *MockSubscriber) Unsubscribe(ctx context.Context) error
```

`Unsubscribe` looks up the first unfulfilled `ExpectedUnsubscribe` for the topic, removes the handler from `m.handlers[topic]`, and marks the expectation as triggered. Returns an error if no matching expectation is found.

### MockMessage

Exported type implementing `broker.Message`. Used for constructing messages passed to `InjectMessage`. Tracks whether `Ack` was called, allowing tests to assert handler acknowledgement behaviour.

```go
type MockMessage struct {
    ctx   context.Context
    topic string
    hdr   metadata.Metadata
    body  []byte
    c     codec.Codec
    err   error
    acked bool
}

func (m *MockMessage) Ack() error    // sets m.acked = true
func (m *MockMessage) Acked() bool   // returns whether Ack was called

// NewMockMessage creates a MockMessage with pre-marshaled body.
// Pass a non-nil codec to enable Unmarshal in the handler under test.
func NewMockMessage(ctx context.Context, topic string, hdr metadata.Metadata, body []byte, c codec.Codec) *MockMessage
```

## Behaviour per Method

| Method | Behaviour |
|---|---|
| `Connect` | Finds first unfulfilled `ExpectedConnect`. Missing → error. Sets internal `connected = true` unless `e.err != nil`. |
| `Disconnect` | Finds first unfulfilled `ExpectedDisconnect`. Missing → error. Sets `connected = false`. |
| `Publish` | Finds first unfulfilled `ExpectedPublish` matching topic. Missing → error. Applies delay. Returns `e.err`. Does **not** call handlers. |
| `Subscribe` | Finds first unfulfilled `ExpectedSubscribe` matching topic. Missing → error. If `e.err != nil`, returns `nil, e.err` without registering handler. Otherwise registers handler in `m.handlers[topic]` and returns `MockSubscriber, nil`. |
| `Unsubscribe` | Finds first unfulfilled `ExpectedUnsubscribe` matching topic. Missing → error. Removes handler from `m.handlers[topic]`. |
| `InjectMessage` | Calls all handlers in `m.handlers[topic]` with `msg`. Supports both `func(broker.Message) error` and `func([]broker.Message) error`. Returns first error. |
| `NewMessage` | Creates a `MockMessage` via codec marshal; does not require a prior expectation. |
| `Name/String/Live/Ready/Health/Address/Options` | Simple getters; no expectation required. |

## Error Cases

- Unexpected call (no matching unfulfilled expectation): `fmt.Errorf("unexpected call to %s for topic %q", methodName, topic)`
- `ExpectationsWereMet`: returns error listing all unfulfilled expectations with their `String()` description
- `InjectMessage` with no registered handlers: returns `nil` (no subscribers is valid)
- `Publish`/`Subscribe` called before `Connect`: returns `broker.ErrNotConnected`

## Usage Examples

### Test publisher

```go
mock := mock.NewMockBroker()
mock.ExpectConnect()
mock.ExpectPublish("orders").WillReturnError(nil)
mock.ExpectDisconnect()

_ = mock.Connect(ctx)
err := myService.SendOrder(ctx, order)
_ = mock.Disconnect(ctx)

assert.NoError(t, err)
assert.NoError(t, mock.ExpectationsWereMet())
```

### Test subscriber handler

```go
mock := mock.NewMockBroker()
mock.ExpectConnect()
mock.ExpectSubscribe("orders")
mock.ExpectUnsubscribe("orders")
mock.ExpectDisconnect()

_ = mock.Connect(ctx)
myService.StartConsumer(ctx, mock)

msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{"key": "val"}, []byte(`{}`), nil)
err := mock.InjectMessage(ctx, "orders", msg)
assert.NoError(t, err)

myService.Stop(ctx)
assert.NoError(t, mock.ExpectationsWereMet())
```

### Full lifecycle

```go
mock := mock.NewMockBroker()
mock.ExpectConnect()
mock.ExpectSubscribe("orders")
mock.ExpectPublish("results")
mock.ExpectUnsubscribe("orders")
mock.ExpectDisconnect()

myService.Run(ctx, mock)

assert.NoError(t, mock.ExpectationsWereMet())
```

## Godoc Examples

A file `example_test.go` in the `mock` package provides runnable `Example_*` functions visible in godoc. Required examples:

- `Example()` — full lifecycle (connect → subscribe → inject → unsubscribe → disconnect)
- `ExampleMockBroker_ExpectPublish()` — testing a publisher
- `ExampleMockBroker_InjectMessage()` — testing a subscriber handler

Each example follows the standard Go testable-example pattern (package `mock_test`, no `// Output:` needed since broker calls are side-effectful).

## Out of Scope

- Message body validation in `ExpectPublish` (not needed per requirements)
- Pass-through delivery from Publish to Subscribe handlers
- Batch-publish expectation (can be added later)
