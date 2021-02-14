package register

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options holds options for register
type Options struct {
	Name      string
	Addrs     []string
	Timeout   time.Duration
	TLSConfig *tls.Config

	// Logger that will be used
	Logger logger.Logger
	// Meter that will be used
	Meter meter.Meter
	// Tracer
	Tracer tracer.Tracer
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// NewOptions returns options that filled by opts
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// RegisterOptions holds options for register method
type RegisterOptions struct {
	TTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
	// Domain to register the service in
	Domain string
	// Attempts specify attempts for register
	Attempts int
}

// NewRegisterOptions returns register options struct filled by opts
func NewRegisterOptions(opts ...RegisterOption) RegisterOptions {
	options := RegisterOptions{
		Domain:  DefaultDomain,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WatchOptions holds watch options
type WatchOptions struct {
	// Specify a service to watch
	// If blank, the watch is for all services
	Service string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
	// Domain to watch
	Domain string
}

// NewWatchOptions returns watch options filled by opts
func NewWatchOptions(opts ...WatchOption) WatchOptions {
	options := WatchOptions{
		Domain:  DefaultDomain,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// DeregisterOptions holds options for deregister method
type DeregisterOptions struct {
	Context context.Context
	// Domain the service was registered in
	Domain string
	// Atempts specify max attempts for deregister
	Attempts int
}

// NewDeregisterOptions returns options for deregister filled by opts
func NewDeregisterOptions(opts ...DeregisterOption) DeregisterOptions {
	options := DeregisterOptions{
		Domain:  DefaultDomain,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// LookupOptions holds lookup options
type LookupOptions struct {
	Context context.Context
	// Domain to scope the request to
	Domain string
}

// NewLookupOptions returns lookup options filled by opts
func NewLookupOptions(opts ...LookupOption) LookupOptions {
	options := LookupOptions{
		Domain:  DefaultDomain,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ListOptions holds the list options for list method
type ListOptions struct {
	Context context.Context
	// Domain to scope the request to
	Domain string
}

// NewListOptions returns list options filled by opts
func NewListOptions(opts ...ListOption) ListOptions {
	options := ListOptions{
		Domain:  DefaultDomain,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Addrs is the register addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Timeout sets the timeout
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Logger sets the logger
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

// Tracer sets the tracer
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Context sets the context
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// TLSConfig Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// RegisterAttempts specifies register atempts count
func RegisterAttempts(t int) RegisterOption {
	return func(o *RegisterOptions) {
		o.Attempts = t
	}
}

// RegisterTTL specifies register ttl
func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// RegisterContext sets the register context
func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// RegisterDomain secifies register domain
func RegisterDomain(d string) RegisterOption {
	return func(o *RegisterOptions) {
		o.Domain = d
	}
}

// WatchService name
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

// WatchContext sets the context for watch method
func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

// WatchDomain sets the domain for watch
func WatchDomain(d string) WatchOption {
	return func(o *WatchOptions) {
		o.Domain = d
	}
}

// DeregisterAttempts specifies deregister atempts count
func DeregisterAttempts(t int) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Attempts = t
	}
}

// DeregisterContext sets the context for deregister method
func DeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}

// DeregisterDomain specifies deregister domain
func DeregisterDomain(d string) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Domain = d
	}
}

// LookupContext sets the context for lookup method
func LookupContext(ctx context.Context) LookupOption {
	return func(o *LookupOptions) {
		o.Context = ctx
	}
}

// LookupDomain sets the domain for lookup
func LookupDomain(d string) LookupOption {
	return func(o *LookupOptions) {
		o.Domain = d
	}
}

// ListContext specifies context for list method
func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}

// ListDomain sets the domain for list method
func ListDomain(d string) ListOption {
	return func(o *ListOptions) {
		o.Domain = d
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
