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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", false
	}
	traceID, traceOk := md.Get(traceIDKey)
	microID, microOk := md.Get("Micro-Id")
	if !traceOk && !microOk {
		isFound = false
		return
	}
	if !traceOk {
		traceID = microID
	}
	parentSpanID, ok = md.Get(spanIDKey)
	return traceID, parentSpanID, ok
}

// NewContext saves the trace and span ids in the context
func NewContext(ctx context.Context, traceID, parentSpanID string) context.Context {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.New(2)
	}
	md.Set(traceIDKey, traceID)
	md.Set(spanIDKey, parentSpanID)
	return metadata.NewContext(ctx, md)
}
