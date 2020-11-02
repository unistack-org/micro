package tracer

// Options struct
type Options struct {
	// Size is the size of ring buffer
	Size int
}

// Option func
type Option func(o *Options)

type ReadOptions struct {
	// Trace id
	Trace string
}

// ReadOption func
type ReadOption func(o *ReadOptions)

// ReadTracer read the given trace
func ReadTrace(t string) ReadOption {
	return func(o *ReadOptions) {
		o.Trace = t
	}
}

const (
	// DefaultSize of the buffer
	DefaultSize = 64
)

// DefaultOptions returns default options
func DefaultOptions() Options {
	return Options{
		Size: DefaultSize,
	}
}
