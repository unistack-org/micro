package flow

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWaitForEvent verifies that a WaitForEvent step blocks until
// Flow.PublishEvent delivers a payload on the topic, and that the payload
// becomes the step response (and the input of dependent steps).
func TestWaitForEvent(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	wait := WaitForEvent("approval")
	approve := NewStep("approve",
		StepRequires("approval"),
		StepExecute(func(ctx context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
			return req, nil
		}),
	)

	w, err := f.WorkflowCreate(ctx, "wf-wait", wait, approve)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	// Wait until the wait step is running.
	require.Eventually(t, func() bool {
		st, err := ss.StepLoad(ctx, eid, "approval")
		return err == nil && st.Status == StatusRunning
	}, 5*time.Second, 10*time.Millisecond)

	// Publish the event; the wait step must complete with the payload.
	require.NoError(t, f.PublishEvent(ctx, "approval", &Message{Body: []byte(`{"ok":true}`)}))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusSuccess, st.Status)

	wst, err := ss.StepLoad(ctx, eid, "approval")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, wst.Status)
	assert.JSONEq(t, `{"ok":true}`, string(wst.Rsp))

	// The dependent step received the event payload as its input.
	ast, err := ss.StepLoad(ctx, eid, "approve")
	require.NoError(t, err)
	assert.Equal(t, StatusSuccess, ast.Status)
	assert.JSONEq(t, `{"ok":true}`, string(ast.Req))
}

// TestWaitForEventAbort verifies that a WaitForEvent step is released when the
// execution is aborted while waiting.
func TestWaitForEventAbort(t *testing.T) {
	f, ss := newTestFlow(t)
	ctx := context.Background()

	wait := WaitForEvent("approval")
	w, err := f.WorkflowCreate(ctx, "wf-wait-abort", wait)
	require.NoError(t, err)
	eid, err := w.Execute(ctx, nil, ExecuteAsync(true))
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		st, err := ss.StepLoad(ctx, eid, "approval")
		return err == nil && st.Status == StatusRunning
	}, 5*time.Second, 10*time.Millisecond)

	require.NoError(t, w.Abort(ctx, eid))

	st := waitFinished(t, ss, eid)
	assert.Equal(t, StatusAborted, st.Status)
}

// TestPublishEventNoWaiters verifies PublishEvent is a no-op when no step is
// waiting on the topic.
func TestPublishEventNoWaiters(t *testing.T) {
	f, _ := newTestFlow(t)
	require.NoError(t, f.PublishEvent(context.Background(), "nobody", &Message{}))
}
