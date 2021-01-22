// Package metadata is a way of defining message headers
package metadata

import (
	"context"
	"net/textproto"
	"sort"
)

type metadataKey struct{}

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

var (
	// DefaultMetadataSize used when need to init new Metadata
	DefaultMetadataSize = 6
)

type Iterator struct {
	cur  int
	cnt  int
	keys []string
	md   Metadata
}

func (iter *Iterator) Next(k, v *string) bool {
	if iter.cur+1 > iter.cnt {
		return false
	}

	*k = iter.keys[iter.cur]
	*v = iter.md[*k]
	iter.cur++
	return true
}

// Iterate returns run user func with map key, val sorted by key
func (md Metadata) Iterator() *Iterator {
	iter := &Iterator{md: md, cnt: len(md)}
	iter.keys = make([]string, 0, iter.cnt)
	for k := range md {
		iter.keys = append(iter.keys, k)
	}
	sort.Strings(iter.keys)
	return iter
}

// Get returns value from metadata by key
func (md Metadata) Get(key string) (string, bool) {
	// fast path
	val, ok := md[key]
	if !ok {
		// slow path
		val, ok = md[textproto.CanonicalMIMEHeaderKey(key)]
	}
	return val, ok
}

// Set is used to store value in metadata
func (md Metadata) Set(key, val string) {
	md[textproto.CanonicalMIMEHeaderKey(key)] = val
}

// Del is used to remove value from metadata
func (md Metadata) Del(key string) {
	// fast path
	if _, ok := md[key]; ok {
		delete(md, key)
	} else {
		// slow path
		delete(md, textproto.CanonicalMIMEHeaderKey(key))
	}
}

// Copy makes a copy of the metadata
func Copy(md Metadata) Metadata {
	nmd := New(len(md))
	for key, val := range md {
		nmd.Set(key, val)
	}
	return nmd
}

// Del deletes key from metadata
func Del(ctx context.Context, key string) context.Context {
	md, ok := FromContext(ctx)
	if !ok {
		md = New(0)
	}
	md.Del(key)
	return context.WithValue(ctx, metadataKey{}, md)
}

// Set add key with val to metadata
func Set(ctx context.Context, key, val string) context.Context {
	md, ok := FromContext(ctx)
	if !ok {
		md = New(0)
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

// New return new sized metadata
func New(size int) Metadata {
	if size == 0 {
		size = DefaultMetadataSize
	}
	return make(Metadata, size)
}

// Merge merges metadata to existing metadata, overwriting if specified
func Merge(omd Metadata, mmd Metadata, overwrite bool) Metadata {
	nmd := Copy(omd)
	for key, val := range mmd {
		if _, ok := nmd[key]; ok && !overwrite {
			// skip
		} else if val != "" {
			nmd.Set(key, val)
		} else {
			nmd.Del(key)
		}
	}
	return nmd
}
