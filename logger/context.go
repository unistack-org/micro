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
