package micro

import (
	"context"
	"fmt"
	"time"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/config"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/server"
	"go.unistack.org/micro/v3/store"
	"go.unistack.org/micro/v3/tracer"
)

// Options for micro service
type Options struct {
	// Context holds external options or cancel stuff
	Context context.Context
	// Metadata holds service metadata
	Metadata metadata.Metadata
	// Version holds service version
	Version string
	// Name holds service name
	Name string
	// Brokers holds brokers
	Brokers []broker.Broker
	// Loggers holds loggers
	Loggers []logger.Logger
	// Meters holds meter
	Meters []meter.Meter
	// Configs holds config
	Configs []config.Config
	// Clients holds clients
	Clients []client.Client
	// Stores holds stores
	Stores []store.Store
	// Registers holds registers
	Registers []register.Register
	// Tracers holds tracers
	Tracers []tracer.Tracer
	// Routers holds routers
	Routers []router.Router
	// BeforeStart holds funcs that runs before service starts
	BeforeStart []func(context.Context) error
	// BeforeStop holds funcs that runs before service stops
	BeforeStop []func(context.Context) error
	// AfterStart holds funcs that runs after service starts
	AfterStart []func(context.Context) error
	// AfterStop holds funcs that runs after service stops
	AfterStop []func(context.Context) error
	// Servers holds servers
	Servers []server.Server
}

// NewOptions returns new Options filled with defaults and overrided by provided opts
func NewOptions(opts ...Option) Options {
	options := Options{
		Context:   context.Background(),
		Servers:   []server.Server{server.DefaultServer},
		Clients:   []client.Client{client.DefaultClient},
		Brokers:   []broker.Broker{broker.DefaultBroker},
		Registers: []register.Register{register.DefaultRegister},
		Routers:   []router.Router{router.DefaultRouter},
		Loggers:   []logger.Logger{logger.DefaultLogger},
		Tracers:   []tracer.Tracer{tracer.DefaultTracer},
		Meters:    []meter.Meter{meter.DefaultMeter},
		Configs:   []config.Config{config.DefaultConfig},
		Stores:    []store.Store{store.DefaultStore},
		// Runtime   runtime.Runtime
		// Profile   profile.Profile
	}

	for _, o := range opts {
		//nolint:errcheck
		o(&options)
	}

	return options
}

// Option func
type Option func(*Options) error

// Broker to be used for client and server
func Broker(b broker.Broker, opts ...BrokerOption) Option {
	return func(o *Options) error {
		var err error
		bopts := brokerOptions{}
		for _, opt := range opts {
			opt(&bopts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, srv := range o.Servers {
			for _, os := range bopts.servers {
				if srv.Name() == os || all {
					if err = srv.Init(server.Broker(b)); err != nil {
						return err
					}
				}
			}
		}
		for _, cli := range o.Clients {
			for _, oc := range bopts.clients {
				if cli.Name() == oc || all {
					if err = cli.Init(client.Broker(b)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

type brokerOptions struct {
	servers []string
	clients []string
}

// BrokerOption func signature
type BrokerOption func(*brokerOptions)

// BrokerClient specifies clients for broker
func BrokerClient(n string) BrokerOption {
	return func(o *brokerOptions) {
		o.clients = append(o.clients, n)
	}
}

// BrokerServer specifies servers for broker
func BrokerServer(n string) BrokerOption {
	return func(o *brokerOptions) {
		o.servers = append(o.servers, n)
	}
}

// Client to be used for service
func Client(c ...client.Client) Option {
	return func(o *Options) error {
		o.Clients = c
		return nil
	}
}

// Clients to be used for service
func Clients(c ...client.Client) Option {
	return func(o *Options) error {
		o.Clients = c
		return nil
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service and for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) error {
		// TODO: Pass context to underline stuff ?
		o.Context = ctx
		return nil
	}
}

/*
// Profile to be used for debug profile
func Profile(p profile.Profile) Option {
	return func(o *Options) {
		o.Profile = p
	}
}
*/

// Server to be used for service
func Server(s ...server.Server) Option {
	return func(o *Options) error {
		o.Servers = s
		return nil
	}
}

// Servers to be used for service
func Servers(s ...server.Server) Option {
	return func(o *Options) error {
		o.Servers = s
		return nil
	}
}

// Store sets the store to use
func Store(s ...store.Store) Option {
	return func(o *Options) error {
		o.Stores = s
		return nil
	}
}

// Stores sets the store to use
func Stores(s ...store.Store) Option {
	return func(o *Options) error {
		o.Stores = s
		return nil
	}
}

// Logger set the logger to use
//nolint:gocyclo
func Logger(l logger.Logger, opts ...LoggerOption) Option {
	return func(o *Options) error {
		var err error
		lopts := loggerOptions{}
		for _, opt := range opts {
			opt(&lopts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, srv := range o.Servers {
			for _, os := range lopts.servers {
				if srv.Name() == os || all {
					if err = srv.Init(server.Logger(l)); err != nil {
						return err
					}
				}
			}
		}
		for _, cli := range o.Clients {
			for _, oc := range lopts.clients {
				if cli.Name() == oc || all {
					if err = cli.Init(client.Logger(l)); err != nil {
						return err
					}
				}
			}
		}
		for _, brk := range o.Brokers {
			for _, ob := range lopts.brokers {
				if brk.Name() == ob || all {
					if err = brk.Init(broker.Logger(l)); err != nil {
						return err
					}
				}
			}
		}
		for _, reg := range o.Registers {
			for _, or := range lopts.registers {
				if reg.Name() == or || all {
					if err = reg.Init(register.Logger(l)); err != nil {
						return err
					}
				}
			}
		}
		for _, str := range o.Stores {
			for _, or := range lopts.stores {
				if str.Name() == or || all {
					if err = str.Init(store.Logger(l)); err != nil {
						return err
					}
				}
			}
		}
		for _, mtr := range o.Meters {
			for _, or := range lopts.meters {
				if mtr.Name() == or || all {
					if err = mtr.Init(meter.Logger(l)); err != nil {
						return err
					}
				}
			}
		}
		for _, trc := range o.Tracers {
			for _, ot := range lopts.tracers {
				if trc.Name() == ot || all {
					if err = trc.Init(tracer.Logger(l)); err != nil {
						return err
					}
				}
			}
		}

		return nil
	}
}

// LoggerOption func signature
type LoggerOption func(*loggerOptions)

// loggerOptions
type loggerOptions struct {
	servers   []string
	clients   []string
	brokers   []string
	registers []string
	stores    []string
	meters    []string
	tracers   []string
}

/*
func LoggerServer(n string) LoggerOption {

}
*/

// Meter set the meter to use
func Meter(m ...meter.Meter) Option {
	return func(o *Options) error {
		o.Meters = m
		return nil
	}
}

// Meters set the meter to use
func Meters(m ...meter.Meter) Option {
	return func(o *Options) error {
		o.Meters = m
		return nil
	}
}

// Register sets the register for the service
// and the underlying components
//nolint:gocyclo
func Register(r register.Register, opts ...RegisterOption) Option {
	return func(o *Options) error {
		var err error
		ropts := registerOptions{}
		for _, opt := range opts {
			opt(&ropts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, rtr := range o.Routers {
			for _, os := range ropts.routers {
				if rtr.Name() == os || all {
					if err = rtr.Init(router.Register(r)); err != nil {
						return err
					}
				}
			}
		}
		for _, srv := range o.Servers {
			for _, os := range ropts.servers {
				if srv.Name() == os || all {
					if err = srv.Init(server.Register(r)); err != nil {
						return err
					}
				}
			}
		}
		for _, brk := range o.Brokers {
			for _, os := range ropts.brokers {
				if brk.Name() == os || all {
					if err = brk.Init(broker.Register(r)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

type registerOptions struct {
	routers []string
	servers []string
	brokers []string
}

// RegisterOption func signature
type RegisterOption func(*registerOptions)

// RegisterRouter speciefies routers for register
func RegisterRouter(n string) RegisterOption {
	return func(o *registerOptions) {
		o.routers = append(o.routers, n)
	}
}

// RegisterServer specifies servers for register
func RegisterServer(n string) RegisterOption {
	return func(o *registerOptions) {
		o.servers = append(o.servers, n)
	}
}

// RegisterBroker specifies broker for register
func RegisterBroker(n string) RegisterOption {
	return func(o *registerOptions) {
		o.brokers = append(o.brokers, n)
	}
}

// Tracer sets the tracer
//nolint:gocyclo
func Tracer(t tracer.Tracer, opts ...TracerOption) Option {
	return func(o *Options) error {
		var err error
		topts := tracerOptions{}
		for _, opt := range opts {
			opt(&topts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, srv := range o.Servers {
			for _, os := range topts.servers {
				if srv.Name() == os || all {
					if err = srv.Init(server.Tracer(t)); err != nil {
						return err
					}
				}
			}
		}
		for _, cli := range o.Clients {
			for _, os := range topts.clients {
				if cli.Name() == os || all {
					if err = cli.Init(client.Tracer(t)); err != nil {
						return err
					}
				}
			}
		}
		for _, str := range o.Stores {
			for _, os := range topts.stores {
				if str.Name() == os || all {
					if err = str.Init(store.Tracer(t)); err != nil {
						return err
					}
				}
			}
		}
		for _, brk := range o.Brokers {
			for _, os := range topts.brokers {
				if brk.Name() == os || all {
					if err = brk.Init(broker.Tracer(t)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

type tracerOptions struct {
	clients []string
	servers []string
	brokers []string
	stores  []string
}

// TracerOption func signature
type TracerOption func(*tracerOptions)

// TracerClient sets the clients for tracer
func TracerClient(n string) TracerOption {
	return func(o *tracerOptions) {
		o.clients = append(o.clients, n)
	}
}

// TracerServer sets the servers for tracer
func TracerServer(n string) TracerOption {
	return func(o *tracerOptions) {
		o.servers = append(o.servers, n)
	}
}

// TracerBroker sets the broker for tracer
func TracerBroker(n string) TracerOption {
	return func(o *tracerOptions) {
		o.brokers = append(o.brokers, n)
	}
}

// TracerStore sets the store for tracer
func TracerStore(n string) TracerOption {
	return func(o *tracerOptions) {
		o.stores = append(o.stores, n)
	}
}

// Config sets the config for the service
func Config(c ...config.Config) Option {
	return func(o *Options) error {
		o.Configs = c
		return nil
	}
}

// Configs sets the configs for the service
func Configs(c ...config.Config) Option {
	return func(o *Options) error {
		o.Configs = c
		return nil
	}
}

/*
// Selector sets the selector for the service client
func Selector(s selector.Selector) Option {
	return func(o *Options) error {
		if o.Client != nil {
			o.Client.Init(client.Selector(s))
		}
		return nil
	}
}
*/
/*
// Runtime sets the runtime
func Runtime(r runtime.Runtime) Option {
	return func(o *Options) {
		o.Runtime = r
	}
}
*/

// Router sets the router
func Router(r router.Router, opts ...RouterOption) Option {
	return func(o *Options) error {
		var err error
		ropts := routerOptions{}
		for _, opt := range opts {
			opt(&ropts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, cli := range o.Clients {
			for _, os := range ropts.clients {
				if cli.Name() == os || all {
					if err = cli.Init(client.Router(r)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

type routerOptions struct {
	clients []string
}

// RouterOption func signature
type RouterOption func(*routerOptions)

// RouterClient sets the clients for router
func RouterClient(n string) RouterOption {
	return func(o *routerOptions) {
		o.clients = append(o.clients, n)
	}
}

// Address sets the address of the server
func Address(addr string) Option {
	return func(o *Options) error {
		switch len(o.Servers) {
		case 0:
			return fmt.Errorf("cant set address on nil server")
		case 1:
			break
		default:
			return fmt.Errorf("cant set same address for multiple servers")
		}
		return o.Servers[0].Init(server.Address(addr))
	}
}

// Name of the service
func Name(n string) Option {
	return func(o *Options) error {
		o.Name = n
		return nil
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) error {
		o.Version = v
		return nil
	}
}

// Metadata associated with the service
func Metadata(md metadata.Metadata) Option {
	return func(o *Options) error {
		o.Metadata = metadata.Copy(md)
		return nil
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(td time.Duration, opts ...RegisterOption) Option {
	return func(o *Options) error {
		var err error
		ropts := registerOptions{}
		for _, opt := range opts {
			opt(&ropts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, srv := range o.Servers {
			for _, os := range ropts.servers {
				if srv.Name() == os || all {
					if err = srv.Init(server.RegisterTTL(td)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(td time.Duration, opts ...RegisterOption) Option {
	return func(o *Options) error {
		var err error
		ropts := registerOptions{}
		for _, opt := range opts {
			opt(&ropts)
		}
		all := false
		if len(opts) == 0 {
			all = true
		}
		for _, srv := range o.Servers {
			for _, os := range ropts.servers {
				if srv.Name() == os || all {
					if err = srv.Init(server.RegisterInterval(td)); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

// BeforeStart run funcs before service starts
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *Options) error {
		o.BeforeStart = append(o.BeforeStart, fn)
		return nil
	}
}

// BeforeStop run funcs before service stops
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *Options) error {
		o.BeforeStop = append(o.BeforeStop, fn)
		return nil
	}
}

// AfterStart run funcs after service starts
func AfterStart(fn func(context.Context) error) Option {
	return func(o *Options) error {
		o.AfterStart = append(o.AfterStart, fn)
		return nil
	}
}

// AfterStop run funcs after service stops
func AfterStop(fn func(context.Context) error) Option {
	return func(o *Options) error {
		o.AfterStop = append(o.AfterStop, fn)
		return nil
	}
}
