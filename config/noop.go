// Package config is an interface for dynamic configuration.
package config

import "context"

type noopConfig struct {
	opts Options
}

func (c *noopConfig) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	return nil
}

func (c *noopConfig) Load(ctx context.Context) error {
	return nil
}

func (c *noopConfig) Save(ctx context.Context) error {
	return nil
}

func (c *noopConfig) Options() Options {
	return c.opts
}

func (c *noopConfig) String() string {
	return "noop"
}

// NewConfig returns new noop config
func NewConfig(opts ...Option) Config {
	return &noopConfig{opts: NewOptions(opts...)}
}
