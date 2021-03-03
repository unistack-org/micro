// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"
)

type tracerKey struct{}

// FromContext returns a tracer from context
func FromContext(ctx context.Context) Tracer {
	if ctx == nil {
		return DefaultTracer
	}
	if tracer, ok := ctx.Value(tracerKey{}).(Tracer); ok {
		return tracer
	}
	return DefaultTracer
}

// NewContext saves the tracer in the context
func NewContext(ctx context.Context, tracer Tracer) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, tracerKey{}, tracer)
}

type spanKey struct{}

// SpanFromContext returns a span from context
func SpanFromContext(ctx context.Context) Span {
	if ctx == nil {
		return &noopSpan{}
	}
	if span, ok := ctx.Value(spanKey{}).(Span); ok {
		return span
	}
	return &noopSpan{}
}

// NewSpanContext saves the span in the context
func NewSpanContext(ctx context.Context, span Span) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, spanKey{}, span)
}
