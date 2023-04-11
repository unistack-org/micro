package fsm // import "go.unistack.org/micro/v4/fsm"

import (
	"context"
	"errors"
)

var (
	ErrInvalidState = errors.New("does not exists")
	StateEnd        = "end"
)

type State interface {
	Name() string
	Body() interface{}
}

// StateWrapper wraps the StateFunc and returns the equivalent
type StateWrapper func(StateFunc) StateFunc

// StateFunc called on state transition and return next step and error
type StateFunc func(ctx context.Context, state State, opts ...StateOption) (State, error)

type FSM interface {
	Start(context.Context, interface{}, ...Option) (interface{}, error)
	Current() string
	Reset()
	State(string, StateFunc)
}
