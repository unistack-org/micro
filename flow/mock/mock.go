package mock

import (
	"context"
	"fmt"
	"sync"

	"go.unistack.org/micro/v4/flow"
)

// ExpectedWorkflowCreate represents an expected WorkflowCreate operation
type ExpectedWorkflowCreate struct {
	id     string
	steps  []flow.Step
	times  int
	called int
	mutex  sync.Mutex
	err    error
	result flow.Workflow
}

func (e *ExpectedWorkflowCreate) match(id string, steps []flow.Step) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.id != "" && e.id != id {
		return false
	}

	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedWorkflowSave represents an expected WorkflowSave operation
type ExpectedWorkflowSave struct {
	workflow flow.Workflow
	times    int
	called   int
	mutex    sync.Mutex
	err      error
}

func (e *ExpectedWorkflowSave) match(w flow.Workflow) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.workflow != nil && e.workflow.ID() != w.ID() {
		return false
	}

	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedWorkflowLoad represents an expected WorkflowLoad operation
type ExpectedWorkflowLoad struct {
	id     string
	times  int
	called int
	mutex  sync.Mutex
	err    error
	result flow.Workflow
}

func (e *ExpectedWorkflowLoad) match(id string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.id != "" && e.id != id {
		return false
	}

	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedWorkflowList represents an expected WorkflowList operation
type ExpectedWorkflowList struct {
	times  int
	called int
	mutex  sync.Mutex
	err    error
	result []flow.Workflow
}

func (e *ExpectedWorkflowList) match() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// ExpectedExecute represents an expected Execute operation on a workflow
type ExpectedExecute struct {
	id     string
	times  int
	called int
	mutex  sync.Mutex
	err    error
	result string
}

func (e *ExpectedExecute) match(id string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.id != "" && e.id != id {
		return false
	}

	if e.times > 0 && e.called >= e.times {
		return false
	}

	e.called++
	return true
}

// Flow is a mock implementation of the flow.Flow interface for testing
type Flow struct {
	expectedCreates []*ExpectedWorkflowCreate
	expectedSaves   []*ExpectedWorkflowSave
	expectedLoads   []*ExpectedWorkflowLoad
	expectedLists   []*ExpectedWorkflowList
	expectedExecutes []*ExpectedExecute

	workflows map[string]flow.Workflow
	opts      flow.Options
	err       error
	mutex     sync.RWMutex
}

// NewFlow creates a new mock flow
func NewFlow(opts ...flow.Option) *Flow {
	options := flow.NewOptions(opts...)
	return &Flow{
		workflows: make(map[string]flow.Workflow),
		opts:      options,
	}
}

// ExpectWorkflowCreate creates an expectation for a WorkflowCreate operation
func (m *Flow) ExpectWorkflowCreate(id string) *ExpectedWorkflowCreate {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedWorkflowCreate{id: id}
	m.expectedCreates = append(m.expectedCreates, exp)
	return exp
}

// ExpectWorkflowSave creates an expectation for a WorkflowSave operation
func (m *Flow) ExpectWorkflowSave() *ExpectedWorkflowSave {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedWorkflowSave{}
	m.expectedSaves = append(m.expectedSaves, exp)
	return exp
}

// ExpectWorkflowLoad creates an expectation for a WorkflowLoad operation
func (m *Flow) ExpectWorkflowLoad(id string) *ExpectedWorkflowLoad {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedWorkflowLoad{id: id}
	m.expectedLoads = append(m.expectedLoads, exp)
	return exp
}

// ExpectWorkflowList creates an expectation for a WorkflowList operation
func (m *Flow) ExpectWorkflowList() *ExpectedWorkflowList {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedWorkflowList{}
	m.expectedLists = append(m.expectedLists, exp)
	return exp
}

// ExpectExecute creates an expectation for an Execute operation
func (m *Flow) ExpectExecute(id string) *ExpectedExecute {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	exp := &ExpectedExecute{id: id}
	m.expectedExecutes = append(m.expectedExecutes, exp)
	return exp
}

// WithResult sets the result to return for expected WorkflowCreate operations
func (e *ExpectedWorkflowCreate) WithResult(w flow.Workflow) *ExpectedWorkflowCreate {
	e.result = w
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedWorkflowCreate) Times(n int) *ExpectedWorkflowCreate {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedWorkflowCreate) WillReturnError(err error) *ExpectedWorkflowCreate {
	e.err = err
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedWorkflowSave) Times(n int) *ExpectedWorkflowSave {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedWorkflowSave) WillReturnError(err error) *ExpectedWorkflowSave {
	e.err = err
	return e
}

// WithResult sets the result to return for expected WorkflowLoad operations
func (e *ExpectedWorkflowLoad) WithResult(w flow.Workflow) *ExpectedWorkflowLoad {
	e.result = w
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedWorkflowLoad) Times(n int) *ExpectedWorkflowLoad {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedWorkflowLoad) WillReturnError(err error) *ExpectedWorkflowLoad {
	e.err = err
	return e
}

// WithResult sets the result to return for expected WorkflowList operations
func (e *ExpectedWorkflowList) WithResult(workflows []flow.Workflow) *ExpectedWorkflowList {
	e.result = workflows
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedWorkflowList) Times(n int) *ExpectedWorkflowList {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedWorkflowList) WillReturnError(err error) *ExpectedWorkflowList {
	e.err = err
	return e
}

// WithResult sets the execution ID to return for expected Execute operations
func (e *ExpectedExecute) WithResult(execID string) *ExpectedExecute {
	e.result = execID
	return e
}

// Times sets how many times the expectation should be called
func (e *ExpectedExecute) Times(n int) *ExpectedExecute {
	e.times = n
	return e
}

// WillReturnError sets an error to return for the expected operation
func (e *ExpectedExecute) WillReturnError(err error) *ExpectedExecute {
	e.err = err
	return e
}

// Options returns the current options
func (m *Flow) Options() flow.Options {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.opts
}

// Init initializes the mock flow
func (m *Flow) Init(opts ...flow.Option) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.err != nil {
		return m.err
	}

	for _, o := range opts {
		o(&m.opts)
	}
	return nil
}

// WorkflowCreate creates a new workflow with specific id and steps
func (m *Flow) WorkflowCreate(ctx context.Context, id string, steps ...flow.Step) (flow.Workflow, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.err != nil {
		return nil, m.err
	}

	// Find matching expectation
	for _, exp := range m.expectedCreates {
		if exp.match(id, steps) {
			if exp.err != nil {
				return nil, exp.err
			}
			if exp.result != nil {
				m.workflows[id] = exp.result
				return exp.result, nil
			}
			// Create default mock workflow if no result specified
			w := NewWorkflow(id)
			m.workflows[id] = w
			return w, nil
		}
	}

	// If no expectation matched, use default behavior
	w := NewWorkflow(id)
	m.workflows[id] = w
	return w, nil
}

// WorkflowSave saves workflow
func (m *Flow) WorkflowSave(ctx context.Context, w flow.Workflow) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.err != nil {
		return m.err
	}

	// Find matching expectation
	for _, exp := range m.expectedSaves {
		if exp.match(w) {
			if exp.err != nil {
				return exp.err
			}
			m.workflows[w.ID()] = w
			return nil
		}
	}

	// If no expectation matched, use default behavior
	m.workflows[w.ID()] = w
	return nil
}

// WorkflowLoad loads workflow with specific id
func (m *Flow) WorkflowLoad(ctx context.Context, id string) (flow.Workflow, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.err != nil {
		return nil, m.err
	}

	// Find matching expectation
	for _, exp := range m.expectedLoads {
		if exp.match(id) {
			if exp.err != nil {
				return nil, exp.err
			}
			if exp.result != nil {
				return exp.result, nil
			}
			// Return stored workflow if exists
			if w, ok := m.workflows[id]; ok {
				return w, nil
			}
			return nil, fmt.Errorf("workflow not found: %s", id)
		}
	}

	// If no expectation matched, use default behavior
	if w, ok := m.workflows[id]; ok {
		return w, nil
	}
	return nil, fmt.Errorf("workflow not found: %s", id)
}

// WorkflowList lists all workflows
func (m *Flow) WorkflowList(ctx context.Context) ([]flow.Workflow, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.err != nil {
		return nil, m.err
	}

	// Find matching expectation
	for _, exp := range m.expectedLists {
		if exp.match() {
			if exp.err != nil {
				return nil, exp.err
			}
			if exp.result != nil {
				return exp.result, nil
			}
			// Return all stored workflows
			result := make([]flow.Workflow, 0, len(m.workflows))
			for _, w := range m.workflows {
				result = append(result, w)
			}
			return result, nil
		}
	}

	// If no expectation matched, use default behavior
	result := make([]flow.Workflow, 0, len(m.workflows))
	for _, w := range m.workflows {
		result = append(result, w)
	}
	return result, nil
}

// ExpectationsWereMet checks that all expected operations were called the expected number of times
func (m *Flow) ExpectationsWereMet() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, exp := range m.expectedCreates {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected workflow create for id %s to be called %d times, but was called %d times", exp.id, exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedSaves {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected workflow save to be called %d times, but was called %d times", exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedLoads {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected workflow load for id %s to be called %d times, but was called %d times", exp.id, exp.times, exp.called)
		}
	}

	for _, exp := range m.expectedLists {
		if exp.times > 0 && exp.called != exp.times {
			return fmt.Errorf("expected workflow list to be called %d times, but was called %d times", exp.times, exp.called)
		}
	}

	return nil
}

// Workflow is a mock implementation of the flow.Workflow interface for testing
type Workflow struct {
	id        string
	status    flow.Status
	steps     map[string]flow.Step
	executed  bool
	executeErr error
	mutex     sync.RWMutex
}

// NewWorkflow creates a new mock workflow
func NewWorkflow(id string) *Workflow {
	return &Workflow{
		id:    id,
		steps: make(map[string]flow.Step),
	}
}

// ID returns the workflow id
func (w *Workflow) ID() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.id
}

// Status returns the workflow status
func (w *Workflow) Status() flow.Status {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.status
}

// SetStatus sets the workflow status
func (w *Workflow) SetStatus(status flow.Status) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = status
}

// SetExecuteError sets an error to return on Execute
func (w *Workflow) SetExecuteError(err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.executeErr = err
}

// AppendSteps appends steps to the workflow
func (w *Workflow) AppendSteps(steps ...flow.Step) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, s := range steps {
		w.steps[s.String()] = s
	}
	return nil
}

// RemoveSteps removes steps from the workflow
func (w *Workflow) RemoveSteps(steps ...flow.Step) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, s := range steps {
		delete(w.steps, s.String())
	}
	return nil
}

// Execute executes the workflow
func (w *Workflow) Execute(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) (string, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if w.executeErr != nil {
		return "", w.executeErr
	}

	w.executed = true
	return w.id, nil
}

// Suspend suspends the workflow
func (w *Workflow) Suspend(ctx context.Context, id string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = flow.StatusSuspend
	return nil
}

// Resume resumes the workflow
func (w *Workflow) Resume(ctx context.Context, id string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = flow.StatusRunning
	return nil
}

// Abort aborts the workflow
func (w *Workflow) Abort(ctx context.Context, id string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.status = flow.StatusAborted
	return nil
}

// Step is a mock implementation of the flow.Step interface for testing
type Step struct {
	id           string
	endpoint     string
	requires     []string
	executeFunc  func(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) (*flow.Message, error)
	compensateFunc func(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) error
	executeErr   error
	compensateErr error
	request      *flow.Message
	response     *flow.Message
	status       flow.Status
	options      flow.StepOptions
	mutex        sync.RWMutex
}

// NewStep creates a new mock step
func NewStep(id string) *Step {
	return &Step{
		id:       id,
		requires: make([]string, 0),
		status:   flow.StatusPending,
	}
}

// ID returns the step id
func (s *Step) ID() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.id
}

// Endpoint returns the step endpoint
func (s *Step) Endpoint() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.endpoint
}

// SetEndpoint sets the step endpoint
func (s *Step) SetEndpoint(endpoint string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.endpoint = endpoint
}

// Requires returns the required steps
func (s *Step) Requires() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.requires
}

// Require adds required steps
func (s *Step) Require(steps ...flow.Step) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, step := range steps {
		s.requires = append(s.requires, step.String())
	}
	return nil
}

// Options returns the step options
func (s *Step) Options() flow.StepOptions {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.options
}

// SetOptions sets the step options
func (s *Step) SetOptions(opts flow.StepOptions) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.options = opts
}

// String returns the step string representation
func (s *Step) String() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.id
}

// GetStatus returns the step status
func (s *Step) GetStatus() flow.Status {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.status
}

// SetStatus sets the step status
func (s *Step) SetStatus(status flow.Status) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.status = status
}

// Request returns the step request message
func (s *Step) Request() *flow.Message {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.request
}

// SetRequest sets the step request message
func (s *Step) SetRequest(req *flow.Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.request = req
}

// Response returns the step response message
func (s *Step) Response() *flow.Message {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.response
}

// SetResponse sets the step response message
func (s *Step) SetResponse(rsp *flow.Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.response = rsp
}

// SetExecuteError sets an error to return on Execute
func (s *Step) SetExecuteError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.executeErr = err
}

// SetCompensateError sets an error to return on Compensate
func (s *Step) SetCompensateError(err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.compensateErr = err
}

// SetExecuteFunc sets a custom execute function
func (s *Step) SetExecuteFunc(fn func(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) (*flow.Message, error)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.executeFunc = fn
}

// SetCompensateFunc sets a custom compensate function
func (s *Step) SetCompensateFunc(fn func(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.compensateFunc = fn
}

// Execute executes the step
func (s *Step) Execute(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) (*flow.Message, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.request = req

	if s.executeFunc != nil {
		s.mutex.Unlock()
		rsp, err := s.executeFunc(ctx, req, opts...)
		s.mutex.Lock()
		s.response = rsp
		return rsp, err
	}

	if s.executeErr != nil {
		s.status = flow.StatusFailure
		return nil, s.executeErr
	}

	s.status = flow.StatusSuccess
	return s.response, nil
}

// Compensate performs rollback for this step
func (s *Step) Compensate(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) error {
	s.mutex.Lock()

	if s.compensateFunc != nil {
		fn := s.compensateFunc
		s.mutex.Unlock()
		return fn(ctx, req, opts...)
	}

	if s.compensateErr != nil {
		s.mutex.Unlock()
		return s.compensateErr
	}

	s.mutex.Unlock()
	return nil
}
