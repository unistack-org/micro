package metrics

import (
	"time"
)

// NoopReporter is an noop implementation of Reporter:
type noopReporter struct {
	opts Options
}

// NewReporter returns a configured noop reporter:
func NewReporter(opts ...Option) Reporter {
	return &noopReporter{
		opts: NewOptions(opts...),
	}
}

// Init initialize options
func (r *noopReporter) Init(opts ...Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

// Count implements the Reporter interface Count method:
func (r *noopReporter) Count(metricName string, value int64, tags Tags) error {
	return nil
}

// Gauge implements the Reporter interface Gauge method:
func (r *noopReporter) Gauge(metricName string, value float64, tags Tags) error {
	return nil
}

// Timing implements the Reporter interface Timing method:
func (r *noopReporter) Timing(metricName string, value time.Duration, tags Tags) error {
	return nil
}

// Options implements the Reporter interface Optios method:
func (r *noopReporter) Options() Options {
	return r.opts
}
