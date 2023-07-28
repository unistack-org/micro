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

// NewContext stores logger into passed context
func NewContext(ctx context.Context, l Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, loggerKey{}, l)
}
