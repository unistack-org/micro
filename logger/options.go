package logger

import (
	"context"
	"io"
	"os"

	"go.unistack.org/micro/v4/options"
)

// Options holds logger options
type Options struct {
	// Out holds the output writer
	Out io.Writer
	// Context holds exernal options
	Context context.Context
	// Fields holds additional metadata
	Fields []interface{}
	// Name holds the logger name
	Name string
	// The logging level the logger should log
	Level Level
	// CallerSkipCount number of frmaes to skip
	CallerSkipCount int
}

// NewOptions creates new options struct
func NewOptions(opts ...options.Option) Options {
	options := Options{
		Level:           DefaultLevel,
		Fields:          make([]interface{}, 0, 6),
		Out:             os.Stderr,
		CallerSkipCount: DefaultCallerSkipCount,
		Context:         context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WithFields set default fields for the logger
func WithFields(fields ...interface{}) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fields, ".Fields")
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
