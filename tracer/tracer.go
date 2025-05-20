// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"

	"go.unistack.org/micro/v3/logger"
)

var (
	// DefaultTracer is the global default tracer
	DefaultTracer Tracer = NewTracer() //nolint:revive
	// TraceIDKey is the key used for the trace id in the log call
	TraceIDKey = "trace-id"
	// SpanIDKey is the key used for the span id in the log call
	SpanIDKey = "span-id"
	// DefaultSkipEndpoints is the slice of endpoint that must not be traced
	DefaultSkipEndpoints = []string{
		"MeterService.Metrics",
		"HealthService.Live",
		"HealthService.Ready",
		"HealthService.Version",
	}
	DefaultContextAttrFuncs []ContextAttrFunc
)

type ContextAttrFunc func(ctx context.Context) []interface{}

func init() {
	logger.DefaultContextAttrFuncs = append(logger.DefaultContextAttrFuncs,
		func(ctx context.Context) []interface{} {
			if span, ok := SpanFromContext(ctx); ok {
				return []interface{}{
					TraceIDKey, span.TraceID(),
					SpanIDKey, span.SpanID(),
				}
			}
			return nil
		})
}

// Tracer is an interface for distributed tracing
type Tracer interface {
	// Name return tracer name
	Name() string
	// Init tracer with options
	Init(...Option) error
	// Start a trace
	Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span)
	// Extract get span metadata from context
	// Extract(ctx context.Context)
	// Flush flushes spans
	Flush(ctx context.Context) error
	// Enabled returns tracer status
	Enabled() bool
}

type Span interface {
	// Tracer return underlining tracer
	Tracer() Tracer
	// Finish complete and send span
	Finish(opts ...SpanOption)
	// Context return context with span
	Context() context.Context
	// SetName set the span name
	SetName(name string)
	// SetStatus set the span status code and msg
	SetStatus(status SpanStatus, msg string)
	// Status returns span status and msg
	Status() (SpanStatus, string)
	// AddLabels append labels to span
	AddLabels(kv ...interface{})
	// AddEvent append event to span
	AddEvent(name string, opts ...EventOption)
	// AddEvent append event to span
	AddLogs(kv ...interface{})
	// Kind returns span kind
	Kind() SpanKind
	// TraceID returns trace id
	TraceID() string
	// SpanID returns span id
	SpanID() string
	// IsRecording returns the recording state of the Span.
	IsRecording() bool
}
