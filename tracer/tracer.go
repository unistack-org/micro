// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"
)

var (
	// DefaultTracer is the global default tracer
	DefaultTracer Tracer = NewTracer()
)

// Tracer is an interface for distributed tracing
type Tracer interface {
	// Name return tracer name
	Name() string
	// Init tracer with options
	Init(...Option) error
	// Start a trace
	Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span)
}

type Span interface {
	// Tracer return underlining tracer
	Tracer() Tracer
	// Finish complete and send span
	Finish(opts ...SpanOption)
	// AddEvent add event to span
	AddEvent(name string, opts ...EventOption)
	// Context return context with span
	Context() context.Context
	// SetName set the span name
	SetName(name string)
	// SetLabels set the span labels
	SetLabels(labels ...Label)
}

type Label struct {
	key string
	val interface{}
}

func Any(k string, v interface{}) Label {
	return Label{k, v}
}

func String(k string, v string) Label {
	return Label{k, v}
}

func Int(k string, v int) Label {
	return Label{k, v}
}

func Int64(k string, v int64) Label {
	return Label{k, v}
}

func Float64(k string, v float64) Label {
	return Label{k, v}
}

func Bool(k string, v bool) Label {
	return Label{k, v}
}
