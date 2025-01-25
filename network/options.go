package network

import (
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/network/tunnel"
	"go.unistack.org/micro/v4/proxy"
	"go.unistack.org/micro/v4/router"
	"go.unistack.org/micro/v4/tracer"
	"go.unistack.org/micro/v4/util/id"
)

// Option func
type Option func(*Options)

// Options configure network
type Options struct {
	// Router used for routing
	Router router.Router
	// Proxy holds the proxy
	Proxy proxy.Proxy
	// Logger used for logging
	Logger logger.Logger
	// Meter used for metrics
	Meter meter.Meter
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Tunnel used for transfer data
	Tunnel tunnel.Tunnel
	// ID of the node
	ID string
	// Name of the network
	Name string
	// Address to bind to
	Address string
	// Advertise sets the address to advertise
	Advertise string
	// Nodes is a list of nodes to connect to
	Nodes []string
}

// ID sets the id of the network node
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// Name sets the network name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Address sets the network address
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Advertise sets the address to advertise
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// Nodes is a list of nodes to connect to
func Nodes(n ...string) Option {
	return func(o *Options) {
		o.Nodes = n
	}
}

// Tunnel sets the network tunnel
func Tunnel(t tunnel.Tunnel) Option {
	return func(o *Options) {
		o.Tunnel = t
	}
}

// Router sets the network router
func Router(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// Proxy sets the network proxy
func Proxy(p proxy.Proxy) Option {
	return func(o *Options) {
		o.Proxy = p
	}
}

// Logger sets the network logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter sets the meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Tracer to be used for tracing
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// NewOptions returns network default options
func NewOptions(opts ...Option) Options {
	options := Options{
		ID:      id.MustNew(),
		Name:    "go.micro",
		Address: ":0",
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
