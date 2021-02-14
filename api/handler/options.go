package handler

import (
	"github.com/unistack-org/micro/v3/api/router"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
)

var (
	// DefaultMaxRecvSize specifies max recv size for handler
	DefaultMaxRecvSize int64 = 1024 * 1024 * 100 // 10Mb
)

// Options struct holds handler options
type Options struct {
	MaxRecvSize int64
	Namespace   string
	Router      router.Router
	Client      client.Client
	Logger      logger.Logger
}

// Option func signature
type Option func(o *Options)

// NewOptions creates new options struct and fills it
func NewOptions(opts ...Option) Options {
	options := Options{
		Client:      client.DefaultClient,
		Router:      router.DefaultRouter,
		Logger:      logger.DefaultLogger,
		MaxRecvSize: DefaultMaxRecvSize,
	}
	for _, o := range opts {
		o(&options)
	}

	// set namespace if blank
	if len(options.Namespace) == 0 {
		WithNamespace("go.micro.api")(&options)
	}

	return options
}

// WithNamespace specifies the namespace for the handler
func WithNamespace(s string) Option {
	return func(o *Options) {
		o.Namespace = s
	}
}

// WithRouter specifies a router to be used by the handler
func WithRouter(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// WithClient specifies client to be used by the handler
func WithClient(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// WithMaxRecvSize specifies max body size
func WithMaxRecvSize(size int64) Option {
	return func(o *Options) {
		o.MaxRecvSize = size
	}
}
