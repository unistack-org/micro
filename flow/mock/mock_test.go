package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.unistack.org/micro/v5/flow"
)

func TestFlow(t *testing.T) {
	ctx := context.Background()
	f := NewFlow()

	// Test WorkflowCreate with expectation
	f.ExpectWorkflowCreate("test_workflow")
	w, err := f.WorkflowCreate(ctx, "test_workflow")
	require.NoError(t, err)
	require.NotNil(t, w)
	assert.Equal(t, "test_workflow", w.ID())

	// Test WorkflowSave with expectation
	f.ExpectWorkflowSave()
	err = f.WorkflowSave(ctx, w)
	require.NoError(t, err)

	// Test WorkflowLoad with expectation
	f.ExpectWorkflowLoad("test_workflow").WithResult(w)
	loadedW, err := f.WorkflowLoad(ctx, "test_workflow")
	require.NoError(t, err)
	require.NotNil(t, loadedW)
	assert.Equal(t, "test_workflow", loadedW.ID())

	// Test WorkflowList with expectation
	f.ExpectWorkflowList().WithResult([]flow.Workflow{w})
	workflows, err := f.WorkflowList(ctx)
	require.NoError(t, err)
	assert.Len(t, workflows, 1)

	// Verify all expectations were met
	require.NoError(t, f.ExpectationsWereMet())
}

func TestWorkflow(t *testing.T) {
	ctx := context.Background()
	w := NewWorkflow("test_workflow")

	// Test ID
	assert.Equal(t, "test_workflow", w.ID())

	// Test Status
	assert.Equal(t, flow.StatusPending, w.Status())

	// Test SetStatus
	w.SetStatus(flow.StatusRunning)
	assert.Equal(t, flow.StatusRunning, w.Status())

	// Test AppendSteps
	step := NewStep("step1")
	err := w.AppendSteps(step)
	require.NoError(t, err)

	// Test Execute
	execID, err := w.Execute(ctx, &flow.Message{})
	require.NoError(t, err)
	assert.Equal(t, "test_workflow", execID)

	// Test Suspend
	err = w.Suspend(ctx, "test_workflow")
	require.NoError(t, err)
	assert.Equal(t, flow.StatusSuspend, w.Status())

	// Test Resume
	err = w.Resume(ctx, "test_workflow")
	require.NoError(t, err)
	assert.Equal(t, flow.StatusRunning, w.Status())

	// Test Abort
	err = w.Abort(ctx, "test_workflow")
	require.NoError(t, err)
	assert.Equal(t, flow.StatusAborted, w.Status())

	// Test SetExecuteError
	w.SetExecuteError(assert.AnError)
	_, err = w.Execute(ctx, &flow.Message{})
	assert.Error(t, err)
}

func TestStep(t *testing.T) {
	ctx := context.Background()
	s := NewStep("test_step")

	// Test ID
	assert.Equal(t, "test_step", s.ID())

	// Test String
	assert.Equal(t, "test_step", s.String())

	// Test Status
	assert.Equal(t, flow.StatusPending, s.GetStatus())

	// Test SetStatus
	s.SetStatus(flow.StatusRunning)
	assert.Equal(t, flow.StatusRunning, s.GetStatus())

	// Test Endpoint
	s.SetEndpoint("service.method")
	assert.Equal(t, "service.method", s.Endpoint())

	// Test Requires
	assert.Empty(t, s.Requires())

	// Test Require
	step2 := NewStep("step2")
	err := s.Require(step2)
	require.NoError(t, err)
	assert.Contains(t, s.Requires(), "step2")

	// Test Execute success
	req := &flow.Message{Body: []byte("test")}
	rsp, err := s.Execute(ctx, req)
	require.NoError(t, err)
	assert.Nil(t, rsp)
	assert.Equal(t, flow.StatusSuccess, s.GetStatus())
	assert.Equal(t, req, s.Request())

	// Test Execute error
	s2 := NewStep("test_step2")
	s2.SetExecuteError(assert.AnError)
	_, err = s2.Execute(ctx, &flow.Message{})
	assert.Error(t, err)
	assert.Equal(t, flow.StatusFailure, s2.GetStatus())

	// Test custom execute function
	s3 := NewStep("test_step3")
	expectedRsp := &flow.Message{Body: []byte("response")}
	s3.SetExecuteFunc(func(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) (*flow.Message, error) {
		return expectedRsp, nil
	})
	rsp, err = s3.Execute(ctx, &flow.Message{})
	require.NoError(t, err)
	assert.Equal(t, expectedRsp, rsp)
	assert.Equal(t, expectedRsp, s3.Response())

	// Test Compensate
	err = s.Compensate(ctx, &flow.Message{})
	require.NoError(t, err)

	// Test Compensate error
	s4 := NewStep("test_step4")
	s4.SetCompensateError(assert.AnError)
	err = s4.Compensate(ctx, &flow.Message{})
	assert.Error(t, err)

	// Test custom compensate function
	s5 := NewStep("test_step5")
	s5.SetCompensateFunc(func(ctx context.Context, req *flow.Message, opts ...flow.ExecuteOption) error {
		return assert.AnError
	})
	err = s5.Compensate(ctx, &flow.Message{})
	assert.Error(t, err)
}

func TestFlowExpectations(t *testing.T) {
	ctx := context.Background()
	f := NewFlow()

	// Test expectation with error
	f.ExpectWorkflowCreate("error_workflow").WillReturnError(assert.AnError)
	_, err := f.WorkflowCreate(ctx, "error_workflow")
	assert.Error(t, err)

	// Test expectation with custom result
	customWorkflow := NewWorkflow("custom")
	f.ExpectWorkflowCreate("custom_workflow").WithResult(customWorkflow)
	w, err := f.WorkflowCreate(ctx, "custom_workflow")
	require.NoError(t, err)
	assert.Equal(t, customWorkflow, w)

	// Test expectation not met
	f2 := NewFlow()
	f2.ExpectWorkflowCreate("not_called").Times(1)
	err = f2.ExpectationsWereMet()
	assert.Error(t, err)
}

func TestWorkflowRemoveSteps(t *testing.T) {
	w := NewWorkflow("test_workflow")

	step1 := NewStep("step1")
	step2 := NewStep("step2")

	err := w.AppendSteps(step1, step2)
	require.NoError(t, err)

	err = w.RemoveSteps(step1)
	require.NoError(t, err)
}
