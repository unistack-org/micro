package http

import "go.unistack.org/micro/v3/router"

// Options struct
type Options struct {
	Router router.Router
}

// Option func
type Option func(*Options)

// WithRouter sets the router.Router option
func WithRouter(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}
