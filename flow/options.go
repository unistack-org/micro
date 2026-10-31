package flow

import (
	"context"
	"time"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/store"
	"go.unistack.org/micro/v5/tracer"
)

// Option func
type Option func(*Options)

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
	// StateStore stores workflow execution state.
	// When Store is set and StateStore is nil, it defaults to a KV
	// adapter (NewKVStateStore) over Store.
	StateStore StateStore
	// CleanupInterval is the interval between cleanup sweeps of finished
	// executions. Zero disables cleanup.
	CleanupInterval time.Duration
	// CleanupTTL is how long a finished execution is kept before deletion.
	CleanupTTL time.Duration
	// PollInterval is the interval between state store polls for
	// cross-process abort/suspend signals of executions running in this
	// process. Zero uses DefaultPollInterval; negative disables polling
	// (useful when the state store pushes status changes via
	// WorkflowWatcher).
	PollInterval time.Duration
	// PoolSize is the maximum number of concurrently executing workflow steps (0 = default: runtime.NumCPU()*2)
	PoolSize int
	// Hooks holds optional lifecycle callbacks invoked by the engine.
	Hooks *Hooks
}

// Hooks are optional lifecycle callbacks invoked by the flow engine.
// They mirror the lifecycle events Cloudflare Workflows emits for rollbacks.
// All callbacks are optional; nil callbacks are skipped.
type Hooks struct {
	// OnCompensateStart is called before a step's compensation begins.
	OnCompensateStart func(eid, stepID string)
	// OnCompensateSuccess is called after a step's compensation succeeds.
	OnCompensateSuccess func(eid, stepID string)
	// OnCompensateFailure is called after a step's compensation fails
	// (after its configured retries).
	OnCompensateFailure func(eid, stepID string, err error)
	// OnCompensateComplete is called when the compensation phase finishes.
	// failed is true if any step's compensation failed.
	OnCompensateComplete func(eid string, failed bool)
}

// NewOptions returns new options struct with default or passed values
func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
		Client:  client.DefaultClient,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Logger sets the logger option
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter sets the meter option
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Client to use for sync/async communication
func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the flow
// or can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Tracer mechanism for distributed tracking
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Store used for intermediate results
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// StateStorage sets a custom state store for workflow execution state.
func StateStorage(ss StateStore) Option {
	return func(o *Options) {
		o.StateStore = ss
	}
}

// Cleanup enables periodic cleanup of finished executions: every interval,
// executions that have been finished for longer than ttl are deleted.
func Cleanup(interval, ttl time.Duration) Option {
	return func(o *Options) {
		o.CleanupInterval = interval
		o.CleanupTTL = ttl
	}
}

// PollInterval sets the interval between state store polls for
// cross-process abort/suspend signals. Zero uses DefaultPollInterval;
// negative disables polling (rely on WorkflowWatcher push instead).
func PollInterval(d time.Duration) Option {
	return func(o *Options) {
		o.PollInterval = d
	}
}

// PoolSize sets the maximum number of concurrently executing workflow steps.
// A value of 0 uses the default: runtime.NumCPU() * 2.
func PoolSize(n int) Option {
	return func(o *Options) {
		o.PoolSize = n
	}
}

// WithHooks sets the lifecycle callbacks invoked by the engine (see Hooks).
func WithHooks(h *Hooks) Option {
	return func(o *Options) {
		o.Hooks = h
	}
}

// WorkflowOption func signature
type WorkflowOption func(*WorkflowOptions)

// WorkflowOptions holds workflow options
type WorkflowOptions struct {
	Context context.Context
	ID      string
}

// WorkflowID set workflow id
func WorkflowID(id string) WorkflowOption {
	return func(o *WorkflowOptions) {
		o.ID = id
	}
}

// ExecuteOptions holds execute options
type ExecuteOptions struct {
	// Client holds the client.Client
	Client client.Client
	// Tracer holds the tracer
	Tracer tracer.Tracer
	// Logger holds the logger
	Logger logger.Logger
	// Meter holds the meter
	Meter meter.Meter
	// Context can be used to abort execution or pass additional opts
	Context context.Context
	// Start step
	Start string
	// Timeout for execution
	Timeout time.Duration
	// Async enables async execution
	Async bool
	// EID is the execution id to use. When empty, a new id is generated.
	// Reusing an existing execution id resumes that execution: steps that
	// already completed successfully are skipped (see handleWorkflow).
	EID string
	// Compensate controls automatic rollback on failure.
	//   true  (default) — on a step failure the engine unwinds all
	//                      successful steps in reverse start order.
	//   false           — on failure the engine stops, keeps every
	//                      successful step's effect, and persists the
	//                      execution as failed. The caller can later resume
	//                      (retry) or compensate it explicitly (see
	//                      Workflow.Compensate / CompensateSteps).
	Compensate bool
}

// ExecuteOption func signature
type ExecuteOption func(*ExecuteOptions)

// ExecuteClient pass client.Client to ExecuteOption
func ExecuteClient(c client.Client) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Client = c
	}
}

// ExecuteTracer pass tracer.Tracer to ExecuteOption
func ExecuteTracer(t tracer.Tracer) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Tracer = t
	}
}

// ExecuteLogger pass logger.Logger to ExecuteOption
func ExecuteLogger(l logger.Logger) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Logger = l
	}
}

// ExecuteMeter pass meter.Meter to ExecuteOption
func ExecuteMeter(m meter.Meter) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Meter = m
	}
}

// ExecuteContext pass context.Context ot ExecuteOption
func ExecuteContext(ctx context.Context) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Context = ctx
	}
}

// ExecuteTimeout pass timeout time.Duration for execution
func ExecuteTimeout(td time.Duration) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Timeout = td
	}
}

// ExecuteAsync says that caller does not wait for execution complete
func ExecuteAsync(b bool) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Async = b
	}
}

// ExecuteEID sets the execution id to use. When empty, a new id is generated.
// Reusing an existing execution id resumes that execution: steps that already
// completed successfully are skipped.
func ExecuteEID(eid string) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.EID = eid
	}
}

// ExecuteCompensate controls automatic rollback on failure. When false, a
// failed execution is left with its successful steps intact (no automatic
// rollback) so the caller can later resume it (retry) or compensate it
// explicitly (see Workflow.Compensate / CompensateSteps). The default is
// true (automatic rollback on failure).
func ExecuteCompensate(b bool) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Compensate = b
	}
}

// NewExecuteOptions create new ExecuteOptions struct
func NewExecuteOptions(opts ...ExecuteOption) ExecuteOptions {
	options := ExecuteOptions{
		Client:     client.DefaultClient,
		Logger:     logger.DefaultLogger,
		Tracer:     tracer.DefaultTracer,
		Meter:      meter.DefaultMeter,
		Context:    context.Background(),
		Compensate: true,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// RetryPolicy configures step-level retries.
type RetryPolicy struct {
	// MaxAttempts is the total number of attempts, including the first
	// one. Values below 1 mean no retries.
	MaxAttempts int
	// Backoff is the initial delay between attempts. The delay doubles
	// with each attempt (exponential backoff). Zero means no delay.
	Backoff time.Duration
}

// StepOptions holds step options
type StepOptions struct {
	Context  context.Context
	ID       string
	Requires []string
	// Retry configures step-level retries. Nil means a single attempt.
	Retry *RetryPolicy
	// CompensateRetry configures retries of the step's compensation
	// (rollback). Nil means a single compensation attempt.
	CompensateRetry *RetryPolicy
	// CompensateTimeout bounds a single compensation attempt. Zero means no
	// timeout (the step context governs).
	CompensateTimeout time.Duration
	// ExecuteFunc is the step's execution logic, set via StepExecute. It is
	// used by steps created with NewStep.
	ExecuteFunc func(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error)
	// CompensateFunc is the step's compensation logic, set via StepCompensate.
	// It is used by steps created with NewStep.
	CompensateFunc func(ctx context.Context, in *CompensateInput, opts ...ExecuteOption) error
}

// StepOption func signature
type StepOption func(*StepOptions)

// NewStepOptions create new StepOptions struct
func NewStepOptions(opts ...StepOption) StepOptions {
	options := StepOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// StepID sets the step id for dag
func StepID(id string) StepOption {
	return func(o *StepOptions) {
		o.ID = id
	}
}

// StepRequires specifies required steps
func StepRequires(steps ...string) StepOption {
	return func(o *StepOptions) {
		o.Requires = steps
	}
}

// StepRetry sets the retry policy of the step.
func StepRetry(p RetryPolicy) StepOption {
	return func(o *StepOptions) {
		o.Retry = &p
	}
}

// StepCompensateRetry sets the retry policy of the step's compensation.
func StepCompensateRetry(p RetryPolicy) StepOption {
	return func(o *StepOptions) {
		o.CompensateRetry = &p
	}
}

// StepCompensateTimeout sets the timeout of a single compensation attempt.
func StepCompensateTimeout(d time.Duration) StepOption {
	return func(o *StepOptions) {
		o.CompensateTimeout = d
	}
}

// StepExecute sets the execution logic of a step created with NewStep.
func StepExecute(fn func(ctx context.Context, req *Message, opts ...ExecuteOption) (*Message, error)) StepOption {
	return func(o *StepOptions) {
		o.ExecuteFunc = fn
	}
}

// StepCompensate sets the compensation (rollback) logic of a step created
// with NewStep.
func StepCompensate(fn func(ctx context.Context, in *CompensateInput, opts ...ExecuteOption) error) StepOption {
	return func(o *StepOptions) {
		o.CompensateFunc = fn
	}
}
