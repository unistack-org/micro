# Flow ants/v2 Worker Pool Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a shared `ants/v2` goroutine pool to the `flow` package so that workflow steps are executed through a bounded pool, dispatched only when all their dependencies are met (dispatch-on-ready).

**Architecture:** A single `*ants.Pool` lives on `microFlow` and is passed to each `microWorkflow` on creation. `handleWorkflow` is rewritten: instead of one goroutine-per-step with a busy-wait loop, each step gets a lightweight watcher goroutine (outside the pool) that blocks on dependency `chan struct{}` channels, then calls `pool.Submit` when ready. `wg.Done()` and `close(doneChan[id])` are called inside the pool task via `defer`, ensuring no pool slot is held by a waiting goroutine.

**Tech Stack:** Go 1.25, `github.com/panjf2000/ants/v2`, `go.unistack.org/micro/v4/store/memory` (tests), `sync/atomic.Bool`.

---

## File Map

| File | Change |
|------|--------|
| `micro/go.mod` | Add `github.com/panjf2000/ants/v2` |
| `micro/flow/options.go` | Add `PoolSize int` to `Options`, add `PoolSize(int) Option` func |
| `micro/flow/flow.go` | Add `Close() error` to `Flow` interface |
| `micro/flow/mock/mock.go` | Add `Close() error` no-op to `*Flow` mock |
| `micro/flow/default.go` | Add `pool *ants.Pool` to both `microFlow` and `microWorkflow`; update `NewFlow`, `Init`, `WorkflowCreate`, `WorkflowLoad`; add `Close()`; replace `handleWorkflow` |
| `micro/flow/pool_test.go` | New file: pool size, close, and concurrency tests |

---

## Task 1: Add ants/v2 dependency

**Files:**
- Modify: `micro/go.mod`, `micro/go.sum`

- [ ] **Step 1: Add ants/v2 to the module**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go get github.com/panjf2000/ants/v2
```

Expected: lines like `go: added github.com/panjf2000/ants/v2 v2.x.x` in output, no error.

- [ ] **Step 2: Verify go.mod contains ants**

```bash
grep ants /Users/vtolstov/Projects/unistack/micro/micro/go.mod
```

Expected: `github.com/panjf2000/ants/v2 v2.x.x`

- [ ] **Step 3: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add go.mod go.sum
git commit -m "chore(flow): add ants/v2 goroutine pool dependency"
```

---

## Task 2: Add PoolSize option

**Files:**
- Modify: `micro/flow/options.go`
- Test: `micro/flow/pool_test.go`

- [ ] **Step 1: Create pool_test.go with failing option test**

Create `/Users/vtolstov/Projects/unistack/micro/micro/flow/pool_test.go`:

```go
package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoolSizeOption(t *testing.T) {
	opts := NewOptions(PoolSize(8))
	assert.Equal(t, 8, opts.PoolSize)
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run TestPoolSizeOption -v
```

Expected: compile error — `PoolSize` undefined.

- [ ] **Step 3: Add PoolSize field and option func to options.go**

In `micro/flow/options.go`, add `PoolSize int` to the `Options` struct after the `Store` field:

```go
// Options server struct
type Options struct {
	// Context holds the external options and can be used for flow shutdown
	Context context.Context
	// Client holds the client.Client
	Client client.Client
	// Tracer holds the tracer
	Tracer tracer.Tracer
	// Logger holds the logger
	Logger logger.Logger
	// Meter holds the meter
	Meter meter.Meter
	// Store used for intermediate results
	Store store.Store
	// PoolSize is the maximum number of concurrently executing workflow steps (0 = default: runtime.NumCPU()*2)
	PoolSize int
}
```

Then add the option func after the existing `Store` func:

```go
// PoolSize sets the maximum number of concurrently executing workflow steps.
// A value of 0 uses the default: runtime.NumCPU() * 2.
func PoolSize(n int) Option {
	return func(o *Options) {
		o.PoolSize = n
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run TestPoolSizeOption -v
```

Expected: `PASS`.

- [ ] **Step 5: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add flow/options.go flow/pool_test.go
git commit -m "feat(flow): add PoolSize option for goroutine pool configuration"
```

---

## Task 3: Add Close() to Flow interface and mock

**Files:**
- Modify: `micro/flow/flow.go`
- Modify: `micro/flow/mock/mock.go`

- [ ] **Step 1: Add Close() to the Flow interface in flow.go**

In `micro/flow/flow.go`, add `Close() error` to the `Flow` interface after `WorkflowList`:

```go
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
	// Close releases resources held by the flow (goroutine pool)
	Close() error
}
```

- [ ] **Step 2: Verify the build breaks (mock doesn't implement Close yet)**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go build ./flow/... 2>&1 | head -20
```

Expected: compile error about `mock.Flow` not implementing `flow.Flow` (missing `Close` method).

- [ ] **Step 3: Add Close() no-op to mock/mock.go**

In `micro/flow/mock/mock.go`, add after the `Init` method:

```go
// Close releases resources (no-op for mock)
func (m *Flow) Close() error {
	return nil
}
```

- [ ] **Step 4: Verify build passes**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go build ./flow/...
```

Expected: no output (success).

- [ ] **Step 5: Run existing tests to verify nothing broke**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/... -v 2>&1 | tail -20
```

Expected: all existing tests PASS.

- [ ] **Step 6: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add flow/flow.go flow/mock/mock.go
git commit -m "feat(flow): add Close() to Flow interface and mock implementation"
```

---

## Task 4: Add pool fields and update struct constructors

**Files:**
- Modify: `micro/flow/default.go`
- Modify: `micro/flow/pool_test.go`

- [ ] **Step 1: Add failing test for pool creation**

Append to `micro/flow/pool_test.go`:

```go
func TestNewFlowCreatesPool(t *testing.T) {
	f := NewFlow(PoolSize(4))
	mf, ok := f.(*microFlow)
	require.True(t, ok)
	assert.NotNil(t, mf.pool)
	assert.Equal(t, 4, mf.pool.Cap())
	require.NoError(t, f.Close())
}

func TestNewFlowDefaultPoolSize(t *testing.T) {
	import_runtime_numcpu_placeholder := 1 // compile will fail - see step 2
	_ = import_runtime_numcpu_placeholder
}
```

Actually skip the placeholder test — write only what compiles once pool is added. Add this import block at the top of pool_test.go instead:

Replace the full content of `micro/flow/pool_test.go` with:

```go
package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
```

- [ ] **Step 2: Run tests to verify they fail (pool not on microFlow yet)**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run "TestNewFlowCreatesPool" -v
```

Expected: compile error — `mf.pool undefined`.

- [ ] **Step 3: Add pool to microFlow and microWorkflow structs, update NewFlow, WorkflowCreate, WorkflowLoad**

In `micro/flow/default.go`, update the imports to add `ants` and `runtime` (do NOT add `sync/atomic` yet — it is only needed in Task 7 when `handleWorkflow` is replaced):

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	ants "github.com/panjf2000/ants/v2"
	"github.com/heimdalr/dag"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/util/id"
)
```

Update `microFlow` struct to add `pool`:

```go
type microFlow struct {
	opts Options
	pool *ants.Pool
}
```

Update `microWorkflow` struct to add `pool`:

```go
type microWorkflow struct {
	opts   Options
	g      *dag.DAG
	steps  map[string]Step
	id     string
	status Status
	sync.RWMutex
	init       bool
	cancelFunc context.CancelFunc
	execCtx    context.Context
	pool       *ants.Pool
}
```

Replace `NewFlow`:

```go
// NewFlow create new flow
func NewFlow(opts ...Option) Flow {
	options := NewOptions(opts...)
	size := options.PoolSize
	if size == 0 {
		size = runtime.NumCPU() * 2
	}
	p, _ := ants.NewPool(size)
	return &microFlow{opts: options, pool: p}
}
```

Update `WorkflowCreate` to pass pool to workflow:

```go
func (f *microFlow) WorkflowCreate(ctx context.Context, id string, steps ...Step) (Workflow, error) {
	w := &microWorkflow{opts: f.opts, pool: f.pool, id: id, g: &dag.DAG{}, steps: make(map[string]Step, len(steps))}

	for _, s := range steps {
		w.steps[s.String()] = s
		if _, err := w.g.AddVertex(s); err != nil {
			return nil, fmt.Errorf("failed to add vertex %s: %w", s.String(), err)
		}
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return nil, ErrStepNotExists
			}
			if err := w.g.AddEdge(src.String(), dst.String()); err != nil {
				return nil, fmt.Errorf("failed to add edge %s -> %s: %w", src.String(), dst.String(), err)
			}
		}
	}

	w.g.ReduceTransitively()
	w.init = true

	return w, nil
}
```

Update `WorkflowLoad` to pass pool:

```go
func (f *microFlow) WorkflowLoad(ctx context.Context, id string) (Workflow, error) {
	workflowStore := store.NewNamespaceStore(f.opts.Store, filepath.Join("workflows", id))

	statusFrame := &codec.Frame{}
	if err := workflowStore.Read(ctx, "status", statusFrame); err != nil {
		return nil, err
	}

	status := StringStatus[string(statusFrame.Data)]

	w := &microWorkflow{
		opts:   f.opts,
		pool:   f.pool,
		id:     id,
		g:      &dag.DAG{},
		steps:  make(map[string]Step),
		status: status,
	}

	stepFrame := &codec.Frame{}
	if err := workflowStore.Read(ctx, "steps", stepFrame); err != nil {
		f.opts.Logger.Warn(ctx, "failed to read steps for workflow %s: %v", id, err)
		return nil, err
	}

	var stepsData map[string][]string
	if err := json.Unmarshal(stepFrame.Data, &stepsData); err != nil {
		f.opts.Logger.Warn(ctx, "failed to unmarshal steps for workflow %s: %v", id, err)
		return nil, err
	}

	w.init = true

	f.opts.Logger.Info(ctx, "workflow %s loaded with status %s", id, status)
	return w, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run "TestNewFlowCreatesPool|TestPoolSizeOption" -v
```

Expected: both PASS.

- [ ] **Step 5: Run all flow tests**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/... -v 2>&1 | tail -30
```

Expected: all tests PASS.

- [ ] **Step 6: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add flow/default.go flow/pool_test.go
git commit -m "feat(flow): add ants pool to microFlow and microWorkflow structs"
```

---

## Task 5: Implement microFlow.Close() and update Init

**Files:**
- Modify: `micro/flow/default.go`
- Modify: `micro/flow/pool_test.go`

- [ ] **Step 1: Add failing Close test**

Append to `micro/flow/pool_test.go`:

```go
func TestFlowClose(t *testing.T) {
	f := NewFlow(PoolSize(2))
	require.NoError(t, f.Close())
	// Second close must not panic
	require.NoError(t, f.Close())
}
```

- [ ] **Step 2: Run to verify it fails**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run TestFlowClose -v
```

Expected: compile error — `*microFlow` doesn't implement `Flow` (missing `Close`).

- [ ] **Step 3: Add Close() to microFlow in default.go**

Add after the existing `Init` method:

```go
// Close releases the goroutine pool.
func (f *microFlow) Close() error {
	if f.pool != nil {
		f.pool.Release()
		f.pool = nil
	}
	return nil
}
```

- [ ] **Step 4: Update Init to recreate pool when PoolSize changes**

Replace the existing `Init` method:

```go
func (f *microFlow) Init(opts ...Option) error {
	for _, o := range opts {
		o(&f.opts)
	}

	if f.pool != nil {
		f.pool.Release()
	}
	size := f.opts.PoolSize
	if size == 0 {
		size = runtime.NumCPU() * 2
	}
	var err error
	f.pool, err = ants.NewPool(size)
	if err != nil {
		return err
	}

	if err := f.opts.Client.Init(); err != nil {
		return err
	}
	if err := f.opts.Tracer.Init(); err != nil {
		return err
	}
	if err := f.opts.Logger.Init(); err != nil {
		return err
	}
	if err := f.opts.Meter.Init(); err != nil {
		return err
	}
	if err := f.opts.Store.Init(); err != nil {
		return err
	}
	return nil
}
```

- [ ] **Step 5: Run Close test**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run TestFlowClose -v
```

Expected: PASS.

- [ ] **Step 6: Run all tests**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/... -v 2>&1 | tail -20
```

Expected: all tests PASS.

- [ ] **Step 7: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add flow/default.go flow/pool_test.go
git commit -m "feat(flow): implement Flow.Close() and pool recreation in Init"
```

---

## Task 6: Write pool concurrency test (TDD — before replacing handleWorkflow)

**Files:**
- Modify: `micro/flow/pool_test.go`

- [ ] **Step 1: Add imports and poolTestStep helper to pool_test.go**

Replace the full content of `micro/flow/pool_test.go` with:

```go
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

func (s *poolTestStep) ID() string          { return s.name }
func (s *poolTestStep) Endpoint() string    { return s.name }
func (s *poolTestStep) String() string      { return s.name }
func (s *poolTestStep) Hashcode() interface{} { return s.name }
func (s *poolTestStep) Requires() []string  { return s.requires }
func (s *poolTestStep) Options() StepOptions { return StepOptions{} }
func (s *poolTestStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.requires = append(s.requires, step.String())
	}
	return nil
}
func (s *poolTestStep) GetStatus() Status { s.mu.Lock(); defer s.mu.Unlock(); return s.status }
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
	require.NoError(t, f.Close())
}

// TestWorkflowPoolConcurrency verifies that PoolSize(1) serializes parallel steps.
func TestWorkflowPoolConcurrency(t *testing.T) {
	ms := memory.NewStore()
	ctx := context.Background()
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

	// Three independent steps (no dependencies) — should run one at a time with PoolSize(1).
	step1 := makeStep("step1")
	step2 := makeStep("step2")
	step3 := makeStep("step3")

	wf, err := f.WorkflowCreate(ctx, "test-pool-concurrency", step1, step2, step3)
	require.NoError(t, err)

	_, err = wf.Execute(ctx, &Message{Body: []byte(`{}`)})
	require.NoError(t, err)

	assert.Equal(t, int32(1), atomic.LoadInt32(&maxConcurrent),
		"PoolSize(1) must serialize step execution: maxConcurrent should be 1")
}
```

- [ ] **Step 2: Run TestWorkflowPoolConcurrency to verify it fails (old handleWorkflow uses raw goroutines)**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run TestWorkflowPoolConcurrency -v -count=1
```

Expected: test runs but `maxConcurrent` is likely > 1 because current `handleWorkflow` spawns raw goroutines not limited by the pool. Test FAILS with assertion error like `expected 1, got 3`.

- [ ] **Step 3: Commit the test (red state)**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add flow/pool_test.go
git commit -m "test(flow): add pool concurrency test (red — handleWorkflow not yet using pool)"
```

---

## Task 7: Replace handleWorkflow with dispatch-on-ready + pool

**Files:**
- Modify: `micro/flow/default.go`

- [ ] **Step 1: Add `"sync/atomic"` to imports in default.go**

In `micro/flow/default.go`, add `"sync/atomic"` to the import block (after `"sync"`):

```go
import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	ants "github.com/panjf2000/ants/v2"
	"github.com/heimdalr/dag"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/util/id"
)
```

- [ ] **Step 2: Replace handleWorkflow**

In `micro/flow/default.go`, replace the entire `handleWorkflow` method with:

```go
func (w *microWorkflow) handleWorkflow(startID string, opts ...ExecuteOption) error {
	options := NewExecuteOptions(opts...)

	eid := w.id
	workflowStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("workflows", eid))
	stepStore := store.NewNamespaceStore(w.opts.Store, filepath.Join("steps", eid))

	w.Lock()
	if w.execCtx == nil || w.cancelFunc == nil {
		execCtx, cancel := context.WithCancel(context.Background())
		w.execCtx = execCtx
		w.cancelFunc = cancel
	}
	execCtx := w.execCtx
	w.Unlock()

	var executedSteps []string
	var execMu sync.Mutex

	failWorkflow := func(err error) error {
		w.opts.Logger.Error(options.Context, "workflow failed: %v", err)

		if werr := workflowStore.Write(options.Context, "status", &codec.Frame{Data: []byte(StatusFailure.String())}); werr != nil {
			w.opts.Logger.Error(options.Context, "store error: %v", werr)
		}

		execMu.Lock()
		stepsToCompensate := make([]string, len(executedSteps))
		copy(stepsToCompensate, executedSteps)
		execMu.Unlock()

		for i := len(stepsToCompensate) - 1; i >= 0; i-- {
			stepID := stepsToCompensate[i]
			step, ok := w.steps[stepID]
			if !ok {
				continue
			}
			w.opts.Logger.Info(options.Context, "compensating step: %s", stepID)
			reqFrame := &codec.Frame{}
			if rerr := stepStore.Read(options.Context, filepath.Join(stepID, "req"), reqFrame); rerr != nil {
				w.opts.Logger.Error(options.Context, "failed to read request for compensation: %v", rerr)
				continue
			}
			req := &Message{Body: reqFrame.Data}
			if cerr := step.Compensate(options.Context, req, opts...); cerr != nil {
				w.opts.Logger.Error(options.Context, "compensation failed for step %s: %v", stepID, cerr)
			} else {
				if werr := stepStore.Write(options.Context, filepath.Join(stepID, "status"), &codec.Frame{Data: []byte(StatusPending.String())}); werr != nil {
					w.opts.Logger.Error(options.Context, "store error: %v", werr)
				}
			}
		}

		return err
	}

	stepResults := make(map[string]*Message)
	var resultsMu sync.RWMutex

	vertices := w.g.GetVertices()
	if len(vertices) == 0 {
		return failWorkflow(fmt.Errorf("no steps to execute"))
	}

	// One done-channel per step; closed when the step finishes (success or failure).
	// Downstream watcher goroutines block on these channels outside the pool.
	doneChan := make(map[string]chan struct{}, len(vertices))
	for stepID := range vertices {
		doneChan[stepID] = make(chan struct{})
	}

	errChan := make(chan error, len(vertices))
	var wg sync.WaitGroup
	var aborted atomic.Bool

	for stepID := range vertices {
		wg.Add(1)
		go func(id string) {
			step, ok := w.steps[id]
			if !ok {
				errChan <- ErrStepNotExists
				close(doneChan[id])
				wg.Done()
				return
			}

			// Wait for all dependency channels to close (outside the pool).
			for _, depID := range step.Requires() {
				select {
				case <-doneChan[depID]:
				case <-execCtx.Done():
					errChan <- execCtx.Err()
					close(doneChan[id])
					wg.Done()
					return
				}
			}

			// Fast-exit if a prior step already failed.
			if aborted.Load() {
				close(doneChan[id])
				wg.Done()
				return
			}

			// Submit to pool. Blocks until a worker slot is available.
			if err := w.pool.Submit(func() {
				defer wg.Done()
				defer close(doneChan[id])

				if aborted.Load() {
					return
				}

				// Resume support: skip steps that already completed successfully.
				stepStatusFrame := &codec.Frame{}
				if rerr := stepStore.Read(execCtx, filepath.Join(id, "status"), stepStatusFrame); rerr == nil {
					if StringStatus[string(stepStatusFrame.Data)] == StatusSuccess {
						rspFrame := &codec.Frame{}
						if rrerr := stepStore.Read(execCtx, filepath.Join(id, "rsp"), rspFrame); rrerr == nil && len(rspFrame.Data) > 0 {
							resultsMu.Lock()
							stepResults[id] = &Message{Body: rspFrame.Data}
							resultsMu.Unlock()
						}
						return
					}
				}

				// Collect output from dependency steps as input.
				inputMsg := &Message{Body: []byte{}, Header: metadata.Metadata{}}
				requires := step.Requires()
				if len(requires) > 0 {
					resultsMu.RLock()
					for _, reqID := range requires {
						if res, exists := stepResults[reqID]; exists && len(res.Body) > 0 {
							inputMsg.Body = res.Body
							inputMsg.Header = res.Header
						}
					}
					resultsMu.RUnlock()
				}

				if werr := stepStore.Write(execCtx, filepath.Join(id, "req"), &codec.Frame{Data: inputMsg.Body}); werr != nil {
					aborted.Store(true)
					errChan <- werr
					return
				}

				step.SetStatus(StatusRunning)
				if werr := stepStore.Write(execCtx, filepath.Join(id, "status"), &codec.Frame{Data: []byte(StatusRunning.String())}); werr != nil {
					aborted.Store(true)
					errChan <- werr
					return
				}

				w.opts.Logger.Info(execCtx, "executing step: %s", id)

				rsp, execErr := step.Execute(execCtx, inputMsg, opts...)
				if execErr != nil {
					step.SetStatus(StatusFailure)
					_ = stepStore.Write(execCtx, filepath.Join(id, "rsp"), &codec.Frame{Data: []byte(execErr.Error())})
					_ = stepStore.Write(execCtx, filepath.Join(id, "status"), &codec.Frame{Data: []byte(StatusFailure.String())})
					aborted.Store(true)
					errChan <- execErr
					return
				}

				step.SetStatus(StatusSuccess)
				if rsp != nil {
					if werr := stepStore.Write(execCtx, filepath.Join(id, "rsp"), &codec.Frame{Data: rsp.Body}); werr != nil {
						aborted.Store(true)
						errChan <- werr
						return
					}
					resultsMu.Lock()
					stepResults[id] = rsp
					resultsMu.Unlock()
				}

				if werr := stepStore.Write(execCtx, filepath.Join(id, "status"), &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
					aborted.Store(true)
					errChan <- werr
					return
				}

				_ = workflowStore.Write(execCtx, "last_step", &codec.Frame{Data: []byte(id)})

				execMu.Lock()
				executedSteps = append(executedSteps, id)
				execMu.Unlock()

				w.opts.Logger.Info(execCtx, "step completed: %s", id)
			}); err != nil {
				// Pool was released or encountered an error.
				errChan <- err
				close(doneChan[id])
				wg.Done()
			}
		}(stepID)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return failWorkflow(err)
		}
	}

	if werr := workflowStore.Write(options.Context, "status", &codec.Frame{Data: []byte(StatusSuccess.String())}); werr != nil {
		w.opts.Logger.Error(options.Context, "store error: %v", werr)
	}

	w.opts.Logger.Info(options.Context, "workflow completed successfully")
	return nil
}
```

- [ ] **Step 3: Run TestWorkflowPoolConcurrency (should now PASS)**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/ -run TestWorkflowPoolConcurrency -v -count=1
```

Expected: PASS — `maxConcurrent == 1`.

- [ ] **Step 4: Run all flow tests**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/... -v -count=1 2>&1 | tail -40
```

Expected: all tests PASS (including existing dag tests and mock tests).

- [ ] **Step 5: Run go vet**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go vet ./flow/...
```

Expected: no output (no issues).

- [ ] **Step 6: Commit**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && git add flow/default.go
git commit -m "feat(flow): replace handleWorkflow busy-wait with ants pool dispatch-on-ready"
```

---

## Task 8: Final verification

- [ ] **Step 1: Run full module test suite**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./... 2>&1 | grep -E "FAIL|ok" | head -40
```

Expected: all packages show `ok`, no `FAIL`.

- [ ] **Step 2: Run race detector on flow package**

```bash
cd /Users/vtolstov/Projects/unistack/micro/micro && go test ./flow/... -race -count=3 -timeout 60s
```

Expected: all tests PASS with no data race warnings.

- [ ] **Step 3: Final commit if any cleanup needed**

If the previous steps produced no additional changes, no commit is needed here.
