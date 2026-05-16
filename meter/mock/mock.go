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

// Init records an unexpected call and returns an error.
func (m *MockMeter) Init(_ ...meter.Option) error {
	return fmt.Errorf("unexpected call to meter.Init")
}

// Write records an unexpected call and returns an error.
func (m *MockMeter) Write(_ io.Writer, _ ...meter.Option) error {
	return fmt.Errorf("unexpected call to meter.Write")
}

// Counter records the unexpected call and returns a stub counter.
func (m *MockMeter) Counter(name string, _ ...string) meter.Counter {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Counter(%q)", name))
	m.mu.Unlock()
	return &mockCounter{}
}

// FloatCounter records the unexpected call and returns a stub float counter.
func (m *MockMeter) FloatCounter(name string, _ ...string) meter.FloatCounter {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("FloatCounter(%q)", name))
	m.mu.Unlock()
	return &mockFloatCounter{}
}

// Gauge records the unexpected call and returns a stub gauge.
func (m *MockMeter) Gauge(name string, _ func() float64, _ ...string) meter.Gauge {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Gauge(%q)", name))
	m.mu.Unlock()
	return &mockGauge{}
}

// Histogram records the unexpected call and returns a stub histogram.
func (m *MockMeter) Histogram(name string, _ ...string) meter.Histogram {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Histogram(%q)", name))
	m.mu.Unlock()
	return &mockHistogram{}
}

// HistogramExt records the unexpected call and returns a stub histogram.
func (m *MockMeter) HistogramExt(name string, _ []float64, _ ...string) meter.Histogram {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("HistogramExt(%q)", name))
	m.mu.Unlock()
	return &mockHistogram{}
}

// Summary records the unexpected call and returns a stub summary.
func (m *MockMeter) Summary(name string, _ ...string) meter.Summary {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Summary(%q)", name))
	m.mu.Unlock()
	return &mockSummary{}
}

// SummaryExt records the unexpected call and returns a stub summary.
func (m *MockMeter) SummaryExt(name string, _ time.Duration, _ []float64, _ ...string) meter.Summary {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("SummaryExt(%q)", name))
	m.mu.Unlock()
	return &mockSummary{}
}

// Unregister records the unexpected call and returns false.
func (m *MockMeter) Unregister(name string, _ ...string) bool {
	m.mu.Lock()
	m.unexpected = append(m.unexpected, fmt.Sprintf("Unregister(%q)", name))
	m.mu.Unlock()
	return false
}

// ExpectationsWereMet verifies that all registered expectations were satisfied.
// Full implementation is added in Task 8; currently always returns nil.
func (m *MockMeter) ExpectationsWereMet() error {
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
