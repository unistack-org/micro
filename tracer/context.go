// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"
)

type tracerKey struct{}

// FromContext returns a tracer from context
func FromContext(ctx context.Context) (Tracer, bool) {
	if ctx == nil {
		return nil, false
	}
	if tracer, ok := ctx.Value(tracerKey{}).(Tracer); ok {
		return tracer, true
	}
	return nil, false
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
func SpanFromContext(ctx context.Context) (Span, bool) {
	if ctx == nil {
		return nil, false
	}
	if span, ok := ctx.Value(spanKey{}).(Span); ok {
		return span, true
	}
	return nil, false
}

// NewSpanContext saves the span in the context
func NewSpanContext(ctx context.Context, span Span) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, spanKey{}, span)
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
