// Package tracer provides an interface for distributed tracing
package tracer

import (
	"context"

	"github.com/unistack-org/micro/v3/metadata"
)

const (
	traceIDKey = "Micro-Trace-Id"
	spanIDKey  = "Micro-Span-Id"
)

// FromContext returns a span from context
func FromContext(ctx context.Context) (traceID string, parentSpanID string, isFound bool) {
	if ctx == nil {
		return "", "", false
	}
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

// NewContext saves the trace and span ids in the context
func NewContext(ctx context.Context, traceID, parentSpanID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	md := metadata.New(2)
	md.Set(traceIDKey, traceID)
	md.Set(spanIDKey, parentSpanID)
	return metadata.MergeContext(ctx, md, true)
}
