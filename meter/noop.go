package meter

import (
	"io"
	"time"

	"go.unistack.org/micro/v4/options"
)

// NoopMeter is an noop implementation of Meter
type noopMeter struct {
	opts Options
}

// NewMeter returns a configured noop reporter:
func NewMeter(opts ...options.Option) Meter {
	return &noopMeter{opts: NewOptions(opts...)}
}

// Clone return old meter with new options
func (r *noopMeter) Clone(opts ...options.Option) Meter {
	options := r.opts
	for _, o := range opts {
		o(&options)
	}
	return &noopMeter{opts: options}
}

func (r *noopMeter) Name() string {
	return r.opts.Name
}

// Init initialize options
func (r *noopMeter) Init(opts ...options.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

// Counter implements the Meter interface
func (r *noopMeter) Counter(name string, labels ...string) Counter {
	return &noopCounter{labels: labels}
}

// FloatCounter implements the Meter interface
func (r *noopMeter) FloatCounter(name string, labels ...string) FloatCounter {
	return &noopFloatCounter{labels: labels}
}

// Gauge implements the Meter interface
func (r *noopMeter) Gauge(name string, f func() float64, labels ...string) Gauge {
	return &noopGauge{labels: labels}
}

// Summary implements the Meter interface
func (r *noopMeter) Summary(name string, labels ...string) Summary {
	return &noopSummary{labels: labels}
}

// SummaryExt implements the Meter interface
func (r *noopMeter) SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) Summary {
	return &noopSummary{labels: labels}
}

// Histogram implements the Meter interface
func (r *noopMeter) Histogram(name string, labels ...string) Histogram {
	return &noopHistogram{labels: labels}
}

// Set implements the Meter interface
func (r *noopMeter) Set(opts ...options.Option) Meter {
	m := &noopMeter{opts: r.opts}

	for _, o := range opts {
		o(&m.opts)
	}

	return m
}

func (r *noopMeter) Write(_ io.Writer, _ ...options.Option) error {
	return nil
}

// Options implements the Meter interface
func (r *noopMeter) Options() Options {
	return r.opts
}

// String implements the Meter interface
func (r *noopMeter) String() string {
	return "noop"
}

type noopCounter struct {
	labels []string
}

func (r *noopCounter) Add(int) {
}

func (r *noopCounter) Dec() {
}

func (r *noopCounter) Get() uint64 {
	return 0
}

func (r *noopCounter) Inc() {
}

func (r *noopCounter) Set(uint64) {
}

type noopFloatCounter struct {
	labels []string
}

func (r *noopFloatCounter) Add(float64) {
}

func (r *noopFloatCounter) Get() float64 {
	return 0
}

func (r *noopFloatCounter) Set(float64) {
}

func (r *noopFloatCounter) Sub(float64) {
}

type noopGauge struct {
	labels []string
}

func (r *noopGauge) Get() float64 {
	return 0
}

type noopSummary struct {
	labels []string
}

func (r *noopSummary) Update(float64) {
}

func (r *noopSummary) UpdateDuration(time.Time) {
}

type noopHistogram struct {
	labels []string
}

func (r *noopHistogram) Reset() {
}

func (r *noopHistogram) Update(float64) {
}

func (r *noopHistogram) UpdateDuration(time.Time) {
}

// func (r *noopHistogram) VisitNonZeroBuckets(f func(vmrange string, count uint64)) {}
