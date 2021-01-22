package meter

import (
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

// NoopMeter is an noop implementation of Meter
type noopMeter struct {
	opts Options
	md   metadata.Metadata
}

// NewMeter returns a configured noop reporter:
func NewMeter(opts ...Option) Meter {
	return &noopMeter{
		opts: NewOptions(opts...),
	}
}

// Init initialize options
func (r *noopMeter) Init(opts ...Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

// Counter implements the Meter interface
func (r *noopMeter) Counter(name string, md metadata.Metadata) Counter {
	return &noopCounter{}
}

// FloatCounter implements the Meter interface
func (r *noopMeter) FloatCounter(name string, md metadata.Metadata) FloatCounter {
	return &noopFloatCounter{}
}

// Gauge implements the Meter interface
func (r *noopMeter) Gauge(name string, f func() float64, md metadata.Metadata) Gauge {
	return &noopGauge{}
}

// Summary implements the Meter interface
func (r *noopMeter) Summary(name string, md metadata.Metadata) Summary {
	return &noopSummary{}
}

// SummaryExt implements the Meter interface
func (r *noopMeter) SummaryExt(name string, window time.Duration, quantiles []float64, md metadata.Metadata) Summary {
	return &noopSummary{}
}

// Histogram implements the Meter interface
func (r *noopMeter) Histogram(name string, md metadata.Metadata) Histogram {
	return &noopHistogram{}
}

// Set implements the Meter interface
func (r *noopMeter) Set(md metadata.Metadata) Meter {
	return &noopMeter{opts: r.opts, md: metadata.Copy(md)}
}

// Options implements the Meter interface
func (r *noopMeter) Options() Options {
	return r.opts
}

// String implements the Meter interface
func (r *noopMeter) String() string {
	return "noop"
}

type noopCounter struct{}

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

type noopFloatCounter struct{}

func (r *noopFloatCounter) Add(float64) {

}

func (r *noopFloatCounter) Get() float64 {
	return 0
}

func (r *noopFloatCounter) Set(float64) {

}

func (r *noopFloatCounter) Sub(float64) {

}

type noopGauge struct{}

func (r *noopGauge) Get() float64 {
	return 0
}

type noopSummary struct{}

func (r *noopSummary) Update(float64) {

}

func (r *noopSummary) UpdateDuration(time.Time) {

}

type noopHistogram struct{}

func (r *noopHistogram) Reset() {

}

func (r *noopHistogram) Update(float64) {

}

func (r *noopHistogram) UpdateDuration(time.Time) {

}

//func (r *noopHistogram) VisitNonZeroBuckets(f func(vmrange string, count uint64)) {}
