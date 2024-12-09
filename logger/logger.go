// Package logger provides a log interface
package logger

import (
	"context"
)

type ContextAttrFunc func(ctx context.Context) []interface{}

var DefaultContextAttrFuncs []ContextAttrFunc

var (
	// DefaultLogger variable
	DefaultLogger = NewLogger()
	// DefaultLevel used by logger
	DefaultLevel = InfoLevel
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
	Info(ctx context.Context, msg string, args ...interface{})
	// Trace level message
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
	// Name returns broker instance name
	Name() string
	// String returns the type of logger
	String() string
}

// Field contains keyval pair
type Field interface{}
