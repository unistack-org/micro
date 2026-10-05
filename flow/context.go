package flow

import (
	"context"
	"sync"
)

type flowKey struct{}

var flowKeyVal = flowKey{}

// FromContext returns Flow from context
func FromContext(ctx context.Context) (Flow, bool) {
	if ctx == nil {
		return nil, false
	}
	c, ok := ctx.Value(flowKeyVal).(Flow)
	return c, ok
}

// NewContext stores Flow to context
func NewContext(ctx context.Context, f Flow) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, flowKeyVal, f)
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v any) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

type stepCheckpointKey struct{}

// stepCheckpoint carries the checkpoint state of a running step through the
// step context (heartbeat-style). The engine creates one per step attempt:
// it is seeded with the checkpoint persisted by a previous attempt, updated
// in-memory by StepCheckpoint and persisted to the state store.
type stepCheckpoint struct {
	mu   sync.Mutex
	last []byte
	save func(ctx context.Context, data []byte) error
}

func (c *stepCheckpoint) checkpoint(ctx context.Context, data []byte) error {
	c.mu.Lock()
	c.last = append([]byte(nil), data...)
	c.mu.Unlock()
	if c.save != nil {
		return c.save(ctx, data)
	}
	return nil
}

func (c *stepCheckpoint) lastValue() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}

// StepCheckpoint saves data as a checkpoint of the currently running step.
// The engine persists the checkpoint to the state store, so a retried or
// restarted step can continue from the last saved point (similar to
// activity heartbeats). It is a no-op when the step does not run under the
// flow engine.
//
// Example:
//
//	func (s *myStep) Execute(ctx context.Context, req *Message, opts ...flow.ExecuteOption) (*flow.Message, error) {
//		// continue from the last checkpoint on retry or resume
//		if saved := flow.StepLastCheckpoint(ctx); len(saved) > 0 {
//			// ...
//		}
//		for i, item := range items {
//			if err := process(item); err != nil {
//				return nil, err
//			}
//			_ = flow.StepCheckpoint(ctx, []byte(strconv.Itoa(i)))
//		}
//		// ...
//	}
func StepCheckpoint(ctx context.Context, data []byte) error {
	if cp, ok := ctx.Value(stepCheckpointKey{}).(*stepCheckpoint); ok {
		return cp.checkpoint(ctx, data)
	}
	return nil
}

// StepLastCheckpoint returns the last checkpoint of the currently running
// step: the value saved by StepCheckpoint during the current attempt, or
// the checkpoint persisted by a previous attempt (seeded by the engine on
// retry or resume). It returns nil when the step does not run under the
// flow engine or no checkpoint was saved.
func StepLastCheckpoint(ctx context.Context) []byte {
	if cp, ok := ctx.Value(stepCheckpointKey{}).(*stepCheckpoint); ok {
		return cp.lastValue()
	}
	return nil
}
