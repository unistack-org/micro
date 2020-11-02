// Package logger provides a log interface
package logger

var (
	// DefaultLogger variable
	DefaultLogger Logger = NewHelper(NewLogger())
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
	// Log writes a log entry
	Log(level Level, v ...interface{})
	// Logf writes a formatted log entry
	Logf(level Level, format string, v ...interface{})
	// String returns the name of logger
	String() string
}

// Init initialize logger
func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

// Fields create logger with specific fields
func Fields(fields map[string]interface{}) Logger {
	return DefaultLogger.Fields(fields)
}

// Log writes log with specific level
func Log(level Level, v ...interface{}) {
	DefaultLogger.Log(level, v...)
}

// Logf writes formatted log with specific level
func Logf(level Level, format string, v ...interface{}) {
	DefaultLogger.Logf(level, format, v...)
}

// String return logger name
func String() string {
	return DefaultLogger.String()
}
