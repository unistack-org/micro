package store

import (
	"context"
)

type storeKey struct{}

// FromContext get store from context
func FromContext(ctx context.Context) (Store, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(storeKey{}).(Store)
	return c, ok
}

// NewContext put store in context
func NewContext(ctx context.Context, c Store) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, storeKey{}, c)
}
