package logger

import "context"

type loggerKey struct{}

// FromContext returns logger from passed context
func FromContext(ctx context.Context) (Logger, bool) {
	l, ok := ctx.Value(loggerKey{}).(Logger)
	return l, ok
}

// NewContext stores logger into passed context
func NewContext(ctx context.Context, l Logger) context.Context {
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
