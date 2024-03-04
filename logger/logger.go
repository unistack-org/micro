// Package logger provides a log interface
package logger

import (
	"context"

	"go.unistack.org/micro/v4/options"
)

type ContextAttrFunc func(ctx context.Context) []interface{}

var DefaultContextAttrFuncs []ContextAttrFunc

var (
	// DefaultLogger variable
	DefaultLogger = NewLogger()
	// DefaultLevel used by logger
	DefaultLevel = InfoLevel
	// DefaultCallerSkipCount used by logger
	DefaultCallerSkipCount = 2
)

// Logger is a generic logging interface
type Logger interface {
	// Init initialises options
	Init(opts ...options.Option) error
	// Clone create logger copy with new options
	Clone(opts ...options.Option) Logger
	// V compare provided verbosity level with current log level
	V(level Level) bool
	// Level sets the log level for logger
	Level(level Level)
	// The Logger options
	Options() Options
	// Attrs set attrs to always be logged with keyval pairs
	Attrs(attrs ...interface{}) Logger
	// Info level message
	Info(ctx context.Context, msg string, attrs ...interface{})
	// Tracef level message
	Trace(ctx context.Context, msg string, attrs ...interface{})
	// Debug level message
	Debug(ctx context.Context, msg string, attrs ...interface{})
	// Warn level message
	Warn(ctx context.Context, msg string, attrs ...interface{})
	// Error level message
	Error(ctx context.Context, msg string, attrs ...interface{})
	// Fatal level message
	Fatal(ctx context.Context, msg string, attrs ...interface{})
	// Log logs message with needed level
	Log(ctx context.Context, level Level, msg string, attrs ...interface{})
	// String returns the type name of logger
	String() string
	// String returns the name of logger
	Name() string
}

// Info writes formatted msg to default logger on info level
func Info(ctx context.Context, msg string, attrs ...interface{}) {
	DefaultLogger.Info(ctx, msg, attrs...)
}

// Error writes formatted msg to default logger on error level
func Error(ctx context.Context, msg string, attrs ...interface{}) {
	DefaultLogger.Error(ctx, msg, attrs...)
}

// Debugf writes formatted msg to default logger on debug level
func Debugf(ctx context.Context, msg string, attrs ...interface{}) {
	DefaultLogger.Debug(ctx, msg, attrs...)
}

// Warn writes formatted msg to default logger on warn level
func Warn(ctx context.Context, msg string, attrs ...interface{}) {
	DefaultLogger.Warn(ctx, msg, attrs...)
}

// Trace writes formatted msg to default logger on trace level
func Trace(ctx context.Context, msg string, attrs ...interface{}) {
	DefaultLogger.Trace(ctx, msg, attrs...)
}

// Fatal writes formatted msg to default logger on fatal level
func Fatal(ctx context.Context, msg string, attrs ...interface{}) {
	DefaultLogger.Fatal(ctx, msg, attrs...)
}

// V returns true if passed level enabled in default logger
func V(level Level) bool {
	return DefaultLogger.V(level)
}

// Init initialize default logger
func Init(opts ...options.Option) error {
	return DefaultLogger.Init(opts...)
}

// Attrs create default logger with specific attrs
func Attrs(attrs ...interface{}) Logger {
	return DefaultLogger.Attrs(attrs...)
}
