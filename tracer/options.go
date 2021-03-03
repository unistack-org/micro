package tracer

import "github.com/unistack-org/micro/v3/logger"

type SpanOptions struct {
}

type SpanOption func(o *SpanOptions)

type EventOptions struct {
}

type EventOption func(o *EventOptions)

// Options struct
type Options struct {
	// Name of the tracer
	Name string
	// Logger is the logger for messages
	Logger logger.Logger
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
