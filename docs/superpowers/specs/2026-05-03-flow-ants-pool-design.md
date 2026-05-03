# Flow: ants/v2 Worker Pool for Workflow Execution

**Date:** 2026-05-03  
**Status:** Approved

## Summary

Add `github.com/panjf2000/ants/v2` goroutine pool to `go.unistack.org/micro/v4/flow` to bound the number of concurrently executing workflow steps. The pool is shared across all workflows on a single `Flow` instance (Flow-level, not Workflow-level). Steps are dispatched to the pool only when all their dependencies have completed (dispatch-on-ready), so pool slots are never occupied by waiting goroutines.

## Architecture

### Pool lifecycle

- Pool is created in `NewFlow`. If `Options.PoolSize == 0`, the default is `runtime.NumCPU() * 2`.
- Pool is stored on `microFlow.pool *ants.Pool` and shared by all concurrent `handleWorkflow` calls.
- `Init` recreates the pool if `PoolSize` changes (releases old pool first).
- `Close() error` is added to `microFlow` and to the `Flow` interface; it calls `pool.Release()`.

### Dispatch model — per-step readiness watcher

Current implementation spawns one goroutine per step with a busy-wait loop checking dependency completion. This is replaced with:

```
// one done-channel per step (closed when step completes)
doneChan := map[string]chan struct{}{}
for each stepID: doneChan[stepID] = make(chan struct{})

// one lightweight watcher goroutine per step (lives outside the pool)
for each stepID:
  go func watcher(id string):
    for _, depID := range step[id].Requires():
      <-doneChan[depID]          // block until dependency channel is closed
    pool.Submit(func():          // blocks until pool slot is available
      executeStep(id)
      close(doneChan[id])        // broadcast completion to dependent watchers
    )
```

Root steps (no dependencies) skip the dependency-wait and call `pool.Submit` immediately. The watcher goroutines are lightweight (blocked on channel receive, not spinning) and live outside the pool.

## Components

### `micro/flow/options.go`

- Add `PoolSize int` to `Options`.
- Add option constructor:
  ```go
  func PoolSize(n int) Option {
      return func(o *Options) { o.PoolSize = n }
  }
  ```

### `micro/flow/flow.go`

- Add `Close() error` to the `Flow` interface.

### `micro/flow/default.go`

- Add `pool *ants.Pool` field to `microFlow`.
- `NewFlow`: create pool after `NewOptions`; default size = `runtime.NumCPU() * 2`.
- `Init`: if `PoolSize` changed, release old pool and create new one.
- `microFlow.Close()`: call `pool.Release()`, return nil.
- `handleWorkflow`: replace busy-wait goroutines with watcher goroutines + `pool.Submit`.

### `micro/go.mod`

- Add `github.com/panjf2000/ants/v2`.

## Data flow

```
Execute(ctx, req) 
  → handleWorkflow(startID)
      → build doneChan map (one chan per vertex)
      → for each vertex: wg.Add(1), spawn watcher goroutine
          watcher: wait on dependency doneChans
          watcher: pool.Submit(func():       ← blocks if pool full, then returns
              defer wg.Done()
              executeStep(id)
              close(doneChan[id])            ← broadcast to dependent watchers
          )
          watcher exits (wg.Done is in pool task, not watcher)
      → wg.Wait()                            ← waits for all pool tasks to finish
      → collect errors from errChan
```

## Error handling

- If a step fails, it sends the error to `errChan` and closes its `doneChan`. Dependent watchers unblock and call `pool.Submit`, where `executeStep` detects the upstream failure via `errChan` (or a shared `atomic.Bool aborted` flag) and returns early without executing, then closes their own `doneChan`. This prevents downstream steps from hanging.
- Compensation logic (`failWorkflow`) is unchanged.
- If the pool itself returns an error from `Submit` (e.g. pool released), the watcher propagates it to `errChan`.

## Testing

- Existing tests in `flow/` must pass unchanged (pool with default size is transparent).
- Add a test with `PoolSize(1)` and a workflow with independent parallel steps: verify steps execute sequentially (pool serializes them) and all complete correctly.
- Add a test confirming `Close()` releases the pool without panic on subsequent workflow creation.
