// Package profiler is for profilers
package profiler // import "go.unistack.org/micro/v4/profiler"

// Profiler interface
type Profiler interface {
	// Start the profiler
	Start() error
	// Stop the profiler
	Stop() error
	// Name of the profiler
	String() string
}

// DefaultProfiler holds the default profiler
var DefaultProfiler Profiler = NewProfiler()

// Options holds the options for profiler
type Options struct {
	// Name to use for the profile
	Name string
}

// Option func signature
type Option func(o *Options)

// Name of the profile
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
