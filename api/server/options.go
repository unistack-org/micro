package server

import (
	"crypto/tls"
	"net/http"

	"github.com/unistack-org/micro/v3/api/resolver"
	"github.com/unistack-org/micro/v3/api/server/acme"
)

// Option func
type Option func(o *Options)

// Options for api server
type Options struct {
	EnableACME   bool
	EnableCORS   bool
	ACMEProvider acme.Provider
	EnableTLS    bool
	ACMEHosts    []string
	TLSConfig    *tls.Config
	Resolver     resolver.Resolver
	Wrappers     []Wrapper
}

// NewOptions returns new Options
func NewOptions(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

type Wrapper func(h http.Handler) http.Handler

func WrapHandler(w ...Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w...)
	}
}

func EnableCORS(b bool) Option {
	return func(o *Options) {
		o.EnableCORS = b
	}
}

func EnableACME(b bool) Option {
	return func(o *Options) {
		o.EnableACME = b
	}
}

func ACMEHosts(hosts ...string) Option {
	return func(o *Options) {
		o.ACMEHosts = hosts
	}
}

func ACMEProvider(p acme.Provider) Option {
	return func(o *Options) {
		o.ACMEProvider = p
	}
}

func EnableTLS(b bool) Option {
	return func(o *Options) {
		o.EnableTLS = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func Resolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
