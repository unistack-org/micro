// Package trace provides an interface for distributed tracing
package tracer

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	DefaultTracer Tracer = newTracer()
)

// Tracer is an interface for distributed tracing
type Tracer interface {
	// Start a trace
	Start(ctx context.Context, name string) (context.Context, *Span)
	// Finish the trace
	Finish(*Span) error
	// Read the traces
	Read(...ReadOption) ([]*Span, error)
}

// SpanType describe the nature of the trace span
type SpanType int

const (
	// SpanTypeRequestInbound is a span created when serving a request
	SpanTypeRequestInbound SpanType = iota
	// SpanTypeRequestOutbound is a span created when making a service call
	SpanTypeRequestOutbound
)

// Span is used to record an entry
type Span struct {
	// Id of the trace
	Trace string
	// name of the span
	Name string
	// id of the span
	Id string
	// parent span id
	Parent string
	// Start time
	Started time.Time
	// Duration in nano seconds
	Duration time.Duration
	// associated data
	Metadata map[string]string
	// Type
	Type SpanType
}

const (
	traceIDKey = "Micro-Trace-Id"
	spanIDKey  = "Micro-Span-Id"
)

// FromContext returns a span from context
func FromContext(ctx context.Context) (traceID string, parentSpanID string, isFound bool) {
	traceID, traceOk := metadata.Get(ctx, traceIDKey)
	microID, microOk := metadata.Get(ctx, "Micro-Id")
	if !traceOk && !microOk {
		isFound = false
		return
	}
	if !traceOk {
		traceID = microID
	}
	parentSpanID, ok := metadata.Get(ctx, spanIDKey)
	return traceID, parentSpanID, ok
}

// ToContext saves the trace and span ids in the context
func ToContext(ctx context.Context, traceID, parentSpanID string) context.Context {
	return metadata.MergeContext(ctx, map[string]string{
		traceIDKey: traceID,
		spanIDKey:  parentSpanID,
	}, true)
}
