// Package config is an interface for dynamic configuration.
package config

import (
	"context"
	"errors"
)

// DefaultConfig default config
var DefaultConfig Config = NewConfig()

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
	// Name returns name of config
	Name() string
	// Init the config
	Init(opts ...Option) error
	// Options in the config
	Options() Options
	// Load config from sources
	Load(context.Context, ...LoadOption) error
	// Save config to sources
	Save(context.Context, ...SaveOption) error
	// Watch a value for changes
	//Watch(context.Context) (Watcher, error)
	// String returns config type name
	String() string
}

// Watcher is the config watcher
type Watcher interface {
	// Next() (, error)
	Stop() error
}

// Load loads config from config sources
func Load(ctx context.Context, cs []Config, opts ...LoadOption) error {
	var err error
	for _, c := range cs {
		if err = c.Init(); err != nil {
			return err
		}
		if err = c.Load(ctx, opts...); err != nil {
			return err
		}
	}
	return nil
}
