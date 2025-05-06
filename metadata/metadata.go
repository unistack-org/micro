package metadata

import (
	"fmt"
	"net/textproto"
	"strings"
)

// defaultMetadataSize is used when initializing new Metadata.
var defaultMetadataSize = 2

// Metadata maps keys to values. Use the New, NewWithMetadata and Pairs functions to create it.
type Metadata map[string][]string

// New creates a zero-value Metadata with the specified size.
func New(l int) Metadata {
	if l == 0 {
		l = defaultMetadataSize
	}
	md := make(Metadata, l)
	return md
}

// NewWithMetadata creates a Metadata from the provided key-value map.
func NewWithMetadata(m map[string]string) Metadata {
	md := make(Metadata, len(m))
	for key, val := range m {
		md[key] = append(md[key], val)
	}
	return md
}

// Pairs returns a Metadata formed from the key-value mapping. It panics if the length of kv is odd.
func Pairs(kv ...string) Metadata {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := make(Metadata, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		md[kv[i]] = append(md[kv[i]], kv[i+1])
	}
	return md
}

// Join combines multiple Metadatas into a single Metadata.
// The order of values for each key is determined by the order in which the Metadatas are provided to Join.
func Join(mds ...Metadata) Metadata {
	out := Metadata{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return out
}

// Copy returns a deep copy of Metadata.
func Copy(src Metadata) Metadata {
	out := make(Metadata, len(src))
	for k, v := range src {
		out[k] = copyOf(v)
	}
	return out
}

// Copy returns a deep copy of Metadata.
func (md Metadata) Copy() Metadata {
	out := make(Metadata, len(md))
	for k, v := range md {
		out[k] = copyOf(v)
	}
	return out
}

// CopyTo performs a deep copy of Metadata to the out.
func (md Metadata) CopyTo(out Metadata) {
	for k, v := range md {
		out[k] = copyOf(v)
	}
}

// Len returns the number of items in Metadata.
func (md Metadata) Len() int {
	return len(md)
}

// AsMap returns a deep copy of Metadata as a map[string]string
func (md Metadata) AsMap() map[string]string {
	out := make(map[string]string, len(md))
	for k, v := range md {
		out[k] = strings.Join(v, ",")
	}
	return out
}

// AsHTTP1 returns a deep copy of Metadata with keys converted to canonical MIME header key format.
func (md Metadata) AsHTTP1() map[string][]string {
	out := make(map[string][]string, len(md))
	for k, v := range md {
		out[textproto.CanonicalMIMEHeaderKey(k)] = copyOf(v)
	}
	return out
}

// AsHTTP2 returns a deep copy of Metadata with keys converted to lowercase.
func (md Metadata) AsHTTP2() map[string][]string {
	out := make(map[string][]string, len(md))
	for k, v := range md {
		out[strings.ToLower(k)] = copyOf(v)
	}
	return out
}

// Get retrieves the values for a given key, checking the key in three formats:
// - exact case,
// - lower case,
// - canonical MIME header key format.
func (md Metadata) Get(k string) []string {
	v, ok := md[k]
	if !ok {
		v, ok = md[strings.ToLower(k)]
	}
	if !ok {
		v = md[textproto.CanonicalMIMEHeaderKey(k)]
	}
	return v
}

// GetJoined retrieves the values for a given key and joins them into a single string, separated by commas.
func (md Metadata) GetJoined(k string) string {
	return strings.Join(md.Get(k), ",")
}

// Set assigns the values to the given key.
func (md Metadata) Set(key string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	md[key] = vals
}

// Append adds values to the existing values for the given key.
func (md Metadata) Append(key string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	md[key] = append(md[key], vals...)
}

// Del removes the values for the given keys k. It checks and removes the keys in the following formats:
// - exact case,
// - lower case,
// - canonical MIME header key format.
func (md Metadata) Del(k ...string) {
	for i := range k {
		delete(md, k[i])
		delete(md, strings.ToLower(k[i]))
		delete(md, textproto.CanonicalMIMEHeaderKey(k[i]))
	}
}
