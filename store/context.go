package store

import (
	"context"
)

type storeKey struct{}

// FromContext get store from context
func FromContext(ctx context.Context) (Store, bool) {
	c, ok := ctx.Value(storeKey{}).(Store)
	return c, ok
}

// NewContext put store in context
func NewContext(ctx context.Context, c Store) context.Context {
	return context.WithValue(ctx, storeKey{}, c)
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
