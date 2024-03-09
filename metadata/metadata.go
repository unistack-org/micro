// Package metadata is a way of defining message headers
package metadata

import (
	"net/textproto"
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
type Metadata map[string][]string

// Get returns value from metadata by key
func (md Metadata) Get(key string) (string, bool) {
	// fast path
	val, ok := md[key]
	if !ok {
		// slow path
		val, ok = md[textproto.CanonicalMIMEHeaderKey(key)]
	}
	return strings.Join(val, ","), ok
}

// Set is used to store value in metadata
func (md Metadata) Set(kv ...string) {
	if len(kv)%2 == 1 {
		kv = kv[:len(kv)-1]
	}
	for idx := 0; idx < len(kv); idx += 2 {
		md[textproto.CanonicalMIMEHeaderKey(kv[idx])] = []string{kv[idx+1]}
	}
}

// Append is used to append value in metadata
func (md Metadata) Append(k string, v ...string) {
	if len(v) == 0 {
		return
	}
	k = textproto.CanonicalMIMEHeaderKey(k)
	md[k] = append(md[k], v...)
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
	nmd := make(Metadata, len(md))
	for k, v := range md {
		nmd[k] = v
	}
	return nmd
}

// New return new sized metadata
func New(size int) Metadata {
	if size == 0 {
		size = 2
	}
	return make(Metadata, size)
}

// Merge merges metadata to existing metadata, overwriting if specified
func Merge(omd Metadata, mmd Metadata, overwrite bool) Metadata {
	nmd := Copy(omd)
	for key, val := range mmd {
		nval, ok := nmd[key]
		switch {
		case ok && overwrite:
			nmd[key] = nval
			continue
		case ok && !overwrite:
			continue
		case !ok:
			for _, v := range val {
				if v != "" {
					nval = append(nval, v)
				}
			}
			nmd[key] = nval
		}
	}
	return nmd
}

// Pairs from which metadata created
func Pairs(kv ...string) Metadata {
	if len(kv)%2 == 1 {
		kv = kv[:len(kv)-1]
	}
	md := make(Metadata, len(kv)/2)
	md.Set(kv...)
	return md
}
