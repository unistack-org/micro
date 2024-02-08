package logger

import (
	"context"
)

type noopLogger struct {
	opts Options
}

func NewLogger(opts ...Option) Logger {
	options := NewOptions(opts...)
	return &noopLogger{opts: options}
}

func (l *noopLogger) V(lvl Level) bool {
	return false
}

func (l *noopLogger) Level(lvl Level) {
}

func (l *noopLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&l.opts)
	}
	return nil
}

func (l *noopLogger) Clone(opts ...Option) Logger {
	nl := &noopLogger{opts: l.opts}
	for _, o := range opts {
		o(&nl.opts)
	}
	return nl
}

func (l *noopLogger) Fields(attrs ...interface{}) Logger {
	return l
}

func (l *noopLogger) Options() Options {
	return l.opts
}

func (l *noopLogger) String() string {
	return "noop"
}

func (l *noopLogger) Log(ctx context.Context, lvl Level, attrs ...interface{}) {
}

func (l *noopLogger) Info(ctx context.Context, attrs ...interface{}) {
}

func (l *noopLogger) Debug(ctx context.Context, attrs ...interface{}) {
}

func (l *noopLogger) Error(ctx context.Context, attrs ...interface{}) {
}

func (l *noopLogger) Trace(ctx context.Context, attrs ...interface{}) {
}

func (l *noopLogger) Warn(ctx context.Context, attrs ...interface{}) {
}

func (l *noopLogger) Fatal(ctx context.Context, attrs ...interface{}) {
}

func (l *noopLogger) Logf(ctx context.Context, lvl Level, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Infof(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Debugf(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Errorf(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Tracef(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Warnf(ctx context.Context, msg string, attrs ...interface{}) {
}

func (l *noopLogger) Fatalf(ctx context.Context, msg string, attrs ...interface{}) {
}
