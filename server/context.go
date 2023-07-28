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

// NewContext stores Server to context
func NewContext(ctx context.Context, s Server) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, serverKey{}, s)
}
