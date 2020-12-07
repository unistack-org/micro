package client

import (
	"context"
)

type clientKey struct{}

// FromContext get client from context
func FromContext(ctx context.Context) (Client, bool) {
	c, ok := ctx.Value(clientKey{}).(Client)
	return c, ok
}

// NewContext put client in context
func NewContext(ctx context.Context, c Client) context.Context {
	return context.WithValue(ctx, clientKey{}, c)
}

// SetPublishOption returns a function to setup a context with given value
func SetPublishOption(k, v interface{}) PublishOption {
	return func(o *PublishOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
