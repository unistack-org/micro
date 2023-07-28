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

// NewContext put client in context
func NewContext(ctx context.Context, c Client) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, clientKey{}, c)
}
