package codec

import (
	"context"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/tracer"
)

// Option func
type Option func(*Options)

// Options contains codec options
type Options struct {
	// Meter used for metrics
	Meter meter.Meter
	// Logger used for logging
	Logger logger.Logger
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Context stores additional codec options
	Context context.Context
	// TagName specifies tag name in struct to control codec
	TagName string
	// Flatten specifies that struct must be analyzed for flatten tag
	Flatten bool
}

// TagName sets the codec tag name in struct
func TagName(n string) Option {
	return func(o *Options) {
		o.TagName = n
	}
}

// Flatten enables checking for flatten tag name
func Flatten(b bool) Option {
	return func(o *Options) {
		o.Flatten = b
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Tracer to be used for tracing
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Meter sets the meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// NewOptions returns new options
func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
		TagName: DefaultTagName,
		Flatten: false,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
