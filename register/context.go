package register

import (
	"context"
)

type registerKey struct{}

// FromContext get register from context
func FromContext(ctx context.Context) (Register, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(registerKey{}).(Register)
	return c, ok
}

// NewContext put register in context
func NewContext(ctx context.Context, c Register) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, registerKey{}, c)
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
