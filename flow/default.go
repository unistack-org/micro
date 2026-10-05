package flow

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/heimdalr/dag"
	ants "github.com/panjf2000/ants/v2"
	codecpb "go.unistack.org/micro-proto/v5/codec"
	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/util/id"
)

// defaultStatusPollInterval is how often the engine checks the state store
// for cross-process abort/suspend signals of executions running in this
// process.
const defaultStatusPollInterval = 5 * time.Second

// microFlow implements Flow.
type microFlow struct {
	opts Options
	pool *ants.Pool

	mu         sync.Mutex
	executions map[string]*microExecution

	cancel context.CancelFunc
}

// microExecution is an execution currently running in this process.
type microExecution struct {
	eid    string
	cancel context.CancelFunc
}

type microWorkflow struct {
	f      *microFlow
	opts   Options
	g      *dag.DAG
	steps  map[string]Step
	id     string
	status Status
	sync.RWMutex
	init bool
}

func (w *microWorkflow) ID() string {
	return w.id
}

func (w *microWorkflow) Status() Status {
	return w.status
}

func (w *microWorkflow) AppendSteps(steps ...Step) error {
	var err error
	w.Lock()
	defer w.Unlock()

	for _, s := range steps {
		w.steps[s.String()] = s
		if _, err = w.g.AddVertex(s); err != nil {
			return err
		}
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			if err = w.g.AddEdge(src.String(), dst.String()); err != nil {
				return err
			}
		}
	}

	w.g.ReduceTransitively()

	return nil
}

func (w *microWorkflow) RemoveSteps(steps ...Step) error {
	// TODO: handle case when some step requires or required by removed step

	w.Lock()
	defer w.Unlock()

	for _, s := range steps {
		delete(w.steps, s.String())
		if err := w.g.DeleteVertex(s.String()); err != nil {
			return fmt.Errorf("failed to delete vertex %s: %w", s.String(), err)
		}
	}

	for _, dst := range steps {
		for _, req := range dst.Requires() {
			src, ok := w.steps[req]
			if !ok {
				return ErrStepNotExists
			}
			if err := w.g.AddEdge(src.String(), dst.String()); err != nil {
				return fmt.Errorf("failed to add edge %s -> %s: %w", src.String(), dst.String(), err)
			}
		}
	}

	w.g.ReduceTransitively()

	return nil
}

// Abort aborts execution eid: the status is written to the state store (so
// it is visible to every process) and, if the execution runs in this
// process, its context is canceled so the current step stops immediately.
func (w *microWorkflow) Abort(ctx context.Context, eid string) error {
	if err := w.f.stateStore().WorkflowSetStatus(ctx, eid, StatusAborted); err != nil {
		return err
	}
	w.f.cancelExecution(eid)
	return nil
}

// Suspend suspends execution eid (see Abort). A suspended execution can be
// resumed later with Resume from any process.
func (w *microWorkflow) Suspend(ctx context.Context, eid string) error {
	if err := w.f.stateStore().WorkflowSetStatus(ctx, eid, StatusSuspend); err != nil {
		return err
	}
	w.f.cancelExecution(eid)
	return nil
}

// Resume resumes execution eid from any process: the status is reset to
// Running and the DAG is re-executed, skipping steps that already completed
// successfully (detected via the state store).
func (w *microWorkflow) Resume(ctx context.Context, eid string) error {
	ss := w.f.stateStore()

	st, err := ss.WorkflowLoad(ctx, eid)
	if err != nil {
		return err
	}

	// Only non-terminal executions can be resumed.
	switch st.Status {
	case StatusRunning, StatusSuccess, StatusAborted:
		return nil
	}

	if err := ss.WorkflowSetStatus(ctx, eid, StatusRunning); err != nil {
		return err
	}

	go func() {
		if err := w.handleWorkflow(context.Background(), eid, nil); err != nil {
			w.opts.Logger.Error(context.Background(), "resume workflow execution failed", "error", err, "eid", eid)
		}
	}()

	return nil
}

func (w *microWorkflow) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (string, error) {
	w.Lock()
	if !w.init {
		w.g.ReduceTransitively()
		w.init = true
	}
	w.Unlock()

	eid, err := id.New()
	if err != nil {
		return "", err
	}

	options := NewExecuteOptions(opts...)

	nopts := make([]ExecuteOption, 0, len(opts)+5)

	nopts = append(nopts,
		ExecuteClient(w.opts.Client),
		ExecuteTracer(w.opts.Tracer),
		ExecuteLogger(w.opts.Logger),
		ExecuteMeter(w.opts.Meter),
	)
	nopts = append(nopts, opts...)

	// Persist the execution record. The record is keyed by the execution id;
	// WorkflowID references the workflow definition.
	if err := w.f.stateStore().WorkflowSave(ctx, eid, &WorkflowState{
		EID:        eid,
		WorkflowID: w.id,
		Status:     StatusRunning,
		Graph:      w.graph(),
		StartedAt:  time.Now(),
	}); err != nil {
		return eid, err
	}

	if options.Async {
		go func() {
			if err := w.handleWorkflow(ctx, eid, req, nopts...); err != nil {
				w.opts.Logger.Error(context.Background(), "async workflow execution failed", "error", err, "eid", eid)
			}
		}()
		return eid, nil
	}

	return eid, w.handleWorkflow(ctx, eid, req, nopts...)
}

// graph returns the step dependency graph: step id -> required step ids.
func (w *microWorkflow) graph() map[string][]string {
	w.RLock()
	defer w.RUnlock()
	g := make(map[string][]string, len(w.steps))
	for id, step := range w.steps {
		g[id] = step.Requires()
	}
	return g
}

// handleWorkflow executes the DAG of the given execution, persisting step and
// workflow state to the state store. Steps that already completed
// successfully (StatusSuccess) are skipped, which makes it safe to re-invoke
// the same execution id to resume a previous execution.
func (w *microWorkflow) handleWorkflow(ctx context.Context, eid string, req *Message, opts ...ExecuteOption) error {
	options := NewExecuteOptions(opts...)
	ss := w.f.stateStore()

	nopts := make([]ExecuteOption, 0, len(opts)+5)
	nopts = append(nopts,
		ExecuteClient(w.opts.Client),
		ExecuteTracer(w.opts.Tracer),
		ExecuteLogger(w.opts.Logger),
		ExecuteMeter(w.opts.Meter),
	)
	nopts = append(nopts, opts...)

	// Execution context: canceled on abort/suspend (in-process directly,
	// cross-process by the status poller) or when the base context is done.
	execCtx, cancel := context.WithCancel(options.Context)
	defer cancel()

	w.f.registerExecution(&microExecution{eid: eid, cancel: cancel})
	defer w.f.unregisterExecution(eid)

	// Load step states persisted by previous attempts of this execution so
	// completed steps can be skipped (resume semantics).
	prevSteps, err := ss.StepList(ctx, eid)
	if err != nil {
		return err
	}

	var executedSteps []string
	var execMu sync.Mutex

	// compensateSteps runs Compensate for all executed steps in reverse
	// order (Saga pattern).
	compensateSteps := func() {
		execMu.Lock()
		stepsToCompensate := make([]string, len(executedSteps))
		copy(stepsToCompensate, executedSteps)
		execMu.Unlock()

		for i := len(stepsToCompensate) - 1; i >= 0; i-- {
			w.compensateStep(ctx, eid, stepsToCompensate[i])
		}
	}

	stepResults := make(map[string]*Message)
	var resultsMu sync.RWMutex

	vertices := w.g.GetVertices()
	if len(vertices) == 0 {
		return w.finishWorkflow(ctx, eid, fmt.Errorf("no steps to execute"))
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

	// Skip steps that completed successfully in a previous attempt of this
	// execution: seed their results and mark them done.
	for stepID := range vertices {
		st := prevSteps[stepID]
		if st == nil || st.Status != StatusSuccess {
			continue
		}
		if step, ok := w.steps[stepID]; ok {
			step.SetStatus(StatusSuccess)
		}
		resultsMu.Lock()
		stepResults[stepID] = &Message{Body: st.Rsp}
		resultsMu.Unlock()
		execMu.Lock()
		executedSteps = append(executedSteps, stepID)
		execMu.Unlock()
		close(doneChan[stepID])
	}

	for stepID := range vertices {
		// Steps skipped above were already marked done and seeded, so they
		// need no execution task.
		if ps, ok := prevSteps[stepID]; ok && ps.Status == StatusSuccess {
			continue
		}
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

			// Cross-process abort/suspend check: the status may have been
			// updated by another process while we were waiting.
			if st, lerr := ss.WorkflowStatus(ctx, eid); lerr == nil &&
				(st == StatusAborted || st == StatusSuspend) {
				aborted.Store(true)
				close(doneChan[id])
				wg.Done()
				return
			}

			// Submit to pool. Blocks until a worker slot is available.
			if err := w.f.pool.Submit(func() {
				defer wg.Done()
				defer close(doneChan[id])

				if aborted.Load() {
					return
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
				} else if req != nil {
					// The workflow input message feeds the root steps.
					inputMsg = req
				}

				maxAttempts := 1
				backoff := time.Duration(0)
				if rp := step.Options().Retry; rp != nil {
					if rp.MaxAttempts > 1 {
						maxAttempts = rp.MaxAttempts
					}
					backoff = rp.Backoff
				}

				// Step checkpoint: seeded from the persisted state of a
				// previous attempt of this execution (resume) and shared
				// across retries in this run.
				cp := &stepCheckpoint{}
				if ps, ok := prevSteps[id]; ok {
					cp.last = append([]byte(nil), ps.Checkpoint...)
				}
				cp.save = func(ctx context.Context, data []byte) error {
					return ss.StepSetCheckpoint(ctx, eid, id, data)
				}

				var (
					attempt int
					rsp     *Message
					execErr error
				)

				for a := 0; a < maxAttempts; a++ {
					attempt = a + 1

					attemptCtx := context.WithValue(execCtx, stepCheckpointKey{}, cp)

					step.SetStatus(StatusRunning)
					if serr := ss.StepSave(ctx, eid, id, &StepState{
						EID: eid, StepID: id, Status: StatusRunning,
						Req: inputMsg.Body, Attempt: attempt,
						StartedAt: time.Now(),
					}); serr != nil {
						w.opts.Logger.Error(ctx, "failed to persist step status", "error", serr, "eid", eid, "step", id)
					}

					w.opts.Logger.Info(attemptCtx, "executing step: %s (attempt %d)", id, attempt)

					rsp, execErr = step.Execute(attemptCtx, inputMsg, nopts...)
					if execErr == nil {
						break
					}
					if execCtx.Err() != nil {
						// Execution aborted/suspended or the base context is
						// done: no point retrying.
						break
					}

					// Persist the failed attempt; the checkpoint (if any)
					// survives for the next attempt.
					_ = ss.StepSave(ctx, eid, id, &StepState{
						EID: eid, StepID: id, Status: StatusFailure,
						Error: execErr.Error(), Attempt: attempt,
						FinishedAt: time.Now(),
					})

					if a < maxAttempts-1 && backoff > 0 {
						select {
						case <-time.After(backoff << a):
						case <-execCtx.Done():
						}
					}
				}

				if execErr != nil {
					step.SetStatus(StatusFailure)
					_ = ss.StepSave(ctx, eid, id, &StepState{
						EID: eid, StepID: id, Status: StatusFailure,
						Error: execErr.Error(), Attempt: attempt,
						FinishedAt: time.Now(),
					})
					aborted.Store(true)
					errChan <- execErr
					return
				}

				step.SetStatus(StatusSuccess)
				if rsp != nil {
					_ = ss.StepSave(ctx, eid, id, &StepState{
						EID: eid, StepID: id, Status: StatusSuccess,
						Rsp: rsp.Body, Attempt: attempt,
						FinishedAt: time.Now(),
					})
					resultsMu.Lock()
					stepResults[id] = rsp
					resultsMu.Unlock()
				}

				execMu.Lock()
				executedSteps = append(executedSteps, id)
				execMu.Unlock()

				_ = ss.WorkflowSetLastStep(ctx, eid, id)

				w.opts.Logger.Info(ctx, "step completed: %s", id)
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

	var runErr error
	for err := range errChan {
		if runErr == nil {
			runErr = err
		}
	}

	if runErr != nil {
		compensateSteps()
	}

	return w.finishWorkflow(ctx, eid, runErr)
}

// compensateStep runs Compensate for a single step (Saga pattern), reading
// the step's persisted request from the state store. The step is reset to
// Pending so a resumed execution re-runs it.
func (w *microWorkflow) compensateStep(ctx context.Context, eid, stepID string) {
	step, ok := w.steps[stepID]
	if !ok {
		return
	}
	ss := w.f.stateStore()

	w.opts.Logger.Info(ctx, "compensating step: %s", stepID)
	st, err := ss.StepLoad(ctx, eid, stepID)
	if err != nil || st == nil {
		w.opts.Logger.Error(ctx, "failed to load step state for compensation", "error", err, "eid", eid, "step", stepID)
		return
	}

	req := &Message{Body: st.Req}
	if cerr := step.Compensate(ctx, req); cerr != nil {
		w.opts.Logger.Error(ctx, "compensation failed for step %s: %v", stepID, cerr)
		return
	}

	if werr := ss.StepSave(ctx, eid, stepID, &StepState{EID: eid, StepID: stepID, Status: StatusPending}); werr != nil {
		w.opts.Logger.Error(ctx, "failed to reset step status after compensation", "error", werr, "eid", eid, "step", stepID)
	}
}

// finishWorkflow writes the terminal status of the execution. A status
// already set by Abort/Suspend (in this or another process) is kept.
func (w *microWorkflow) finishWorkflow(ctx context.Context, eid string, runErr error) error {
	ss := w.f.stateStore()

	st, lerr := ss.WorkflowStatus(ctx, eid)
	if lerr != nil {
		st = StatusRunning
	}
	switch {
	case runErr != nil && st == StatusRunning:
		st = StatusFailure
	case st == StatusRunning:
		st = StatusSuccess
	}

	w.setStatus(st)

	if state, lerr := ss.WorkflowLoad(ctx, eid); lerr == nil && state != nil {
		state.Status = st
		state.FinishedAt = time.Now()
		if err := ss.WorkflowSave(ctx, eid, state); err != nil {
			return err
		}
	} else if lerr != nil {
		if err := ss.WorkflowSetStatus(ctx, eid, st); err != nil {
			return err
		}
	}

	if runErr != nil {
		w.opts.Logger.Error(ctx, "workflow failed: %v", runErr)
		return runErr
	}
	w.opts.Logger.Info(ctx, "workflow completed successfully")
	return nil
}

// setStatus updates the in-process status cache.
func (w *microWorkflow) setStatus(s Status) {
	w.Lock()
	w.status = s
	w.Unlock()
}

// NewFlow create new flow
func NewFlow(opts ...Option) Flow {
	options := NewOptions(opts...)
	size := options.PoolSize
	if size == 0 {
		size = runtime.NumCPU() * 2
	}
	p, _ := ants.NewPool(size)
	return &microFlow{
		opts:       options,
		pool:       p,
		executions: make(map[string]*microExecution),
	}
}

func (f *microFlow) Options() Options {
	return f.opts
}

// Close releases the goroutine pool and stops background goroutines.
func (f *microFlow) Close() error {
	if f.cancel != nil {
		f.cancel()
	}
	if f.pool != nil {
		f.pool.Release()
		f.pool = nil
	}
	return nil
}

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
	if f.opts.Store != nil {
		if err := f.opts.Store.Init(); err != nil {
			return err
		}
	}

	// State store for workflow execution state. Defaults to a KV adapter
	// over Store; when neither is set, a no-op store is used.
	if f.opts.StateStore == nil && f.opts.Store != nil {
		f.opts.StateStore = NewKVStateStore(f.opts.Store)
	}
	if f.opts.StateStore == nil {
		f.opts.StateStore = noopStateStore{}
	}
	if f.executions == nil {
		f.executions = make(map[string]*microExecution)
	}

	if f.cancel == nil {
		ctx, cancel := context.WithCancel(f.opts.Context)
		f.cancel = cancel

		go f.pollStatuses(ctx)
		if f.opts.CleanupInterval > 0 {
			go f.cleanup(ctx)
		}
	}

	return nil
}

// stateStore returns the state store used by this flow (never nil).
func (f *microFlow) stateStore() StateStore {
	if f.opts.StateStore != nil {
		return f.opts.StateStore
	}
	return noopStateStore{}
}

// registerExecution registers an execution running in this process.
func (f *microFlow) registerExecution(e *microExecution) {
	f.mu.Lock()
	f.executions[e.eid] = e
	f.mu.Unlock()
}

// unregisterExecution removes an execution from the in-process registry.
func (f *microFlow) unregisterExecution(eid string) {
	f.mu.Lock()
	delete(f.executions, eid)
	f.mu.Unlock()
}

// cancelExecution cancels the execution running in this process (if any), so
// the current step stops immediately.
func (f *microFlow) cancelExecution(eid string) {
	f.mu.Lock()
	e, ok := f.executions[eid]
	f.mu.Unlock()
	if ok {
		e.cancel()
	}
}

// pollStatuses periodically checks the state store for cross-process
// abort/suspend signals of executions running in this process and cancels
// their execution contexts.
func (f *microFlow) pollStatuses(ctx context.Context) {
	ticker := time.NewTicker(defaultStatusPollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.mu.Lock()
			execs := make([]*microExecution, 0, len(f.executions))
			for _, e := range f.executions {
				execs = append(execs, e)
			}
			f.mu.Unlock()

			for _, e := range execs {
				st, err := f.stateStore().WorkflowStatus(ctx, e.eid)
				if err != nil {
					continue
				}
				if st == StatusAborted || st == StatusSuspend {
					e.cancel()
				}
			}
		}
	}
}

// cleanup periodically deletes executions that have been in a terminal
// status for longer than opts.CleanupTTL.
func (f *microFlow) cleanup(ctx context.Context) {
	ticker := time.NewTicker(f.opts.CleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			entries, err := f.stateStore().WorkflowList(ctx, &WorkflowFilter{
				Statuses: []Status{StatusSuccess, StatusFailure, StatusAborted},
			})
			if err != nil {
				f.opts.Logger.Error(ctx, "cleanup failed to list workflows", "error", err)
				continue
			}
			deadline := time.Now().Add(-f.opts.CleanupTTL)
			for _, entry := range entries {
				if entry.State.FinishedAt.IsZero() || entry.State.FinishedAt.After(deadline) {
					continue
				}
				if err := f.stateStore().WorkflowDelete(ctx, entry.EID); err != nil {
					f.opts.Logger.Error(ctx, "cleanup failed to delete workflow", "error", err, "eid", entry.EID)
				}
			}
		}
	}
}

func (f *microFlow) WorkflowRemove(ctx context.Context, id string) error {
	return f.stateStore().WorkflowDelete(ctx, id)
}

func (f *microFlow) WorkflowCreate(ctx context.Context, id string, steps ...Step) (Workflow, error) {
	w := &microWorkflow{f: f, opts: f.opts, id: id, g: dag.NewDAG(), steps: make(map[string]Step, len(steps))}

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

func (f *microFlow) WorkflowSave(ctx context.Context, w Workflow) error {
	mw, ok := w.(*microWorkflow)
	if !ok {
		return fmt.Errorf("invalid workflow type")
	}

	// The workflow definition is stored as a record keyed by its id, with the
	// step dependency graph.
	if err := f.stateStore().WorkflowSave(ctx, mw.id, &WorkflowState{
		EID:        mw.id,
		WorkflowID: mw.id,
		Status:     mw.Status(),
		Graph:      mw.graph(),
	}); err != nil {
		return err
	}

	f.opts.Logger.Info(ctx, "workflow %s saved", mw.id)
	return nil
}

func (f *microFlow) WorkflowLoad(ctx context.Context, id string) (Workflow, error) {
	st, err := f.stateStore().WorkflowLoad(ctx, id)
	if err != nil {
		return nil, err
	}

	// Steps are Go objects (callables) and cannot be restored from the store;
	// they must be (re)registered via WorkflowCreate/AppendSteps. The record
	// carries the dependency graph and the status.
	w := &microWorkflow{
		f:      f,
		opts:   f.opts,
		id:     st.EID,
		g:      dag.NewDAG(),
		steps:  make(map[string]Step),
		status: st.Status,
		init:   true,
	}

	f.opts.Logger.Info(ctx, "workflow %s loaded with status %s", id, st.Status.String())
	return w, nil
}

func (f *microFlow) WorkflowList(ctx context.Context) ([]Workflow, error) {
	entries, err := f.stateStore().WorkflowList(ctx, nil)
	if err != nil {
		return nil, err
	}

	workflows := make([]Workflow, 0, len(entries))
	for _, entry := range entries {
		workflows = append(workflows, &microWorkflow{
			f:      f,
			opts:   f.opts,
			id:     entry.EID,
			g:      dag.NewDAG(),
			steps:  make(map[string]Step),
			status: entry.State.Status,
			init:   true,
		})
	}

	return workflows, nil
}

type microCallStep struct {
	rsp     *Message
	req     *Message
	service string
	method  string
	opts    StepOptions
	status  Status
}

func (s *microCallStep) Request() *Message {
	return s.req
}

func (s *microCallStep) Response() *Message {
	return s.rsp
}

func (s *microCallStep) ID() string {
	return s.String()
}

func (s *microCallStep) Options() StepOptions {
	return s.opts
}

func (s *microCallStep) Endpoint() string {
	return s.method
}

func (s *microCallStep) Requires() []string {
	return s.opts.Requires
}

func (s *microCallStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.opts.Requires = append(s.opts.Requires, step.String())
	}
	return nil
}

func (s *microCallStep) String() string {
	if s.opts.ID != "" {
		return s.opts.ID
	}
	return fmt.Sprintf("%s.%s", s.service, s.method)
}

func (s *microCallStep) Name() string {
	return s.String()
}

func (s *microCallStep) Hashcode() interface{} {
	return s.String()
}

func (s *microCallStep) GetStatus() Status {
	return s.status
}

func (s *microCallStep) SetStatus(status Status) {
	s.status = status
}

func (s *microCallStep) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error) {
	options := NewExecuteOptions(opts...)
	if options.Client == nil {
		return nil, ErrMissingClient
	}
	rsp := &codecpb.Frame{}
	copts := []client.CallOption{client.WithRetries(0)}
	if options.Timeout > 0 {
		copts = append(copts,
			client.WithRequestTimeout(options.Timeout),
			client.WithDialTimeout(options.Timeout))
	}
	nctx := metadata.NewOutgoingContext(ctx, req.Header)
	err := options.Client.Call(nctx, options.Client.NewRequest(s.service, s.method, &codecpb.Frame{Data: req.Body}), rsp, copts...)
	if err != nil {
		return nil, err
	}
	md, _ := metadata.FromOutgoingContext(nctx)
	return &Message{Header: md, Body: rsp.Data}, err
}

// Compensate performs rollback for this step (default implementation returns nil)
func (s *microCallStep) Compensate(ctx context.Context, req *Message, opts ...ExecuteOption) error {
	// Default implementation does nothing - override in custom steps if compensation is needed
	return nil
}

type microPublishStep struct {
	req    *Message
	rsp    *Message
	topic  string
	opts   StepOptions
	status Status
}

func (s *microPublishStep) Request() *Message {
	return s.req
}

func (s *microPublishStep) Response() *Message {
	return s.rsp
}

func (s *microPublishStep) ID() string {
	return s.String()
}

func (s *microPublishStep) Options() StepOptions {
	return s.opts
}

func (s *microPublishStep) Endpoint() string {
	return s.topic
}

func (s *microPublishStep) Requires() []string {
	return s.opts.Requires
}

func (s *microPublishStep) Require(steps ...Step) error {
	for _, step := range steps {
		s.opts.Requires = append(s.opts.Requires, step.String())
	}
	return nil
}

func (s *microPublishStep) String() string {
	if s.opts.ID != "" {
		return s.opts.ID
	}
	return s.topic
}

func (s *microPublishStep) Name() string {
	return s.String()
}

func (s *microPublishStep) Hashcode() interface{} {
	return s.String()
}

func (s *microPublishStep) GetStatus() Status {
	return s.status
}

func (s *microPublishStep) SetStatus(status Status) {
	s.status = status
}

func (s *microPublishStep) Execute(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error) {
	return nil, nil
}

// Compensate performs rollback for this step (default implementation returns nil)
func (s *microPublishStep) Compensate(ctx context.Context, req *Message, opts ...ExecuteOption) error {
	// Default implementation does nothing - override in custom steps if compensation is needed
	return nil
}

// NewCallStep create new step with client.Call
func NewCallStep(service string, name string, method string, opts ...StepOption) Step {
	options := NewStepOptions(opts...)
	return &microCallStep{service: service, method: name + "." + method, opts: options}
}

// NewPublishStep create new step with client.Publish
func NewPublishStep(topic string, opts ...StepOption) Step {
	options := NewStepOptions(opts...)
	return &microPublishStep{topic: topic, opts: options}
}
