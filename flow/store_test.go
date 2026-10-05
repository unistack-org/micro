package flow

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.unistack.org/micro/v5/store/memory"
)

// StateStoreConformance verifies that a StateStore implementation
// satisfies the flow.StateStore contract. It is exported so custom
// implementations can be conformance-tested from their own package:
//
//	func TestMyStateStore(t *testing.T) {
//		flow.StateStoreConformance(t, func() flow.StateStore { return newMyStore() })
//	}
func StateStoreConformance(t *testing.T, newStore func() StateStore) {
	ctx := context.Background()
	ss := newStore()
	other := newStore()

	t.Run("WorkflowRoundtrip", func(t *testing.T) {
		started := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
		require.NoError(t, ss.WorkflowSave(ctx, "eid-1", &WorkflowState{
			EID:        "eid-1",
			WorkflowID: "wf",
			Status:     StatusRunning,
			Graph:      map[string][]string{"b": {"a"}},
			StartedAt:  started,
		}))

		st, err := ss.WorkflowLoad(ctx, "eid-1")
		require.NoError(t, err)
		assert.Equal(t, "eid-1", st.EID)
		assert.Equal(t, "wf", st.WorkflowID)
		assert.Equal(t, StatusRunning, st.Status)
		assert.Equal(t, map[string][]string{"b": {"a"}}, st.Graph)
		assert.Equal(t, started, st.StartedAt)

		status, err := ss.WorkflowStatus(ctx, "eid-1")
		require.NoError(t, err)
		assert.Equal(t, StatusRunning, status)

		require.NoError(t, ss.WorkflowSetStatus(ctx, "eid-1", StatusAborted))
		status, err = ss.WorkflowStatus(ctx, "eid-1")
		require.NoError(t, err)
		assert.Equal(t, StatusAborted, status)

		require.NoError(t, ss.WorkflowSetLastStep(ctx, "eid-1", "b"))
		st, err = ss.WorkflowLoad(ctx, "eid-1")
		require.NoError(t, err)
		assert.Equal(t, "b", st.LastStep)
	})

	t.Run("WorkflowNotFound", func(t *testing.T) {
		_, err := other.WorkflowLoad(ctx, "missing")
		assert.Error(t, err)
	})

	t.Run("StepRoundtrip", func(t *testing.T) {
		require.NoError(t, ss.WorkflowSave(ctx, "eid-2", &WorkflowState{
			EID: "eid-2", WorkflowID: "wf", Status: StatusRunning,
		}))
		require.NoError(t, ss.StepSave(ctx, "eid-2", "s1", &StepState{
			EID:     "eid-2",
			StepID:  "s1",
			Req:     []byte(`{"in":1}`),
			Rsp:     []byte(`{"out":2}`),
			Error:   "step error",
			Attempt: 2,
		}))

		st, err := ss.StepLoad(ctx, "eid-2", "s1")
		require.NoError(t, err)
		assert.Equal(t, "eid-2", st.EID)
		assert.Equal(t, "s1", st.StepID)
		assert.Equal(t, []byte(`{"in":1}`), st.Req)
		assert.Equal(t, []byte(`{"out":2}`), st.Rsp)
		assert.Equal(t, "step error", st.Error)
		assert.Equal(t, 2, st.Attempt)

		// partial update: only the status changes
		require.NoError(t, ss.StepSave(ctx, "eid-2", "s1", &StepState{
			EID: "eid-2", StepID: "s1", Status: StatusSuccess,
		}))
		st, err = ss.StepLoad(ctx, "eid-2", "s1")
		require.NoError(t, err)
		assert.Equal(t, StatusSuccess, st.Status)
		assert.Equal(t, []byte(`{"in":1}`), st.Req)
		assert.Equal(t, []byte(`{"out":2}`), st.Rsp)
		assert.Equal(t, "step error", st.Error)
		assert.Equal(t, 2, st.Attempt)
	})

	t.Run("StepList", func(t *testing.T) {
		require.NoError(t, ss.WorkflowSave(ctx, "eid-3", &WorkflowState{
			EID: "eid-3", WorkflowID: "wf", Status: StatusRunning,
		}))
		require.NoError(t, ss.StepSave(ctx, "eid-3", "s1", &StepState{EID: "eid-3", StepID: "s1", Status: StatusSuccess}))
		require.NoError(t, ss.StepSave(ctx, "eid-3", "s2", &StepState{EID: "eid-3", StepID: "s2", Status: StatusRunning}))

		steps, err := ss.StepList(ctx, "eid-3")
		require.NoError(t, err)
		assert.Len(t, steps, 2)
		assert.Equal(t, StatusSuccess, steps["s1"].Status)
		assert.Equal(t, StatusRunning, steps["s2"].Status)
	})

	t.Run("Checkpoint", func(t *testing.T) {
		require.NoError(t, ss.WorkflowSave(ctx, "eid-4", &WorkflowState{
			EID: "eid-4", WorkflowID: "wf", Status: StatusRunning,
		}))
		require.NoError(t, ss.StepSave(ctx, "eid-4", "s", &StepState{
			EID: "eid-4", StepID: "s", Status: StatusRunning,
		}))

		require.NoError(t, ss.StepSetCheckpoint(ctx, "eid-4", "s", []byte("cp-1")))
		st, err := ss.StepLoad(ctx, "eid-4", "s")
		require.NoError(t, err)
		assert.Equal(t, []byte("cp-1"), st.Checkpoint)

		// a save without checkpoint clears it
		require.NoError(t, ss.StepSave(ctx, "eid-4", "s", &StepState{
			EID: "eid-4", StepID: "s", Status: StatusSuccess,
		}))
		st, err = ss.StepLoad(ctx, "eid-4", "s")
		require.NoError(t, err)
		assert.Nil(t, st.Checkpoint)
	})

	t.Run("WorkflowList", func(t *testing.T) {
		lss := newStore()
		require.NoError(t, lss.WorkflowSave(ctx, "l-eid-1", &WorkflowState{
			EID: "l-eid-1", WorkflowID: "wf-a", Status: StatusRunning,
			StartedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		}))
		require.NoError(t, lss.WorkflowSave(ctx, "l-eid-2", &WorkflowState{
			EID: "l-eid-2", WorkflowID: "wf-a", Status: StatusSuccess,
			StartedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		}))
		require.NoError(t, lss.WorkflowSave(ctx, "l-eid-3", &WorkflowState{
			EID: "l-eid-3", WorkflowID: "wf-b", Status: StatusFailure,
			StartedAt: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
		}))

		// all entries, deduplicated, newest first
		list, err := lss.WorkflowList(ctx, &WorkflowFilter{})
		require.NoError(t, err)
		require.Len(t, list, 3)
		assert.Equal(t, "l-eid-3", list[0].EID)
		assert.Equal(t, StatusFailure, list[0].State.Status)
		assert.Equal(t, "l-eid-2", list[1].EID)
		assert.Equal(t, StatusSuccess, list[1].State.Status)
		assert.Equal(t, "l-eid-1", list[2].EID)
		assert.Equal(t, StatusRunning, list[2].State.Status)

		// filter by workflow id
		list, err = lss.WorkflowList(ctx, &WorkflowFilter{WorkflowID: "wf-a"})
		require.NoError(t, err)
		require.Len(t, list, 2)
		assert.Equal(t, "l-eid-2", list[0].EID)
		assert.Equal(t, "l-eid-1", list[1].EID)

		// filter by status
		list, err = lss.WorkflowList(ctx, &WorkflowFilter{Statuses: []Status{StatusRunning}})
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, "l-eid-1", list[0].EID)

		// limit
		list, err = lss.WorkflowList(ctx, &WorkflowFilter{Limit: 1})
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, "l-eid-3", list[0].EID)

		// offset + limit
		list, err = lss.WorkflowList(ctx, &WorkflowFilter{Offset: 1, Limit: 1})
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, "l-eid-2", list[0].EID)
	})

	t.Run("WorkflowDelete", func(t *testing.T) {
		dss := newStore()
		require.NoError(t, dss.WorkflowSave(ctx, "d-eid", &WorkflowState{
			EID: "d-eid", WorkflowID: "wf", Status: StatusRunning,
		}))
		require.NoError(t, dss.StepSave(ctx, "d-eid", "s", &StepState{
			EID: "d-eid", StepID: "s", Status: StatusRunning,
		}))

		require.NoError(t, dss.WorkflowDelete(ctx, "d-eid"))

		_, err := dss.WorkflowLoad(ctx, "d-eid")
		assert.Error(t, err)
		_, err = dss.WorkflowStatus(ctx, "d-eid")
		assert.Error(t, err)

		steps, err := dss.StepList(ctx, "d-eid")
		require.NoError(t, err)
		assert.Empty(t, steps)
	})
}

func TestKVStateStore(t *testing.T) {
	newStore := func() StateStore {
		s := memory.NewStore()
		require.NoError(t, s.Init())
		require.NoError(t, s.Connect(context.Background()))
		return NewKVStateStore(s)
	}
	StateStoreConformance(t, newStore)
}

func TestKVStateStoreWatchWorkflow(t *testing.T) {
	s := memory.NewStore()
	require.NoError(t, s.Init())
	require.NoError(t, s.Connect(context.Background()))
	ss := NewKVStateStore(s)

	w, ok := ss.(WorkflowWatcher)
	require.True(t, ok)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	ch, err := w.WatchWorkflow(ctx, "watch-eid")
	require.NoError(t, err)

	// No initial value for a non-existent execution.
	select {
	case st := <-ch:
		t.Fatalf("unexpected initial status: %v", st)
	case <-time.After(50 * time.Millisecond):
	}

	require.NoError(t, ss.WorkflowSetStatus(ctx, "watch-eid", StatusRunning))
	select {
	case st := <-ch:
		assert.Equal(t, StatusRunning, st)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for pushed status")
	}

	require.NoError(t, ss.WorkflowSetStatus(ctx, "watch-eid", StatusAborted))
	select {
	case st := <-ch:
		assert.Equal(t, StatusAborted, st)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for pushed status")
	}

	// Cancellation of the subscription context closes the channel.
	cancel()
	select {
	case _, ok := <-ch:
		assert.False(t, ok)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}
