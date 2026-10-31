// Package flow is an interface used for saga pattern microservice workflow.
//
// Idempotency: steps and their compensations may be re-run (step retries,
// execution resume after a crash), so both must be idempotent. Pass an
// idempotency key derived from the execution id and step id to external
// systems so that a re-execution does not duplicate side effects:
//
//	idempotencyKey := eid + ":" + stepID
//
// Compensation handlers receive the step's persisted response
// (CompensateInput.Rsp) and should use it (e.g. a transaction or charge id)
// to reverse exactly the effect the step produced.
package flow // import "go.unistack.org/micro/v5/flow"

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"go.unistack.org/micro/v5/metadata"
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

// FlowResult represents the result of a step execution
type FlowResult struct {
	StepID string
	Result *Message
	Error  error
}

// CompensateInput carries the context passed to Step.Compensate.
//
// It mirrors the rollback context of Cloudflare Workflows: the compensation
// handler receives the step's persisted output (Rsp), the error that started
// the rollback (Err) and the step's last request (Req).
type CompensateInput struct {
	// Req is the step's last request message.
	Req *Message
	// Rsp is the step's persisted response message. It is nil when the step
	// failed before producing a response, in which case the compensation may
	// be a no-op (the step likely created no external effect to undo).
	Rsp *Message
	// Err is the error that triggered the compensation.
	Err error
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
	// Compensate performs rollback for this step (Saga pattern).
	// It is invoked in reverse start order when a subsequent step fails.
	// in carries the step's last request (Req), its persisted response (Rsp,
	// nil when the step failed before producing one) and the error that
	// triggered the compensation (Err). Return nil if the step does not
	// require compensation or the compensation succeeded.
	Compensate(ctx context.Context, in *CompensateInput, opts ...ExecuteOption) error
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
	// StatusSuspended step suspended
	StatusSuspended
	// StatusCompensated step was compensated (rolled back) after a failure.
	// A compensated step is re-run on resume, since its effect was undone.
	StatusCompensated
	// StatusSkipped step was skipped by a conditional branch (its Condition
	// returned false) or pruned as a dependent of a skipped step. A skipped
	// step has no effect and is not compensated.
	StatusSkipped
)

var (
	// StatusString contains map status => string
	StatusString = map[Status]string{
		StatusPending:     "StatusPending",
		StatusRunning:     "StatusRunning",
		StatusFailure:     "StatusFailure",
		StatusSuccess:     "StatusSuccess",
		StatusAborted:     "StatusAborted",
		StatusSuspended:   "StatusSuspended",
		StatusCompensated: "StatusCompensated",
		StatusSkipped:     "StatusSkipped",
	}
	// StringStatus contains map string => status
	StringStatus = map[string]Status{
		"StatusPending":     StatusPending,
		"StatusRunning":     StatusRunning,
		"StatusFailure":     StatusFailure,
		"StatusSuccess":     StatusSuccess,
		"StatusAborted":     StatusAborted,
		"StatusSuspended":   StatusSuspended,
		"StatusCompensated": StatusCompensated,
		"StatusSkipped":     StatusSkipped,
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
	// Suspend suspends execution
	Suspend(ctx context.Context, id string) error
	// Resume resumes execution
	Resume(ctx context.Context, id string) error
	// Abort abort execution
	Abort(ctx context.Context, id string) error
	// Compensate runs the rollback of execution id on demand: it unwinds all
	// successful steps in reverse start order (the same semantics as the
	// automatic rollback on failure). It is idempotent (already-compensated
	// steps are skipped) and is the way to compensate an execution that was
	// left uncompensated (see ExecuteCompensate(false)).
	Compensate(ctx context.Context, id string) error
	// CompensateSteps runs the rollback of a subset of execution id's steps.
	// When ids is empty, all successful steps are compensated (same as
	// Compensate). When ids is non-empty, only the given steps are
	// compensated, in reverse start order. Idempotent like Compensate.
	CompensateSteps(ctx context.Context, id string, ids ...string) error
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
	// PublishEvent delivers payload to all steps currently waiting on topic
	// (see WaitForEvent). Steps that start waiting after the event was
	// published do not receive it (live delivery only).
	PublishEvent(ctx context.Context, topic string, payload *Message) error
	// Close releases resources held by the flow (goroutine pool)
	Close() error
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
