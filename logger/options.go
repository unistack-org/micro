package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Option func signature
type Option func(*Options)

// Options holds logger options
type Options struct {
	// Out holds the output writer
	Out io.Writer
	// Context holds exernal options
	Context context.Context
	// Name holds the logger name
	Name string
	// Fields holds additional metadata
	Fields []interface{}
	// CallerSkipCount number of frmaes to skip
	CallerSkipCount int
	// ContextAttrFuncs contains funcs that executed before log func on context
	ContextAttrFuncs []ContextAttrFunc
	// TimeKey is the key used for the time of the log call
	TimeKey string
	// LevelKey is the key used for the level of the log call
	LevelKey string
	// MessageKey is the key used for the message of the log call
	MessageKey string
	// SourceKey is the key used for the source file and line of the log call
	SourceKey string
	// StacktraceKey is the key used for the stacktrace
	StacktraceKey string
	// Stacktrace controls writing of stacktaces on error
	Stacktrace bool
	// The logging level the logger should log
	Level Level
}

// NewOptions creates new options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Level:            DefaultLevel,
		Fields:           make([]interface{}, 0, 6),
		Out:              os.Stderr,
		CallerSkipCount:  DefaultCallerSkipCount,
		Context:          context.Background(),
		ContextAttrFuncs: DefaultContextAttrFuncs,
	}

	WithMicroKeys()(&options)

	for _, o := range opts {
		o(&options)
	}
	return options
}

// WithContextAttrFuncs appends default funcs for the context arrts filler
func WithContextAttrFuncs(fncs ...ContextAttrFunc) Option {
	return func(o *Options) {
		o.ContextAttrFuncs = append(o.ContextAttrFuncs, fncs...)
	}
}

// WithFields set default fields for the logger
func WithFields(fields ...interface{}) Option {
	return func(o *Options) {
		o.Fields = fields
	}
}

// WithLevel set default level for the logger
func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithOutput set default output writer for the logger
func WithOutput(out io.Writer) Option {
	return func(o *Options) {
		o.Out = out
	}
}

// WithStacktrace controls writing stacktrace on error
func WithStacktrace(v bool) Option {
	return func(o *Options) {
		o.Stacktrace = v
	}
}

// WithCallerSkipCount set frame count to skip
func WithCallerSkipCount(c int) Option {
	return func(o *Options) {
		o.CallerSkipCount = c
	}
}

// WithContext set context
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// WithName sets the name
func WithName(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func WithZapKeys() Option {
	return func(o *Options) {
		o.TimeKey = "@timestamp"
		o.LevelKey = "level"
		o.MessageKey = "msg"
		o.SourceKey = "caller"
		o.StacktraceKey = "stacktrace"
	}
}

func WithZerologKeys() Option {
	return func(o *Options) {
		o.TimeKey = "time"
		o.LevelKey = "level"
		o.MessageKey = "message"
		o.SourceKey = "caller"
		o.StacktraceKey = "stacktrace"
	}
}

func WithSlogKeys() Option {
	return func(o *Options) {
		o.TimeKey = slog.TimeKey
		o.LevelKey = slog.LevelKey
		o.MessageKey = slog.MessageKey
		o.SourceKey = slog.SourceKey
		o.StacktraceKey = "stacktrace"
	}
}

func WithMicroKeys() Option {
	return func(o *Options) {
		o.TimeKey = "timestamp"
		o.LevelKey = "level"
		o.MessageKey = "msg"
		o.SourceKey = "caller"
		o.StacktraceKey = "stacktrace"
	}
}
