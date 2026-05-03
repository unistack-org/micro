package flow

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	memory "go.unistack.org/micro/v4/store/memory"
)

// poolTestStep is a Step implementation for pool concurrency tests.
type poolTestStep struct {
	name          string
	concurrent    *int32
	maxConcurrent *int32
	requires      []string
	status        Status
	mu            sync.Mutex
}

func (s *poolTestStep) ID() string               { return s.name }
func (s *poolTestStep) Endpoint() string         { return s.name }
func (s *poolTestStep) String() string           { return s.name }
func (s *poolTestStep) Hashcode() interface{}    { return s.name }
func (s *poolTestStep) Requires() []string       { return s.requires }
func (s *poolTestStep) Options() StepOptions     { return StepOptions{} }
func (s *poolTestStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.requires = append(s.requires, step.String())
	}
	return nil
}
func (s *poolTestStep) GetStatus() Status  { s.mu.Lock(); defer s.mu.Unlock(); return s.status }
func (s *poolTestStep) SetStatus(st Status) { s.mu.Lock(); defer s.mu.Unlock(); s.status = st }
func (s *poolTestStep) Request() *Message  { return nil }
func (s *poolTestStep) Response() *Message { return nil }
func (s *poolTestStep) Compensate(_ context.Context, _ *Message, _ ...ExecuteOption) error {
	return nil
}
func (s *poolTestStep) Execute(_ context.Context, req *Message, _ ...ExecuteOption) (*Message, error) {
	c := atomic.AddInt32(s.concurrent, 1)
	for {
		cur := atomic.LoadInt32(s.maxConcurrent)
		if c <= cur {
			break
		}
		if atomic.CompareAndSwapInt32(s.maxConcurrent, cur, c) {
			break
		}
	}
	time.Sleep(30 * time.Millisecond)
	atomic.AddInt32(s.concurrent, -1)
	return &Message{Body: []byte(`{}`)}, nil
}

func TestPoolSizeOption(t *testing.T) {
	opts := NewOptions(PoolSize(8))
	assert.Equal(t, 8, opts.PoolSize)
}

func TestNewFlowCreatesPool(t *testing.T) {
	f := NewFlow(PoolSize(4))
	mf, ok := f.(*microFlow)
	require.True(t, ok)
	assert.NotNil(t, mf.pool)
	assert.Equal(t, 4, mf.pool.Cap())
	require.NoError(t, f.Close())
}

func TestFlowClose(t *testing.T) {
	f := NewFlow(PoolSize(2))
	require.NoError(t, f.Close())
	// Second close must not panic
	require.NoError(t, f.Close())
}

// TestWorkflowPoolConcurrency verifies that PoolSize(1) serializes parallel steps.
func TestWorkflowPoolConcurrency(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
	require.NoError(t, ms.Init())
	require.NoError(t, ms.Connect(ctx))

	var concurrent int32
	var maxConcurrent int32

	makeStep := func(name string) *poolTestStep {
		return &poolTestStep{
			name:          name,
			concurrent:    &concurrent,
			maxConcurrent: &maxConcurrent,
		}
	}

	f := NewFlow(PoolSize(1), Store(ms))
	defer f.Close()

	// Diamond: start → (step_a, step_b) → end.
	// With PoolSize(1), step_a and step_b must run serially — maxConcurrent stays 1.
	start := makeStep("start")
	stepA := makeStep("step_a")
	stepA.requires = []string{"start"}
	stepB := makeStep("step_b")
	stepB.requires = []string{"start"}
	end := makeStep("end")
	end.requires = []string{"step_a", "step_b"}

	wf, err := f.WorkflowCreate(ctx, "test-pool-concurrency", start, stepA, stepB, end)
	require.NoError(t, err)

	_, err = wf.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.NoError(t, err)

	assert.Equal(t, int32(1), atomic.LoadInt32(&maxConcurrent),
		"PoolSize(1) must serialize step execution: maxConcurrent should be 1")
}
