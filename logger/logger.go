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

// Info writes msg to default logger on info level
func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

// Error writes msg to default logger on error level
func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

// Debug writes msg to default logger on debug level
func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

// Warn writes msg to default logger on warn level
func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

// Trace writes msg to default logger on trace level
func Trace(args ...interface{}) {
	DefaultLogger.Trace(args...)
}

// Fatal writes msg to default logger on fatal level
func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

// Infof writes formatted msg to default logger on info level
func Infof(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

// Errorf writes formatted msg to default logger on error level
func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Errorf(msg, args...)
}

// Debugf writes formatted msg to default logger on debug level
func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Debugf(msg, args...)
}

// Warnf writes formatted msg to default logger on warn level
func Warnf(msg string, args ...interface{}) {
	DefaultLogger.Warnf(msg, args...)
}

// Tracef writes formatted msg to default logger on trace level
func Tracef(msg string, args ...interface{}) {
	DefaultLogger.Tracef(msg, args...)
}

// Fatalf writes formatted msg to default logger on fatal level
func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Fatalf(msg, args...)
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
func Fields(fields map[string]interface{}) Logger {
	return DefaultLogger.Fields(fields)
}
