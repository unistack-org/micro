// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"
)

type tracerKey struct{}

var tracerKeyVal = tracerKey{}

// FromContext returns a tracer from context
func FromContext(ctx context.Context) (Tracer, bool) {
	if ctx == nil {
		return nil, false
	}
	if tracer, ok := ctx.Value(tracerKeyVal).(Tracer); ok {
		return tracer, true
	}
	return nil, false
}

// MustContext returns a tracer from context
func MustContext(ctx context.Context) Tracer {
	t, ok := FromContext(ctx)
	if !ok {
		panic("missing tracer")
	}
	return t
}

// NewContext saves the tracer in the context
func NewContext(ctx context.Context, tracer Tracer) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, tracerKeyVal, tracer)
}

type spanKey struct{}

var spanKeyVal = spanKey{}

// SpanFromContext returns a span from context
func SpanMustContext(ctx context.Context) Span {
	sp, ok := SpanFromContext(ctx)
	if !ok {
		panic("missing span")
	}
	return sp
}

// SpanFromContext returns a span from context
func SpanFromContext(ctx context.Context) (Span, bool) {
	if ctx == nil {
		return nil, false
	}
	if span, ok := ctx.Value(spanKeyVal).(Span); ok {
		return span, true
	}
	return nil, false
}

// NewSpanContext saves the span in the context
func NewSpanContext(ctx context.Context, span Span) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, spanKeyVal, span)
}

// SetOption returns a function to setup a context with given value
func SetOption(k, v any) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
