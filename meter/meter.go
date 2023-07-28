// Package meter is for instrumentation
package meter // import "go.unistack.org/micro/v4/meter"

import (
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.unistack.org/micro/v4/options"
)

var (
	// DefaultMeter is the default meter
	DefaultMeter = NewMeter()
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
	// Name returns meter name
	Name() string
	// Init initialize meter
	Init(opts ...options.Option) error
	// Clone create meter copy with new options
	Clone(opts ...options.Option) Meter
	// Counter get or create counter
	Counter(name string, labels ...string) Counter
	// FloatCounter get or create float counter
	FloatCounter(name string, labels ...string) FloatCounter
	// Gauge get or create gauge
	Gauge(name string, fn func() float64, labels ...string) Gauge
	// Set create new meter metrics set
	Set(opts ...options.Option) Meter
	// Histogram get or create histogram
	Histogram(name string, labels ...string) Histogram
	// Summary get or create summary
	Summary(name string, labels ...string) Summary
	// SummaryExt get or create summary with spcified quantiles and window time
	SummaryExt(name string, window time.Duration, quantiles []float64, labels ...string) Summary
	// Write writes metrics to io.Writer
	Write(w io.Writer, opts ...options.Option) error
	// Options returns meter options
	Options() Options
	// String return meter type
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

// sort labels alphabeticaly by label name
type byKey []string

func (k byKey) Len() int           { return len(k) / 2 }
func (k byKey) Less(i, j int) bool { return k[i*2] < k[j*2] }
func (k byKey) Swap(i, j int) {
	k[i*2], k[j*2] = k[j*2], k[i*2]
	k[i*2+1], k[j*2+1] = k[j*2+1], k[i*2+1]
}

// BuildLabels used to sort labels and delete duplicates.
// Last value wins in case of duplicate label keys.
func BuildLabels(labels ...string) []string {
	if len(labels)%2 == 1 {
		labels = labels[:len(labels)-1]
	}
	sort.Sort(byKey(labels))
	return labels
}

// BuildName used to combine metric with labels.
// If labels count is odd, drop last element
func BuildName(name string, labels ...string) string {
	if len(labels)%2 == 1 {
		labels = labels[:len(labels)-1]
	}

	if len(labels) > 2 {
		sort.Sort(byKey(labels))

		idx := 0
		for {
			if labels[idx] == labels[idx+2] {
				copy(labels[idx:], labels[idx+2:])
				labels = labels[:len(labels)-2]
			} else {
				idx += 2
			}
			if idx+2 >= len(labels) {
				break
			}
		}
	}

	var b strings.Builder
	_, _ = b.WriteString(name)
	_, _ = b.WriteRune('{')
	for idx := 0; idx < len(labels); idx += 2 {
		if idx > 0 {
			_, _ = b.WriteRune(',')
		}
		_, _ = b.WriteString(labels[idx])
		_, _ = b.WriteString(`=`)
		_, _ = b.WriteString(strconv.Quote(labels[idx+1]))
	}
	_, _ = b.WriteRune('}')

	return b.String()
}
