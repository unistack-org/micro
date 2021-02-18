// Package meter is for instrumentation
package meter

import (
	"io"
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
	Counter(name string, opts ...Option) Counter
	FloatCounter(name string, opts ...Option) FloatCounter
	Gauge(name string, fn func() float64, opts ...Option) Gauge
	Set(opts ...Option) Meter
	Histogram(name string, opts ...Option) Histogram
	Summary(name string, opts ...Option) Summary
	SummaryExt(name string, window time.Duration, quantiles []float64, opts ...Option) Summary
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

// Labels holds the metrics labels with k, v
type Labels struct {
	keys []string
	vals []string
}

type labels Labels

func (ls labels) sort() {
	sort.Sort(ls)
}

func (ls labels) Len() int {
	return len(ls.keys)
}

func (ls labels) Swap(i, j int) {
	ls.keys[i], ls.keys[j] = ls.keys[j], ls.keys[i]
	ls.vals[i], ls.vals[j] = ls.vals[j], ls.vals[i]
}

func (ls labels) Less(i, j int) bool {
	return ls.vals[i] < ls.vals[j]
}

// Append adds labels to label set
func (ls Labels) Append(nls Labels) Labels {
	for n := range nls.keys {
		ls.keys = append(ls.keys, nls.keys[n])
		ls.vals = append(ls.vals, nls.vals[n])
	}
	return ls
}

func (ls Labels) Len() int {
	return len(ls.keys)
}

// LabelIter holds the
type LabelIter struct {
	labels Labels
	cnt    int
	cur    int
}

// Iter returns labels iterator
func (ls Labels) Iter() *LabelIter {
	labels(ls).sort()
	return &LabelIter{labels: ls, cnt: len(ls.keys)}
}

// Next advance itarator to new pos
func (iter *LabelIter) Next(k, v *string) bool {
	if iter.cur+1 > iter.cnt {
		return false
	}

	*k = iter.labels.keys[iter.cur]
	*v = iter.labels.vals[iter.cur]
	iter.cur++
	return true
}
