package flow

import (
	"context"
	"time"

	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/store"
	"go.unistack.org/micro/v3/tracer"
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
// Can be used for extra option values.
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

// WorflowOption signature
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

type ExecuteOption func(*ExecuteOptions)

func ExecuteClient(c client.Client) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Client = c
	}
}

func ExecuteTracer(t tracer.Tracer) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Tracer = t
	}
}

func ExecuteLogger(l logger.Logger) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Logger = l
	}
}

func ExecuteMeter(m meter.Meter) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Meter = m
	}
}

func ExecuteContext(ctx context.Context) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Context = ctx
	}
}

func ExecuteReverse(b bool) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Reverse = b
	}
}

func ExecuteTimeout(td time.Duration) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Timeout = td
	}
}

func ExecuteAsync(b bool) ExecuteOption {
	return func(o *ExecuteOptions) {
		o.Async = b
	}
}

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

type StepOptions struct {
	Context  context.Context
	Fallback string
	ID       string
	Requires []string
}

type StepOption func(*StepOptions)

func NewStepOptions(opts ...StepOption) StepOptions {
	options := StepOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func StepID(id string) StepOption {
	return func(o *StepOptions) {
		o.ID = id
	}
}

func StepRequires(steps ...string) StepOption {
	return func(o *StepOptions) {
		o.Requires = steps
	}
}

func StepFallback(step string) StepOption {
	return func(o *StepOptions) {
		o.Fallback = step
	}
}
