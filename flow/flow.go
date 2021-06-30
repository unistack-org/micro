// Package flow is an interface used for saga pattern microservice workflow
package flow

import (
	"context"
	"errors"
)

var (
	ErrStepNotExists = errors.New("step not exists")
)

// Step represents dedicated workflow step
type Step interface {
	// ID returns step id
	ID() string
	// Endpoint returns rpc endpoint service_name.service_method or broker topic
	Endpoint() string
	// Execute step run
	Execute(ctx context.Context, req interface{}, opts ...ExecuteOption) error
	// Requires returns dependent steps
	Requires() []string
	// Options returns step options
	Options() StepOptions
	// Require add required steps
	Require(steps ...Step) error
	// String
	String() string
}

// Workflow contains all steps to execute
type Workflow interface {
	// ID returns id of the workflow
	ID() string
	// Steps returns steps slice where parallel steps returned on the same level
	Steps() [][]Step
	// Execute workflow with args, return execution id and error
	Execute(ctx context.Context, req interface{}, opts ...ExecuteOption) (string, error)
	// RemoveSteps remove steps from workflow
	RemoveSteps(ctx context.Context, steps ...Step) error
	// AppendSteps append steps to workflow
	AppendSteps(ctx context.Context, steps ...Step) error
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
