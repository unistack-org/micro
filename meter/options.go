package meter

import (
	"context"

	"github.com/unistack-org/micro/v3/logger"
)

// Option powers the configuration for metrics implementations:
type Option func(*Options)

// Options for metrics implementations:
type Options struct {
	Name    string
	Address string
	Path    string
	Labels  Labels
	//TimingObjectives map[float64]float64
	Logger              logger.Logger
	Context             context.Context
	MetricPrefix        string
	LabelPrefix         string
	WriteProcessMetrics bool
	WriteFDMetrics      bool
}

// NewOptions prepares a set of options:
func NewOptions(opt ...Option) Options {
	opts := Options{
		Address:      DefaultAddress,
		Path:         DefaultPath,
		Context:      context.Background(),
		Logger:       logger.DefaultLogger,
		MetricPrefix: DefaultMetricPrefix,
		LabelPrefix:  DefaultLabelPrefix,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
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

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Label sets the label
func Label(key, val string) Option {
	return func(o *Options) {
		o.Labels.keys = append(o.Labels.keys, key)
		o.Labels.vals = append(o.Labels.vals, val)
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
