package logger

import (
	"context"

	"go.unistack.org/micro/v4/options"
)

type noopLogger struct {
	opts Options
}

func NewLogger(opts ...options.Option) Logger {
	options := NewOptions(opts...)
	return &noopLogger{opts: options}
}

func (l *noopLogger) V(lvl Level) bool {
	return false
}

func (l *noopLogger) Level(lvl Level) {
}

func (l *noopLogger) Init(opts ...options.Option) error {
	for _, o := range opts {
		o(&l.opts)
	}
	return nil
}

func (l *noopLogger) Clone(opts ...options.Option) Logger {
	nl := &noopLogger{opts: l.opts}
	for _, o := range opts {
		o(&nl.opts)
	}
	return nl
}

func (l *noopLogger) Attrs(attrs ...interface{}) Logger {
	return l
}

func (l *noopLogger) Options() Options {
	return l.opts
}

func (l *noopLogger) Name() string {
	return l.opts.Name
}

func (l *noopLogger) String() string {
	return "noop"
}

func (l *noopLogger) Log(ctx context.Context, lvl Level, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Info(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Debug(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Error(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Trace(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Warn(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Fatal(ctx context.Context, msg string, attrs ...interface{}) {
}
