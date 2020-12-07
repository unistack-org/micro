package tracer

var (
	// DefaultSize of the buffer
	DefaultSize = 64
)

// Options struct
type Options struct {
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

// NewOptions returns default options
func NewOptions(opts ...Option) Options {
	options := Options{
		Size: DefaultSize,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
