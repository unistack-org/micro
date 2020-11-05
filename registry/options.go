package registry

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/unistack-org/micro/v3/logger"
)

type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	Logger    logger.Logger
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

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

type DeregisterOptions struct {
	Context context.Context
	// Domain the service was registered in
	Domain string
	// Atempts specify max attempts for deregister
	Attempts int
}

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

type GetOptions struct {
	Context context.Context
	// Domain to scope the request to
	Domain string
}

func NewGetOptions(opts ...GetOption) GetOptions {
	options := GetOptions{
		Domain:  DefaultDomain,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

type ListOptions struct {
	Context context.Context
	// Domain to scope the request to
	Domain string
}

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

// Addrs is the registry addresses to use
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Secure communication with the registry
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Context sets the context
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func RegisterAttempts(t int) RegisterOption {
	return func(o *RegisterOptions) {
		o.Attempts = t
	}
}

func RegisterTTL(t time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.TTL = t
	}
}

func RegisterContext(ctx context.Context) RegisterOption {
	return func(o *RegisterOptions) {
		o.Context = ctx
	}
}

func RegisterDomain(d string) RegisterOption {
	return func(o *RegisterOptions) {
		o.Domain = d
	}
}

// Watch a service
func WatchService(name string) WatchOption {
	return func(o *WatchOptions) {
		o.Service = name
	}
}

func WatchContext(ctx context.Context) WatchOption {
	return func(o *WatchOptions) {
		o.Context = ctx
	}
}

func WatchDomain(d string) WatchOption {
	return func(o *WatchOptions) {
		o.Domain = d
	}
}

func DeregisterTimeout(t int) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Attempts = t
	}
}

func DeregisterContext(ctx context.Context) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Context = ctx
	}
}

func DeregisterDomain(d string) DeregisterOption {
	return func(o *DeregisterOptions) {
		o.Domain = d
	}
}

func GetContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

func GetDomain(d string) GetOption {
	return func(o *GetOptions) {
		o.Domain = d
	}
}

func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}

func ListDomain(d string) ListOption {
	return func(o *ListOptions) {
		o.Domain = d
	}
}
