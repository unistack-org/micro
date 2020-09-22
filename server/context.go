package server

import (
	"context"
)

type serverKey struct{}

func FromContext(ctx context.Context) (Server, bool) {
	c, ok := ctx.Value(serverKey{}).(Server)
	return c, ok
}

func NewContext(ctx context.Context, s Server) context.Context {
	return context.WithValue(ctx, serverKey{}, s)
}

// SetServerSubscriberOption returns a function to setup a context with given value
func SetServerSubscriberOption(k, v interface{}) SubscriberOption {
	return func(o *SubscriberOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
