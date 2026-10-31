package flow

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStepTimeout verifies that a step whose attempt exceeds its configured
// timeout is canceled and the execution fails.
func TestStepTimeout(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := NewStep("slow",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			select {
			case <-time.After(2 * time.Second):
				return &Message{Body: []byte("done")}, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}),
		StepTimeout(100*time.Millisecond),
	)

	w, err := f.WorkflowCreate(ctx, "wf", a)
	require.NoError(t, err)

	// The step exceeds its timeout, so the execution fails and Execute
	// returns the error (the eid is still returned for inspection).
	eid, err := w.Execute(ctx, nil)
	require.Error(t, err)
	require.NotEmpty(t, eid)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusFailure, st.Status)
}

// TestStepTimeoutSucceeds verifies that a step that finishes within its
// timeout succeeds.
func TestStepTimeoutSucceeds(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := NewStep("fast",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			time.Sleep(10 * time.Millisecond)
			return &Message{Body: []byte("ok")}, nil
		}),
		StepTimeout(time.Second),
	)

	w, err := f.WorkflowCreate(ctx, "wf", a)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)
}

// TestConditionalBranch verifies that a step whose Condition returns true is
// executed.
func TestConditionalBranch(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := NewStep("a", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("go")}, nil
	}))
	b := NewStep("b",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte("b")}, nil
		}),
		StepCondition(func(ctx context.Context, req *Message) (bool, error) {
			return string(req.Body) == "go", nil
		}),
	)
	require.NoError(t, b.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf", a, b)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)

	bst, err := ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, bst.Status)
}

// TestConditionalBranchPrune verifies that a step whose Condition returns
// false is skipped and its transitive dependents are pruned, while unrelated
// branches still run.
func TestConditionalBranchPrune(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	// a -> b (only if input == "go") -> d
	// a -> c (always)
	a := NewStep("a", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("stop")}, nil
	}))
	b := NewStep("b",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte("b")}, nil
		}),
		StepCondition(func(ctx context.Context, req *Message) (bool, error) {
			return string(req.Body) == "go", nil
		}),
	)
	require.NoError(t, b.Require(a))
	d := NewStep("d", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("d")}, nil
	}))
	require.NoError(t, d.Require(b))
	c := NewStep("c", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("c")}, nil
	}))
	require.NoError(t, c.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf", a, b, c, d)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)

	// b is skipped (condition false).
	bst, err := ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	assert.Equal(t, StatusSkipped, bst.Status)

	// d is pruned (transitive dependent of skipped b).
	dst, err := ss.StepLoad(ctx, eid, "d")
	require.NoError(t, err)
	assert.Equal(t, StatusSkipped, dst.Status)

	// c still runs (unrelated branch).
	cst, err := ss.StepLoad(ctx, eid, "c")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, cst.Status)
}

// TestConditionalBranchResume verifies that a resumed execution keeps the same
// branch: a previously skipped step stays skipped.
func TestConditionalBranchResume(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := NewStep("a", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("stop")}, nil
	}))
	b := NewStep("b",
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return &Message{Body: []byte("b")}, nil
		}),
		StepCondition(func(ctx context.Context, req *Message) (bool, error) {
			return string(req.Body) == "go", nil
		}),
	)
	require.NoError(t, b.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf", a, b)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)

	// b was skipped.
	bst, err := ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	assert.Equal(t, StatusSkipped, bst.Status)

	// Resume: b must stay skipped (the branch is stable).
	require.NoError(t, w.Resume(ctx, eid))
	st = waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)

	bst, err = ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	assert.Equal(t, StatusSkipped, bst.Status)
}

// TestDelay verifies that a Delay step pauses the execution for the given
// duration.
func TestDelay(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := Delay("wait", 200*time.Millisecond)
	b := NewStep("b", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("b")}, nil
	}))
	require.NoError(t, b.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf", a, b)
	require.NoError(t, err)

	start := time.Now()
	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)
	assert.GreaterOrEqual(t, time.Since(start), 200*time.Millisecond)
}

// TestSubWorkflow verifies that a NewWorkflowStep runs its sub-workflow and
// returns the sub-workflow's last step response.
func TestSubWorkflow(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	// Sub-workflow: s1 -> s2.
	s1 := NewStep("s1", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("s1")}, nil
	}))
	s2 := NewStep("s2", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("s2")}, nil
	}))
	require.NoError(t, s2.Require(s1))

	// Parent workflow: a -> sub.
	a := NewStep("a", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("a")}, nil
	}))
	sub := NewWorkflowStep("sub", s1, s2)
	require.NoError(t, sub.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf", a, sub)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)

	subst, err := ss.StepLoad(ctx, eid, "sub")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, subst.Status)
	assert.Equal(t, []byte("s2"), subst.Rsp)
}

// TestParallelSteps verifies that two steps that depend on the same step run
// concurrently (the DAG scheduler does not serialize independent branches).
func TestParallelSteps(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	var running, maxRunning int32

	makeStep := func(name string) Step {
		return NewStep(name, StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			cur := atomic.AddInt32(&running, 1)
			for {
				m := atomic.LoadInt32(&maxRunning)
				if cur <= m || atomic.CompareAndSwapInt32(&maxRunning, m, cur) {
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
			atomic.AddInt32(&running, -1)
			return &Message{Body: []byte(name)}, nil
		}))
	}

	a := makeStep("a")
	b := makeStep("b")
	c := makeStep("c")
	require.NoError(t, b.Require(a))
	require.NoError(t, c.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf", a, b, c)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)
	// b and c share the dependency a, so they must have overlapped.
	assert.GreaterOrEqual(t, atomic.LoadInt32(&maxRunning), int32(2))
}

// TestWorkflowVersion verifies that the flow's version is stamped onto the
// execution record.
func TestWorkflowVersion(t *testing.T) {
	f, ss := newTestFlow(t, WorkflowVersion("v1"))
	ctx := context.Background()

	a := NewStep("a", StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
		return &Message{Body: []byte("a")}, nil
	}))

	w, err := f.WorkflowCreate(ctx, "wf", a)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)
	assert.Equal(t, "v1", st.Version)
}
