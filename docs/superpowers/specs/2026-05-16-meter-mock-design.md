# Design: meter/mock — strict mock for meter.Meter

**Date:** 2026-05-16
**Module:** `go.unistack.org/micro/v5`
**Path:** `meter/mock/mock.go`, `meter/mock/mock_test.go`

---

## Goal

Implement a mock of the `meter.Meter` interface in package `go.unistack.org/micro/v5/meter/mock`, following the same pattern as `broker/mock`. The mock must:

1. Register expectations on Meter methods (`ExpectCounter(...)`, `ExpectInit()`, etc.)
2. Verify all expectations were fulfilled via `ExpectationsWereMet()`
3. Allow tests to inspect the actual accumulated metric values (Counter, Gauge, etc.)
4. Return an error on unexpected calls to methods that return `error`

---

## meter.Meter interface (reference)

```go
type Meter interface {
    Name() string
    Init(opts ...Option) error
    Clone(opts ...Option) Meter
    Counter(name string, labels ...string) Counter
    FloatCounter(name string, labels ...string) FloatCounter
    Gauge(name string, fn func() float64, labels ...string) Gauge
    Set(opts ...Option) Meter
    Histogram(name string, labels ...string) Histogram
    HistogramExt(name string, quantiles []float64, labels ...string) Histogram
    Summary(name string, labels ...string) Summary
    SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) Summary
    Write(w io.Writer, opts ...Option) error
    Options() Options
    String() string
    Unregister(name string, labels ...string) bool
}
```

---

## Architecture

### Expectation infrastructure (mirroring broker/mock)

```go
// expectation is the common interface for all Expected* types.
type expectation interface {
    fulfilled() bool
    Lock()
    Unlock()
    String() string
}

// commonExpectation is embedded in all Expected* types.
type commonExpectation struct {
    sync.Mutex
    triggered bool
    err       error
}

func (e *commonExpectation) fulfilled() bool { return e.triggered }
```

### Expected types

| Type | Meter method | Fluent methods |
|------|--------------|----------------|
| `ExpectedInit` | `Init()` | `WillReturnError(err)` |
| `ExpectedWrite` | `Write()` | `WillReturnError(err)` |
| `ExpectedUnregister` | `Unregister(name, labels...)` | `WillReturn(bool)` |
| `ExpectedCounter` | `Counter(name, labels...)` | exposes `*mockCounter` |
| `ExpectedFloatCounter` | `FloatCounter(name, labels...)` | exposes `*mockFloatCounter` |
| `ExpectedGauge` | `Gauge(name, fn, labels...)` | exposes `*mockGauge` |
| `ExpectedHistogram` | `Histogram(name, labels...)` | exposes `*mockHistogram` |
| `ExpectedHistogramExt` | `HistogramExt(name, quantiles, labels...)` | exposes `*mockHistogram` |
| `ExpectedSummary` | `Summary(name, labels...)` | exposes `*mockSummary` |
| `ExpectedSummaryExt` | `SummaryExt(name, window, quantiles, labels...)` | exposes `*mockSummary` |

### Metric mock objects (state-tracking)

Each sub-interface mock implements the corresponding `meter.*` interface and accumulates real values:

- `mockCounter` — `value uint64`, guarded by `sync.Mutex`
- `mockFloatCounter` — `value float64`, guarded by `sync.Mutex`
- `mockGauge` — `value float64`, guarded by `sync.Mutex`; `fn func() float64` stored for `Get()` when provided
- `mockHistogram` — `updates []float64`, guarded by `sync.Mutex`; `Reset()` clears the slice
- `mockSummary` — `updates []float64`, guarded by `sync.Mutex`

Tests access the accumulated value through the expectation's field:

```go
exp := m.ExpectCounter("requests", "method", "GET")
// ... run code under test ...
if exp.Counter().Get() != 5 {
    t.Errorf("expected 5, got %d", exp.Counter().Get())
}
```

### MockMeter

```go
type MockMeter struct {
    opts       meter.Options
    mu         sync.Mutex
    expected   []expectation
    unexpected []string  // unexpected call descriptions
}

func NewMockMeter(opts ...meter.Option) *MockMeter
```

**Call dispatch rules:**

- The first untriggered matching expectation fires and is marked `triggered = true`.
- If no matching expectation exists:
  - `Init()`, `Write()` → `fmt.Errorf("unexpected call to ...")`
  - `Counter()`, `FloatCounter()`, `Gauge()`, `Histogram()`, `Summary()` → noop object + append to `unexpected`
  - `Unregister()` → `false` + append to `unexpected`
- `Clone()` and `Set()` → return the same `*MockMeter` (no expectation required)

### ExpectationsWereMet()

```go
func (m *MockMeter) ExpectationsWereMet() error
```

Returns an error if:
1. Any registered expectation was not triggered (`triggered == false`)
2. There were unexpected calls (non-empty `unexpected`)

---

## Files

```
meter/
  mock/
    mock.go       — MockMeter, Expected*, mockCounter/FloatCounter/Gauge/Histogram/Summary
    mock_test.go  — tests: Init, Counter.Inc, Histogram.Update, ExpectationsWereMet, etc.
```

---

## Usage examples

```go
// Basic scenario
m := mock.NewMockMeter()
m.ExpectInit()
m.ExpectCounter("requests", "method", "GET")

_ = m.Init()
c := m.Counter("requests", "method", "GET")
c.Inc()
c.Inc()

exp := m.expected[1].(*mock.ExpectedCounter)
if exp.Counter().Get() != 2 {
    t.Error("expected 2 increments")
}
if err := m.ExpectationsWereMet(); err != nil {
    t.Error(err)
}

// Error injection
m2 := mock.NewMockMeter()
m2.ExpectInit().WillReturnError(errors.New("init failed"))
if err := m2.Init(); err == nil {
    t.Error("expected error")
}

// Unregister
m3 := mock.NewMockMeter()
m3.ExpectUnregister("old_metric").WillReturn(true)
if !m3.Unregister("old_metric") {
    t.Error("expected true")
}
```

---

## Test coverage

`mock_test.go` covers:

1. `TestMockMeter_Init` — expectation, error injection, unexpected call
2. `TestMockMeter_Counter` — `Inc`, `Dec`, `Add`, `Set`, `Get` via `ExpectedCounter.Counter()`
3. `TestMockMeter_FloatCounter` — `Add`, `Sub`, `Set`, `Get`
4. `TestMockMeter_Gauge` — `Inc`, `Dec`, `Set`, `Get`
5. `TestMockMeter_Histogram` — `Update`, `UpdateDuration`, `Reset`
6. `TestMockMeter_HistogramExt` — same as Histogram
7. `TestMockMeter_Summary` — `Update`, `UpdateDuration`
8. `TestMockMeter_SummaryExt` — same as Summary
9. `TestMockMeter_Write` — expectation, error injection
10. `TestMockMeter_Unregister` — `WillReturn(true/false)`
11. `TestMockMeter_ExpectationsWereMet` — unfulfilled expectations, unexpected calls
12. `TestMockMeter_CloneSet` — Clone and Set without expectations

---

## Project standards

- Compile-time interface check: `var _ meter.Meter = (*MockMeter)(nil)`
- All public types and methods documented (required per project CLAUDE.md)
- Package declaration: `package mock`
- Import path: `go.unistack.org/micro/v5/meter`
