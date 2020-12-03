// Package config is an interface for dynamic configuration.
package config

import (
	"context"
	"errors"
)

var (
	DefaultConfig Config = NewConfig()
)

var (
	// ErrWatcherStopped is returned when source watcher has been stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Config is an interface abstraction for dynamic configuration
type Config interface {
	// Init the config
	Init(opts ...Option) error
	// Options in the config
	Options() Options
	// Load config from sources
	Load(context.Context) error
	// Save config to sources
	Save(context.Context) error
	// Watch a value for changes
	//	Watch(interface{}) (Watcher, error)
	String() string
}

// Watcher is the config watcher
//type Watcher interface {
//	Next() (, error)
//	Stop() error
//}
