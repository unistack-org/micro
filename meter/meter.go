// Package meter is for instrumentation
package meter

import (
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	DefaultReporter Meter = NewMeter()
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
