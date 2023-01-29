package fsm

import (
	"context"
	"fmt"
	"sync"
)

type state struct {
	body interface{}
	name string
}

var _ State = &state{}

func (s *state) Name() string {
	return s.name
}

func (s *state) Body() interface{} {
	return s.body
}

// fsm is a finite state machine
type fsm struct {
	statesMap   map[string]StateFunc
	current     string
	statesOrder []string
	opts        Options
	mu          sync.Mutex
}

// NewFSM creates a new finite state machine having the specified initial state
// with specified options
func NewFSM(opts ...Option) *fsm {
	return &fsm{
		statesMap: map[string]StateFunc{},
		opts:      NewOptions(opts...),
	}
}

// Current returns the current state
func (f *fsm) Current() string {
	f.mu.Lock()
	s := f.current
	f.mu.Unlock()
	return s
}

// Current returns the current state
func (f *fsm) Reset() {
	f.mu.Lock()
	f.current = f.opts.Initial
	f.mu.Unlock()
}

// State adds state to fsm
func (f *fsm) State(state string, fn StateFunc) {
	f.mu.Lock()
	f.statesMap[state] = fn
	f.statesOrder = append(f.statesOrder, state)
	f.mu.Unlock()
}

// Start runs state machine with provided data
func (f *fsm) Start(ctx context.Context, args interface{}, opts ...Option) (interface{}, error) {
	var err error

	f.mu.Lock()
	options := f.opts

	for _, opt := range opts {
		opt(&options)
	}

	sopts := []StateOption{StateDryRun(options.DryRun)}

	cstate := options.Initial
	states := make(map[string]StateFunc, len(f.statesMap))
	for k, v := range f.statesMap {
		states[k] = v
	}
	f.current = cstate
	f.mu.Unlock()

	var s State
	s = &state{name: cstate, body: args}
	nstate := s.Name()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			fn, ok := states[nstate]
			if !ok {
				return nil, fmt.Errorf(`state "%s" %w`, nstate, ErrInvalidState)
			}
			f.mu.Lock()
			f.current = nstate
			f.mu.Unlock()

			// wrap the handler func
			for i := len(options.Wrappers); i > 0; i-- {
				fn = options.Wrappers[i-1](fn)
			}

			s, err = fn(ctx, s, sopts...)

			switch {
			case err != nil:
				return s.Body(), err
			case s.Name() == StateEnd:
				return s.Body(), nil
			case s.Name() == "":
				for idx := range f.statesOrder {
					if f.statesOrder[idx] == nstate && len(f.statesOrder) > idx+1 {
						nstate = f.statesOrder[idx+1]
					}
				}
			default:
				nstate = s.Name()
			}
		}
	}
}
