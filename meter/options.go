package meter

import (
	"context"

	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
)

// Option powers the configuration for metrics implementations:
type Option func(*Options)

// Options for metrics implementations:
type Options struct {
	Address  string
	Path     string
	Metadata metadata.Metadata
	//TimingObjectives map[float64]float64
	Logger       logger.Logger
	Context      context.Context
	MetricPrefix string
	LabelPrefix  string
}

// NewOptions prepares a set of options:
func NewOptions(opt ...Option) Options {
	opts := Options{
		Address:      DefaultAddress,
		Metadata:     metadata.New(3), // 3 elements contains service name, version and id
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

// Metadata will be added to every metric
func Metadata(md metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = metadata.Copy(md)
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
