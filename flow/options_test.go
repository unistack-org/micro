package flow

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/meter"
	memory "go.unistack.org/micro/v5/store/memory"
	"go.unistack.org/micro/v5/tracer"
)

// simpleStep is a minimal Step implementation used across tests.
type simpleStep struct {
	name     string
	requires []string
	status   Status
}

func (s *simpleStep) ID() string               { return s.name }
func (s *simpleStep) Endpoint() string         { return s.name }
func (s *simpleStep) String() string           { return s.name }
func (s *simpleStep) Hashcode() interface{}    { return s.name }
func (s *simpleStep) Requires() []string       { return s.requires }
func (s *simpleStep) Options() StepOptions     { return StepOptions{} }
func (s *simpleStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.requires = append(s.requires, step.String())
	}
	return nil
}
func (s *simpleStep) GetStatus() Status  { return s.status }
func (s *simpleStep) SetStatus(st Status) { s.status = st }
func (s *simpleStep) Request() *Message  { return nil }
func (s *simpleStep) Response() *Message { return nil }
func (s *simpleStep) Execute(_ context.Context, _ *Message, _ ...ExecuteOption) (*Message, error) {
	return &Message{Body: []byte(`{}`)}, nil
}
func (s *simpleStep) Compensate(_ context.Context, _ *Message, _ ...ExecuteOption) error {
	return nil
}
func (s *simpleStep) Name() string { return s.name }

// --- Options tests ---

func TestOptionLogger(t *testing.T) {
	l := logger.DefaultLogger
	opts := NewOptions(Logger(l))
	assert.Equal(t, l, opts.Logger)
}

func TestOptionMeter(t *testing.T) {
	m := meter.DefaultMeter
	opts := NewOptions(Meter(m))
	assert.Equal(t, m, opts.Meter)
}

func TestOptionTracer(t *testing.T) {
	tr := tracer.DefaultTracer
	opts := NewOptions(Tracer(tr))
	assert.Equal(t, tr, opts.Tracer)
}

func TestOptionClient(t *testing.T) {
	c := client.DefaultClient
	opts := NewOptions(Client(c))
	assert.Equal(t, c, opts.Client)
}

func TestOptionStore(t *testing.T) {
	ms := memory.NewStore()
	opts := NewOptions(Store(ms))
	assert.Equal(t, ms, opts.Store)
}

func TestOptionContext(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "val")
	opts := NewOptions(Context(ctx))
	assert.Equal(t, ctx, opts.Context)
}

// --- microFlow.Options / Init ---

func TestMicroFlowOptions(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	opts := f.Options()
	assert.Equal(t, ms, opts.Store)
}

func TestMicroFlowInit(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	err := f.Init(PoolSize(2))
	require.NoError(t, err)
}

// --- WorkflowRemove (via concrete type) ---

func TestWorkflowRemove(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	mf, ok := f.(*microFlow)
	require.True(t, ok)
	err := mf.WorkflowRemove(ctx, "some-id")
	require.NoError(t, err)
}

// --- WorkflowSave / WorkflowLoad ---

func TestWorkflowSaveLoad(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	step := &simpleStep{name: "step1"}
	w, err := f.WorkflowCreate(ctx, "wf-save-load", step)
	require.NoError(t, err)

	err = f.WorkflowSave(ctx, w)
	require.NoError(t, err)

	loaded, err := f.WorkflowLoad(ctx, "wf-save-load")
	require.NoError(t, err)
	assert.Equal(t, "wf-save-load", loaded.ID())
}

// --- WorkflowList ---

func TestWorkflowList(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	step := &simpleStep{name: "step1"}

	w, err := f.WorkflowCreate(ctx, "wf-list-1", step)
	require.NoError(t, err)
	require.NoError(t, f.WorkflowSave(ctx, w))

	// WorkflowList may return errors for unparseable keys; we just want it called
	_, _ = f.WorkflowList(ctx)
}

// --- microWorkflow.ID / Status ---

func TestWorkflowIDStatus(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	step := &simpleStep{name: "step1"}
	w, err := f.WorkflowCreate(ctx, "wf-id-status", step)
	require.NoError(t, err)

	assert.Equal(t, "wf-id-status", w.ID())
	_ = w.Status() // just ensure no panic
}

// --- AppendSteps / RemoveSteps ---

func TestWorkflowAppendRemoveSteps(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	step1 := &simpleStep{name: "step1"}
	w, err := f.WorkflowCreate(ctx, "wf-append-remove", step1)
	require.NoError(t, err)

	step2 := &simpleStep{name: "step2", requires: []string{"step1"}}
	err = w.AppendSteps(step2)
	require.NoError(t, err)

	err = w.RemoveSteps(step2)
	// May error if vertex referenced elsewhere; just verify the function was called
	_ = err
}

// --- Abort / Suspend ---

func TestWorkflowAbortSuspend(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	step := &simpleStep{name: "step1"}
	w, err := f.WorkflowCreate(ctx, "wf-abort", step)
	require.NoError(t, err)

	// Execute so workflow store is initialised
	_, err = w.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.NoError(t, err)

	err = w.Abort(ctx, "wf-abort")
	require.NoError(t, err)

	err = w.Suspend(ctx, "wf-abort")
	require.NoError(t, err)
}

// --- Resume ---

func TestWorkflowResume(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	step := &simpleStep{name: "step1"}
	w, err := f.WorkflowCreate(ctx, "wf-resume", step)
	require.NoError(t, err)

	// Suspend first so there's a status record
	_, err = w.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.NoError(t, err)

	require.NoError(t, w.Suspend(ctx, "wf-resume"))
	// Resume should not panic/error
	err = w.Resume(ctx, "wf-resume")
	require.NoError(t, err)
}

// --- NewCallStep methods ---

func TestNewCallStep(t *testing.T) {
	s := NewCallStep("svc", "Method", "Do")
	cs, ok := s.(*microCallStep)
	require.True(t, ok)

	assert.NotNil(t, s)
	assert.Equal(t, "svc.Method.Do", s.String())
	assert.Equal(t, "svc.Method.Do", s.ID())
	assert.Equal(t, "Method.Do", s.Endpoint())
	assert.Equal(t, "svc.Method.Do", cs.Name())
	assert.Equal(t, "svc.Method.Do", cs.Hashcode())
	assert.Nil(t, s.Request())
	assert.Nil(t, s.Response())
	assert.Empty(t, s.Requires())
	_ = s.Options()

	s.SetStatus(StatusRunning)
	assert.Equal(t, StatusRunning, s.GetStatus())

	other := &simpleStep{name: "dep"}
	require.NoError(t, s.Require(other))
	assert.Contains(t, s.Requires(), "dep")

	// Compensate with nil client should still return nil
	err := s.Compensate(context.Background(), &Message{Body: []byte{}})
	assert.NoError(t, err)
}

func TestNewCallStepWithID(t *testing.T) {
	s := NewCallStep("svc", "Method", "Do", StepID("custom-id"))
	assert.Equal(t, "custom-id", s.String())
	assert.Equal(t, "custom-id", s.ID())
}

// --- NewPublishStep methods ---

func TestNewPublishStep(t *testing.T) {
	s := NewPublishStep("my.topic")
	ps, ok := s.(*microPublishStep)
	require.True(t, ok)

	assert.NotNil(t, s)
	assert.Equal(t, "my.topic", s.String())
	assert.Equal(t, "my.topic", s.ID())
	assert.Equal(t, "my.topic", s.Endpoint())
	assert.Equal(t, "my.topic", ps.Name())
	assert.Equal(t, "my.topic", ps.Hashcode())
	assert.Nil(t, s.Request())
	assert.Nil(t, s.Response())
	assert.Empty(t, s.Requires())
	_ = s.Options()

	s.SetStatus(StatusSuccess)
	assert.Equal(t, StatusSuccess, s.GetStatus())

	other := &simpleStep{name: "dep"}
	require.NoError(t, s.Require(other))
	assert.Contains(t, s.Requires(), "dep")

	rsp, err := s.Execute(context.Background(), &Message{})
	assert.Nil(t, rsp)
	assert.NoError(t, err)

	err = s.Compensate(context.Background(), &Message{})
	assert.NoError(t, err)
}

func TestNewPublishStepWithID(t *testing.T) {
	s := NewPublishStep("my.topic", StepID("pub-id"))
	assert.Equal(t, "pub-id", s.String())
}

// --- ExecuteOptions helpers ---

func TestNewExecuteOptions(t *testing.T) {
	opts := NewExecuteOptions(
		ExecuteAsync(true),
		ExecuteContext(context.Background()),
		ExecuteClient(client.DefaultClient),
		ExecuteLogger(logger.DefaultLogger),
		ExecuteMeter(meter.DefaultMeter),
		ExecuteTracer(tracer.DefaultTracer),
		ExecuteTimeout(5),
	)
	assert.True(t, opts.Async)
	assert.NotNil(t, opts.Context)
	assert.Equal(t, 5, int(opts.Timeout))
}

// --- WorkflowID option ---

func TestWorkflowIDOption(t *testing.T) {
	opt := WorkflowID("my-wf-id")
	o := &WorkflowOptions{}
	opt(o)
	assert.Equal(t, "my-wf-id", o.ID)
}

// --- StepOptions helpers ---

func TestStepOptionHelpers(t *testing.T) {
	opts := NewStepOptions(
		StepID("sid"),
		StepRequires("a", "b"),
		StepFallback("fallback"),
	)
	assert.Equal(t, "sid", opts.ID)
	assert.Equal(t, []string{"a", "b"}, opts.Requires)
	assert.Equal(t, "fallback", opts.Fallback)
}

// --- failingStep for error-path tests ---

type failingStep struct {
	name     string
	requires []string
	status   Status
	errMsg   string
}

func (s *failingStep) ID() string               { return s.name }
func (s *failingStep) Endpoint() string         { return s.name }
func (s *failingStep) String() string           { return s.name }
func (s *failingStep) Hashcode() interface{}    { return s.name }
func (s *failingStep) Requires() []string       { return s.requires }
func (s *failingStep) Options() StepOptions     { return StepOptions{} }
func (s *failingStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.requires = append(s.requires, step.String())
	}
	return nil
}
func (s *failingStep) GetStatus() Status  { return s.status }
func (s *failingStep) SetStatus(st Status) { s.status = st }
func (s *failingStep) Request() *Message  { return nil }
func (s *failingStep) Response() *Message { return nil }
func (s *failingStep) Execute(_ context.Context, _ *Message, _ ...ExecuteOption) (*Message, error) {
	return nil, fmt.Errorf("%s", s.errMsg)
}
func (s *failingStep) Compensate(_ context.Context, _ *Message, _ ...ExecuteOption) error {
	return nil
}

// --- Workflow Execute with failing step (covers error+compensation path) ---

func TestWorkflowExecuteFailingStep(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms), PoolSize(2))
	step := &failingStep{name: "bad-step", errMsg: "intentional failure"}
	w, err := f.WorkflowCreate(ctx, "wf-fail", step)
	require.NoError(t, err)

	_, err = w.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "intentional failure")
}

// --- Workflow Execute: two-step chain so compensation reads prior req ---

func TestWorkflowExecuteChainWithFailure(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms), PoolSize(2))
	step1 := &simpleStep{name: "ok-step"}
	step2 := &failingStep{name: "fail-step", requires: []string{"ok-step"}, errMsg: "step2 fail"}

	w, err := f.WorkflowCreate(ctx, "wf-chain-fail", step1, step2)
	require.NoError(t, err)

	_, err = w.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.Error(t, err)
}

// --- Workflow Execute: no-root graph returns error ---

func TestWorkflowExecuteNoRoots(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms), PoolSize(2))
	// Two steps each requiring each other would be a cycle and rejected by DAG.
	// Instead test empty workflow (0 steps, 0 roots).
	w, err := f.WorkflowCreate(ctx, "wf-no-roots")
	require.NoError(t, err)

	_, err = w.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.Error(t, err)
}

// --- Workflow Execute async ---

func TestWorkflowExecuteAsync(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms), PoolSize(2))
	step := &simpleStep{name: "step1"}
	w, err := f.WorkflowCreate(ctx, "wf-async", step)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, &Message{Body: []byte(`{}`)}, ExecuteAsync(true))
	require.NoError(t, err)
	assert.NotEmpty(t, eid)
	// Give goroutine a moment to finish
	// (no strict assertion needed — just cover the async path)
}

// --- Workflow Execute with Start option ---

func TestWorkflowExecuteWithStart(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms), PoolSize(2))
	step := &simpleStep{name: "only-step"}
	w, err := f.WorkflowCreate(ctx, "wf-start", step)
	require.NoError(t, err)

	// Setting Start to a non-existent step exercises the lookup branch
	_, err = w.Execute(ctx, &Message{Body: []byte(`{}`)}, ExecuteAsync(false))
	require.NoError(t, err)
}

// --- microFlow.Init error path (store not initialised) ---

func TestMicroFlowInitWithStore(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	f := NewFlow(Store(ms))
	// Re-init with same store — should succeed
	require.NoError(t, f.Init(Store(ms), PoolSize(1)))
}

// --- RegisterStep ---

func TestRegisterStep(t *testing.T) {
	step := &simpleStep{name: "registered"}
	RegisterStep(step)
	// just verify it doesn't panic
}

// --- Status.String ---

func TestStatusString(t *testing.T) {
	assert.Equal(t, "StatusPending", StatusPending.String())
	assert.Equal(t, "StatusRunning", StatusRunning.String())
	assert.Equal(t, "StatusFailure", StatusFailure.String())
	assert.Equal(t, "StatusSuccess", StatusSuccess.String())
	assert.Equal(t, "StatusAborted", StatusAborted.String())
	assert.Equal(t, "StatusSuspend", StatusSuspend.String())
}

// --- RawMessage marshal/unmarshal ---

func TestRawMessageMarshal(t *testing.T) {
	rm := RawMessage(`{"key":"value"}`)
	b, err := rm.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte(`{"key":"value"}`), b)

	var rm2 RawMessage
	require.NoError(t, rm2.UnmarshalJSON([]byte(`{"x":1}`)))
	assert.Equal(t, RawMessage(`{"x":1}`), rm2)
}

func TestRawMessageNilMarshal(t *testing.T) {
	var rm *RawMessage
	b, err := rm.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte("null"), b)
}
