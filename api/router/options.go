package router

import (
	"context"

	"github.com/unistack-org/micro/v3/api/resolver"
	"github.com/unistack-org/micro/v3/api/resolver/vpath"
	"github.com/unistack-org/micro/v3/registry"
)

type Options struct {
	Handler  string
	Registry registry.Registry
	Resolver resolver.Resolver
	Context  context.Context
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Context: context.Background(),
		Handler: "meta",
	}

	for _, o := range opts {
		o(&options)
	}

	if options.Resolver == nil {
		options.Resolver = vpath.NewResolver(
			resolver.WithHandler(options.Handler),
		)
	}

	return options
}

// WithContext sets the context
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// WithHandler sets the handler
func WithHandler(h string) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

// WithRegistry sets the registry
func WithRegistry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// WithResolver sets the resolver
func WithResolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
