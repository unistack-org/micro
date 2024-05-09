package meter

import (
	"context"
)

// Option powers the configuration for metrics implementations:
type Option func(*Options)

// Options for metrics implementations
type Options struct {
	// Context holds external options
	Context context.Context
	// Name holds the meter name
	Name string
	// Address holds the address that serves metrics
	Address string
	// Path holds the path for metrics
	Path string
	// MetricPrefix holds the prefix for all metrics
	MetricPrefix string
	// LabelPrefix holds the prefix for all labels
	LabelPrefix string
	// Labels holds the default labels
	Labels []string
	// WriteProcessMetrics flag to write process metrics
	WriteProcessMetrics bool
	// WriteFDMetrics flag to write fd metrics
	WriteFDMetrics bool
}

// NewOptions prepares a set of options:
func NewOptions(opt ...Option) Options {
	opts := Options{
		Address:      DefaultAddress,
		Path:         DefaultPath,
		Context:      context.Background(),
		MetricPrefix: DefaultMetricPrefix,
		LabelPrefix:  DefaultLabelPrefix,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// LabelPrefix sets the labels prefix
func LabelPrefix(pref string) Option {
	return func(o *Options) {
		o.LabelPrefix = pref
	}
}

// MetricPrefix sets the metric prefix
func MetricPrefix(pref string) Option {
	return func(o *Options) {
		o.MetricPrefix = pref
	}
}

// Context sets the metrics context
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Path used to serve metrics over HTTP
func Path(value string) Option {
	return func(o *Options) {
		o.Path = value
	}
}

// Address is the listen address to serve metrics
func Address(value string) Option {
	return func(o *Options) {
		o.Address = value
	}
}

/*
// TimingObjectives defines the desired spread of statistics for histogram / timing metrics:
func TimingObjectives(value map[float64]float64) Option {
	return func(o *Options) {
		o.TimingObjectives = value
	}
}
*/

// Labels sets the meter labels
func Labels(ls ...string) Option {
	return func(o *Options) {
		o.Labels = append(o.Labels, ls...)
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// WriteProcessMetrics enable process metrics output for write
func WriteProcessMetrics(b bool) Option {
	return func(o *Options) {
		o.WriteProcessMetrics = b
	}
}

// WriteFDMetrics enable fd metrics output for write
func WriteFDMetrics(b bool) Option {
	return func(o *Options) {
		o.WriteFDMetrics = b
	}
}
