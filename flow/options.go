package flow

import (
	"context"
	"time"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/tracer"
)

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
func NewOptions(opts ...options.Option) Options {
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

// WorkflowOptions holds workflow options
type WorkflowOptions struct {
	Context context.Context
	ID      string
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

// Reverse says that dag must be run in reverse order
func Reverse(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".Reverse")
	}
}

// Async says that caller does not wait for execution complete
func Async(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".Async")
	}
}

// NewExecuteOptions create new ExecuteOptions struct
func NewExecuteOptions(opts ...options.Option) ExecuteOptions {
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

// NewStepOptions create new StepOptions struct
func NewStepOptions(opts ...options.Option) StepOptions {
	options := StepOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Requires specifies required steps
func Requires(steps ...string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, steps, ".Requires")
	}
}

// Fallback set the step to run on error
func Fallback(step string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, step, ".Fallback")
	}
}

// ID sets the step ID
func StepID(id string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, id, ".ID")
	}
}
