# meter/mock Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement `go.unistack.org/micro/v5/meter/mock` — a strict expectations-based mock of the `meter.Meter` interface for use in tests.

**Architecture:** `MockMeter` holds an ordered slice of `expectation` values registered via `ExpectXxx()` methods. On each method call the mock scans for the first untriggered matching expectation, marks it triggered, and returns configured values. Unexpected calls (no matching expectation) are recorded and reported by `ExpectationsWereMet()`. Sub-interface types (`mockCounter`, `mockGauge`, etc.) accumulate real values so tests can inspect them via `exp.Counter().Get()`.

**Tech Stack:** Go 1.26, `go.unistack.org/micro/v5/meter`, `sync`, `fmt`, `io`, `time`

---

## File map

| File | Action | Responsibility |
|------|--------|----------------|
| `meter/mock/mock.go` | Create | `MockMeter`, all `Expected*` types, all `mock*` sub-interface types |
| `meter/mock/mock_test.go` | Create | All tests (grows with each task) |

---

### Task 1: Package scaffold

Create the two files. `mock.go` contains `MockMeter` satisfying `meter.Meter` with all methods stubbed (no expectation dispatch yet), plus all five sub-interface mock types. `mock_test.go` has only a compile-time assertion test.

**Files:**
- Create: `meter/mock/mock.go`
- Create: `meter/mock/mock_test.go`

- [ ] **Step 1: Write the failing test**

Create `meter/mock/mock_test.go`:

```go
package mock

import (
	"testing"

	"go.unistack.org/micro/v5/meter"
)

func TestMockMeter_Implements(t *testing.T) {
	var _ meter.Meter = NewMockMeter()
}
```

- [ ] **Step 2: Run test — expect failure**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro
go test ./meter/mock/... -v -count=1
```

Expected: `cannot find package` or `undefined: NewMockMeter`

- [ ] **Step 3: Create `meter/mock/mock.go` with the scaffold**

```go
// Package mock provides a strict expectations-based mock of meter.Meter.
package mock

import (
	"fmt"
	"io"
	"sync"
	"time"

	"go.unistack.org/micro/v5/meter"
)

// Compile-time interface check.
var _ meter.Meter = (*MockMeter)(nil)

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

// MockMeter is a mock implementation of meter.Meter for use in tests.
//
// Register expectations with ExpectInit, ExpectCounter, etc. before exercising
// the code under test. After the test call ExpectationsWereMet to verify all
// declared expectations were fulfilled.
type MockMeter struct {
	opts       meter.Options
	mu         sync.Mutex
	expected   []expectation
	unexpected []string
}

// NewMockMeter creates a new MockMeter.
func NewMockMeter(opts ...meter.Option) *MockMeter {
	return &MockMeter{opts: meter.NewOptions(opts...)}
}

// Name implements meter.Meter.
func (m *MockMeter) Name() string { return m.opts.Name }

// Options implements meter.Meter.
func (m *MockMeter) Options() meter.Options { return m.opts }

// String implements meter.Meter.
func (m *MockMeter) String() string { return "mock" }

// Clone implements meter.Meter. Returns the same MockMeter without requiring an expectation.
func (m *MockMeter) Clone(_ ...meter.Option) meter.Meter { return m }

// Set implements meter.Meter. Returns the same MockMeter without requiring an expectation.
func (m *MockMeter) Set(_ ...meter.Option) meter.Meter { return m }

// Init implements meter.Meter.
func (m *MockMeter) Init(_ ...meter.Option) error {
	return fmt.Errorf("unexpected call to meter.Init")
}

// Write implements meter.Meter.
func (m *MockMeter) Write(_ io.Writer, _ ...meter.Option) error {
	return fmt.Errorf("unexpected call to meter.Write")
}

// Counter implements meter.Meter.
func (m *MockMeter) Counter(name string, _ ...string) meter.Counter {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Counter(%q)", name))
	m.mu.Unlock()
	return &mockCounter{}
}

// FloatCounter implements meter.Meter.
func (m *MockMeter) FloatCounter(name string, _ ...string) meter.FloatCounter {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("FloatCounter(%q)", name))
	m.mu.Unlock()
	return &mockFloatCounter{}
}

// Gauge implements meter.Meter.
func (m *MockMeter) Gauge(name string, _ func() float64, _ ...string) meter.Gauge {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Gauge(%q)", name))
	m.mu.Unlock()
	return &mockGauge{}
}

// Histogram implements meter.Meter.
func (m *MockMeter) Histogram(name string, _ ...string) meter.Histogram {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Histogram(%q)", name))
	m.mu.Unlock()
	return &mockHistogram{}
}

// HistogramExt implements meter.Meter.
func (m *MockMeter) HistogramExt(name string, _ []float64, _ ...string) meter.Histogram {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("HistogramExt(%q)", name))
	m.mu.Unlock()
	return &mockHistogram{}
}

// Summary implements meter.Meter.
func (m *MockMeter) Summary(name string, _ ...string) meter.Summary {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Summary(%q)", name))
	m.mu.Unlock()
	return &mockSummary{}
}

// SummaryExt implements meter.Meter.
func (m *MockMeter) SummaryExt(name string, _ time.Duration, _ []float64, _ ...string) meter.Summary {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("SummaryExt(%q)", name))
	m.mu.Unlock()
	return &mockSummary{}
}

// Unregister implements meter.Meter.
func (m *MockMeter) Unregister(name string, _ ...string) bool {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Unregister(%q)", name))
	m.mu.Unlock()
	return false
}

// ExpectationsWereMet returns an error if any declared expectation was not fulfilled
// or if there were unexpected calls.
func (m *MockMeter) ExpectationsWereMet() error {
	return nil // implemented in Task 8
}

// --- Sub-interface mock types ---

// mockCounter is a state-tracking implementation of meter.Counter.
type mockCounter struct {
	mu    sync.Mutex
	value uint64
}

func (c *mockCounter) Add(v int) {
	c.mu.Lock()
	if v >= 0 {
		c.value += uint64(v)
	} else {
		u := uint64(-v)
		if u <= c.value {
			c.value -= u
		} else {
			c.value = 0
		}
	}
	c.mu.Unlock()
}

func (c *mockCounter) Dec() { c.mu.Lock(); if c.value > 0 { c.value-- }; c.mu.Unlock() }
func (c *mockCounter) Inc() { c.mu.Lock(); c.value++; c.mu.Unlock() }
func (c *mockCounter) Get() uint64 { c.mu.Lock(); defer c.mu.Unlock(); return c.value }
func (c *mockCounter) Set(v uint64) { c.mu.Lock(); c.value = v; c.mu.Unlock() }

// mockFloatCounter is a state-tracking implementation of meter.FloatCounter.
type mockFloatCounter struct {
	mu    sync.Mutex
	value float64
}

func (c *mockFloatCounter) Add(v float64) { c.mu.Lock(); c.value += v; c.mu.Unlock() }
func (c *mockFloatCounter) Sub(v float64) { c.mu.Lock(); c.value -= v; c.mu.Unlock() }
func (c *mockFloatCounter) Get() float64  { c.mu.Lock(); defer c.mu.Unlock(); return c.value }
func (c *mockFloatCounter) Set(v float64) { c.mu.Lock(); c.value = v; c.mu.Unlock() }

// mockGauge is a state-tracking implementation of meter.Gauge.
type mockGauge struct {
	mu    sync.Mutex
	value float64
}

func (g *mockGauge) Add(v float64) { g.mu.Lock(); g.value += v; g.mu.Unlock() }
func (g *mockGauge) Set(v float64) { g.mu.Lock(); g.value = v; g.mu.Unlock() }
func (g *mockGauge) Inc()          { g.mu.Lock(); g.value++; g.mu.Unlock() }
func (g *mockGauge) Dec()          { g.mu.Lock(); g.value--; g.mu.Unlock() }
func (g *mockGauge) Get() float64  { g.mu.Lock(); defer g.mu.Unlock(); return g.value }

// mockHistogram is a state-tracking implementation of meter.Histogram.
type mockHistogram struct {
	mu      sync.Mutex
	updates []float64
}

func (h *mockHistogram) Reset()                     { h.mu.Lock(); h.updates = nil; h.mu.Unlock() }
func (h *mockHistogram) Update(v float64)           { h.mu.Lock(); h.updates = append(h.updates, v); h.mu.Unlock() }
func (h *mockHistogram) UpdateDuration(t time.Time) { h.Update(time.Since(t).Seconds()) }

// Updates returns a copy of the recorded observation values.
func (h *mockHistogram) Updates() []float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]float64, len(h.updates))
	copy(out, h.updates)
	return out
}

// mockSummary is a state-tracking implementation of meter.Summary.
type mockSummary struct {
	mu      sync.Mutex
	updates []float64
}

func (s *mockSummary) Update(v float64)           { s.mu.Lock(); s.updates = append(s.updates, v); s.mu.Unlock() }
func (s *mockSummary) UpdateDuration(t time.Time) { s.Update(time.Since(t).Seconds()) }

// Updates returns a copy of the recorded observation values.
func (s *mockSummary) Updates() []float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]float64, len(s.updates))
	copy(out, s.updates)
	return out
}
```

- [ ] **Step 4: Run test — expect pass**

```bash
go test ./meter/mock/... -v -count=1
```

Expected: `PASS TestMockMeter_Implements`

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add package scaffold with sub-interface mocks"
```

---

### Task 2: Init expectation

Add `ExpectedInit`, `ExpectInit()`, and update `Init()` dispatch.

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
import (
	"errors"
	"testing"

	"go.unistack.org/micro/v5/meter"
)

func TestMockMeter_Init_Success(t *testing.T) {
	m := NewMockMeter()
	m.ExpectInit()
	if err := m.Init(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMockMeter_Init_Error(t *testing.T) {
	m := NewMockMeter()
	m.ExpectInit().WillReturnError(errors.New("init failed"))
	err := m.Init()
	if err == nil || err.Error() != "init failed" {
		t.Fatalf("expected 'init failed', got %v", err)
	}
}

func TestMockMeter_Init_Unexpected(t *testing.T) {
	m := NewMockMeter()
	if err := m.Init(); err == nil {
		t.Fatal("expected error for unexpected Init call")
	}
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run TestMockMeter_Init
```

Expected: `undefined: ExpectInit` (compile error)

- [ ] **Step 3: Add `ExpectedInit` and `ExpectInit()` to `mock.go`, update `Init()`**

After the `commonExpectation` definition, add:

```go
// ExpectedInit holds an expectation for meter.Init.
type ExpectedInit struct {
	commonExpectation
}

// WillReturnError configures the error returned by Init.
func (e *ExpectedInit) WillReturnError(err error) *ExpectedInit { e.err = err; return e }

func (e *ExpectedInit) String() string {
	if e.err != nil {
		return fmt.Sprintf("ExpectedInit => should return error: %s", e.err)
	}
	return "ExpectedInit"
}
```

Add `ExpectInit()` method to `MockMeter`:

```go
// ExpectInit registers an expectation that meter.Init will be called.
func (m *MockMeter) ExpectInit() *ExpectedInit {
	e := &ExpectedInit{}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}
```

Replace the `Init` method body:

```go
// Init implements meter.Meter.
func (m *MockMeter) Init(_ ...meter.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		ei, ok := e.(*ExpectedInit)
		if !ok {
			continue
		}
		ei.Lock()
		if ei.triggered {
			ei.Unlock()
			continue
		}
		ei.triggered = true
		err := ei.err
		ei.Unlock()
		return err
	}
	return fmt.Errorf("unexpected call to meter.Init")
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1 -run TestMockMeter_Init
```

Expected: all three `TestMockMeter_Init_*` tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add Init expectation"
```

---

### Task 3: Counter, FloatCounter, Gauge expectations

Add `ExpectedCounter`, `ExpectedFloatCounter`, `ExpectedGauge` types plus their `ExpectXxx()` methods and dispatch updates.

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
func TestMockMeter_Counter(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectCounter("requests", "method", "GET")

	c := m.Counter("requests", "method", "GET")
	c.Inc()
	c.Inc()
	c.Add(3)
	c.Dec()
	c.Set(10)

	if got := exp.Counter().Get(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestMockMeter_Counter_Unexpected(t *testing.T) {
	m := NewMockMeter()
	_ = m.Counter("hits") // no expectation — goes to unexpected
	if len(m.unexpected) == 0 {
		t.Fatal("expected unexpected call to be recorded")
	}
}

func TestMockMeter_FloatCounter(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectFloatCounter("bytes")

	fc := m.FloatCounter("bytes")
	fc.Add(1.5)
	fc.Add(2.5)
	fc.Sub(1.0)
	fc.Set(5.0)

	if got := exp.FloatCounter().Get(); got != 5.0 {
		t.Fatalf("expected 5.0, got %f", got)
	}
}

func TestMockMeter_Gauge(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectGauge("temperature")

	g := m.Gauge("temperature", nil)
	g.Set(42.0)
	g.Inc()
	g.Dec()
	g.Add(0.5)

	if got := exp.Gauge().Get(); got != 42.5 {
		t.Fatalf("expected 42.5, got %f", got)
	}
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_(Counter|FloatCounter|Gauge)'
```

Expected: `undefined: ExpectCounter` (compile error)

- [ ] **Step 3: Add Expected types and Expect methods to `mock.go`**

After `ExpectedInit`, add:

```go
// ExpectedCounter holds an expectation for meter.Counter.
type ExpectedCounter struct {
	commonExpectation
	name    string
	labels  []string
	counter *mockCounter
}

// Counter returns the mock counter for value inspection.
func (e *ExpectedCounter) Counter() *mockCounter { return e.counter }

func (e *ExpectedCounter) String() string {
	return fmt.Sprintf("ExpectedCounter(%q)", e.name)
}

// ExpectedFloatCounter holds an expectation for meter.FloatCounter.
type ExpectedFloatCounter struct {
	commonExpectation
	name         string
	labels       []string
	floatCounter *mockFloatCounter
}

// FloatCounter returns the mock float counter for value inspection.
func (e *ExpectedFloatCounter) FloatCounter() *mockFloatCounter { return e.floatCounter }

func (e *ExpectedFloatCounter) String() string {
	return fmt.Sprintf("ExpectedFloatCounter(%q)", e.name)
}

// ExpectedGauge holds an expectation for meter.Gauge.
type ExpectedGauge struct {
	commonExpectation
	name   string
	labels []string
	gauge  *mockGauge
}

// Gauge returns the mock gauge for value inspection.
func (e *ExpectedGauge) Gauge() *mockGauge { return e.gauge }

func (e *ExpectedGauge) String() string {
	return fmt.Sprintf("ExpectedGauge(%q)", e.name)
}
```

Add `ExpectCounter`, `ExpectFloatCounter`, `ExpectGauge` methods to `MockMeter`:

```go
// ExpectCounter registers an expectation that meter.Counter will be called with name.
func (m *MockMeter) ExpectCounter(name string, labels ...string) *ExpectedCounter {
	e := &ExpectedCounter{name: name, labels: labels, counter: &mockCounter{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectFloatCounter registers an expectation that meter.FloatCounter will be called with name.
func (m *MockMeter) ExpectFloatCounter(name string, labels ...string) *ExpectedFloatCounter {
	e := &ExpectedFloatCounter{name: name, labels: labels, floatCounter: &mockFloatCounter{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectGauge registers an expectation that meter.Gauge will be called with name.
func (m *MockMeter) ExpectGauge(name string, labels ...string) *ExpectedGauge {
	e := &ExpectedGauge{name: name, labels: labels, gauge: &mockGauge{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}
```

Replace `Counter`, `FloatCounter`, `Gauge` method bodies:

```go
// Counter implements meter.Meter.
func (m *MockMeter) Counter(name string, _ ...string) meter.Counter {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		ec, ok := e.(*ExpectedCounter)
		if !ok {
			continue
		}
		ec.Lock()
		if ec.triggered || ec.name != name {
			ec.Unlock()
			continue
		}
		ec.triggered = true
		c := ec.counter
		ec.Unlock()
		return c
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("Counter(%q)", name))
	return &mockCounter{}
}

// FloatCounter implements meter.Meter.
func (m *MockMeter) FloatCounter(name string, _ ...string) meter.FloatCounter {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		ec, ok := e.(*ExpectedFloatCounter)
		if !ok {
			continue
		}
		ec.Lock()
		if ec.triggered || ec.name != name {
			ec.Unlock()
			continue
		}
		ec.triggered = true
		fc := ec.floatCounter
		ec.Unlock()
		return fc
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("FloatCounter(%q)", name))
	return &mockFloatCounter{}
}

// Gauge implements meter.Meter.
func (m *MockMeter) Gauge(name string, _ func() float64, _ ...string) meter.Gauge {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		eg, ok := e.(*ExpectedGauge)
		if !ok {
			continue
		}
		eg.Lock()
		if eg.triggered || eg.name != name {
			eg.Unlock()
			continue
		}
		eg.triggered = true
		g := eg.gauge
		eg.Unlock()
		return g
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("Gauge(%q)", name))
	return &mockGauge{}
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_(Counter|FloatCounter|Gauge)'
```

Expected: all five tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add Counter, FloatCounter, Gauge expectations"
```

---

### Task 4: Histogram and HistogramExt expectations

Add `ExpectedHistogram`, `ExpectedHistogramExt` and their dispatch.

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
func TestMockMeter_Histogram(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectHistogram("latency")

	h := m.Histogram("latency")
	h.Update(1.0)
	h.Update(2.5)
	h.Reset()
	h.Update(3.0)

	updates := exp.Histogram().Updates()
	if len(updates) != 1 || updates[0] != 3.0 {
		t.Fatalf("expected [3.0], got %v", updates)
	}
}

func TestMockMeter_HistogramExt(t *testing.T) {
	m := NewMockMeter()
	quantiles := []float64{0.5, 0.9, 0.99}
	exp := m.ExpectHistogramExt("latency_ext", quantiles)

	h := m.HistogramExt("latency_ext", quantiles)
	h.Update(1.0)
	h.Update(2.0)

	updates := exp.Histogram().Updates()
	if len(updates) != 2 {
		t.Fatalf("expected 2 updates, got %d", len(updates))
	}
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Histogram'
```

Expected: `undefined: ExpectHistogram` (compile error)

- [ ] **Step 3: Add Expected types and Expect methods to `mock.go`**

After `ExpectedGauge`, add:

```go
// ExpectedHistogram holds an expectation for meter.Histogram.
type ExpectedHistogram struct {
	commonExpectation
	name      string
	labels    []string
	histogram *mockHistogram
}

// Histogram returns the mock histogram for value inspection.
func (e *ExpectedHistogram) Histogram() *mockHistogram { return e.histogram }

func (e *ExpectedHistogram) String() string {
	return fmt.Sprintf("ExpectedHistogram(%q)", e.name)
}

// ExpectedHistogramExt holds an expectation for meter.HistogramExt.
type ExpectedHistogramExt struct {
	commonExpectation
	name      string
	quantiles []float64
	labels    []string
	histogram *mockHistogram
}

// Histogram returns the mock histogram for value inspection.
func (e *ExpectedHistogramExt) Histogram() *mockHistogram { return e.histogram }

func (e *ExpectedHistogramExt) String() string {
	return fmt.Sprintf("ExpectedHistogramExt(%q)", e.name)
}
```

Add `ExpectHistogram` and `ExpectHistogramExt` to `MockMeter`:

```go
// ExpectHistogram registers an expectation that meter.Histogram will be called with name.
func (m *MockMeter) ExpectHistogram(name string, labels ...string) *ExpectedHistogram {
	e := &ExpectedHistogram{name: name, labels: labels, histogram: &mockHistogram{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectHistogramExt registers an expectation that meter.HistogramExt will be called with name.
func (m *MockMeter) ExpectHistogramExt(name string, quantiles []float64, labels ...string) *ExpectedHistogramExt {
	e := &ExpectedHistogramExt{name: name, quantiles: quantiles, labels: labels, histogram: &mockHistogram{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}
```

Replace `Histogram` and `HistogramExt` method bodies:

```go
// Histogram implements meter.Meter.
func (m *MockMeter) Histogram(name string, _ ...string) meter.Histogram {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		eh, ok := e.(*ExpectedHistogram)
		if !ok {
			continue
		}
		eh.Lock()
		if eh.triggered || eh.name != name {
			eh.Unlock()
			continue
		}
		eh.triggered = true
		h := eh.histogram
		eh.Unlock()
		return h
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("Histogram(%q)", name))
	return &mockHistogram{}
}

// HistogramExt implements meter.Meter.
func (m *MockMeter) HistogramExt(name string, _ []float64, _ ...string) meter.Histogram {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		eh, ok := e.(*ExpectedHistogramExt)
		if !ok {
			continue
		}
		eh.Lock()
		if eh.triggered || eh.name != name {
			eh.Unlock()
			continue
		}
		eh.triggered = true
		h := eh.histogram
		eh.Unlock()
		return h
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("HistogramExt(%q)", name))
	return &mockHistogram{}
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Histogram'
```

Expected: both Histogram tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add Histogram and HistogramExt expectations"
```

---

### Task 5: Summary and SummaryExt expectations

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
func TestMockMeter_Summary(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectSummary("response_size")

	s := m.Summary("response_size")
	s.Update(100.0)
	s.Update(200.0)

	updates := exp.Summary().Updates()
	if len(updates) != 2 || updates[0] != 100.0 || updates[1] != 200.0 {
		t.Fatalf("expected [100.0, 200.0], got %v", updates)
	}
}

func TestMockMeter_SummaryExt(t *testing.T) {
	m := NewMockMeter()
	exp := m.ExpectSummaryExt("response_size_ext", 5*time.Minute, []float64{0.5, 0.9})

	s := m.SummaryExt("response_size_ext", 5*time.Minute, []float64{0.5, 0.9})
	s.Update(42.0)

	updates := exp.Summary().Updates()
	if len(updates) != 1 || updates[0] != 42.0 {
		t.Fatalf("expected [42.0], got %v", updates)
	}
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Summary'
```

Expected: `undefined: ExpectSummary` (compile error)

- [ ] **Step 3: Add Expected types and Expect methods to `mock.go`**

After `ExpectedHistogramExt`, add:

```go
// ExpectedSummary holds an expectation for meter.Summary.
type ExpectedSummary struct {
	commonExpectation
	name    string
	labels  []string
	summary *mockSummary
}

// Summary returns the mock summary for value inspection.
func (e *ExpectedSummary) Summary() *mockSummary { return e.summary }

func (e *ExpectedSummary) String() string {
	return fmt.Sprintf("ExpectedSummary(%q)", e.name)
}

// ExpectedSummaryExt holds an expectation for meter.SummaryExt.
type ExpectedSummaryExt struct {
	commonExpectation
	name      string
	window    time.Duration
	quantiles []float64
	labels    []string
	summary   *mockSummary
}

// Summary returns the mock summary for value inspection.
func (e *ExpectedSummaryExt) Summary() *mockSummary { return e.summary }

func (e *ExpectedSummaryExt) String() string {
	return fmt.Sprintf("ExpectedSummaryExt(%q)", e.name)
}
```

Add `ExpectSummary` and `ExpectSummaryExt` to `MockMeter`:

```go
// ExpectSummary registers an expectation that meter.Summary will be called with name.
func (m *MockMeter) ExpectSummary(name string, labels ...string) *ExpectedSummary {
	e := &ExpectedSummary{name: name, labels: labels, summary: &mockSummary{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectSummaryExt registers an expectation that meter.SummaryExt will be called with name.
func (m *MockMeter) ExpectSummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) *ExpectedSummaryExt {
	e := &ExpectedSummaryExt{name: name, window: window, quantiles: quantiles, labels: labels, summary: &mockSummary{}}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}
```

Replace `Summary` and `SummaryExt` method bodies:

```go
// Summary implements meter.Meter.
func (m *MockMeter) Summary(name string, _ ...string) meter.Summary {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		es, ok := e.(*ExpectedSummary)
		if !ok {
			continue
		}
		es.Lock()
		if es.triggered || es.name != name {
			es.Unlock()
			continue
		}
		es.triggered = true
		s := es.summary
		es.Unlock()
		return s
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("Summary(%q)", name))
	return &mockSummary{}
}

// SummaryExt implements meter.Meter.
func (m *MockMeter) SummaryExt(name string, _ time.Duration, _ []float64, _ ...string) meter.Summary {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		es, ok := e.(*ExpectedSummaryExt)
		if !ok {
			continue
		}
		es.Lock()
		if es.triggered || es.name != name {
			es.Unlock()
			continue
		}
		es.triggered = true
		s := es.summary
		es.Unlock()
		return s
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("SummaryExt(%q)", name))
	return &mockSummary{}
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Summary'
```

Expected: both Summary tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add Summary and SummaryExt expectations"
```

---

### Task 6: Write expectation

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
func TestMockMeter_Write_Success(t *testing.T) {
	m := NewMockMeter()
	m.ExpectWrite()
	if err := m.Write(io.Discard); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMockMeter_Write_Error(t *testing.T) {
	m := NewMockMeter()
	m.ExpectWrite().WillReturnError(errors.New("write failed"))
	if err := m.Write(io.Discard); err == nil || err.Error() != "write failed" {
		t.Fatalf("expected 'write failed', got %v", err)
	}
}

func TestMockMeter_Write_Unexpected(t *testing.T) {
	m := NewMockMeter()
	if err := m.Write(io.Discard); err == nil {
		t.Fatal("expected error for unexpected Write call")
	}
}
```

Add `"io"` to the import block in `mock_test.go`:

```go
import (
	"errors"
	"io"
	"testing"
	"time"

	"go.unistack.org/micro/v5/meter"
)
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Write'
```

Expected: `undefined: ExpectWrite` (compile error)

- [ ] **Step 3: Add `ExpectedWrite` and update `Write()` in `mock.go`**

After `ExpectedInit`, add:

```go
// ExpectedWrite holds an expectation for meter.Write.
type ExpectedWrite struct {
	commonExpectation
}

// WillReturnError configures the error returned by Write.
func (e *ExpectedWrite) WillReturnError(err error) *ExpectedWrite { e.err = err; return e }

func (e *ExpectedWrite) String() string {
	if e.err != nil {
		return fmt.Sprintf("ExpectedWrite => should return error: %s", e.err)
	}
	return "ExpectedWrite"
}
```

Add `ExpectWrite()` to `MockMeter`:

```go
// ExpectWrite registers an expectation that meter.Write will be called.
func (m *MockMeter) ExpectWrite() *ExpectedWrite {
	e := &ExpectedWrite{}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}
```

Replace the `Write` method body:

```go
// Write implements meter.Meter.
func (m *MockMeter) Write(_ io.Writer, _ ...meter.Option) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		ew, ok := e.(*ExpectedWrite)
		if !ok {
			continue
		}
		ew.Lock()
		if ew.triggered {
			ew.Unlock()
			continue
		}
		ew.triggered = true
		err := ew.err
		ew.Unlock()
		return err
	}
	return fmt.Errorf("unexpected call to meter.Write")
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Write'
```

Expected: all three Write tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add Write expectation"
```

---

### Task 7: Unregister expectation

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
func TestMockMeter_Unregister_True(t *testing.T) {
	m := NewMockMeter()
	m.ExpectUnregister("old_metric").WillReturn(true)
	if !m.Unregister("old_metric") {
		t.Fatal("expected true")
	}
}

func TestMockMeter_Unregister_False(t *testing.T) {
	m := NewMockMeter()
	m.ExpectUnregister("missing").WillReturn(false)
	if m.Unregister("missing") {
		t.Fatal("expected false")
	}
}

func TestMockMeter_Unregister_Unexpected(t *testing.T) {
	m := NewMockMeter()
	if m.Unregister("anything") {
		t.Fatal("expected false for unexpected call")
	}
	if len(m.unexpected) == 0 {
		t.Fatal("expected unexpected call to be recorded")
	}
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Unregister'
```

Expected: `undefined: ExpectUnregister` (compile error)

- [ ] **Step 3: Add `ExpectedUnregister` and update `Unregister()` in `mock.go`**

After `ExpectedWrite`, add:

```go
// ExpectedUnregister holds an expectation for meter.Unregister.
type ExpectedUnregister struct {
	commonExpectation
	name   string
	labels []string
	result bool
}

// WillReturn configures the bool returned by Unregister.
func (e *ExpectedUnregister) WillReturn(v bool) *ExpectedUnregister { e.result = v; return e }

func (e *ExpectedUnregister) String() string {
	return fmt.Sprintf("ExpectedUnregister(%q)", e.name)
}
```

Add `ExpectUnregister()` to `MockMeter`:

```go
// ExpectUnregister registers an expectation that meter.Unregister will be called with name.
func (m *MockMeter) ExpectUnregister(name string, labels ...string) *ExpectedUnregister {
	e := &ExpectedUnregister{name: name, labels: labels}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}
```

Replace the `Unregister` method body:

```go
// Unregister implements meter.Meter.
func (m *MockMeter) Unregister(name string, _ ...string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.expected {
		eu, ok := e.(*ExpectedUnregister)
		if !ok {
			continue
		}
		eu.Lock()
		if eu.triggered || eu.name != name {
			eu.Unlock()
			continue
		}
		eu.triggered = true
		result := eu.result
		eu.Unlock()
		return result
	}
	m.unexpected = append(m.unexpected, fmt.Sprintf("Unregister(%q)", name))
	return false
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_Unregister'
```

Expected: all three Unregister tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): add Unregister expectation"
```

---

### Task 8: ExpectationsWereMet

Implement `ExpectationsWereMet()` to report unfulfilled expectations and unexpected calls.

**Files:**
- Modify: `meter/mock/mock.go`
- Modify: `meter/mock/mock_test.go`

- [ ] **Step 1: Add tests to `mock_test.go`**

Append to `mock_test.go`:

```go
func TestMockMeter_ExpectationsWereMet_OK(t *testing.T) {
	m := NewMockMeter()
	m.ExpectInit()
	m.ExpectCounter("hits")

	_ = m.Init()
	_ = m.Counter("hits")

	if err := m.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMockMeter_ExpectationsWereMet_Unfulfilled(t *testing.T) {
	m := NewMockMeter()
	m.ExpectInit()
	// Never call Init — expectation remains unfulfilled.

	if err := m.ExpectationsWereMet(); err == nil {
		t.Fatal("expected error for unfulfilled expectation")
	}
}

func TestMockMeter_ExpectationsWereMet_Unexpected(t *testing.T) {
	m := NewMockMeter()
	_ = m.Counter("unregistered") // no expectation

	if err := m.ExpectationsWereMet(); err == nil {
		t.Fatal("expected error for unexpected call")
	}
}

func TestMockMeter_CloneSet(t *testing.T) {
	m := NewMockMeter(meter.Name("base"))
	c := m.Clone()
	if c == nil {
		t.Fatal("Clone returned nil")
	}
	s := m.Set()
	if s == nil {
		t.Fatal("Set returned nil")
	}
}
```

- [ ] **Step 2: Run tests — expect failure**

```bash
go test ./meter/mock/... -v -count=1 -run 'TestMockMeter_ExpectationsWereMet|TestMockMeter_CloneSet'
```

Expected: `TestMockMeter_ExpectationsWereMet_Unfulfilled` and `TestMockMeter_ExpectationsWereMet_Unexpected` fail because `ExpectationsWereMet` currently returns `nil`.

- [ ] **Step 3: Replace `ExpectationsWereMet()` body in `mock.go`**

```go
// ExpectationsWereMet returns an error if any declared expectation was not fulfilled
// or if there were unexpected calls.
func (m *MockMeter) ExpectationsWereMet() error {
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
	if len(m.unexpected) > 0 {
		return fmt.Errorf("there were unexpected calls: %v", m.unexpected)
	}
	return nil
}
```

- [ ] **Step 4: Run all tests — expect pass**

```bash
go test ./meter/mock/... -v -count=1
```

Expected: all tests PASS

- [ ] **Step 5: Commit**

```bash
git add meter/mock/mock.go meter/mock/mock_test.go
git commit -m "feat(meter/mock): implement ExpectationsWereMet"
```

---

## Self-review against spec

**Spec coverage:**

| Spec requirement | Covered by |
|-----------------|------------|
| `MockMeter` implementing `meter.Meter` | Task 1 |
| `ExpectInit` + `WillReturnError` | Task 2 |
| `ExpectWrite` + `WillReturnError` | Task 6 |
| `ExpectUnregister` + `WillReturn(bool)` | Task 7 |
| `ExpectCounter` + `Counter()` inspection | Task 3 |
| `ExpectFloatCounter` + `FloatCounter()` inspection | Task 3 |
| `ExpectGauge` + `Gauge()` inspection | Task 3 |
| `ExpectHistogram` + `Histogram()` inspection | Task 4 |
| `ExpectHistogramExt` + `Histogram()` inspection | Task 4 |
| `ExpectSummary` + `Summary()` inspection | Task 5 |
| `ExpectSummaryExt` + `Summary()` inspection | Task 5 |
| Unexpected calls recorded + reported | Tasks 1, 8 |
| `ExpectationsWereMet()` | Task 8 |
| `Clone()` / `Set()` without expectation | Task 8 |
| Compile-time `var _ meter.Meter` check | Task 1 |

**No placeholders found.**

**Type consistency:** `exp.Counter()` returns `*mockCounter`, `exp.Histogram()` returns `*mockHistogram`, `exp.Summary()` returns `*mockSummary` — consistent throughout all tasks. `Updates()` helper exists on both `mockHistogram` and `mockSummary`.
