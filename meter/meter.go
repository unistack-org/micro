// Package meter is for instrumentation
package meter

import (
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	// DefaultMeter is the default meter
	DefaultMeter Meter = NewMeter()
	// DefaultAddress data will be made available on this host:port
	DefaultAddress = ":9090"
	// DefaultPath the meter endpoint where the Meter data will be made available
	DefaultPath = "/metrics"
	// timingObjectives is the default spread of stats we maintain for timings / histograms:
	//defaultTimingObjectives = map[float64]float64{0.0: 0, 0.5: 0.05, 0.75: 0.04, 0.90: 0.03, 0.95: 0.02, 0.98: 0.001, 1: 0}
	// default metric prefix
	DefaultMetricPrefix = "micro_"
	// default label prefix
	DefaultLabelPrefix = "micro_"
)

// Meter is an interface for collecting and instrumenting metrics
type Meter interface {
	Init(...Option) error
	Counter(string, metadata.Metadata) Counter
	FloatCounter(string, metadata.Metadata) FloatCounter
	Gauge(string, func() float64, metadata.Metadata) Gauge
	Set(metadata.Metadata) Meter
	Histogram(string, metadata.Metadata) Histogram
	Summary(string, metadata.Metadata) Summary
	SummaryExt(string, time.Duration, []float64, metadata.Metadata) Summary
	Options() Options
	String() string
}

// Counter is a counter
type Counter interface {
	Add(int)
	Dec()
	Get() uint64
	Inc()
	Set(uint64)
}

// FloatCounter is a float64 counter
type FloatCounter interface {
	Add(float64)
	Get() float64
	Set(float64)
	Sub(float64)
}

// Gauge is a float64 gauge
type Gauge interface {
	Get() float64
}

// Histogram is a histogram for non-negative values with automatically created buckets
type Histogram interface {
	Reset()
	Update(float64)
	UpdateDuration(time.Time)
	// VisitNonZeroBuckets(f func(vmrange string, count uint64))
}

// Summary is the summary
type Summary interface {
	Update(float64)
	UpdateDuration(time.Time)
}
