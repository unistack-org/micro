package fsm // import "go.unistack.org/micro/v3/fsm"

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrInvalidState = errors.New("does not exists")
	StateEnd        = "end"
)

// Options struct holding fsm options
type Options struct {
	// DryRun mode
	DryRun bool
	// Initial state
	Initial string
	// HooksBefore func slice runs in order before state
	HooksBefore []HookBeforeFunc
	// HooksAfter func slice runs in order after state
	HooksAfter []HookAfterFunc
}

// HookBeforeFunc func signature
type HookBeforeFunc func(ctx context.Context, state string, args interface{})

// HookAfterFunc func signature
type HookAfterFunc func(ctx context.Context, state string, args interface{})

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

// HookBefore provides hook func slice
func HookBefore(fns ...HookBeforeFunc) Option {
	return func(o *Options) {
		o.HooksBefore = fns
	}
}

// HookAfter provides hook func slice
func HookAfter(fns ...HookAfterFunc) Option {
	return func(o *Options) {
		o.HooksAfter = fns
	}
}

// StateFunc called on state transition and return next step and error
type StateFunc func(ctx context.Context, args interface{}, opts...StateOption) (string, interface{}, error)

// FSM is a finite state machine
type FSM struct {
	mu          sync.Mutex
	statesMap   map[string]StateFunc
	statesOrder []string
	opts        *Options
	current     string
}

// New creates a new finite state machine having the specified initial state
// with specified options
func New(opts ...Option) *FSM {
	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	return &FSM{
		statesMap: map[string]StateFunc{},
		opts:      options,
	}
}

// Current returns the current state
func (f *FSM) Current() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.current
}

// Current returns the current state
func (f *FSM) Reset() {
	f.mu.Lock()
	f.current = f.opts.Initial
	f.mu.Unlock()
}

// State adds state to fsm
func (f *FSM) State(state string, fn StateFunc) {
	f.mu.Lock()
	f.statesMap[state] = fn
	f.statesOrder = append(f.statesOrder, state)
	f.mu.Unlock()
}

// Init initialize fsm and check states

// Start runs state machine with provided data
func (f *FSM) Start(ctx context.Context, args interface{}, opts ...Option) (interface{}, error) {
	var err error
	var ok bool
	var fn StateFunc
	var nstate string

	f.mu.Lock()
	options := f.opts

	for _, opt := range opts {
		opt(options)
	}

	sopts := []StateOption{StateDryRun(options.DryRun)}
	
	cstate := options.Initial
	states := make(map[string]StateFunc, len(f.statesMap))
	for k, v := range f.statesMap {
		states[k] = v
	}
	f.current = cstate
	f.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			fn, ok = states[cstate]
			if !ok {
				return nil, fmt.Errorf(`state "%s" %w`, cstate, ErrInvalidState)
			}
			f.mu.Lock()
			f.current = cstate
			f.mu.Unlock()
			for _, fn := range options.HooksBefore {
				fn(ctx, cstate, args)
			}
			nstate, args, err = fn(ctx, args, sopts...)
			for _, fn := range options.HooksAfter {
				fn(ctx, cstate, args)
			}
			switch {
			case err != nil:
				return args, err
			case nstate == StateEnd:
				return args, nil
			case nstate == "":
				for idx := range f.statesOrder {
					if f.statesOrder[idx] == cstate && len(f.statesOrder) > idx+1 {
						nstate = f.statesOrder[idx+1]
					}
				}
			}
			cstate = nstate
		}
	}
}