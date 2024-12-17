package codec

import (
	"context"
)

type codecKey struct{}

// FromContext returns codec from context
func FromContext(ctx context.Context) (Codec, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(codecKey{}).(Codec)
	return c, ok
}

// MustContext returns codec from context
func MustContext(ctx context.Context) Codec {
	c, ok := FromContext(ctx)
	if !ok {
		panic("missing codec")
	}
	return c
}

// NewContext put codec in context
func NewContext(ctx context.Context, c Codec) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, codecKey{}, c)
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
