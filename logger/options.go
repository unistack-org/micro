package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"reflect"

	"go.unistack.org/micro/v4/options"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

// Options holds logger options
type Options struct {
	// Out holds the output writer
	Out io.Writer
	// Context holds exernal options
	Context context.Context
	// Attrs holds additional attributes
	Attrs []interface{}
	// Name holds the logger name
	Name string
	// The logging level the logger should log
	Level Level
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
}

// NewOptions creates new options struct
func NewOptions(opts ...options.Option) Options {
	options := Options{
		Level:            DefaultLevel,
		Attrs:            make([]interface{}, 0, 6),
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
func WithContextAttrFuncs(fncs ...ContextAttrFunc) options.Option {
	return func(src interface{}) error {
		v, err := options.Get(src, ".ContextAttrFuncs")
		if err != nil {
			return err
		} else if rutil.IsZero(v) {
			v = reflect.MakeSlice(reflect.TypeOf(v), 0, len(fncs)).Interface()
		}
		cv := reflect.ValueOf(v)
		for _, l := range fncs {
			cv = reflect.Append(cv, reflect.ValueOf(l))
		}
		return options.Set(src, cv.Interface(), ".ContextAttrFuncs")
	}
}

// WithAttrs set default fields for the logger
func WithAttrs(attrs ...interface{}) options.Option {
	return func(src interface{}) error {
		return options.Set(src, attrs, ".Attrs")
	}
}

// WithLevel set default level for the logger
func WithLevel(lvl Level) options.Option {
	return func(src interface{}) error {
		return options.Set(src, lvl, ".Level")
	}
}

// WithOutput set default output writer for the logger
func WithOutput(out io.Writer) options.Option {
	return func(src interface{}) error {
		return options.Set(src, out, ".Out")
	}
}

// WithCallerSkipCount set frame count to skip
func WithCallerSkipCount(c int) options.Option {
	return func(src interface{}) error {
		return options.Set(src, c, ".CallerSkipCount")
	}
}

func WithZapKeys() options.Option {
	return func(src interface{}) error {
		var err error
		if err = options.Set(src, "@timestamp", ".TimeKey"); err != nil {
			return err
		}
		if err = options.Set(src, "level", ".LevelKey"); err != nil {
			return err
		}
		if err = options.Set(src, "msg", ".MessageKey"); err != nil {
			return err
		}
		if err = options.Set(src, "caller", ".SourceKey"); err != nil {
			return err
		}
		return nil
	}
}

func WithZerologKeys() options.Option {
	return func(src interface{}) error {
		var err error
		if err = options.Set(src, "time", ".TimeKey"); err != nil {
			return err
		}
		if err = options.Set(src, "level", ".LevelKey"); err != nil {
			return err
		}
		if err = options.Set(src, "message", ".MessageKey"); err != nil {
			return err
		}
		if err = options.Set(src, "caller", ".SourceKey"); err != nil {
			return err
		}
		return nil
	}
}

func WithSlogKeys() options.Option {
	return func(src interface{}) error {
		var err error
		if err = options.Set(src, slog.TimeKey, ".TimeKey"); err != nil {
			return err
		}
		if err = options.Set(src, slog.LevelKey, ".LevelKey"); err != nil {
			return err
		}
		if err = options.Set(src, slog.MessageKey, ".MessageKey"); err != nil {
			return err
		}
		if err = options.Set(src, slog.SourceKey, ".SourceKey"); err != nil {
			return err
		}
		return nil
	}
}

func WithMicroKeys() options.Option {
	return func(src interface{}) error {
		var err error
		if err = options.Set(src, "timestamp", ".TimeKey"); err != nil {
			return err
		}
		if err = options.Set(src, "level", ".LevelKey"); err != nil {
			return err
		}
		if err = options.Set(src, "msg", ".MessageKey"); err != nil {
			return err
		}
		if err = options.Set(src, "caller", ".SourceKey"); err != nil {
			return err
		}
		return nil
	}
}
