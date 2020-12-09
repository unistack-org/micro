package config

import (
	"context"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
)

type Options struct {
	BeforeLoad []func(context.Context, Config) error
	AfterLoad  []func(context.Context, Config) error
	BeforeSave []func(context.Context, Config) error
	AfterSave  []func(context.Context, Config) error
	// Struct that holds config data
	Struct interface{}
	// struct tag name
	StructTag string
	// logger that will be used
	Logger logger.Logger
	// codec that used for load/save
	Codec codec.Codec
	// for alternative data
	Context context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}

	return options
}

func BeforeLoad(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.BeforeLoad = fn
	}
}

func AfterLoad(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.AfterLoad = fn
	}
}

func BeforeSave(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.BeforeSave = fn
	}
}

func AfterSave(fn ...func(context.Context, Config) error) Option {
	return func(o *Options) {
		o.AfterSave = fn
	}
}

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

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
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
