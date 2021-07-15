// Package flow is an interface used for saga pattern microservice workflow
package flow

import (
	"context"
	"errors"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	ErrStepNotExists = errors.New("step not exists")
	ErrMissingClient = errors.New("client not set")
)

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can be used to delay decoding or precompute a encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return *m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("RawMessage UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

type Message struct {
	Header metadata.Metadata
	Body   RawMessage
}

// Step represents dedicated workflow step
type Step interface {
	// ID returns step id
	ID() string
	// Endpoint returns rpc endpoint service_name.service_method or broker topic
	Endpoint() string
	// Execute step run
	Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error)
	// Requires returns dependent steps
	Requires() []string
	// Options returns step options
	Options() StepOptions
	// Require add required steps
	Require(steps ...Step) error
	// String
	String() string
	// GetStatus returns step status
	GetStatus() Status
	// SetStatus sets the step status
	SetStatus(Status)
	// Request returns step request message
	Request() *Message
	// Response returns step response message
	Response() *Message
}

type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusFailure
	StatusSuccess
	StatusAborted
	StatusSuspend
)

// Workflow contains all steps to execute
type Workflow interface {
	// ID returns id of the workflow
	ID() string
	// Execute workflow with args, return execution id and error
	Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (string, error)
	// RemoveSteps remove steps from workflow
	RemoveSteps(steps ...Step) error
	// AppendSteps append steps to workflow
	AppendSteps(steps ...Step) error
	// Status returns workflow status
	Status() Status
	// Steps returns steps slice where parallel steps returned on the same level
	Steps() ([][]Step, error)
}

// Flow the base interface to interact with workflows
type Flow interface {
	// Options returns options
	Options() Options
	// Init initialize
	Init(...Option) error
	// WorkflowCreate creates new workflow with specific id and steps
	WorkflowCreate(ctx context.Context, id string, steps ...Step) (Workflow, error)
	// WorkflowSave saves workflow
	WorkflowSave(ctx context.Context, w Workflow) error
	// WorkflowLoad loads workflow with specific id
	WorkflowLoad(ctx context.Context, id string) (Workflow, error)
	// WorkflowList lists all workflows
	WorkflowList(ctx context.Context) ([]Workflow, error)
}
