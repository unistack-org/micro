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
	Info(args ...interface{})
	// Trace level message
	Trace(args ...interface{})
	// Debug level message
	Debug(args ...interface{})
	// Warn level message
	Warn(args ...interface{})
	// Error level message
	Error(args ...interface{})
	// Fatal level message
	Fatal(args ...interface{})
	// Infof level message
	Infof(msg string, args ...interface{})
	// Tracef level message
	Tracef(msg string, args ...interface{})
	// Debug level message
	Debugf(msg string, args ...interface{})
	// Warn level message
	Warnf(msg string, args ...interface{})
	// Error level message
	Errorf(msg string, args ...interface{})
	// Fatal level message
	Fatalf(msg string, args ...interface{})
	// String returns the name of logger
	String() string
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Trace(args ...interface{}) {
	DefaultLogger.Trace(args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

func Infof(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Errorf(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Debugf(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	DefaultLogger.Warnf(msg, args...)
}

func Tracef(msg string, args ...interface{}) {
	DefaultLogger.Tracef(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Fatalf(msg, args...)
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
