// Package store is an interface for distributed data storage.
package store // import "go.unistack.org/micro/v4/store"

import (
	"context"
	"errors"

	"go.unistack.org/micro/v4/options"
)

var (
	// ErrNotFound is returned when a key doesn't exist
	ErrNotFound = errors.New("not found")
	// ErrInvalidKey is returned when a key has empty or have invalid format
	ErrInvalidKey = errors.New("invalid key")
	// DefaultStore is the global default store
	DefaultStore = NewStore()
	// DefaultSeparator is the gloabal default key parts separator
	DefaultSeparator = "/"
)

// Store is a data storage interface
type Store interface {
	Name() string
	// Init initialises the store
	Init(opts ...options.Option) error
	// Connect is used when store needs to be connected
	Connect(ctx context.Context) error
	// Options allows you to view the current options.
	Options() Options
	// Exists check that key exists in store
	Exists(ctx context.Context, key string, opts ...options.Option) error
	// Read reads a single key name to provided value with optional options
	Read(ctx context.Context, key string, val interface{}, opts ...options.Option) error
	// Write writes a value to key name to the store with optional options
	Write(ctx context.Context, key string, val interface{}, opts ...options.Option) error
	// Delete removes the record with the corresponding key from the store with optional options
	Delete(ctx context.Context, key string, opts ...options.Option) error
	// List returns any keys that match, or an empty list with no error if none matched with optional options
	List(ctx context.Context, opts ...options.Option) ([]string, error)
	// Disconnect the store
	Disconnect(ctx context.Context) error
	// String returns the name of the implementation.
	String() string
}
