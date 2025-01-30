package client

import (
	"context"
)

type clientKey struct{}

// FromContext get client from context
func FromContext(ctx context.Context) (Client, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(clientKey{}).(Client)
	return c, ok
}

// MustContext get client from context
func MustContext(ctx context.Context) Client {
	c, ok := FromContext(ctx)
	if !ok {
		panic("missing client")
	}
	return c
}

// NewContext put client in context
func NewContext(ctx context.Context, c Client) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, clientKey{}, c)
}

// SetCallOption returns a function to setup a context with given value
func SetCallOption(k, v interface{}) CallOption {
	return func(o *CallOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
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
