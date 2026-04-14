package config

import (
	"context"
)

type configKey struct{}

var configKeyVal = configKey{}

// FromContext returns store from context
func FromContext(ctx context.Context) (Config, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(configKeyVal).(Config)
	return c, ok
}

// MustContext returns store from context
func MustContext(ctx context.Context) Config {
	c, ok := FromContext(ctx)
	if !ok {
		panic("missing config")
	}
	return c
}

// NewContext put store in context
func NewContext(ctx context.Context, c Config) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, configKeyVal, c)
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v any) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetSaveOption returns a function to setup a context with given value
func SetSaveOption(k, v any) SaveOption {
	return func(o *SaveOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetLoadOption returns a function to setup a context with given value
func SetLoadOption(k, v any) LoadOption {
	return func(o *LoadOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetWatchOption returns a function to setup a context with given value
func SetWatchOption(k, v any) WatchOption {
	return func(o *WatchOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
