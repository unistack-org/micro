package tracer

import "go.unistack.org/micro/v3/logger"

// SpanOptions contains span option
type SpanOptions struct{}

// SpanOption func signature
type SpanOption func(o *SpanOptions)

// EventOptions contains event options
type EventOptions struct{}

// EventOption func signature
type EventOption func(o *EventOptions)

// Options struct
type Options struct {
	// Logger used for logging
	Logger logger.Logger
	// Name of the tracer
	Name string
}

// Option func signature
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
