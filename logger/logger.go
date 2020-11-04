// Package logger provides a log interface
package logger

var (
	// DefaultLogger variable
	DefaultLogger Logger = NewLogger()
)

// Logger is a generic logging interface
type Logger interface {
	// Init initialises options
	Init(opts ...Option) error
	// V compare provided verbosity level with current log level
	V(level Level) bool
	// The Logger options
	Options() Options
	// Fields set fields to always be logged
	Fields(fields map[string]interface{}) Logger
	// Info level message
	Info(msg string, args ...interface{})
	// Trace level message
	Trace(msg string, args ...interface{})
	// Debug level message
	Debug(msg string, args ...interface{})
	// Warn level message
	Warn(msg string, args ...interface{})
	// Error level message
	Error(msg string, args ...interface{})
	// Fatal level message
	Fatal(msg string, args ...interface{})
	// String returns the name of logger
	String() string
}

func Info(msg string, args ...interface{}) {
	DefaultLogger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	DefaultLogger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	DefaultLogger.Debug(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	DefaultLogger.Warn(msg, args...)
}

func Trace(msg string, args ...interface{}) {
	DefaultLogger.Trace(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	DefaultLogger.Fatal(msg, args...)
}

func V(level Level) bool {
	return DefaultLogger.V(level)
}

// Init initialize logger
func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

// Fields create logger with specific fields
func Fields(fields map[string]interface{}) Logger {
	return DefaultLogger.Fields(fields)
}
