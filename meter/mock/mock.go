// Package mock provides a mock implementation of the meter.Meter interface
// for use in tests. It records unexpected calls and supports expectation-based
// verification in later tasks.
package mock

import (
	"fmt"
	"io"
	"sync"
	"time"

	"go.unistack.org/micro/v5/meter"
)

// Compile-time check that MockMeter satisfies the meter.Meter interface.
var _ meter.Meter = (*MockMeter)(nil)

// expectation is the interface that all expected call objects must satisfy.
type expectation interface {
	// fulfilled reports whether this expectation has been satisfied.
	fulfilled() bool
	// Lock locks the expectation's mutex.
	Lock()
	// Unlock unlocks the expectation's mutex.
	Unlock()
	// String returns a human-readable description of the expectation.
	String() string
}

// commonExpectation holds the shared state for all expectation types.
type commonExpectation struct {
	sync.Mutex
	triggered bool
	err       error
}

// fulfilled reports whether this expectation has been triggered at least once.
func (e *commonExpectation) fulfilled() bool {
	return e.triggered
}

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

// MockMeter is a mock implementation of meter.Meter for use in tests.
type MockMeter struct {
	opts       meter.Options
	mu         sync.Mutex
	expected   []expectation
	unexpected []string
}

// NewMockMeter creates and returns a new MockMeter configured with the given options.
func NewMockMeter(opts ...meter.Option) *MockMeter {
	return &MockMeter{
		opts: meter.NewOptions(opts...),
	}
}

// Name returns the meter name from its options.
func (m *MockMeter) Name() string {
	return m.opts.Name
}

// Options returns the current meter options.
func (m *MockMeter) Options() meter.Options {
	return m.opts
}

// String returns the string identifier for this meter implementation.
func (m *MockMeter) String() string {
	return "mock"
}

// Clone returns m unchanged. Expectation dispatch is added in a later task.
func (m *MockMeter) Clone(_ ...meter.Option) meter.Meter {
	return m
}

// Set returns m unchanged. Expectation dispatch is added in a later task.
func (m *MockMeter) Set(_ ...meter.Option) meter.Meter {
	return m
}

// ExpectWrite registers an expectation that meter.Write will be called.
func (m *MockMeter) ExpectWrite() *ExpectedWrite {
	e := &ExpectedWrite{}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

// ExpectInit registers an expectation that meter.Init will be called.
func (m *MockMeter) ExpectInit() *ExpectedInit {
	e := &ExpectedInit{}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

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

// ExpectUnregister registers an expectation that meter.Unregister will be called with name.
func (m *MockMeter) ExpectUnregister(name string, labels ...string) *ExpectedUnregister {
	e := &ExpectedUnregister{name: name, labels: labels}
	m.mu.Lock()
	m.expected = append(m.expected, e)
	m.mu.Unlock()
	return e
}

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

// mockCounter is a thread-safe stub implementation of meter.Counter.
type mockCounter struct {
	mu    sync.Mutex
	value uint64
}

// Add increments the counter by delta. Negative delta is treated as zero.
func (c *mockCounter) Add(delta int) {
	c.mu.Lock()
	if delta > 0 {
		c.value += uint64(delta)
	}
	c.mu.Unlock()
}

// Dec decrements the counter by one, guarding against underflow.
func (c *mockCounter) Dec() {
	c.mu.Lock()
	if c.value > 0 {
		c.value--
	}
	c.mu.Unlock()
}

// Inc increments the counter by one.
func (c *mockCounter) Inc() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

// Get returns the current counter value.
func (c *mockCounter) Get() uint64 {
	c.mu.Lock()
	v := c.value
	c.mu.Unlock()
	return v
}

// Set assigns an explicit value to the counter.
func (c *mockCounter) Set(v uint64) {
	c.mu.Lock()
	c.value = v
	c.mu.Unlock()
}

// mockFloatCounter is a thread-safe stub implementation of meter.FloatCounter.
type mockFloatCounter struct {
	mu    sync.Mutex
	value float64
}

// Add adds delta to the counter.
func (c *mockFloatCounter) Add(delta float64) {
	c.mu.Lock()
	c.value += delta
	c.mu.Unlock()
}

// Sub subtracts delta from the counter.
func (c *mockFloatCounter) Sub(delta float64) {
	c.mu.Lock()
	c.value -= delta
	c.mu.Unlock()
}

// Get returns the current counter value.
func (c *mockFloatCounter) Get() float64 {
	c.mu.Lock()
	v := c.value
	c.mu.Unlock()
	return v
}

// Set assigns an explicit value to the counter.
func (c *mockFloatCounter) Set(v float64) {
	c.mu.Lock()
	c.value = v
	c.mu.Unlock()
}

// mockGauge is a thread-safe stub implementation of meter.Gauge.
type mockGauge struct {
	mu    sync.Mutex
	value float64
}

// Add adds delta to the gauge.
func (g *mockGauge) Add(delta float64) {
	g.mu.Lock()
	g.value += delta
	g.mu.Unlock()
}

// Set assigns an explicit value to the gauge.
func (g *mockGauge) Set(v float64) {
	g.mu.Lock()
	g.value = v
	g.mu.Unlock()
}

// Inc increments the gauge by one.
func (g *mockGauge) Inc() {
	g.mu.Lock()
	g.value++
	g.mu.Unlock()
}

// Dec decrements the gauge by one.
func (g *mockGauge) Dec() {
	g.mu.Lock()
	g.value--
	g.mu.Unlock()
}

// Get returns the current gauge value.
func (g *mockGauge) Get() float64 {
	g.mu.Lock()
	v := g.value
	g.mu.Unlock()
	return v
}

// mockHistogram is a thread-safe stub implementation of meter.Histogram.
type mockHistogram struct {
	mu      sync.Mutex
	updates []float64
}

// Reset clears all recorded updates.
func (h *mockHistogram) Reset() {
	h.mu.Lock()
	h.updates = nil
	h.mu.Unlock()
}

// Update records a new observation.
func (h *mockHistogram) Update(v float64) {
	h.mu.Lock()
	h.updates = append(h.updates, v)
	h.mu.Unlock()
}

// UpdateDuration records the elapsed seconds since t as a new observation.
func (h *mockHistogram) UpdateDuration(t time.Time) {
	h.Update(time.Since(t).Seconds())
}

// Updates returns a copy of all recorded observations.
func (h *mockHistogram) Updates() []float64 {
	h.mu.Lock()
	cp := make([]float64, len(h.updates))
	copy(cp, h.updates)
	h.mu.Unlock()
	return cp
}

// mockSummary is a thread-safe stub implementation of meter.Summary.
type mockSummary struct {
	mu      sync.Mutex
	updates []float64
}

// Update records a new observation.
func (s *mockSummary) Update(v float64) {
	s.mu.Lock()
	s.updates = append(s.updates, v)
	s.mu.Unlock()
}

// UpdateDuration records the elapsed seconds since t as a new observation.
func (s *mockSummary) UpdateDuration(t time.Time) {
	s.Update(time.Since(t).Seconds())
}

// Updates returns a copy of all recorded observations.
func (s *mockSummary) Updates() []float64 {
	s.mu.Lock()
	cp := make([]float64, len(s.updates))
	copy(cp, s.updates)
	s.mu.Unlock()
	return cp
}
