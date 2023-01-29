package fsm

// Options struct holding fsm options
type Options struct {
	// Initial state
	Initial string
	// Wrappers runs before state
	Wrappers []StateWrapper
	// DryRun mode
	DryRun bool
}

// Option func signature
type Option func(*Options)

// StateOptions holds state options
type StateOptions struct {
	DryRun bool
}

// StateDryRun says that state executes in dry run mode
func StateDryRun(b bool) StateOption {
	return func(o *StateOptions) {
		o.DryRun = b
	}
}

// StateOption func signature
type StateOption func(*StateOptions)

// InitialState sets init state for state machine
func InitialState(initial string) Option {
	return func(o *Options) {
		o.Initial = initial
	}
}

// WrapState adds a state Wrapper to a list of options passed into the fsm
func WrapState(w StateWrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

// NewOptions returns new Options struct filled by passed Option
func NewOptions(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return options
}
