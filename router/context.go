package router

import (
	"context"
)

type routerKey struct{}

// FromContext get router from context
func FromContext(ctx context.Context) (Router, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(routerKey{}).(Router)
	return c, ok
}

// MustContext get router from context
func MustContext(ctx context.Context) Router {
	r, ok := FromContext(ctx)
	if !ok {
		panic("missing router")
	}
	return r
}

// NewContext put router in context
func NewContext(ctx context.Context, c Router) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, routerKey{}, c)
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
