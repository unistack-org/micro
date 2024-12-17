package logger

import "context"

type loggerKey struct{}

// FromContext returns logger from passed context
func FromContext(ctx context.Context) (Logger, bool) {
	if ctx == nil {
		return nil, false
	}
	l, ok := ctx.Value(loggerKey{}).(Logger)
	return l, ok
}

// MustContext returns logger from passed context or DefaultLogger if empty
func MustContext(ctx context.Context) Logger {
	l, ok := FromContext(ctx)
	if !ok {
		panic("missing logger")
	}
	return l
}

// NewContext stores logger into passed context
func NewContext(ctx context.Context, l Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, loggerKey{}, l)
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
