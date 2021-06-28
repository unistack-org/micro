package flow

import (
	"context"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/tracer"
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
	ID      string
	Context context.Context
}

// WorkflowID set workflow id
func WorkflowID(id string) WorkflowOption {
	return func(o *WorkflowOptions) {
		o.ID = id
	}
}

type ExecuteOptions struct {
	Context context.Context
}

type ExecuteOption func(*ExecuteOptions)

type StepOptions struct {
	ID      string
	Name    string
	Context context.Context
}

type StepOption func(*StepOptions)
