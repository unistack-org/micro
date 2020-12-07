package registry

import (
	"context"
)

type registryKey struct{}

// FromContext get registry from context
func FromContext(ctx context.Context) (Registry, bool) {
	c, ok := ctx.Value(registryKey{}).(Registry)
	return c, ok
}

// NewContext put registry in context
func NewContext(ctx context.Context, c Registry) context.Context {
	return context.WithValue(ctx, registryKey{}, c)
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
