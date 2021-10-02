package router

import (
	"context"

	"go.unistack.org/micro/v3/api/resolver"
	"go.unistack.org/micro/v3/api/resolver/vpath"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/register"
)

// Options holds the options for api router
type Options struct {
	// Register for service lookup
	Register register.Register
	// Resolver to use
	Resolver resolver.Resolver
	// Logger micro logger
	Logger logger.Logger
	// Context is for external options
	Context context.Context
	// Handler name
	Handler string
}

// Option func signature
type Option func(o *Options)

// NewOptions returns options struct filled by opts
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
