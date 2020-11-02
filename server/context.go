package server

import (
	"context"
)

type serverKey struct{}

// FromContext returns Server from context
func FromContext(ctx context.Context) (Server, bool) {
	c, ok := ctx.Value(serverKey{}).(Server)
	return c, ok
}

// NewContext stores Server to context
func NewContext(ctx context.Context, s Server) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}

// Setoption returns a function to setup a context with given value
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
