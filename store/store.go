// Package store is an interface for distributed data storage.
package store

import (
	"context"
	"errors"
	"time"
)

var (
	ErrWatcherStopped = errors.New("watcher stopped")
	// ErrNotConnected is returned when a store is not connected
	ErrNotConnected = errors.New("not connected")
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
	// Name returns store name
	Name() string
	// Init initialises the store
	Init(opts ...Option) error
	// Connect is used when store needs to be connected
	Connect(ctx context.Context) error
	// Options allows you to view the current options.
	Options() Options
	// Exists check that key exists in store
	Exists(ctx context.Context, key string, opts ...ExistsOption) error
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
	// Watch returns events watcher
	Watch(ctx context.Context, opts ...WatchOption) (Watcher, error)
	// Live returns store liveness
	Live() bool
	// Ready returns store readiness
	Ready() bool
	// Health returns store health
	Health() bool
}

type (
	FuncExists func(ctx context.Context, key string, opts ...ExistsOption) error
	HookExists func(next FuncExists) FuncExists
	FuncRead   func(ctx context.Context, key string, val interface{}, opts ...ReadOption) error
	HookRead   func(next FuncRead) FuncRead
	FuncWrite  func(ctx context.Context, key string, val interface{}, opts ...WriteOption) error
	HookWrite  func(next FuncWrite) FuncWrite
	FuncDelete func(ctx context.Context, key string, opts ...DeleteOption) error
	HookDelete func(next FuncDelete) FuncDelete
	FuncList   func(ctx context.Context, opts ...ListOption) ([]string, error)
	HookList   func(next FuncList) FuncList
)

type EventType int

const (
	EventTypeUnknown = iota
	EventTypeConnect
	EventTypeDisconnect
	EventTypeOpError
)

type Event interface {
	Timestamp() time.Time
	Error() error
	Type() EventType
}

type Watcher interface {
	// Next is a blocking call
	Next() (Event, error)
	// Stop stops the watcher
	Stop()
}

type WatchOption func(*WatchOptions) error

type WatchOptions struct{}

func NewWatchOptions(opts ...WatchOption) (WatchOptions, error) {
	options := WatchOptions{}
	var err error
	for _, o := range opts {
		if err = o(&options); err != nil {
			break
		}
	}

	return options, err
}

func Watch(context.Context) (Watcher, error) {
	return nil, nil
}
