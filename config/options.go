package config

import (
	"context"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options hold the config options
type Options struct {
	Name       string
	AllowFail  bool
	BeforeLoad []func(context.Context, Config) error
	AfterLoad  []func(context.Context, Config) error
	BeforeSave []func(context.Context, Config) error
	AfterSave  []func(context.Context, Config) error
	// Struct that holds config data
	Struct interface{}
	// StructTag name
	StructTag string
	// Logger that will be used
	Logger logger.Logger
	// Meter that will be used
	Meter meter.Meter
	// Tracer used for trace
	Tracer tracer.Tracer
	// Codec that used for load/save
	Codec codec.Codec
	// Context for alternative data
	Context context.Context
}

// Option function signature
type Option func(o *Options)

// NewOptions new options struct with filed values
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}

	return options
}

// AllowFail allows config source to fail
func AllowFail(b bool) Option {
	return func(o *Options) {
		o.AllowFail = b
	}
}

// BeforeLoad run funcs before config load
func BeforeLoad(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.BeforeLoad = fn
	}
}

// AfterLoad run funcs after config load
func AfterLoad(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.AfterLoad = fn
	}
}

// BeforeSave run funcs before save
func BeforeSave(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.BeforeSave = fn
	}
}

// AfterSave run fncs after save
func AfterSave(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.AfterSave = fn
	}
}

// Context pass context
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Codec sets the source codec
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
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

// Struct used as config
func Struct(v interface{}) Option {
	return func(o *Options) {
		o.Struct = v
	}
}

// StructTag
func StructTag(name string) Option {
	return func(o *Options) {
		o.StructTag = name
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
