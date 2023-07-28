package flow

import (
	"context"
)

type flowKey struct{}

// FromContext returns Flow from context
func FromContext(ctx context.Context) (Flow, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(flowKey{}).(Flow)
	return c, ok
}

// NewContext stores Flow to context
func NewContext(ctx context.Context, f Flow) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, flowKey{}, f)
}
