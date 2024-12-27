package register

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/tracer"
)

// Options holds options for register
type Options struct {
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Context holds external options
	Context context.Context
	// Logged used for logging
	Logger logger.Logger
	// Meter used for metrics
	Meter meter.Meter
	// TLSConfig holds tls.TLSConfig options
	TLSConfig *tls.Config
	// Name holds the name of register
	Name string
	// Addrs specifies register addrs
	Addrs []string
	// Timeout specifies timeout
	Timeout time.Duration
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
type RegisterOptions struct { // nolint: golint,revive
	Context   context.Context
	Namespace string
	TTL       time.Duration
	Attempts  int
}

// NewRegisterOptions returns register options struct filled by opts
func NewRegisterOptions(opts ...RegisterOption) RegisterOptions {
	options := RegisterOptions{
		Namespace: DefaultNamespace,
		Context:   context.Background(),
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
	// Namespace to watch
	Namespace string
}

// NewWatchOptions returns watch options filled by opts
func NewWatchOptions(opts ...WatchOption) WatchOptions {
	options := WatchOptions{
		Namespace: DefaultNamespace,
		Context:   context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// DeregisterOptions holds options for deregister method
type DeregisterOptions struct {
	Context context.Context
	// Namespace the service was registered in
	Namespace string
	// Atempts specify max attempts for deregister
	Attempts int
}

// NewDeregisterOptions returns options for deregister filled by opts
func NewDeregisterOptions(opts ...DeregisterOption) DeregisterOptions {
	options := DeregisterOptions{
		Namespace: DefaultNamespace,
		Context:   context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// LookupOptions holds lookup options
type LookupOptions struct {
	Context context.Context
	// Namespace to scope the request to
	Namespace string
}

// NewLookupOptions returns lookup options filled by opts
func NewLookupOptions(opts ...LookupOption) LookupOptions {
	options := LookupOptions{
		Namespace: DefaultNamespace,
		Context:   context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ListOptions holds the list options for list method
type ListOptions struct {
	// Context used to store additional options
	Context context.Context
	// Namespace to scope the request to
	Namespace string
	// Name filter services by name
	Name string
}

// NewListOptions returns list options filled by opts
func NewListOptions(opts ...ListOption) ListOptions {
	options := ListOptions{
		Namespace: DefaultNamespace,
		Context:   context.Background(),
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
func RegisterAttempts(t int) RegisterOption { // nolint: golint,revive
	return func(o *RegisterOptions) {
		o.Attempts = t
	}
}

// RegisterTTL specifies register ttl
func RegisterTTL(t time.Duration) RegisterOption { // nolint: golint,revive
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

// RegisterContext sets the register context
func RegisterContext(ctx context.Context) RegisterOption { // nolint: golint,revive
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

// RegisterNamespace secifies register Namespace
func RegisterNamespace(d string) RegisterOption { // nolint: golint,revive
	return func(o *RegisterOptions) {
		o.Namespace = d
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

// WatchNamespace sets the Namespace for watch
func WatchNamespace(d string) WatchOption {
	return func(o *WatchOptions) {
		o.Namespace = d
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

// DeregisterNamespace specifies deregister Namespace
func DeregisterNamespace(d string) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Namespace = d
	}
}

// LookupContext sets the context for lookup method
func LookupContext(ctx context.Context) LookupOption {
	return func(o *LookupOptions) {
		o.Context = ctx
	}
}

// LookupNamespace sets the Namespace for lookup
func LookupNamespace(d string) LookupOption {
	return func(o *LookupOptions) {
		o.Namespace = d
	}
}

// ListContext specifies context for list method
func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}

// ListNamespace sets the Namespace for list method
func ListNamespace(d string) ListOption {
	return func(o *ListOptions) {
		o.Namespace = d
	}
}

// ListName sets the name for list method to filter needed services
func ListName(n string) ListOption {
	return func(o *ListOptions) {
		o.Name = n
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
