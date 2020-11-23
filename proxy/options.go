// Package proxy is a transparent proxy built on the micro/server
package proxy

import (
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/router"
)

// Options for proxy
type Options struct {
	// Specific endpoint to always call
	Endpoint string
	// The default client to use
	Client client.Client
	// The default router to use
	Router router.Router
	// Extra links for different clients
	Links map[string]client.Client
	// Logger
	Logger logger.Logger
}

// Option func signature
type Option func(o *Options)

// WithEndpoint sets a proxy endpoint
func WithEndpoint(e string) Option {
	return func(o *Options) {
		o.Endpoint = e
	}
}

// WithClient sets the client
func WithClient(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// WithRouter specifies the router to use
func WithRouter(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// WithLogger specifies the logger to use
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// WithLink sets a link for outbound requests
func WithLink(name string, c client.Client) Option {
	return func(o *Options) {
		if o.Links == nil {
			o.Links = make(map[string]client.Client)
		}
		o.Links[name] = c
	}
}
