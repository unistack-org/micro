// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/options"
)

// DefaultTracer is the global default tracer
var DefaultTracer = NewTracer()

var (
	// TraceIDKey is the key used for the trace id in the log call
	TraceIDKey = "trace-id"
	// SpanIDKey is the key used for the span id in the log call
	SpanIDKey = "span-id"
)

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
	Init(...options.Option) error
	// Start a trace
	Start(ctx context.Context, name string, opts ...options.Option) (context.Context, Span)
	// Flush flushes spans
	Flush(ctx context.Context) error
}

type Span interface {
	// Tracer return underlining tracer
	Tracer() Tracer
	// Finish complete and send span
	Finish(opts ...options.Option)
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
	AddEvent(name string, opts ...options.Option)
	// AddLogs append logs to span
	AddLogs(kv ...interface{})
	// Kind returns span kind
	Kind() SpanKind
	// TraceID returns trace id
	TraceID() string
	// SpanID returns span id
	SpanID() string
}
