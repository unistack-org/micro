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
	// Initial state
	Initial string
	// HooksBefore func slice runs in order before state
	HooksBefore []HookBeforeFunc
	// HooksAfter func slice runs in order after state
	HooksAfter []HookAfterFunc
}

// HookBeforeFunc func signature
type HookBeforeFunc func(ctx context.Context, state string, args map[string]interface{})

// HookAfterFunc func signature
type HookAfterFunc func(ctx context.Context, state string, args map[string]interface{})

// Option func signature
type Option func(*Options)

// StateInitial sets init state for state machine
func StateInitial(initial string) Option {
	return func(o *Options) {
		o.Initial = initial
	}
}

// StateHookBefore provides hook func slice
func StateHookBefore(fns ...HookBeforeFunc) Option {
	return func(o *Options) {
		o.HooksBefore = fns
	}
}

// StateHookAfter provides hook func slice
func StateHookAfter(fns ...HookAfterFunc) Option {
	return func(o *Options) {
		o.HooksAfter = fns
	}
}

// StateFunc called on state transition and return next step and error
type StateFunc func(ctx context.Context, args map[string]interface{}) (string, map[string]interface{}, error)

// FSM is a finite state machine
type FSM struct {
	mu      sync.Mutex
	states  map[string]StateFunc
	opts    *Options
	current string
}

// New creates a new finite state machine having the specified initial state
// with specified options
func New(opts ...Option) *FSM {
	options := &Options{}

	for _, opt := range opts {
		opt(options)
	}

	return &FSM{
		states: map[string]StateFunc{},
		opts:   options,
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
	f.states[state] = fn
	f.mu.Unlock()
}

// Start runs state machine with provided data
func (f *FSM) Start(ctx context.Context, args map[string]interface{}, opts ...Option) (map[string]interface{}, error) {
	var err error
	var ok bool
	var fn StateFunc
	var nstate string

	f.mu.Lock()
	options := f.opts

	for _, opt := range opts {
		opt(options)
	}

	cstate := options.Initial
	states := make(map[string]StateFunc, len(f.states))
	for k, v := range f.states {
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
			nstate, args, err = fn(ctx, args)
			for _, fn := range options.HooksAfter {
				fn(ctx, cstate, args)
			}
			if err != nil {
				return args, err
			} else if nstate == "" || nstate == StateEnd {
				return args, nil
			}
			cstate = nstate
		}
	}
}
