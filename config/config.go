// Package config is an interface for dynamic configuration.
package config

import (
	"context"
	"errors"
)

var (
	// DefaultConfig default config
	DefaultConfig Config = NewConfig()
)

var (
	// ErrCodecMissing is returned when codec needed and not specified
	ErrCodecMissing = errors.New("codec missing")
	// ErrInvalidStruct is returned when the target struct is invalid
	ErrInvalidStruct = errors.New("invalid struct specified")
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

// Load loads config from config sources
func Load(ctx context.Context, cs ...Config) error {
	var err error
	for _, c := range cs {
		if err = c.Init(); err != nil {
			return err
		}
		if err = c.Load(ctx); err != nil {
			return err
		}
	}
	return nil
}
