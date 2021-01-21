// Package metadata is a way of defining message headers
package metadata

import (
	"context"
)

// FromContext returns metadata from the given context
func FromContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	md, ok := ctx.Value(metadataKey{}).(Metadata)
	if !ok {
		return nil, ok
	}
	nmd := Copy(md)
	return nmd, ok
}

// NewContext creates a new context with the given metadata
func NewContext(ctx context.Context, md Metadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, metadataKey{}, Copy(md))
}

// MergeContext merges metadata to existing metadata, overwriting if specified
func MergeContext(ctx context.Context, pmd Metadata, overwrite bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	md, ok := FromContext(ctx)
	if !ok {
		return context.WithValue(ctx, metadataKey{}, Copy(pmd))
	}
	nmd := Copy(md)
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
