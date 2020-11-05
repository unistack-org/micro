// Package store is an interface for distributed data storage.
// The design document is located at https://github.com/micro/development/blob/master/design/framework/store.md
package store

import (
	"context"
	"errors"
	"time"
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
	// Read takes a single key name and optional ReadOptions. It returns matching []*Record or an error.
	Read(ctx context.Context, key string, opts ...ReadOption) ([]*Record, error)
	// Write() writes a record to the store, and returns an error if the record was not written.
	Write(ctx context.Context, r *Record, opts ...WriteOption) error
	// Delete removes the record with the corresponding key from the store.
	Delete(ctx context.Context, key string, opts ...DeleteOption) error
	// List returns any keys that match, or an empty list with no error if none matched.
	List(ctx context.Context, opts ...ListOption) ([]string, error)
	// Disconnect the store
	Disconnect(ctx context.Context) error
	// String returns the name of the implementation.
	String() string
}

// Record is an item stored or retrieved from a Store
type Record struct {
	// The key to store the record
	Key string `json:"key"`
	// The value within the record
	Value []byte `json:"value"`
	// Any associated metadata for indexing
	Metadata map[string]interface{} `json:"metadata"`
	// Time to expire a record: TODO: change to timestamp
	Expiry time.Duration `json:"expiry,omitempty"`
}
