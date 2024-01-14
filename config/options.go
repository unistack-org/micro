package config

import (
	"context"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/tracer"
)

// Options hold the config options
type Options struct {
	// Struct holds the destination config struct
	Struct interface{}
	// Codec that used for load/save
	Codec codec.Codec
	// Tracer that will be used
	Tracer tracer.Tracer
	// Meter that will be used
	Meter meter.Meter
	// Logger that will be used
	Logger logger.Logger
	// Context used for external options
	Context context.Context
	// Name of the config
	Name string
	// StructTag name
	StructTag string
	// BeforeSave contains slice of funcs that runs before Save
	BeforeSave []func(context.Context, Config) error
	// AfterSave contains slice of funcs that runs after Save
	AfterSave []func(context.Context, Config) error
	// BeforeLoad contains slice of funcs that runs before Load
	BeforeLoad []func(context.Context, Config) error
	// AfterLoad contains slice of funcs that runs after Load
	AfterLoad []func(context.Context, Config) error
	// BeforeInit contains slice of funcs that runs before Init
	BeforeInit []func(context.Context, Config) error
	// AfterInit contains slice of funcs that runs after Init
	AfterInit []func(context.Context, Config) error
	// AllowFail flag to allow fail in config source
	AllowFail bool
	// SkipLoad runs only if condition returns true
	SkipLoad func(context.Context, Config) bool
	// SkipSave runs only if condition returns true
	SkipSave func(context.Context, Config) bool
}

// NewOptions new options struct with filed values
func NewOptions(opts ...options.Option) Options {
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

// LoadOptions struct
type LoadOptions struct {
	Struct   interface{}
	Context  context.Context
	Override bool
	Append   bool
}

// NewLoadOptions create LoadOptions struct with provided opts
func NewLoadOptions(opts ...options.Option) LoadOptions {
	options := LoadOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// LoadOverride override values when load
func LoadOverride(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".Override")
	}
}

// LoadAppend override values when load
func LoadAppend(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".Append")
	}
}

// SaveOptions struct
type SaveOptions struct {
	Struct  interface{}
	Context context.Context
}

// NewSaveOptions fill SaveOptions struct
func NewSaveOptions(opts ...options.Option) SaveOptions {
	options := SaveOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// AllowFail allows config source to fail
func AllowFail(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".AllowFail")
	}
}

// BeforeInit run funcs before config Init
func BeforeInit(fn ...func(context.Context, Config) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".BeforeInit")
	}
}

// AfterInit run funcs after config Init
func AfterInit(fn ...func(context.Context, Config) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".AfterInit")
	}
}

// BeforeLoad run funcs before config load
func BeforeLoad(fn ...func(context.Context, Config) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".BeforeLoad")
	}
}

// AfterLoad run funcs after config load
func AfterLoad(fn ...func(context.Context, Config) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".AfterLoad")
	}
}

// BeforeSave run funcs before save
func BeforeSave(fn ...func(context.Context, Config) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".BeforeSave")
	}
}

// AfterSave run fncs after save
func AfterSave(fn ...func(context.Context, Config) error) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".AfterSave")
	}
}

// Struct used as config
func Struct(v interface{}) options.Option {
	return func(src interface{}) error {
		return options.Set(src, v, ".Struct")
	}
}

// StructTag sets the struct tag that used for filling
func StructTag(name string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, name, ".StructTag")
	}
}

// WatchOptions struuct
type WatchOptions struct {
	// Context used by non default options
	Context context.Context
	// Struct for filling
	Struct interface{}
	// MinInterval specifies the min time.Duration interval for poll changes
	MinInterval time.Duration
	// MaxInterval specifies the max time.Duration interval for poll changes
	MaxInterval time.Duration
	// Coalesce multiple events to one
	Coalesce bool
}

// NewWatchOptions create WatchOptions struct with provided opts
func NewWatchOptions(opts ...options.Option) WatchOptions {
	options := WatchOptions{
		Context:     context.Background(),
		MinInterval: DefaultWatcherMinInterval,
		MaxInterval: DefaultWatcherMaxInterval,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Coalesce controls watch event combining
func WatchCoalesce(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".Coalesce")
	}
}

// WatchInterval specifies min and max time.Duration for pulling changes
func WatchInterval(min, max time.Duration) options.Option {
	return func(src interface{}) error {
		var err error
		if err = options.Set(src, min, ".MinInterval"); err == nil {
			err = options.Set(src, max, ".MaxInterval")
		}
		return err
	}
}
