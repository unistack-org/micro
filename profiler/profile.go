// Package profile is for profilers
package profile

type Profile interface {
	// Start the profiler
	Start() error
	// Stop the profiler
	Stop() error
	// Name of the profiler
	String() string
}

var (
	DefaultProfile Profile = &NoopProfile{}
)

type NoopProfile struct{}

func (p *NoopProfile) Start() error {
	return nil
}

func (p *NoopProfile) Stop() error {
	return nil
}

func (p *NoopProfile) String() string {
	return "noop"
}

type Options struct {
	// Name to use for the profile
	Name string
}

type Option func(o *Options)

// Name of the profile
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
