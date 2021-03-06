// Package meter is for instrumentation
package meter

import (
	"io"
	"reflect"
	"sort"
	"time"
)

var (
	// DefaultMeter is the default meter
	DefaultMeter Meter = NewMeter()
	// DefaultAddress data will be made available on this host:port
	DefaultAddress = ":9090"
	// DefaultPath the meter endpoint where the Meter data will be made available
	DefaultPath = "/metrics"
	// DefaultMetricPrefix holds the string that prepends to all metrics
	DefaultMetricPrefix = "micro_"
	// DefaultLabelPrefix holds the string that prepends to all labels
	DefaultLabelPrefix = "micro_"
	// DefaultSummaryQuantiles is the default spread of stats for summary
	DefaultSummaryQuantiles = []float64{0.5, 0.9, 0.97, 0.99, 1}
	// DefaultSummaryWindow is the default window for summary
	DefaultSummaryWindow = 5 * time.Minute
)

// Meter is an interface for collecting and instrumenting metrics
type Meter interface {
	Name() string
	Init(opts ...Option) error
	Counter(name string, labels ...string) Counter
	FloatCounter(name string, labels ...string) FloatCounter
	Gauge(name string, fn func() float64, labels ...string) Gauge
	Set(opts ...Option) Meter
	Histogram(name string, labels ...string) Histogram
	Summary(name string, labels ...string) Summary
	SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) Summary
	Write(w io.Writer, opts ...Option) error
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

type byKey []string

func (k byKey) Len() int           { return len(k) / 2 }
func (k byKey) Less(i, j int) bool { return k[i*2] < k[j*2] }
func (k byKey) Swap(i, j int) {
	k[i*2], k[i*2+1], k[j*2], k[j*2+1] = k[j*2], k[j*2+1], k[i*2], k[i*2+1]
}

func Sort(slice *[]string) {
	bk := byKey(*slice)
	if bk.Len() <= 1 {
		return
	}
	sort.Sort(bk)
	v := reflect.ValueOf(slice).Elem()
	cnt := 0
	key := 0
	val := 1
	for key < v.Len() {
		if len(bk) > key+2 && bk[key] == bk[key+2] {
			key += 2
			val += 2
			continue
		}
		v.Index(cnt).Set(v.Index(key))
		cnt++
		v.Index(cnt).Set(v.Index(val))
		cnt++
		key += 2
		val += 2
	}
	v.SetLen(cnt)
}
