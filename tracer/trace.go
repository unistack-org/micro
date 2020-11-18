// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	// DefaultTracer is the global default tracer
	DefaultTracer Tracer = NewTracer()
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
	Metadata metadata.Metadata
	// Type
	Type SpanType
}
