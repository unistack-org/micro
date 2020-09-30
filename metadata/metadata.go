// Package metadata is a way of defining message headers
package metadata

import (
	"context"
	"net/textproto"
)

type metadataKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

func (md Metadata) Get(key string) (string, bool) {
	// fast path
	val, ok := md[key]
	if !ok {
		// slow path
		val, ok = md[textproto.CanonicalMIMEHeaderKey(key)]
	}
	return val, ok
}

func (md Metadata) Set(key, val string) {
	md[textproto.CanonicalMIMEHeaderKey(key)] = val
}

func (md Metadata) Del(key string) {
	delete(md, textproto.CanonicalMIMEHeaderKey(key))
}

// Copy makes a copy of the metadata
func Copy(md Metadata) Metadata {
	nmd := make(Metadata, len(md))
	for k, v := range md {
		nmd[k] = v
	}
	return nmd
}

func Del(ctx context.Context, key string) context.Context {
	md, ok := FromContext(ctx)
	if !ok {
		md = make(Metadata)
	}
	md.Del(key)
	return context.WithValue(ctx, metadataKey{}, md)
}

// Set add key with val to metadata
func Set(ctx context.Context, key, val string) context.Context {
	md, ok := FromContext(ctx)
	if !ok {
		md = make(Metadata)
	}
	md.Set(key, val)
	return context.WithValue(ctx, metadataKey{}, md)
}

// Get returns a single value from metadata in the context
func Get(ctx context.Context, key string) (string, bool) {
	md, ok := FromContext(ctx)
	if !ok {
		return "", ok
	}
	return md.Get(key)
}

// FromContext returns metadata from the given context
func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	return md, ok
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataKey{}, md)
}

// MergeContext merges metadata to existing metadata, overwriting if specified
func MergeContext(ctx context.Context, pmd Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	md, ok := FromContext(ctx)
	if !ok {
		md = make(Metadata)
	}
	nmd := make(Metadata, len(md))
	for key, val := range md {
		nmd.Set(key, val)
	}
	for key, val := range pmd {
		if _, ok := nmd[key]; ok && !overwrite {
			// skip
		} else if val != "" {
			nmd.Set(key, val)
		} else {
			nmd.Del(key)
		}
	}
	return context.WithValue(ctx, metadataKey{}, nmd)
}
