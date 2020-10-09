package micro

import (
	"context"
	"time"

	"github.com/micro/cli/v2"
	cmd "github.com/unistack-org/micro-config-cmd"
	"github.com/unistack-org/micro/v3/auth"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/config"
	"github.com/unistack-org/micro/v3/debug/profile"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/registry"
	"github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/runtime"
	"github.com/unistack-org/micro/v3/selector"
	"github.com/unistack-org/micro/v3/server"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options for micro service
type Options struct {
	Auth      auth.Auth
	Broker    broker.Broker
	Logger    logger.Logger
	Cmd       cmd.Cmd
	Config    config.Config
	Client    client.Client
	Server    server.Server
	Store     store.Store
	Registry  registry.Registry
	Router    router.Router
	Runtime   runtime.Runtime
	Transport transport.Transport
	Profile   profile.Profile

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Context:   context.Background(),
		Server:    server.DefaultServer,
		Client:    client.DefaultClient,
		Broker:    broker.DefaultBroker,
		Registry:  registry.DefaultRegistry,
		Router:    router.DefaultRouter,
		Auth:      auth.DefaultAuth,
		Logger:    logger.DefaultLogger,
		Config:    config.DefaultConfig,
		Store:     store.DefaultStore,
		Transport: transport.DefaultTransport,
		//Runtime   runtime.Runtime
		//Profile   profile.Profile
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

type Option func(*Options)

// Broker to be used for service
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
		if o.Client != nil {
			// Update Client and Server
			o.Client.Init(client.Broker(b))
		}
		if o.Server != nil {
			o.Server.Init(server.Broker(b))
		}
	}
}

func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

// Client to be used for service
func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service and for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Profile to be used for debug profile
func Profile(p profile.Profile) Option {
	return func(o *Options) {
		o.Profile = p
	}
}

// Server to be used for service
func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Store sets the store to use
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// Logger set the logger to use
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Registry sets the registry for the service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
		if o.Router != nil {
			// Update router
			o.Router.Init(router.Registry(r))
		}
		if o.Server != nil {
			// Update server
			o.Server.Init(server.Registry(r))
		}
		if o.Broker != nil {
			// Update Broker
			o.Broker.Init(broker.Registry(r))
		}
	}
}

// Tracer sets the tracer for the service
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		if o.Server != nil {
			//todo client trace
			o.Server.Init(server.Tracer(t))
		}
	}
}

// Auth sets the auth for the service
func Auth(a auth.Auth) Option {
	return func(o *Options) {
		o.Auth = a
		if o.Server != nil {
			o.Server.Init(server.Auth(a))
		}
	}
}

// Config sets the config for the service
func Config(c config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// Selector sets the selector for the service client
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		if o.Client != nil {
			o.Client.Init(client.Selector(s))
		}
	}
}

// Transport sets the transport for the service
// and the underlying components
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
		// Update Client and Server
		if o.Client != nil {
			o.Client.Init(client.Transport(t))
		}
		if o.Server != nil {
			o.Server.Init(server.Transport(t))
		}
	}
}

// Runtime sets the runtime
func Runtime(r runtime.Runtime) Option {
	return func(o *Options) {
		o.Runtime = r
	}
}

// Router sets the router
func Router(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
		// Update client
		if o.Client != nil {
			o.Client.Init(client.Router(r))
		}
	}
}

// Convenience options

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) {
		if o.Server != nil {
			o.Server.Init(server.Address(addr))
		}
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) {
		if o.Server != nil {
			o.Server.Init(server.Name(n))
		}
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		if o.Server != nil {
			o.Server.Init(server.Version(v))
		}
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		if o.Server != nil {
			o.Server.Init(server.Metadata(md))
		}
	}
}

// Flags that can be passed to service
func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		if o.Cmd != nil {
			o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
		}
	}
}

// Action can be used to parse user provided cli options
func Action(a func(*cli.Context) error) Option {
	return func(o *Options) {
		if o.Cmd != nil {
			o.Cmd.App().Action = a
		}
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		if o.Server != nil {
			o.Server.Init(server.RegisterTTL(t))
		}
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		if o.Server != nil {
			o.Server.Init(server.RegisterInterval(t))
		}
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		o.Client.Init(client.WrapCall(w...))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// Init once
		o.Server.Init(wrappers...)
	}
}

// Before and Afters

// BeforeStart run funcs before service starts
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart run funcs after service starts
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
