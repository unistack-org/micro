package flow

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewStepEndToEnd verifies the functional step constructor (NewStep) runs
// a workflow and compensates the successful step when a later step fails.
func TestNewStepEndToEnd(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{"id":"a"}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-newstep", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil)
	require.Error(t, err)

	st, err := ss.WorkflowLoad(ctx, eid)
	require.NoError(t, err)
	assert.Equal(t, StatusFailure, st.Status)

	// a succeeded and was compensated.
	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusCompensated, ast.Status)
}

// TestCompensateReceivesResponse verifies the compensation handler receives the
// step's persisted response (Rsp) and the error that triggered the rollback
// (Err), mirroring Cloudflare Workflows' rollback context.
func TestCompensateReceivesResponse(t *testing.T) {
	f, _ := newTestFlow(t)
	ctx := context.Background()

	var (
		mu     sync.Mutex
		gotRsp *Message
		gotErr error
	)

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{"id":"txn-1"}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			mu.Lock()
			gotRsp = in.Rsp
			gotErr = in.Err
			mu.Unlock()
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-rsp", a, b)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotNil(t, gotRsp, "compensation must receive the step response")
	assert.JSONEq(t, `{"id":"txn-1"}`, string(gotRsp.Body))
	require.NotNil(t, gotErr, "compensation must receive the triggering error")
	assert.Contains(t, gotErr.Error(), "boom")
}

// TestCompensateReverseOrder verifies compensation runs in reverse start order
// (the stable source of truth for rollback sequencing).
func TestCompensateReverseOrder(t *testing.T) {
	f, _ := newTestFlow(t)
	ctx := context.Background()

	var (
		mu    sync.Mutex
		order []string
	)
	record := func(name string) {
		mu.Lock()
		order = append(order, name)
		mu.Unlock()
	}

	mk := func(id, dep string, fail bool) Step {
		opts := []StepOption{
			StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
				if fail {
					return nil, errors.New("boom")
				}
				return &Message{Body: []byte(`{}`)}, nil
			}),
			StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
				record(id)
				return nil
			}),
		}
		if dep != "" {
			opts = append(opts, StepRequires(dep))
		}
		return NewStep(id, opts...)
	}

	a := mk("a", "", false)
	b := mk("b", "a", false)
	c := mk("c", "b", true) // c fails

	w, err := f.WorkflowCreate(ctx, "wf-order", a, b, c)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	// c failed (not compensated); b then a in reverse start order.
	assert.Equal(t, []string{"b", "a"}, order)
}

// TestCompensateRetry verifies the compensation retry policy is applied.
func TestCompensateRetry(t *testing.T) {
	f, _ := newTestFlow(t)
	ctx := context.Background()

	var attempts int32

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			if atomic.AddInt32(&attempts, 1) < 3 {
				return errors.New("comp fail")
			}
			return nil
		}),
		StepCompensateRetry(RetryPolicy{MaxAttempts: 3, Backoff: time.Millisecond}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-compretry", a, b)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)

	assert.EqualValues(t, 3, atomic.LoadInt32(&attempts))
}

// TestCompensateStopOnFailure verifies Cloudflare Workflows semantics: when a
// compensation fails (after its retries), the remaining compensations are not
// run.
func TestCompensateStopOnFailure(t *testing.T) {
	f, _ := newTestFlow(t)
	ctx := context.Background()

	var (
		mu          sync.Mutex
		compensated []string
	)
	record := func(name string) {
		mu.Lock()
		compensated = append(compensated, name)
		mu.Unlock()
	}

	mk := func(id, dep string, failExec, failComp bool) Step {
		opts := []StepOption{
			StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
				if failExec {
					return nil, errors.New("boom")
				}
				return &Message{Body: []byte(`{}`)}, nil
			}),
			StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
				record(id)
				if failComp {
					return errors.New("comp fail")
				}
				return nil
			}),
		}
		if dep != "" {
			opts = append(opts, StepRequires(dep))
		}
		return NewStep(id, opts...)
	}

	a := mk("a", "", false, false)
	b := mk("b", "a", false, true) // b's compensation fails
	c := mk("c", "b", true, false) // c fails, triggering rollback

	w, err := f.WorkflowCreate(ctx, "wf-stop", a, b, c)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	// b (higher start order) is compensated first and fails -> a is skipped.
	assert.Equal(t, []string{"b"}, compensated)
}

// TestCompensateHooks verifies the lifecycle hooks fire around compensation.
func TestCompensateHooks(t *testing.T) {
	var (
		mu        sync.Mutex
		started   []string
		succeeded []string
		failed    []string
		completed []bool
	)
	hooks := &Hooks{
		OnCompensateStart: func(eid, stepID string) {
			mu.Lock()
			started = append(started, stepID)
			mu.Unlock()
		},
		OnCompensateSuccess: func(eid, stepID string) {
			mu.Lock()
			succeeded = append(succeeded, stepID)
			mu.Unlock()
		},
		OnCompensateFailure: func(eid, stepID string, err error) {
			mu.Lock()
			failed = append(failed, stepID)
			mu.Unlock()
		},
		OnCompensateComplete: func(eid string, failed bool) {
			mu.Lock()
			completed = append(completed, failed)
			mu.Unlock()
		},
	}
	f, _ := newTestFlow(t, WithHooks(hooks))
	ctx := context.Background()

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-hooks", a, b)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"a"}, started)
	assert.Equal(t, []string{"a"}, succeeded)
	assert.Empty(t, failed)
	require.Len(t, completed, 1)
	assert.False(t, completed[0])
}

// TestCompensateIdempotentResume verifies a compensated step is re-run (not
// skipped) on resume and is not double-compensated within a single pass.
func TestCompensateIdempotentResume(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	var execA, compA int32

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			atomic.AddInt32(&execA, 1)
			return &Message{Body: []byte(`{"id":"a"}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			atomic.AddInt32(&compA, 1)
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-idem", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil)
	require.Error(t, err)

	// First run: a executed once, compensated once.
	assert.EqualValues(t, 1, atomic.LoadInt32(&execA))
	assert.EqualValues(t, 1, atomic.LoadInt32(&compA))

	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusCompensated, ast.Status)

	// Resume: a (StatusCompensated) re-runs; b fails again; a compensated again.
	require.NoError(t, w.Resume(ctx, eid))
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&execA) == 2
	}, 5*time.Second, 10*time.Millisecond)

	assert.EqualValues(t, 2, atomic.LoadInt32(&execA))
	assert.EqualValues(t, 2, atomic.LoadInt32(&compA))
}

// TestCompensateOrderPersisted verifies the start order is persisted to the
// state store and used for deterministic rollback sequencing.
func TestCompensateOrderPersisted(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	mk := func(id, dep string, fail bool) Step {
		opts := []StepOption{
			StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
				if fail {
					return nil, errors.New("boom")
				}
				return &Message{Body: []byte(`{}`)}, nil
			}),
		}
		if dep != "" {
			opts = append(opts, StepRequires(dep))
		}
		return NewStep(id, opts...)
	}

	a := mk("a", "", false)
	b := mk("b", "a", false)
	c := mk("c", "b", true)

	w, err := f.WorkflowCreate(ctx, "wf-orderpersist", a, b, c)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil)
	require.Error(t, err)

	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	bst, err := ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	cst, err := ss.StepLoad(ctx, eid, "c")
	require.NoError(t, err)

	// Start order is assigned sequentially and persisted.
	assert.Equal(t, 1, ast.Order)
	assert.Equal(t, 2, bst.Order)
	assert.Equal(t, 3, cst.Order)
}

// TestNewStepNoExecute verifies a NewStep without StepExecute returns a clear
// error when run.
func TestNewStepNoExecute(t *testing.T) {
	f, _ := newTestFlow(t)
	ctx := context.Background()

	s := NewStep("noexec")
	w, err := f.WorkflowCreate(ctx, "wf-noexec", s)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no Execute function")
}

// TestExecuteCompensateFalse verifies that ExecuteCompensate(false) leaves the
// successful steps intact on failure (no automatic rollback), so the caller can
// later retry or compensate the execution explicitly.
func TestExecuteCompensateFalse(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	var compA int32

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{"id":"a"}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			atomic.AddInt32(&compA, 1)
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-nocomp", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteCompensate(false))
	require.Error(t, err)

	// The execution failed, but a was NOT compensated (no automatic rollback).
	st, err := ss.WorkflowLoad(ctx, eid)
	require.NoError(t, err)
	assert.Equal(t, StatusFailure, st.Status)

	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, ast.Status)
	assert.EqualValues(t, 0, atomic.LoadInt32(&compA), "no automatic compensation")
}

// TestExecuteCompensateFalseThenCompensate verifies the on-demand Compensate
// unwinds the successful steps of an execution that was left uncompensated.
func TestExecuteCompensateFalseThenCompensate(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	var compA int32

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{"id":"a"}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			atomic.AddInt32(&compA, 1)
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-nocomp2", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteCompensate(false))
	require.Error(t, err)
	assert.EqualValues(t, 0, atomic.LoadInt32(&compA))

	// On-demand compensation unwinds a.
	require.NoError(t, w.Compensate(ctx, eid))
	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&compA) == 1
	}, 5*time.Second, 10*time.Millisecond)

	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusCompensated, ast.Status)

	// The execution's terminal status is restored (still failure).
	wst := waitWorkflowStatus(t, ss, eid, StatusFailure)
	assert.Equal(t, StatusFailure, wst.Status)
}

// TestCompensateInputFields verifies CompensateInput carries Rsp and Err.
func TestCompensateInputFields(t *testing.T) {
	f, _ := newTestFlow(t)
	ctx := context.Background()

	var (
		mu     sync.Mutex
		gotRsp *Message
		gotErr error
	)

	a := NewStep("a",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte(`{"out":1}`)}, nil
		}),
		StepCompensate(func(ctx context.Context, in *CompensateInput, _ ...ExecuteOption) error {
			mu.Lock()
			gotRsp = in.Rsp
			gotErr = in.Err
			mu.Unlock()
			return nil
		}),
	)
	b := NewStep("b",
		StepRequires("a"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return nil, errors.New("boom")
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-fields", a, b)
	require.NoError(t, err)
	_, err = w.Execute(ctx, nil)
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	require.NotNil(t, gotRsp)
	assert.JSONEq(t, `{"out":1}`, string(gotRsp.Body))
	require.NotNil(t, gotErr)
	assert.Contains(t, gotErr.Error(), "boom")
}
