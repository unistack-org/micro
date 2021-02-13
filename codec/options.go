package codec

import (
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
)

// Option func
type Option func(*Options)

// Options contains codec options
type Options struct {
	MaxMsgSize int
	Meter      meter.Meter
	Logger     logger.Logger
	Tracer     tracer.Tracer
}

// MaxMsgSize sets the max message size
func MaxMsgSize(n int) Option {
	return func(o *Options) {
		o.MaxMsgSize = n
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
		Logger:     logger.DefaultLogger,
		Meter:      meter.DefaultMeter,
		Tracer:     tracer.DefaultTracer,
		MaxMsgSize: DefaultMaxMsgSize,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
