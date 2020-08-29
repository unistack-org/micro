package router

import (
	"context"

	"github.com/google/uuid"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/registry"
)

// Options are router options
type Options struct {
	// Id is router id
	Id string
	// Address is router address
	Address string
	// Gateway is network gateway
	Gateway string
	// Network is network address
	Network string
	// Registry is the local registry
	Registry registry.Registry
	// Precache routes
	Precache bool
	// Logger
	Logger logger.Logger
	// Context for additional options
	Context context.Context
}

// Id sets Router Id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
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

// Registry sets the local registry
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Precache the routes
func Precache() Option {
	return func(o *Options) {
		o.Precache = true
	}
}

// DefaultOptions returns router default options
func DefaultOptions() Options {
	return Options{
		Id:      uuid.New().String(),
		Network: DefaultNetwork,
		Context: context.Background(),
	}
}
