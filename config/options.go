package config

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
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
	// BeforeSave contains slice of funcs that runs before save
	BeforeSave []func(context.Context, Config) error
	// AfterLoad contains slice of funcs that runs after load
	AfterLoad []func(context.Context, Config) error
	// BeforeLoad contains slice of funcs that runs before load
	BeforeLoad []func(context.Context, Config) error
	// AfterSave contains slice of funcs that runs after save
	AfterSave []func(context.Context, Config) error
	// AllowFail flag to allow fail in config source
	AllowFail bool
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

// LoadOption function signature
type LoadOption func(o *LoadOptions)

// LoadOptions struct
type LoadOptions struct {
	Struct   interface{}
	Override bool
	Append   bool
}

func NewLoadOptions(opts ...LoadOption) LoadOptions {
	options := LoadOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// LoadOverride override values when load
func LoadOverride(b bool) LoadOption {
	return func(o *LoadOptions) {
		o.Override = b
	}
}

// LoadAppend override values when load
func LoadAppend(b bool) LoadOption {
	return func(o *LoadOptions) {
		o.Append = b
	}
}

// LoadStruct override struct for loading
func LoadStruct(src interface{}) LoadOption {
	return func(o *LoadOptions) {
		o.Struct = src
	}
}

// SaveOption function signature
type SaveOption func(o *SaveOptions)

// SaveOptions struct
type SaveOptions struct {
	Struct interface{}
}

// SaveStruct override struct for save to config
func SaveStruct(src interface{}) SaveOption {
	return func(o *SaveOptions) {
		o.Struct = src
	}
}

// NewSaveOptions fill SaveOptions struct
func NewSaveOptions(opts ...SaveOption) SaveOptions {
	options := SaveOptions{}
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

// StructTag sets the struct tag that used for filling
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

// WatchOptions struuct
type WatchOptions struct {
	// Context used by non default options
	Context context.Context
	// Coalesce multiple events to one
	Coalesce bool
	// Interval to periodically pull changes if config source not supports async notify
	Interval time.Duration
	// Struct for filling
	Struct interface{}
}

type WatchOption func(*WatchOptions)

func NewWatchOptions(opts ...WatchOption) WatchOptions {
	options := WatchOptions{
		Context:  context.Background(),
		Interval: DefaultWatcherInterval,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WatchContext pass context
func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

// WatchCoalesce controls watch event combining
func WatchCoalesce(b bool) WatchOption {
	return func(o *WatchOptions) {
		o.Coalesce = b
	}
}

// WatchInterval specifies time.Duration for pulling changes
func WatchInterval(td time.Duration) WatchOption {
	return func(o *WatchOptions) {
		o.Interval = td
	}
}

// WatchStruct overrides struct for fill
func WatchStruct(src interface{}) WatchOption {
	return func(o *WatchOptions) {
		o.Struct = src
	}
}
