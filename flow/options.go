package flow

import (
	"context"
	"time"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/tracer"
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
	// Reverse execution
	Reverse bool
	// Async enables async execution
	Async bool
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

// ExecuteReverse says that dag must be run in reverse order
func ExecuteReverse(b bool) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Reverse = b
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

// NewExecuteOptions create new ExecuteOptions struct
func NewExecuteOptions(opts ...ExecuteOption) ExecuteOptions {
	options := ExecuteOptions{
		Client:  client.DefaultClient,
		Logger:  logger.DefaultLogger,
		Tracer:  tracer.DefaultTracer,
		Meter:   meter.DefaultMeter,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// StepOptions holds step options
type StepOptions struct {
	Context  context.Context
	Fallback string
	ID       string
	Requires []string
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

// StepFallback set the step to run on error
func StepFallback(step string) StepOption {
	return func(o *StepOptions) {
		o.Fallback = step
	}
}
