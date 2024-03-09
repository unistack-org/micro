// Package metadata is a way of defining message headers
package metadata

import (
	"context"
)

type (
	mdIncomingKey struct{}
	mdOutgoingKey struct{}
	mdKey         struct{}
)

// FromIncomingContext returns metadata from incoming ctx
// returned metadata shoud not be modified or race condition happens
func FromIncomingContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	md, ok := ctx.Value(mdIncomingKey{}).(Metadata)
	if !ok || md == nil {
		return nil, false
	}
	return md, ok
}

// FromOutgoingContext returns metadata from outgoing ctx
// returned metadata shoud not be modified or race condition happens
func FromOutgoingContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	md, ok := ctx.Value(mdOutgoingKey{}).(Metadata)
	if !ok || md == nil {
		return nil, false
	}
	return md, ok
}

// FromContext returns metadata from the given context
// returned metadata shoud not be modified or race condition happens
func FromContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	md, ok := ctx.Value(mdKey{}).(Metadata)
	if !ok || md == nil {
		return nil, false
	}
	return md, ok
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, mdKey{}, md)
	return ctx
}

// NewIncomingContext creates a new context with incoming metadata attached
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, mdIncomingKey{}, md)
	return ctx
}

// NewOutgoingContext creates a new context with outcoming metadata attached
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, mdOutgoingKey{}, md)
	return ctx
}

// AppendOutgoingContext apends new md to context
func AppendOutgoingContext(ctx context.Context, kv ...string) context.Context {
	md := Pairs(kv...)
	omd, ok := FromOutgoingContext(ctx)
	if !ok {
		return NewOutgoingContext(ctx, md)
	}
	nmd := Merge(omd, md, true)
	return NewOutgoingContext(ctx, nmd)
}

// AppendIncomingContext apends new md to context
func AppendIncomingContext(ctx context.Context, kv ...string) context.Context {
	md := Pairs(kv...)
	omd, ok := FromIncomingContext(ctx)
	if !ok {
		return NewIncomingContext(ctx, md)
	}
	nmd := Merge(omd, md, true)
	return NewIncomingContext(ctx, nmd)
}
