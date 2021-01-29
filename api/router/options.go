package router

import (
	"context"

	"github.com/unistack-org/micro/v3/api/resolver"
	"github.com/unistack-org/micro/v3/api/resolver/vpath"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/register"
)

type Options struct {
	Handler  string
	Register register.Register
	Resolver resolver.Resolver
	Logger   logger.Logger
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

// WithRegister sets the register
func WithRegister(r register.Register) Option {
	return func(o *Options) {
		o.Register = r
	}
}

// WithResolver sets the resolver
func WithResolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
