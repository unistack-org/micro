package config

import (
	"context"
)

type configKey struct{}

func FromContext(ctx context.Context) (Config, bool) {
	c, ok := ctx.Value(configKey{}).(Config)
	return c, ok
}

func NewContext(ctx context.Context, c Config) context.Context {
	return context.WithValue(ctx, configKey{}, c)
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
