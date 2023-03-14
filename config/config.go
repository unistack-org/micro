// Package config is an interface for dynamic configuration.
package config // import "go.unistack.org/micro/v3/config"

import (
	"context"
	"errors"
	"reflect"
	"time"
)

type Validator interface {
	Validate() error
}

// DefaultConfig default config
var DefaultConfig = NewConfig()

// DefaultWatcherMinInterval default min interval for poll changes
var DefaultWatcherMinInterval = 5 * time.Second

// DefaultWatcherMaxInterval default max interval for poll changes
var DefaultWatcherMaxInterval = 9 * time.Second

var (
	// ErrCodecMissing is returned when codec needed and not specified
	ErrCodecMissing = errors.New("codec missing")
	// ErrInvalidStruct is returned when the target struct is invalid
	ErrInvalidStruct = errors.New("invalid struct specified")
	// ErrWatcherStopped is returned when source watcher has been stopped
	ErrWatcherStopped = errors.New("watcher stopped")
	// ErrWatcherNotImplemented returned when config does not implement watch
	ErrWatcherNotImplemented = errors.New("watcher not implemented")
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
	// Watch a config for changes
	Watch(context.Context, ...WatchOption) (Watcher, error)
	// String returns config type name
	String() string
}

// Watcher is the config watcher
type Watcher interface {
	// Next blocks until update happens or error returned
	Next() (map[string]interface{}, error)
	// Stop stops watcher
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

// Validate runs Validate() error func for each struct field
func Validate(ctx context.Context, cfg interface{}) error {
	if cfg == nil {
		return nil
	}

	if v, ok := cfg.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	sv := reflect.ValueOf(cfg)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		return nil
	}

	typ := sv.Type()
	for idx := 0; idx < typ.NumField(); idx++ {
		fld := typ.Field(idx)
		val := sv.Field(idx)
		if !val.IsValid() || len(fld.PkgPath) != 0 {
			continue
		}

		if v, ok := val.Interface().(Validator); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}

		switch val.Kind() {
		case reflect.Ptr:
			if reflect.Indirect(val).Kind() == reflect.Struct {
				if err := Validate(ctx, val.Interface()); err != nil {
					return err
				}
			}
		case reflect.Struct:
			if err := Validate(ctx, val.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

var (
	// DefaultBeforeLoad default func that runs before config Load
	DefaultBeforeLoad = func(ctx context.Context, c Config) error {
		for _, fn := range c.Options().BeforeLoad {
			if err := fn(ctx, c); err != nil {
				c.Options().Logger.Errorf(ctx, "%s BeforeLoad err: %v", c.String(), err)
				if !c.Options().AllowFail {
					return err
				}
			}
		}
		return nil
	}
	// DefaultAfterLoad default func that runs after config Load
	DefaultAfterLoad = func(ctx context.Context, c Config) error {
		for _, fn := range c.Options().AfterLoad {
			if err := fn(ctx, c); err != nil {
				c.Options().Logger.Errorf(ctx, "%s AfterLoad err: %v", c.String(), err)
				if !c.Options().AllowFail {
					return err
				}
			}
		}
		return nil
	}
	// DefaultBeforeSave default func that runs befora config Save
	DefaultBeforeSave = func(ctx context.Context, c Config) error {
		for _, fn := range c.Options().BeforeSave {
			if err := fn(ctx, c); err != nil {
				c.Options().Logger.Errorf(ctx, "%s BeforeSave err: %v", c.String(), err)
				if !c.Options().AllowFail {
					return err
				}
			}
		}
		return nil
	}
	// DefaultAfterSave default func that runs after config Save
	DefaultAfterSave = func(ctx context.Context, c Config) error {
		for _, fn := range c.Options().AfterSave {
			if err := fn(ctx, c); err != nil {
				c.Options().Logger.Errorf(ctx, "%s AfterSave err: %v", c.String(), err)
				if !c.Options().AllowFail {
					return err
				}
			}
		}
		return nil
	}
	// DefaultBeforeInit default func that runs befora config Init
	DefaultBeforeInit = func(ctx context.Context, c Config) error {
		for _, fn := range c.Options().BeforeInit {
			if err := fn(ctx, c); err != nil {
				c.Options().Logger.Errorf(ctx, "%s BeforeInit err: %v", c.String(), err)
				if !c.Options().AllowFail {
					return err
				}
			}
		}
		return nil
	}
	// DefaultAfterInit default func that runs after config Init
	DefaultAfterInit = func(ctx context.Context, c Config) error {
		for _, fn := range c.Options().AfterSave {
			if err := fn(ctx, c); err != nil {
				c.Options().Logger.Errorf(ctx, "%s AfterInit err: %v", c.String(), err)
				if !c.Options().AllowFail {
					return err
				}
			}
		}
		return nil
	}
)
