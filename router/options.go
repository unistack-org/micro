package router

import (
	"context"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/util/id"
)

// Options are router options
type Options struct {
	Logger   logger.Logger
	Context  context.Context
	Register register.Register
	Name     string
	Gateway  string
	Network  string
	ID       string
	Address  string
	Precache bool
}

// ID sets Router Id
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// Address sets router service address
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Gateway sets network gateway
func Gateway(g string) Option {
	return func(o *Options) {
		o.Gateway = g
	}
}

// Network sets router network
func Network(n string) Option {
	return func(o *Options) {
		o.Network = n
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Register sets the local register
func Register(r register.Register) Option {
	return func(o *Options) {
		o.Register = r
	}
}

// Precache the routes
func Precache() Option {
	return func(o *Options) {
		o.Precache = true
	}
}

// Name of the router
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// NewOptions returns router default options
func NewOptions(opts ...Option) Options {
	options := Options{
		ID:       id.Must(),
		Network:  DefaultNetwork,
		Register: register.DefaultRegister,
		Logger:   logger.DefaultLogger,
		Context:  context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
