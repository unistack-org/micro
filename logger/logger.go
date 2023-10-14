// Package logger provides a log interface
package logger // import "go.unistack.org/micro/v4/logger"

import (
	"context"
	"os"

	"go.unistack.org/micro/v4/options"
)

var (
	// DefaultLogger variable
	DefaultLogger = NewLogger(WithLevel(ParseLevel(os.Getenv("MICRO_LOG_LEVEL"))))
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
	// Fields set fields to always be logged with keyval pairs
	Fields(fields ...interface{}) Logger
	// Info level message
	Info(ctx context.Context, msg string, args ...interface{})
	// Tracef level message
	Trace(ctx context.Context, msg string, args ...interface{})
	// Debug level message
	Debug(ctx context.Context, msg string, args ...interface{})
	// Warn level message
	Warn(ctx context.Context, msg string, args ...interface{})
	// Error level message
	Error(ctx context.Context, msg string, args ...interface{})
	// Fatal level message
	Fatal(ctx context.Context, msg string, args ...interface{})
	// Log logs message with needed level
	Log(ctx context.Context, level Level, msg string, args ...interface{})
	// String returns the name of logger
	String() string
}

// Field contains keyval pair
type Field interface{}

// Info writes formatted msg to default logger on info level
func Info(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Info(ctx, msg, args...)
}

// Error writes formatted msg to default logger on error level
func Error(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Error(ctx, msg, args...)
}

// Debugf writes formatted msg to default logger on debug level
func Debugf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Debug(ctx, msg, args...)
}

// Warn writes formatted msg to default logger on warn level
func Warn(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Warn(ctx, msg, args...)
}

// Trace writes formatted msg to default logger on trace level
func Trace(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Trace(ctx, msg, args...)
}

// Fatal writes formatted msg to default logger on fatal level
func Fatal(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Fatal(ctx, msg, args...)
}

// V returns true if passed level enabled in default logger
func V(level Level) bool {
	return DefaultLogger.V(level)
}

// Init initialize logger
func Init(opts ...options.Option) error {
	return DefaultLogger.Init(opts...)
}

// Fields create logger with specific fields
func Fields(fields ...interface{}) Logger {
	return DefaultLogger.Fields(fields...)
}
