package tracer

import "github.com/unistack-org/micro/v3/logger"

var (
	// DefaultSize of the buffer
	DefaultSize = 64
)

// Options struct
type Options struct {
	Name string
	// Logger is the logger for messages
	Logger logger.Logger
	// Size is the size of ring buffer
	Size int
}

// Option func
type Option func(o *Options)

// ReadOptions struct
type ReadOptions struct {
	// Trace id
	Trace string
}

// ReadOption func
type ReadOption func(o *ReadOptions)

// ReadTrace read the given trace
func ReadTrace(t string) ReadOption {
	return func(o *ReadOptions) {
		o.Trace = t
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// NewOptions returns default options
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger: logger.DefaultLogger,
		Size:   DefaultSize,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
