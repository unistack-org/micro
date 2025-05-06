package metadata

import (
	"context"
	"fmt"
	"strings"
)

// In the metadata package, context and metadata are treated as immutable.
// Deep copies of metadata are made to keep things safe and correct.
// If a user takes a map and changes it across threads, it's their responsibility.
//
// 1. Incoming Context
//
// This context is provided by an external system and populated by the server or broker of the micro framework.
// It should not be modified. The idea is to extract all necessary data from it,
// validate the data, and transfer it into the current context.
// After that, only the current context should be used throughout the code.
//
// 2. Current Context
//
// This is the context used during the execution flow.
// You can add any needed metadata to it and pass it through your code.
//
// 3. Outgoing Context
//
// This context is for sending data to external systems.
// You can add what you need before sending it out.
// But it’s usually better to build and prepare this context right before making the external call,
// instead of changing it in many places.
//
// Execution Flow:
//
// [External System]
//       ↓
// [Incoming Context]
//       ↓
// [Extract & Validate Metadata from Incoming Context]
//       ↓
// [Prepare Current Context]
//       ↓
// [Enrich Current Context]
//       ↓
// [Business Logic]
//       ↓
// [Prepare Outgoing Context]
//       ↓
// [External System Call]

type (
	metadataCurrentKey  struct{}
	metadataIncomingKey struct{}
	metadataOutgoingKey struct{}

	rawMetadata struct {
		md    Metadata
		added [][]string
	}
)

// NewContext creates a new context with the provided Metadata attached.
// The Metadata must not be modified after calling this function.
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataCurrentKey{}, rawMetadata{md: md})
}

// NewIncomingContext creates a new context with the provided incoming Metadata attached.
// The Metadata must not be modified after calling this function.
func NewIncomingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataIncomingKey{}, rawMetadata{md: md})
}

// NewOutgoingContext creates a new context with the provided outgoing Metadata attached.
// The Metadata must not be modified after calling this function.
func NewOutgoingContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metadataOutgoingKey{}, rawMetadata{md: md})
}

// AppendContext returns a new context with the provided key-value pairs (kv)
// merged with any existing metadata in the context. For a description of kv,
// please refer to the Pairs documentation.
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

// AppendOutgoingContext returns a new context with the provided key-value pairs (kv)
// merged with any existing metadata in the context. For a description of kv,
// please refer to the Pairs documentation.
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

// FromContext retrieves a deep copy of the metadata from the context and returns it
// with a boolean indicating if it was found.
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

// MustContext retrieves a deep copy of the metadata from the context and panics
// if the metadata is not found.
func MustContext(ctx context.Context) Metadata {
	md, ok := FromContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// FromIncomingContext retrieves a deep copy of the metadata from the context and returns it
// with a boolean indicating if it was found.
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

// MustIncomingContext retrieves a deep copy of the metadata from the context and panics
// if the metadata is not found.
func MustIncomingContext(ctx context.Context) Metadata {
	md, ok := FromIncomingContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// FromOutgoingContext retrieves a deep copy of the metadata from the context and returns it
// with a boolean indicating if it was found.
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

// MustOutgoingContext retrieves a deep copy of the metadata from the context and panics
// if the metadata is not found.
func MustOutgoingContext(ctx context.Context) Metadata {
	md, ok := FromOutgoingContext(ctx)
	if !ok {
		panic("missing metadata")
	}
	return md
}

// ValueFromCurrentContext retrieves a deep copy of the metadata for the given key
// from the context, performing a case-insensitive search if needed. Returns nil if not found.
func ValueFromCurrentContext(ctx context.Context, key string) []string {
	md, ok := ctx.Value(metadataCurrentKey{}).(rawMetadata)
	if !ok {
		return nil
	}

	if v, ok := md.md[key]; ok {
		return copyOf(v)
	}
	for k, v := range md.md {
		// Case-insensitive comparison: Metadata is a map, and there's no guarantee
		// that the Metadata attached to the context is created using our helper
		// functions.
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}

// ValueFromIncomingContext retrieves a deep copy of the metadata for the given key
// from the context, performing a case-insensitive search if needed. Returns nil if not found.
func ValueFromIncomingContext(ctx context.Context, key string) []string {
	raw, ok := ctx.Value(metadataIncomingKey{}).(rawMetadata)
	if !ok {
		return nil
	}

	if v, ok := raw.md[key]; ok {
		return copyOf(v)
	}
	for k, v := range raw.md {
		// Case-insensitive comparison: Metadata is a map, and there's no guarantee
		// that the Metadata attached to the context is created using our helper
		// functions.
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}

// ValueFromOutgoingContext retrieves a deep copy of the metadata for the given key
// from the context, performing a case-insensitive search if needed. Returns nil if not found.
func ValueFromOutgoingContext(ctx context.Context, key string) []string {
	md, ok := ctx.Value(metadataOutgoingKey{}).(rawMetadata)
	if !ok {
		return nil
	}

	if v, ok := md.md[key]; ok {
		return copyOf(v)
	}
	for k, v := range md.md {
		// Case-insensitive comparison: Metadata is a map, and there's no guarantee
		// that the Metadata attached to the context is created using our helper
		// functions.
		if strings.EqualFold(k, key) {
			return copyOf(v)
		}
	}
	return nil
}
