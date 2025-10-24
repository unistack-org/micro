package logger

// Level means logger level
type Level int8

const (
	// TraceLevel usually used to find bugs, very verbose
	TraceLevel Level = iota - 2
	// DebugLevel used only when enabled debugging
	DebugLevel
	// InfoLevel used for general info about what's going on inside the application
	InfoLevel
	// WarnLevel used for non-critical entries
	WarnLevel
	// ErrorLevel used for errors that should definitely be noted
	ErrorLevel
	// FatalLevel used for critical errors and then calls `os.Exit(1)`
	FatalLevel
	// NoneLevel used to disable logging
	NoneLevel
)

// String returns logger level string representation
func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case NoneLevel:
		return "none"
	}
	return "info"
}

// Enabled returns true if the given level is at or above this level.
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

// ParseLevel converts a level string into a logger Level value.
// returns an InfoLevel if the input string does not match known values.
func ParseLevel(lvl string) Level {
	switch lvl {
	case TraceLevel.String():
		return TraceLevel
	case DebugLevel.String():
		return DebugLevel
	case InfoLevel.String():
		return InfoLevel
	case WarnLevel.String():
		return WarnLevel
	case ErrorLevel.String():
		return ErrorLevel
	case FatalLevel.String():
		return FatalLevel
	case NoneLevel.String():
		return NoneLevel
	}
	return InfoLevel
}
