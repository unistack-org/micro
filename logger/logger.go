// Package logger provides a log interface
package logger

import (
	"context"
)

type ContextAttrFunc func(ctx context.Context) []interface{}

var DefaultContextAttrFuncs []ContextAttrFunc

var (
	// DefaultLogger variable
	DefaultLogger Logger = NewLogger()
	// DefaultLevel used by logger
	DefaultLevel = InfoLevel
	// DefaultCallerSkipCount used by logger
	DefaultCallerSkipCount = 2
)

// Logger is a generic logging interface
type Logger interface {
	// Init initialises options
	Init(opts ...Option) error
	// Clone create logger copy with new options
	Clone(opts ...Option) Logger
	// V compare provided verbosity level with current log level
	V(level Level) bool
	// Level sets the log level for logger
	Level(level Level)
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
	// Name returns broker instance name
	Name() string
	// String returns the type of logger
	String() string
}

// Field contains keyval pair
type Field interface{}

// Info writes msg to default logger on info level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Info(ctx context.Context, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Info(ctx, args...)
}

// Error writes msg to default logger on error level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Error(ctx context.Context, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Error(ctx, args...)
}

// Debug writes msg to default logger on debug level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Debug(ctx context.Context, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Debug(ctx, args...)
}

// Warn writes msg to default logger on warn level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Warn(ctx context.Context, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Warn(ctx, args...)
}

// Trace writes msg to default logger on trace level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Trace(ctx context.Context, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Trace(ctx, args...)
}

// Fatal writes msg to default logger on fatal level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Fatal(ctx context.Context, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Fatal(ctx, args...)
}

// Infof writes formatted msg to default logger on info level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Infof(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Infof(ctx, msg, args...)
}

// Errorf writes formatted msg to default logger on error level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Errorf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Errorf(ctx, msg, args...)
}

// Debugf writes formatted msg to default logger on debug level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Debugf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Debugf(ctx, msg, args...)
}

// Warnf writes formatted msg to default logger on warn level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Warnf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Warnf(ctx, msg, args...)
}

// Tracef writes formatted msg to default logger on trace level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Tracef(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Tracef(ctx, msg, args...)
}

// Fatalf writes formatted msg to default logger on fatal level
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	DefaultLogger.Clone(WithCallerSkipCount(DefaultCallerSkipCount+1)).Fatalf(ctx, msg, args...)
}

// V returns true if passed level enabled in default logger
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func V(level Level) bool {
	return DefaultLogger.V(level)
}

// Init initialize logger
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

// Fields create logger with specific fields
//
// Deprecated: Dont use logger methods directly, use instance of logger to avoid additional allocations
func Fields(fields ...interface{}) Logger {
	return DefaultLogger.Fields(fields...)
}
