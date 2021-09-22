package codec

import (
	"context"

	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
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
	// MaxMsgSize specifies max messages size that reads by codec
	MaxMsgSize int
}

// MaxMsgSize sets the max message size
func MaxMsgSize(n int) Option {
	return func(o *Options) {
		o.MaxMsgSize = n
	}
}

// TagName sets the codec tag name in struct
func TagName(n string) Option {
	return func(o *Options) {
		o.TagName = n
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
		Context:    context.Background(),
		Logger:     logger.DefaultLogger,
		Meter:      meter.DefaultMeter,
		Tracer:     tracer.DefaultTracer,
		MaxMsgSize: DefaultMaxMsgSize,
		TagName:    DefaultTagName,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
