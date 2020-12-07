package transport

import (
	"context"
)

type transportKey struct{}

// FromContext get transport from context
func FromContext(ctx context.Context) (Transport, bool) {
	c, ok := ctx.Value(transportKey{}).(Transport)
	return c, ok
}

// NewContext put transport in context
func NewContext(ctx context.Context, c Transport) context.Context {
	return context.WithValue(ctx, transportKey{}, c)
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
