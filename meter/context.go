package meter

import (
	"context"
)

type meterKey struct{}

var meterKeyVal = meterKey{}

// FromContext get meter from context
func FromContext(ctx context.Context) (Meter, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(meterKeyVal).(Meter)
	return c, ok
}

// MustContext get meter from context
func MustContext(ctx context.Context) Meter {
	m, ok := FromContext(ctx)
	if !ok {
		panic("missing meter")
	}
	return m
}

// NewContext put meter in context
func NewContext(ctx context.Context, c Meter) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, meterKeyVal, c)
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v any) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
