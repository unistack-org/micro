package server

import (
	"context"
)

type serverKey struct{}

var serverKeyVal = serverKey{}

// FromContext returns Server from context
func FromContext(ctx context.Context) (Server, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(serverKeyVal).(Server)
	return c, ok
}

// MustContext returns Server from context
func MustContext(ctx context.Context) Server {
	s, ok := FromContext(ctx)
	if !ok {
		panic("missing server")
	}
	return s
}

// NewContext stores Server to context
func NewContext(ctx context.Context, s Server) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, serverKeyVal, s)
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

// SetHandlerOption returns a function to setup a context with given value
func SetHandlerOption(k, v any) HandlerOption {
	return func(o *HandlerOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
