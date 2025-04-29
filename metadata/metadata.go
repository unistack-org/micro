package metadata

import (
	"context"
	"fmt"
	"net/textproto"
	"sort"
	"strings"
)

// defaultMetadataSize used when need to init new Metadata
var defaultMetadataSize = 2

// Metadata is a mapping from metadata keys to values. Users should use the following
// two convenience functions New and Pairs to generate Metadata.
type Metadata map[string][]string

// New creates an zero Metadata.
func New(l int) Metadata {
	if l == 0 {
		l = defaultMetadataSize
	}
	md := make(Metadata, l)
	return md
}

// NewWithMetadata creates an Metadata from a given key-value map.
func NewWithMetadata(m map[string]string) Metadata {
	md := make(Metadata, len(m))
	for key, val := range m {
		md[key] = append(md[key], val)
	}
	return md
}

// Pairs returns an Metadata formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
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

// Len returns the number of items in Metadata.
func (md Metadata) Len() int {
	return len(md)
}

// Copy returns a copy of Metadata.
func Copy(src Metadata) Metadata {
	out := make(Metadata, len(src))
	for k, v := range src {
		out[k] = copyOf(v)
	}
	return out
}

// Copy returns a copy of Metadata.
func (md Metadata) Copy() Metadata {
	out := make(Metadata, len(md))
	for k, v := range md {
		out[k] = copyOf(v)
	}
	return out
}

// AsMap returns a copy of Metadata with map[string]string.
func (md Metadata) AsMap() map[string]string {
	out := make(map[string]string, len(md))
	for k, v := range md {
		out[k] = strings.Join(v, ",")
	}
	return out
}

// AsHTTP1 returns a copy of Metadata
// with CanonicalMIMEHeaderKey.
func (md Metadata) AsHTTP1() map[string][]string {
	out := make(map[string][]string, len(md))
	for k, v := range md {
		out[textproto.CanonicalMIMEHeaderKey(k)] = copyOf(v)
	}
	return out
}

// AsHTTP1 returns a copy of Metadata
// with strings.ToLower.
func (md Metadata) AsHTTP2() map[string][]string {
	out := make(map[string][]string, len(md))
	for k, v := range md {
		out[strings.ToLower(k)] = copyOf(v)
	}
	return out
}

// CopyTo copies Metadata to out.
func (md Metadata) CopyTo(out Metadata) {
	for k, v := range md {
		out[k] = copyOf(v)
	}
}

// Get obtains the values for a given key.
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

// GetJoined obtains the values for a given key
// with joined values with "," symbol
func (md Metadata) GetJoined(k string) string {
	return strings.Join(md.Get(k), ",")
}

// Set sets the value of a given key with a slice of values.
func (md Metadata) Set(key string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	md[key] = vals
}

// Append adds the values to key k, not overwriting what was already stored at
// that key.
func (md Metadata) Append(key string, vals ...string) {
	if len(vals) == 0 {
		return
	}
	md[key] = append(md[key], vals...)
}

// Del removes the values for a given keys k.
func (md Metadata) Del(k ...string) {
	for i := range k {
		delete(md, k[i])
		delete(md, strings.ToLower(k[i]))
		delete(md, textproto.CanonicalMIMEHeaderKey(k[i]))
	}
}

// Join joins any number of Metadatas into a single Metadata.
//
// The order of values for each key is determined by the order in which the Metadatas
// containing those values are presented to Join.
func Join(mds ...Metadata) Metadata {
	out := Metadata{}
	for _, Metadata := range mds {
		for k, v := range Metadata {
			out[k] = append(out[k], v...)
		}
	}
	return out
}

type (
	metadataIncomingKey struct{}
	metadataOutgoingKey struct{}
	metadataCurrentKey  struct{}
)

// NewContext creates a new context with Metadata attached. Metadata must
// not be modified after calling this function.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataCurrentKey{}, rawMetadata{md: md})
}

// NewIncomingContext creates a new context with incoming Metadata attached. Metadata must
// not be modified after calling this function.
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataIncomingKey{}, rawMetadata{md: md})
}

// NewOutgoingContext creates a new context with outgoing Metadata attached. If used
// in conjunction with AppendOutgoingContext, NewOutgoingContext will
// overwrite any previously-appended metadata. Metadata must not be modified after
// calling this function.
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataOutgoingKey{}, rawMetadata{md: md})
}

// AppendContext returns a new context with the provided kv merged
// with any existing metadata in the context. Please refer to the documentation
// of Pairs for a description of kv.
func AppendContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := ctx.Value(metadataCurrentKey{}).(rawMetadata)
	added := make([][]string, len(md.added)+1)
	copy(added, md.added)
	kvCopy := make([]string, 0, len(kv))
	for i := 0; i < len(kv); i += 2 {
		kvCopy = append(kvCopy, strings.ToLower(kv[i]), kv[i+1])
	}
	added[len(added)-1] = kvCopy
	return context.WithValue(ctx, metadataCurrentKey{}, rawMetadata{md: md.md, added: added})
}

// AppendOutgoingContext returns a new context with the provided kv merged
// with any existing metadata in the context. Please refer to the documentation
// of Pairs for a description of kv.
func AppendOutgoingContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := ctx.Value(metadataOutgoingKey{}).(rawMetadata)
	added := make([][]string, len(md.added)+1)
	copy(added, md.added)
	kvCopy := make([]string, 0, len(kv))
	for i := 0; i < len(kv); i += 2 {
		kvCopy = append(kvCopy, strings.ToLower(kv[i]), kv[i+1])
	}
	added[len(added)-1] = kvCopy
	return context.WithValue(ctx, metadataOutgoingKey{}, rawMetadata{md: md.md, added: added})
}

// FromContext returns the metadata in ctx if it exists.
func FromContext(ctx context.Context) (Metadata, bool) {
	raw, ok := ctx.Value(metadataCurrentKey{}).(rawMetadata)
	if !ok {
		return nil, false
	}
	metadataSize := len(raw.md)
	for i := range raw.added {
		metadataSize += len(raw.added[i]) / 2
	}

	out := make(Metadata, metadataSize)
	for k, v := range raw.md {
		out[k] = copyOf(v)
	}
	for _, added := range raw.added {
		if len(added)%2 == 1 {
			panic(fmt.Sprintf("metadata: FromContext got an odd number of input pairs for metadata: %d", len(added)))
		}

		for i := 0; i < len(added); i += 2 {
			out[added[i]] = append(out[added[i]], added[i+1])
		}
	}
	return out, true
}

// MustContext returns the metadata in ctx.
func MustContext(ctx context.Context) Metadata {
	md, ok := FromContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// FromIncomingContext returns the incoming metadata in ctx if it exists.
func FromIncomingContext(ctx context.Context) (Metadata, bool) {
	raw, ok := ctx.Value(metadataIncomingKey{}).(rawMetadata)
	if !ok {
		return nil, false
	}
	metadataSize := len(raw.md)
	for i := range raw.added {
		metadataSize += len(raw.added[i]) / 2
	}

	out := make(Metadata, metadataSize)
	for k, v := range raw.md {
		out[k] = copyOf(v)
	}
	for _, added := range raw.added {
		if len(added)%2 == 1 {
			panic(fmt.Sprintf("metadata: FromIncomingContext got an odd number of input pairs for metadata: %d", len(added)))
		}

		for i := 0; i < len(added); i += 2 {
			out[added[i]] = append(out[added[i]], added[i+1])
		}
	}
	return out, true
}

// MustIncomingContext returns the incoming metadata in ctx.
func MustIncomingContext(ctx context.Context) Metadata {
	md, ok := FromIncomingContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// ValueFromIncomingContext returns the metadata value corresponding to the metadata
// key from the incoming metadata if it exists. Keys are matched in a case insensitive
// manner.
func ValueFromIncomingContext(ctx context.Context, key string) []string {
	raw, ok := ctx.Value(metadataIncomingKey{}).(rawMetadata)
	if !ok {
		return nil
	}

	if v, ok := raw.md[key]; ok {
		return copyOf(v)
	}
	for k, v := range raw.md {
		// Case insensitive comparison: Metadata is a map, and there's no guarantee
		// that the Metadata attached to the context is created using our helper
		// functions.
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}

// ValueFromCurrentContext returns the metadata value corresponding to the metadata
// key from the incoming metadata if it exists. Keys are matched in a case insensitive
// manner.
func ValueFromCurrentContext(ctx context.Context, key string) []string {
	md, ok := ctx.Value(metadataCurrentKey{}).(rawMetadata)
	if !ok {
		return nil
	}

	if v, ok := md.md[key]; ok {
		return copyOf(v)
	}
	for k, v := range md.md {
		// Case insensitive comparison: Metadata is a map, and there's no guarantee
		// that the Metadata attached to the context is created using our helper
		// functions.
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}

// MustOutgoingContext returns the outgoing metadata in ctx.
func MustOutgoingContext(ctx context.Context) Metadata {
	md, ok := FromOutgoingContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// ValueFromOutgoingContext returns the metadata value corresponding to the metadata
// key from the incoming metadata if it exists. Keys are matched in a case insensitive
// manner.
func ValueFromOutgoingContext(ctx context.Context, key string) []string {
	md, ok := ctx.Value(metadataOutgoingKey{}).(rawMetadata)
	if !ok {
		return nil
	}

	if v, ok := md.md[key]; ok {
		return copyOf(v)
	}
	for k, v := range md.md {
		// Case insensitive comparison: Metadata is a map, and there's no guarantee
		// that the Metadata attached to the context is created using our helper
		// functions.
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}

func copyOf(v []string) []string {
	vals := make([]string, len(v))
	copy(vals, v)
	return vals
}

// FromOutgoingContext returns the outgoing metadata in ctx if it exists.
func FromOutgoingContext(ctx context.Context) (Metadata, bool) {
	raw, ok := ctx.Value(metadataOutgoingKey{}).(rawMetadata)
	if !ok {
		return nil, false
	}

	metadataSize := len(raw.md)
	for i := range raw.added {
		metadataSize += len(raw.added[i]) / 2
	}

	out := make(Metadata, metadataSize)
	for k, v := range raw.md {
		out[k] = copyOf(v)
	}
	for _, added := range raw.added {
		if len(added)%2 == 1 {
			panic(fmt.Sprintf("metadata: FromOutgoingContext got an odd number of input pairs for metadata: %d", len(added)))
		}

		for i := 0; i < len(added); i += 2 {
			out[added[i]] = append(out[added[i]], added[i+1])
		}
	}
	return out, ok
}

type rawMetadata struct {
	md    Metadata
	added [][]string
}

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

	if k != nil && v != nil {
		*k = iter.keys[iter.cur]
		vv := iter.md[*k]
		*v = make([]string, len(vv))
		copy(*v, vv)
		iter.cur++
	}
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
