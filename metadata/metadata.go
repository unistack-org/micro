// Package metadata is a way of defining message headers
package metadata

import (
	"net/textproto"
	"sort"
)

var (
	// HeaderPrefix for all headers passed
	HeaderPrefix = "Micro-"
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

type rawMetadata struct {
	md Metadata
}

var (
	// defaultMetadataSize used when need to init new Metadata
	defaultMetadataSize = 2
)

// Iterator used to iterate over metadata with order
type Iterator struct {
	cur  int
	cnt  int
	keys []string
	md   Metadata
}

// Next advance iterator to next element
func (iter *Iterator) Next(k, v *string) bool {
	if iter.cur+1 > iter.cnt {
		return false
	}

	*k = iter.keys[iter.cur]
	*v = iter.md[*k]
	iter.cur++
	return true
}

// Iterator returns the itarator for metadata in sorted order
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
	delete(md, key)
	// slow path
	delete(md, textproto.CanonicalMIMEHeaderKey(key))
}

// Copy makes a copy of the metadata
func Copy(md Metadata) Metadata {
	nmd := New(len(md))
	for key, val := range md {
		nmd.Set(key, val)
	}
	return nmd
}

// New return new sized metadata
func New(size int) Metadata {
	if size == 0 {
		size = defaultMetadataSize
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
