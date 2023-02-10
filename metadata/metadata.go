// Package metadata is a way of defining message headers
package metadata // import "go.unistack.org/micro/v3/metadata"

import (
	"net/textproto"
	"sort"
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
)

// Metadata is our way of representing request headers internally.
// They're used at the RPC level and translate back and forth
// from Transport headers.
type Metadata map[string][]string

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
func (iter *Iterator) Next(k *string, v *[]string) bool {
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

// Values returns values from metadata by key
func (md Metadata) Values(key string) ([]string, bool) {
	if md == nil {
		return nil, false
	}
	// fast path
	val, ok := md[key]
	if !ok {
		// slow path
		val, ok = md[textproto.CanonicalMIMEHeaderKey(key)]
	}
	return val, ok
}

// Get returns value from metadata by key
func (md Metadata) Get(key string) (string, bool) {
	if md == nil {
		return "", false
	}
	// fast path
	val, ok := md[key]
	if !ok {
		// slow path
		val, ok = md[textproto.CanonicalMIMEHeaderKey(key)]
	}
	if len(val) == 0 {
		return "", false
	}
	return val[0], ok
}

// Set is used to store value in metadata
func (md Metadata) Set(kv ...string) {
	if md == nil {
		return
	}
	if len(kv)%2 == 1 {
		kv = kv[:len(kv)-1]
	}
	for idx := 0; idx < len(kv); idx += 2 {
		md[textproto.CanonicalMIMEHeaderKey(kv[idx])] = []string{kv[idx+1]}
	}
}

// Add is used to append value in metadata
func (md Metadata) Add(k string, v string) {
	kn := textproto.CanonicalMIMEHeaderKey(k)
	md[kn] = append(md[kn], v)
}

// Del is used to remove value from metadata
func (md Metadata) Del(keys ...string) {
	for _, key := range keys {
		// fast path
		delete(md, key)
		// slow path
		delete(md, textproto.CanonicalMIMEHeaderKey(key))
	}
}

// Copy makes a copy of the metadata
func Copy(md Metadata) Metadata {
	nmd := New(len(md))
	for key, val := range md {
		nmd[key] = val
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
	for key, nval := range mmd {
		oval, ok := nmd[key]
		switch {
		case ok && !overwrite:
			continue
		case len(nval) == 1 && nval[0] != "":
			nmd[key] = nval
		case ok && len(nval) > 1:
			sort.Strings(nval)
			sort.Strings(oval)
			for idx, v := range nval {
				if oval[idx] != v {
					oval[idx] = v
				}
			}
			nmd[key] = oval
		case ok && nval[0] == "":
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
	md.Set(kv...)
	return md, true
}
