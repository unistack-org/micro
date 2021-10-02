package tracer

import "go.unistack.org/micro/v3/logger"

type SpanOptions struct{}

type SpanOption func(o *SpanOptions)

type EventOptions struct{}

type EventOption func(o *EventOptions)

// Options struct
type Options struct {
	// Logger used for logging
	Logger logger.Logger
	// Name of the tracer
	Name string
}

// Option func
type Option func(o *Options)

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// NewOptions returns default options
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger: logger.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
