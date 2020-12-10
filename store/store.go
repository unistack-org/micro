// Package store is an interface for distributed data storage.
// The design document is located at https://github.com/micro/development/blob/master/design/framework/store.md
package store

import (
	"context"
	"errors"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	// ErrNotFound is returned when a key doesn't exist
	ErrNotFound = errors.New("not found")
	// DefaultStore is the global default store
	DefaultStore Store = NewStore()
)

// Store is a data storage interface
type Store interface {
	// Init initialises the store
	Init(opts ...Option) error
	// Connect is used when store needs to be connected
	Connect(ctx context.Context) error
	// Options allows you to view the current options.
	Options() Options
	// Exists check that key exists in store
	Exists(ctx context.Context, key string) error
	// Read reads a single key name to provided value with optional ReadOptions
	Read(ctx context.Context, key string, val interface{}, opts ...ReadOption) error
	// Write writes a value to key name to the store with optional WriteOption
	Write(ctx context.Context, key string, val interface{}, opts ...WriteOption) error
	// Delete removes the record with the corresponding key from the store.
	Delete(ctx context.Context, key string, opts ...DeleteOption) error
	// List returns any keys that match, or an empty list with no error if none matched.
	List(ctx context.Context, opts ...ListOption) ([]string, error)
	// Disconnect the store
	Disconnect(ctx context.Context) error
	// String returns the name of the implementation.
	String() string
}

// Value is an item stored or retrieved from a Store
// may be used in store implementations to provide metadata
type Value struct {
	// Data holds underline struct
	Data interface{} `json:"data"`
	// Metadata associated with data for indexing
	Metadata metadata.Metadata `json:"metadata"`
}
