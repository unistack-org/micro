// Package flow is an interface used for saga pattern microservice workflow
package flow // import "go.unistack.org/micro/v4/flow"

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"go.unistack.org/micro/v4/metadata"
)

var (
	// ErrStepNotExists returns when step not found
	ErrStepNotExists = errors.New("step not exists")
	// ErrMissingClient returns when client.Client is missing
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

// Message used to transfer data between steps
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

// Status contains step current status
type Status int

func (status Status) String() string {
	return StatusString[status]
}

const (
	// StatusPending step waiting to start
	StatusPending Status = iota
	// StatusRunning step is running
	StatusRunning
	// StatusFailure step competed with error
	StatusFailure
	// StatusSuccess step completed without error
	StatusSuccess
	// StatusAborted step aborted while it running
	StatusAborted
	// StatusSuspend step suspended
	StatusSuspend
)

var (
	// StatusString contains map status => string
	StatusString = map[Status]string{
		StatusPending: "StatusPending",
		StatusRunning: "StatusRunning",
		StatusFailure: "StatusFailure",
		StatusSuccess: "StatusSuccess",
		StatusAborted: "StatusAborted",
		StatusSuspend: "StatusSuspend",
	}
	// StringStatus contains map string => status
	StringStatus = map[string]Status{
		"StatusPending": StatusPending,
		"StatusRunning": StatusRunning,
		"StatusFailure": StatusFailure,
		"StatusSuccess": StatusSuccess,
		"StatusAborted": StatusAborted,
		"StatusSuspend": StatusSuspend,
	}
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
	// Suspend suspends execution
	Suspend(ctx context.Context, id string) error
	// Resume resumes execution
	Resume(ctx context.Context, id string) error
	// Abort abort execution
	Abort(ctx context.Context, id string) error
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

var (
	flowMu      sync.Mutex
	atomicSteps atomic.Value
)

// RegisterStep register own step with workflow
func RegisterStep(step Step) {
	flowMu.Lock()
	steps, _ := atomicSteps.Load().([]Step)
	atomicSteps.Store(append(steps, step))
	flowMu.Unlock()
}
