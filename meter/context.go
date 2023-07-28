package meter

import (
	"context"
)

type meterKey struct{}

// FromContext get meter from context
func FromContext(ctx context.Context) (Meter, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(meterKey{}).(Meter)
	return c, ok
}

// NewContext put meter in context
func NewContext(ctx context.Context, c Meter) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, meterKey{}, c)
}
