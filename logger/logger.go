// Package logger provides a log interface
package logger

import (
	"context"
	"os"
)

var (
	// DefaultLogger variable
	DefaultLogger Logger = NewLogger(WithLevel(ParseLevel(os.Getenv("MICRO_LOG_LEVEL"))))
	// DefaultLevel used by logger
	DefaultLevel Level = InfoLevel
	// DefaultCallerSkipCount used by logger
	DefaultCallerSkipCount = 2
)

// Logger is a generic logging interface
type Logger interface {
	// Init initialises options
	Init(opts ...Option) error
	// V compare provided verbosity level with current log level
	V(level Level) bool
	// The Logger options
	Options() Options
	// Fields set fields to always be logged with keyval pairs
	Fields(fields ...interface{}) Logger
	// Info level message
	Info(ctx context.Context, args ...interface{})
	// Trace level message
	Trace(ctx context.Context, args ...interface{})
	// Debug level message
	Debug(ctx context.Context, args ...interface{})
	// Warn level message
	Warn(ctx context.Context, args ...interface{})
	// Error level message
	Error(ctx context.Context, args ...interface{})
	// Fatal level message
	Fatal(ctx context.Context, args ...interface{})
	// Infof level message
	Infof(ctx context.Context, msg string, args ...interface{})
	// Tracef level message
	Tracef(ctx context.Context, msg string, args ...interface{})
	// Debug level message
	Debugf(ctx context.Context, msg string, args ...interface{})
	// Warn level message
	Warnf(ctx context.Context, msg string, args ...interface{})
	// Error level message
	Errorf(ctx context.Context, msg string, args ...interface{})
	// Fatal level message
	Fatalf(ctx context.Context, msg string, args ...interface{})
	// Log logs message with needed level
	Log(ctx context.Context, level Level, args ...interface{})
	// Logf logs message with needed level
	Logf(ctx context.Context, level Level, msg string, args ...interface{})
	// String returns the name of logger
	String() string
}

// Field contains keyval pair
type Field interface{}

// Info writes msg to default logger on info level
func Info(ctx context.Context, args ...interface{}) {
	DefaultLogger.Info(ctx, args...)
}

// Error writes msg to default logger on error level
func Error(ctx context.Context, args ...interface{}) {
	DefaultLogger.Error(ctx, args...)
}

// Debug writes msg to default logger on debug level
func Debug(ctx context.Context, args ...interface{}) {
	DefaultLogger.Debug(ctx, args...)
}

// Warn writes msg to default logger on warn level
func Warn(ctx context.Context, args ...interface{}) {
	DefaultLogger.Warn(ctx, args...)
}

// Trace writes msg to default logger on trace level
func Trace(ctx context.Context, args ...interface{}) {
	DefaultLogger.Trace(ctx, args...)
}

// Fatal writes msg to default logger on fatal level
func Fatal(ctx context.Context, args ...interface{}) {
	DefaultLogger.Fatal(ctx, args...)
}

// Infof writes formatted msg to default logger on info level
func Infof(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Infof(ctx, msg, args...)
}

// Errorf writes formatted msg to default logger on error level
func Errorf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Errorf(ctx, msg, args...)
}

// Debugf writes formatted msg to default logger on debug level
func Debugf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Debugf(ctx, msg, args...)
}

// Warnf writes formatted msg to default logger on warn level
func Warnf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Warnf(ctx, msg, args...)
}

// Tracef writes formatted msg to default logger on trace level
func Tracef(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Tracef(ctx, msg, args...)
}

// Fatalf writes formatted msg to default logger on fatal level
func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Fatalf(ctx, msg, args...)
}

// V returns true if passed level enabled in default logger
func V(level Level) bool {
	return DefaultLogger.V(level)
}

// Init initialize logger
func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

// Fields create logger with specific fields
func Fields(fields ...interface{}) Logger {
	return DefaultLogger.Fields(fields...)
}
