package flow

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.unistack.org/micro/v5/store/memory"
)

// engineStep is a Step implementation for engine tests.
type engineStep struct {
	name     string
	requires []string
	opts     StepOptions
	status   Status

	mu   sync.Mutex
	fn   func(ctx context.Context, req *Message) (*Message, error)
	seen []string

	attempts int32
}

func (s *engineStep) ID() string            { return s.name }
func (s *engineStep) Endpoint() string      { return s.name }
func (s *engineStep) String() string        { return s.name }
func (s *engineStep) Hashcode() interface{} { return s.name }
func (s *engineStep) Requires() []string    { return s.requires }
func (s *engineStep) Options() StepOptions  { return s.opts }

func (s *engineStep) Require(steps ...Step) error {
	for _, st := range steps {
		s.requires = append(s.requires, st.String())
	}
	return nil
}

func (s *engineStep) GetStatus() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *engineStep) SetStatus(st Status) {
	s.mu.Lock()
	s.status = st
	s.mu.Unlock()
}

func (s *engineStep) Request() *Message  { return nil }
func (s *engineStep) Response() *Message { return nil }

func (s *engineStep) setFn(fn func(ctx context.Context, req *Message) (*Message, error)) {
	s.mu.Lock()
	s.fn = fn
	s.mu.Unlock()
}

func (s *engineStep) Compensate(context.Context, *Message, ...ExecuteOption) error {
	return nil
}

func (s *engineStep) Execute(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
	atomic.AddInt32(&s.attempts, 1)
	s.mu.Lock()
	fn := s.fn
	s.mu.Unlock()
	if fn == nil {
		s.recordInput(req)
		return &Message{Body: []byte(fmt.Sprintf(`{"step":%q}`, s.name))}, nil
	}
	rsp, err := fn(ctx, req)
	if rsp != nil {
		s.recordInput(req)
	}
	return rsp, err
}

func (s *engineStep) recordInput(req *Message) {
	s.mu.Lock()
	s.seen = append(s.seen, string(req.Body))
	s.mu.Unlock()
}

func (s *engineStep) input(i int) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i < len(s.seen) {
		return s.seen[i]
	}
	return ""
}

// newTestFlow creates an initialized flow backed by a memory store and
// returns the flow with a StateStore view over the same store.
func newTestFlow(t *testing.T, opts ...Option) (Flow, StateStore) {
	t.Helper()
	s := memory.NewStore()
	require.NoError(t, s.Init())
	require.NoError(t, s.Connect(context.Background()))
	f := NewFlow(append([]Option{Store(s)}, opts...)...)
	require.NoError(t, f.Init())
	t.Cleanup(func() { _ = f.Close() })
	return f, NewKVStateStore(s)
}

// waitFinished polls until the execution has reached a terminal write
// (FinishedAt is set) and returns the final record.
func waitFinished(t *testing.T, ss StateStore, eid string) *WorkflowState {
	t.Helper()
	var out *WorkflowState
	require.Eventually(t, func() bool {
		st, err := ss.WorkflowLoad(context.Background(), eid)
		if err != nil {
			return false
		}
		out = st
		return !st.FinishedAt.IsZero()
	}, 5*time.Second, 10*time.Millisecond)
	return out
}

// waitStepStatus polls until the given step record exists with the expected status.
func waitStepStatus(t *testing.T, ss StateStore, eid, stepID string, st Status) *StepState {
	t.Helper()
	var out *StepState
	require.Eventually(t, func() bool {
		var err error
		out, err = ss.StepLoad(context.Background(), eid, stepID)
		return err == nil && out != nil && out.Status == st
	}, 5*time.Second, 10*time.Millisecond)
	return out
}

// waitWorkflowStatus polls until the execution reaches the given status.
func waitWorkflowStatus(t *testing.T, ss StateStore, eid string, st Status) *WorkflowState {
	t.Helper()
	var out *WorkflowState
	require.Eventually(t, func() bool {
		var err error
		out, err = ss.WorkflowLoad(context.Background(), eid)
		return err == nil && out != nil && out.Status == st
	}, 5*time.Second, 10*time.Millisecond)
	return out
}

func TestEngineExecutePersists(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	b := &engineStep{name: "b"}
	require.NoError(t, b.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf-persist", a, b)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, &Message{Body: []byte(`{"in":1}`)})
	require.NoError(t, err)
	assert.NotEmpty(t, eid)

	st, err := ss.WorkflowLoad(ctx, eid)
	require.NoError(t, err)
	assert.Equal(t, "wf-persist", st.WorkflowID)
	assert.Equal(t, StatusSuccess, st.Status)
	assert.False(t, st.FinishedAt.IsZero())

	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, ast.Status)
	assert.Equal(t, `{"in":1}`, string(ast.Req))

	bst, err := ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, bst.Status)
	// b must receive a's response as its input.
	assert.Equal(t, string(ast.Rsp), string(bst.Req))
	assert.EqualValues(t, 1, a.attempts)
	assert.EqualValues(t, 1, b.attempts)
}

func TestEngineFailureCompensates(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	b := &engineStep{name: "b"}
	c := &engineStep{name: "c"}
	require.NoError(t, b.Require(a))
	require.NoError(t, c.Require(b))

	b.setFn(func(context.Context, *Message) (*Message, error) {
		return nil, errors.New("boom")
	})

	w, err := f.WorkflowCreate(ctx, "wf-fail", a, b, c)
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.Error(t, err)

	st, err := ss.WorkflowLoad(ctx, eid)
	require.NoError(t, err)
	assert.Equal(t, StatusFailure, st.Status)

	// a succeeded and was compensated back to Pending.
	ast, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusPending, ast.Status)

	// b failed with its error persisted.
	bst, err := ss.StepLoad(ctx, eid, "b")
	require.NoError(t, err)
	assert.Equal(t, StatusFailure, bst.Status)
	assert.Equal(t, "boom", bst.Error)

	// c never started.
	_, err = ss.StepLoad(ctx, eid, "c")
	assert.Error(t, err)
	assert.EqualValues(t, 0, c.attempts)
}

func TestEngineResumeSkipsCompleted(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	const eid = "eid-resume"
	// Simulate a previously failed execution where a already succeeded.
	require.NoError(t, ss.WorkflowSave(ctx, eid, &WorkflowState{
		EID:        eid,
		WorkflowID: "wf-resume",
		Status:     StatusFailure,
		Graph:      map[string][]string{"b": {"a"}},
		StartedAt:  time.Now(),
	}))
	require.NoError(t, ss.StepSave(ctx, eid, "a", &StepState{
		EID:    eid,
		StepID: "a",
		Status: StatusSuccess,
		Rsp:    []byte(`{"from":"a"}`),
	}))

	a := &engineStep{name: "a"}
	b := &engineStep{name: "b"}
	require.NoError(t, b.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf-resume", a, b)
	require.NoError(t, err)
	require.NoError(t, w.Resume(ctx, eid))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)
	assert.EqualValues(t, 0, a.attempts) // a was skipped: already succeeded
	assert.EqualValues(t, 1, b.attempts)
	assert.Equal(t, `{"from":"a"}`, b.input(0))
}

func TestEngineNoSteps(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	w, err := f.WorkflowCreate(ctx, "wf-empty")
	require.NoError(t, err)

	eid, err := w.Execute(ctx, nil)
	require.Error(t, err)
	assert.NotEmpty(t, eid)

	st, err := ss.WorkflowLoad(ctx, eid)
	require.NoError(t, err)
	assert.Equal(t, StatusFailure, st.Status)
}

func TestEngineCrossProcessRestart(t *testing.T) {
	ctx := context.Background()
	s := memory.NewStore()
	require.NoError(t, s.Init())
	require.NoError(t, s.Connect(ctx))
	ss := NewKVStateStore(s)

	a1 := &engineStep{name: "a"}
	b1 := &engineStep{name: "b"}
	require.NoError(t, b1.Require(a1))
	b1.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		time.Sleep(300 * time.Millisecond)
		return &Message{Body: []byte(`{"b":true}`)}, nil
	})

	flow1 := NewFlow(Store(s))
	require.NoError(t, flow1.Init())
	w1, err := flow1.WorkflowCreate(ctx, "wf-restart", a1, b1)
	require.NoError(t, err)

	eid, err := w1.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	// Wait until a has completed and b is running, then "crash".
	require.Eventually(t, func() bool {
		ast, err := ss.StepLoad(ctx, eid, "a")
		if err != nil || ast.Status != StatusSuccess {
			return false
		}
		bst, err := ss.StepLoad(ctx, eid, "b")
		return err == nil && bst.Status == StatusRunning
	}, 5*time.Second, 10*time.Millisecond)

	require.NoError(t, flow1.Close())
	// A supervisor marks the stale execution as failed so it can be resumed.
	require.NoError(t, ss.WorkflowSetStatus(ctx, eid, StatusFailure))

	// Restart in a "new process": fresh flow and fresh step objects.
	flow2 := NewFlow(Store(s))
	require.NoError(t, flow2.Init())
	t.Cleanup(func() { _ = flow2.Close() })

	a2 := &engineStep{name: "a"}
	b2 := &engineStep{name: "b"}
	require.NoError(t, b2.Require(a2))
	w2, err := flow2.WorkflowCreate(ctx, "wf-restart", a2, b2)
	require.NoError(t, err)
	require.NoError(t, w2.Resume(ctx, eid))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)
	assert.EqualValues(t, 0, a2.attempts) // a was skipped: persisted before the crash
}

func TestEngineStepRetry(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	t.Run("SucceedsAfterRetries", func(t *testing.T) {
		a := &engineStep{name: "a"}
		a.opts = StepOptions{Retry: &RetryPolicy{MaxAttempts: 3, Backoff: 5 * time.Millisecond}}
		var n int32
		a.setFn(func(context.Context, *Message) (*Message, error) {
			if atomic.AddInt32(&n, 1) < 3 {
				return nil, errors.New("flaky")
			}
			return &Message{Body: []byte(`{"ok":true}`)}, nil
		})

		w, err := f.WorkflowCreate(ctx, "wf-retry", a)
		require.NoError(t, err)
		eid, err := w.Execute(ctx, nil)
		require.NoError(t, err)
		assert.EqualValues(t, 3, atomic.LoadInt32(&n))

		st, err := ss.StepLoad(ctx, eid, "a")
		require.NoError(t, err)
		assert.Equal(t, StatusSuccess, st.Status)
		assert.Equal(t, 3, st.Attempt)
	})

	t.Run("ExhaustsRetries", func(t *testing.T) {
		a := &engineStep{name: "a"}
		a.opts = StepOptions{Retry: &RetryPolicy{MaxAttempts: 2, Backoff: 2 * time.Millisecond}}
		a.setFn(func(context.Context, *Message) (*Message, error) {
			return nil, errors.New("always")
		})

		w, err := f.WorkflowCreate(ctx, "wf-retry-exhaust", a)
		require.NoError(t, err)
		eid, err := w.Execute(ctx, nil)
		require.Error(t, err)

		st, err := ss.StepLoad(ctx, eid, "a")
		require.NoError(t, err)
		assert.Equal(t, StatusFailure, st.Status)
		assert.Equal(t, 2, st.Attempt)
		assert.Equal(t, "always", st.Error)
	})
}

func TestEngineCheckpoint(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	a.opts = StepOptions{Retry: &RetryPolicy{MaxAttempts: 2}}
	var resumedFrom string
	a.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		start := 0
		if saved := StepLastCheckpoint(ctx); len(saved) > 0 {
			resumedFrom = string(saved)
			var n int
			n, _ = strconv.Atoi(string(saved))
			start = n + 1
		}
		for i := start; i < 10; i++ {
			if err := StepCheckpoint(ctx, []byte(strconv.Itoa(i))); err != nil {
				return nil, err
			}
			if i == 7 {
				return nil, errors.New("fail-at-7")
			}
		}
		return &Message{Body: []byte(`{"done":true}`)}, nil
	})

	w, err := f.WorkflowCreate(ctx, "wf-checkpoint", a)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	// The second attempt must have resumed from the persisted checkpoint.
	assert.Equal(t, "7", resumedFrom)

	st, err := ss.StepLoad(ctx, eid, "a")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, st.Status)
	assert.Equal(t, 2, st.Attempt)
	assert.Nil(t, st.Checkpoint) // checkpoint cleared on success
}

func TestEngineCleanup(t *testing.T) {
	f, ss := newTestFlow(t, Cleanup(20*time.Millisecond, 5*time.Millisecond))
	ctx := context.Background()

	a := &engineStep{name: "a"}
	w, err := f.WorkflowCreate(ctx, "wf-cleanup", a)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	// The finished execution must be deleted once it is older than the TTL.
	require.Eventually(t, func() bool {
		_, err := ss.WorkflowLoad(ctx, eid)
		return err != nil
	}, 5*time.Second, 10*time.Millisecond)
}

func TestEngineAbort(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	b := &engineStep{name: "b"}
	require.NoError(t, b.Require(a))
	b.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})

	w, err := f.WorkflowCreate(ctx, "wf-abort", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	waitStepStatus(t, ss, eid, "b", StatusRunning)
	require.NoError(t, w.Abort(ctx, eid))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusAborted, st.Status)
}

func TestEngineSuspendResume(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	b := &engineStep{name: "b"}
	require.NoError(t, b.Require(a))
	b.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})

	w, err := f.WorkflowCreate(ctx, "wf-suspend", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	waitStepStatus(t, ss, eid, "b", StatusRunning)
	require.NoError(t, w.Suspend(ctx, eid))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuspend, st.Status)

	// Let b succeed and resume the execution.
	b.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		return &Message{Body: []byte(`{"b":true}`)}, nil
	})
	require.NoError(t, w.Resume(ctx, eid))

	// waitFinished would return the stale record of the first run here;
	// wait for the terminal status of the resumed run instead.
	st = waitWorkflowStatus(t, ss, eid, StatusSuccess)
	assert.False(t, st.FinishedAt.IsZero())
	assert.EqualValues(t, 2, b.attempts)
}

func TestEngineCrossProcessAbort(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	a.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		time.Sleep(300 * time.Millisecond)
		return &Message{Body: []byte(`{"a":true}`)}, nil
	})
	b := &engineStep{name: "b"}
	require.NoError(t, b.Require(a))

	w, err := f.WorkflowCreate(ctx, "wf-xabort", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	// Abort from "another process" while a is still running; b must not start.
	require.Eventually(t, func() bool {
		_, err := ss.WorkflowStatus(ctx, eid)
		return err == nil
	}, 5*time.Second, 10*time.Millisecond)
	require.NoError(t, ss.WorkflowSetStatus(ctx, eid, StatusAborted))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusAborted, st.Status)
	assert.EqualValues(t, 0, b.attempts) // b saw the aborted status before starting
}

func TestEngineCrossProcessAbortInFlight(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	b := &engineStep{name: "b"}
	require.NoError(t, b.Require(a))
	var cancelled int32
	b.setFn(func(ctx context.Context, req *Message) (*Message, error) {
		<-ctx.Done()
		atomic.StoreInt32(&cancelled, 1)
		return nil, ctx.Err()
	})

	w, err := f.WorkflowCreate(ctx, "wf-xabort-flight", a, b)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	waitStepStatus(t, ss, eid, "b", StatusRunning)

	// Set the status directly in the store (another process); the poller
	// must cancel the in-flight execution.
	require.NoError(t, ss.WorkflowSetStatus(ctx, eid, StatusAborted))

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&cancelled) == 1
	}, 15*time.Second, 100*time.Millisecond)

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusAborted, st.Status)
}

func TestEngineWorkflowListRemove(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	a := &engineStep{name: "a"}
	w, err := f.WorkflowCreate(ctx, "wf-list", a)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil)
	require.NoError(t, err)

	found := false
	workflows, err := f.WorkflowList(ctx)
	require.NoError(t, err)
	for _, wf := range workflows {
		if wf.ID() == eid {
			found = true
		}
	}
	assert.True(t, found)

	require.NoError(t, ss.WorkflowDelete(ctx, eid))
	_, err = ss.WorkflowLoad(ctx, eid)
	assert.Error(t, err)
}

func TestEngineNoopStore(t *testing.T) {
	// No Store option: the flow falls back to the no-op state store.
	f := NewFlow()
	require.NoError(t, f.Init())
	t.Cleanup(func() { _ = f.Close() })

	a := &engineStep{name: "a"}
	w, err := f.WorkflowCreate(context.Background(), "wf-noop", a)
	require.NoError(t, err)
	_, err = w.Execute(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, w.Status())
}
