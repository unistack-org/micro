// Package cache is a caching interface
package cache

import "context"

// Cache is an interface for caching
type Cache interface {
	// Initialise options
	Init(...Option) error
	// Get a value
	Get(ctx context.Context, key string) (interface{}, error)
	// Set a value
	Set(ctx context.Context, key string, val interface{}) error
	// Delete a value
	Delete(ctx context.Context, key string) error
	// Name of the implementation
	String() string
}

// Options struct
type Options struct {
	Nodes   []string
	Context context.Context
}

// Option func
type Option func(o *Options)

// Nodes sets the nodes for the cache
func Nodes(v ...string) Option {
	return func(o *Options) {
		o.Nodes = v
	}
}
