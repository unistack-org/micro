//go:build !exclude

// Package metadata is a way of defining message headers
package metadata

import (
	"net/textproto"
	"sort"
	"strings"
)

var (
	// HeaderTopic is the header name that contains topic name
	HeaderTopic = "Micro-Topic"
	// HeaderContentType specifies content type of message
	HeaderContentType = "Content-Type"
	// HeaderEndpoint specifies endpoint in service
	HeaderEndpoint = "Micro-Endpoint"
	// HeaderService specifies service
	HeaderService = "Micro-Service"
	// HeaderTimeout specifies timeout of operation
	HeaderTimeout = "Micro-Timeout"
	// HeaderAuthorization specifies Authorization header
	HeaderAuthorization = "Authorization"
	// HeaderXRequestID specifies request id
	HeaderXRequestID = "X-Request-Id"
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string]string

type rawMetadata struct {
	md Metadata
}

// defaultMetadataSize used when need to init new Metadata
var defaultMetadataSize = 2

// Iterator used to iterate over metadata with order
type Iterator struct {
	md   Metadata
	keys []string
	cur  int
	cnt  int
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

func (md Metadata) MustGet(key string) string {
	val, ok := md.Get(key)
	if !ok {
		panic("missing metadata key")
	}
	return val
}

// Get returns value from metadata by key
func (md Metadata) Get(key string) (string, bool) {
	// fast path
	val, ok := md[key]
	if !ok {
		// slow path
		val, ok = md[textproto.CanonicalMIMEHeaderKey(key)]
		if !ok {
			val, ok = md[strings.ToLower(key)]
		}
	}
	return val, ok
}

// Set is used to store value in metadata
func (md Metadata) Set(kv ...string) {
	if len(kv)%2 == 1 {
		kv = kv[:len(kv)-1]
	}
	for idx := 0; idx < len(kv); idx += 2 {
		md[textproto.CanonicalMIMEHeaderKey(kv[idx])] = kv[idx+1]
	}
}

// Del is used to remove value from metadata
func (md Metadata) Del(keys ...string) {
	for _, key := range keys {
		// fast path
		delete(md, key)
		// slow path
		delete(md, textproto.CanonicalMIMEHeaderKey(key))
		// very slow path
		delete(md, strings.ToLower(key))
	}
}

// Copy makes a copy of the metadata
func (md Metadata) CopyTo(dst Metadata) {
	for k, v := range md {
		dst[k] = v
	}
}

// Copy makes a copy of the metadata
func Copy(md Metadata, exclude ...string) Metadata {
	nmd := New(len(md))
	for k, v := range md {
		nmd[k] = v
	}
	nmd.Del(exclude...)
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
	var ok bool
	nmd := Copy(omd)
	for key, val := range mmd {
		_, ok = nmd[key]
		switch {
		case ok && !overwrite:
			continue
		case val != "":
			nmd[key] = val
		case ok && val == "":
			nmd.Del(key)
		}
	}
	return nmd
}

// Pairs from which metadata created
func Pairs(kv ...string) (Metadata, bool) {
	if len(kv)%2 == 1 {
		return nil, false
	}
	md := New(len(kv) / 2)
	for idx := 0; idx < len(kv); idx += 2 {
		md[kv[idx]] = kv[idx+1]
	}
	return md, true
}
