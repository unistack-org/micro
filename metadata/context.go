//go:build !exclude

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
	md, ok := ctx.Value(mdIncomingKey{}).(*rawMetadata)
	if !ok || md.md == nil {
		return nil, false
	}
	return md.md, ok
}

// MustIncomingContext returns metadata from incoming ctx
// returned metadata shoud not be modified or race condition happens.
// If metadata not exists panics.
func MustIncomingContext(ctx context.Context) Metadata {
	md, ok := FromIncomingContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// FromOutgoingContext returns metadata from outgoing ctx
// returned metadata shoud not be modified or race condition happens
func FromOutgoingContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	md, ok := ctx.Value(mdOutgoingKey{}).(*rawMetadata)
	if !ok || md.md == nil {
		return nil, false
	}
	return md.md, ok
}

// MustOutgoingContext returns metadata from outgoing ctx
// returned metadata shoud not be modified or race condition happens.
// If metadata not exists panics.
func MustOutgoingContext(ctx context.Context) Metadata {
	md, ok := FromOutgoingContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// FromContext returns metadata from the given context
// returned metadata shoud not be modified or race condition happens
func FromContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	md, ok := ctx.Value(mdKey{}).(*rawMetadata)
	if !ok || md.md == nil {
		return nil, false
	}
	return md.md, ok
}

// MustContext returns metadata from the given context
// returned metadata shoud not be modified or race condition happens
func MustContext(ctx context.Context) Metadata {
	md, ok := FromContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, mdKey{}, &rawMetadata{md})
}

// SetOutgoingContext modify outgoing context with given metadata
func SetOutgoingContext(ctx context.Context, md Metadata) bool {
	if ctx == nil {
		return false
	}
	if omd, ok := ctx.Value(mdOutgoingKey{}).(*rawMetadata); ok {
		omd.md = md
		return true
	}
	return false
}

// SetIncomingContext modify incoming context with given metadata
func SetIncomingContext(ctx context.Context, md Metadata) bool {
	if ctx == nil {
		return false
	}
	if omd, ok := ctx.Value(mdIncomingKey{}).(*rawMetadata); ok {
		omd.md = md
		return true
	}
	return false
}

// NewIncomingContext creates a new context with incoming metadata attached
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, mdIncomingKey{}, &rawMetadata{md})
}

// NewOutgoingContext creates a new context with outcoming metadata attached
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, mdOutgoingKey{}, &rawMetadata{md})
}

// AppendOutgoingContext apends new md to context
func AppendOutgoingContext(ctx context.Context, kv ...string) context.Context {
	md, ok := Pairs(kv...)
	if !ok {
		return ctx
	}
	omd, ok := FromOutgoingContext(ctx)
	if !ok {
		return NewOutgoingContext(ctx, md)
	}
	for k, v := range md {
		omd.Set(k, v)
	}
	return ctx
}

// AppendIncomingContext apends new md to context
func AppendIncomingContext(ctx context.Context, kv ...string) context.Context {
	md, ok := Pairs(kv...)
	if !ok {
		return ctx
	}
	omd, ok := FromIncomingContext(ctx)
	if !ok {
		return NewIncomingContext(ctx, md)
	}
	for k, v := range md {
		omd.Set(k, v)
	}
	return ctx
}

// AppendContext apends new md to context
func AppendContext(ctx context.Context, kv ...string) context.Context {
	md, ok := Pairs(kv...)
	if !ok {
		return ctx
	}
	omd, ok := FromContext(ctx)
	if !ok {
		return NewContext(ctx, md)
	}
	for k, v := range md {
		omd.Set(k, v)
	}
	return ctx
}
