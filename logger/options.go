package logger

import (
	"context"
	"io"
	"os"
)

// Option func
type Option func(*Options)

// Options holds logger options
type Options struct {
	// Out holds the output writer
	Out io.Writer
	// Context holds exernal options
	Context context.Context
	// Fields holds additional metadata
	Fields map[string]interface{}
	// Name holds the logger name
	Name string
	// CallerSkipCount number of frmaes to skip
	CallerSkipCount int
	// The logging level the logger should log
	Level Level
}

// NewOptions creates new options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Level:           DefaultLevel,
		Fields:          make(map[string]interface{}),
		Out:             os.Stderr,
		CallerSkipCount: 0,
		Context:         context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WithFields set default fields for the logger
func WithFields(fields map[string]interface{}) Option {
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
func withName(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
