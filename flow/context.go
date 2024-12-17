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

// MustContext returns Flow from context
func MustContext(ctx context.Context) Flow {
	f, ok := FromContext(ctx)
	if !ok {
		panic("missing flow")
	}
	return f
}

// NewContext stores Flow to context
func NewContext(ctx context.Context, f Flow) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, flowKey{}, f)
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
