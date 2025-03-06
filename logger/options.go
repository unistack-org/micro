package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"slices"
	"time"

	"go.unistack.org/micro/v3/meter"
)

// Option func signature
type Option func(*Options)

// Options holds logger options
type Options struct {
	// TimeKey is the key used for the time of the log call
	TimeKey string
	// LevelKey is the key used for the level of the log call
	LevelKey string
	// ErroreKey is the key used for the error of the log call
	ErrorKey string
	// MessageKey is the key used for the message of the log call
	MessageKey string
	// SourceKey is the key used for the source file and line of the log call
	SourceKey string
	// StacktraceKey is the key used for the stacktrace
	StacktraceKey string
	// Name holds the logger name
	Name string
	// Out holds the output writer
	Out io.Writer
	// Context holds exernal options
	Context context.Context
	// Meter used to count logs for specific level
	Meter meter.Meter
	// TimeFunc used to obtain current time
	TimeFunc func() time.Time
	// Fields holds additional metadata
	Fields []interface{}
	// ContextAttrFuncs contains funcs that executed before log func on context
	ContextAttrFuncs []ContextAttrFunc
	// callerSkipCount number of frmaes to skip
	CallerSkipCount int
	// The logging level the logger should log
	Level Level
	// AddSource enabled writing source file and position in log
	AddSource bool
	// AddStacktrace controls writing of stacktaces on error
	AddStacktrace bool
	// DedupKeys deduplicate keys in log output
	DedupKeys bool
}

// NewOptions creates new options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Level:            DefaultLevel,
		Fields:           make([]interface{}, 0, 6),
		Out:              os.Stderr,
		Context:          context.Background(),
		ContextAttrFuncs: DefaultContextAttrFuncs,
		AddSource:        true,
		TimeFunc:         time.Now,
		Meter:            meter.DefaultMeter,
	}

	WithMicroKeys()(&options)

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithContextAttrFuncs appends default funcs for the context attrs filler
func WithContextAttrFuncs(fncs ...ContextAttrFunc) Option {
	return func(o *Options) {
		o.ContextAttrFuncs = append(o.ContextAttrFuncs, fncs...)
	}
}

// WithDedupKeys dont log duplicate keys
func WithDedupKeys(b bool) Option {
	return func(o *Options) {
		o.DedupKeys = b
	}
}

// WithAddFields add fields for the logger
func WithAddFields(fields ...interface{}) Option {
	return func(o *Options) {
		if o.DedupKeys {
			for i := 0; i < len(o.Fields); i += 2 {
				for j := 0; j < len(fields); j += 2 {
					iv, iok := o.Fields[i].(string)
					jv, jok := fields[j].(string)
					if iok && jok && iv == jv {
						o.Fields[i+1] = fields[j+1]
						fields = slices.Delete(fields, j, j+2)
					}
				}
			}
			if len(fields) > 0 {
				o.Fields = append(o.Fields, fields...)
			}
		} else {
			o.Fields = append(o.Fields, fields...)
		}
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

// WithAddStacktrace controls writing stacktrace on error
func WithAddStacktrace(v bool) Option {
	return func(o *Options) {
		o.AddStacktrace = v
	}
}

// WithAddSource controls writing source file and pos in log
func WithAddSource(v bool) Option {
	return func(o *Options) {
		o.AddSource = v
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

// WithMeter sets the meter
func WithMeter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// WithTimeFunc sets the func to obtain current time
func WithTimeFunc(fn func() time.Time) Option {
	return func(o *Options) {
		o.TimeFunc = fn
	}
}

func WithZapKeys() Option {
	return func(o *Options) {
		o.TimeKey = "@timestamp"
		o.LevelKey = slog.LevelKey
		o.MessageKey = slog.MessageKey
		o.SourceKey = "caller"
		o.StacktraceKey = "stacktrace"
		o.ErrorKey = "error"
	}
}

func WithZerologKeys() Option {
	return func(o *Options) {
		o.TimeKey = slog.TimeKey
		o.LevelKey = slog.LevelKey
		o.MessageKey = "message"
		o.SourceKey = "caller"
		o.StacktraceKey = "stacktrace"
		o.ErrorKey = "error"
	}
}

func WithSlogKeys() Option {
	return func(o *Options) {
		o.TimeKey = slog.TimeKey
		o.LevelKey = slog.LevelKey
		o.MessageKey = slog.MessageKey
		o.SourceKey = slog.SourceKey
		o.StacktraceKey = "stacktrace"
		o.ErrorKey = "error"
	}
}

func WithMicroKeys() Option {
	return func(o *Options) {
		o.TimeKey = "timestamp"
		o.LevelKey = slog.LevelKey
		o.MessageKey = slog.MessageKey
		o.SourceKey = "caller"
		o.StacktraceKey = "stacktrace"
		o.ErrorKey = "error"
	}
}

// WithAddCallerSkipCount add skip count for copy logger
func WithAddCallerSkipCount(n int) Option {
	return func(o *Options) {
		if n > 0 {
			o.CallerSkipCount += n
		}
	}
}
