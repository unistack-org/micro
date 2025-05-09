package server

import (
	"context"
)

type serverKey struct{}

// FromContext returns Server from context
func FromContext(ctx context.Context) (Server, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(serverKey{}).(Server)
	return c, ok
}

// MustContext returns Server from context
func MustContext(ctx context.Context) Server {
	s, ok := FromContext(ctx)
	if !ok {
		panic("missing server")
	}
	return s
}

// NewContext stores Server to context
func NewContext(ctx context.Context, s Server) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, serverKey{}, s)
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

// SetSubscriberOption returns a function to setup a context with given value
func SetSubscriberOption(k, v interface{}) SubscriberOption {
	return func(o *SubscriberOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

// SetHandlerOption returns a function to setup a context with given value
func SetHandlerOption(k, v interface{}) HandlerOption {
	return func(o *HandlerOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
