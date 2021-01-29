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
	// timingObjectives is the default spread of stats we maintain for timings / histograms:
	//defaultTimingObjectives = map[float64]float64{0.0: 0, 0.5: 0.05, 0.75: 0.04, 0.90: 0.03, 0.95: 0.02, 0.98: 0.001, 1: 0}
	// default metric prefix
	DefaultMetricPrefix = "micro_"
	// default label prefix
	DefaultLabelPrefix = "micro_"
)

// Meter is an interface for collecting and instrumenting metrics
type Meter interface {
	Name() string
	Init(...Option) error
	Counter(string, ...Option) Counter
	FloatCounter(string, ...Option) FloatCounter
	Gauge(string, func() float64, ...Option) Gauge
	Set(...Option) Meter
	Histogram(string, ...Option) Histogram
	Summary(string, ...Option) Summary
	SummaryExt(string, time.Duration, []float64, ...Option) Summary
	Write(io.Writer, bool) error
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

type Labels struct {
	keys []string
	vals []string
}

func (ls Labels) Len() int {
	return len(ls.keys)
}

func (ls Labels) Swap(i, j int) {
	ls.keys[i], ls.keys[j] = ls.keys[j], ls.keys[i]
	ls.vals[i], ls.vals[j] = ls.vals[j], ls.vals[i]
}

func (ls Labels) Less(i, j int) bool {
	return ls.vals[i] < ls.vals[j]
}

func (ls Labels) Sort() {
	sort.Sort(ls)
}

func (ls Labels) Append(nls Labels) Labels {
	for n := range nls.keys {
		ls.keys = append(ls.keys, nls.keys[n])
		ls.vals = append(ls.vals, nls.vals[n])
	}
	return ls
}

type LabelIter struct {
	labels Labels
	cnt    int
	cur    int
}

func (ls Labels) Iter() *LabelIter {
	ls.Sort()
	return &LabelIter{labels: ls, cnt: len(ls.keys)}
}

func (iter *LabelIter) Next(k, v *string) bool {
	if iter.cur+1 > iter.cnt {
		return false
	}

	*k = iter.labels.keys[iter.cur]
	*v = iter.labels.vals[iter.cur]
	iter.cur++
	return true
}
