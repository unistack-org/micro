package server

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/unistack-org/micro/v3/auth"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/registry"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options server struct
type Options struct {
	Codecs       map[string]codec.Codec
	Broker       broker.Broker
	Registry     registry.Registry
	Tracer       tracer.Tracer
	Auth         auth.Auth
	Logger       logger.Logger
	Transport    transport.Transport
	Metadata     metadata.Metadata
	Name         string
	Address      string
	Advertise    string
	Id           string
	Namespace    string
	Version      string
	HdlrWrappers []HandlerWrapper
	SubWrappers  []SubscriberWrapper

	// RegisterCheck runs a check function before registering the service
	RegisterCheck func(context.Context) error
	// The register expiry time
	RegisterTTL time.Duration
	// The interval on which to register
	RegisterInterval time.Duration
	// RegisterAttempts specify how many times try to register
	RegisterAttempts int
	// DeegisterAttempts specify how many times try to deregister
	DeregisterAttempts int

	// The router for requests
	Router Router

	// TLSConfig specifies tls.Config for secure serving
	TLSConfig *tls.Config

	Wait *sync.WaitGroup
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// NewOptions returns new options struct with default or passed values
func NewOptions(opts ...Option) Options {
	options := Options{
		Auth:             auth.DefaultAuth,
		Codecs:           make(map[string]codec.Codec),
		Context:          context.Background(),
		Metadata:         metadata.New(0),
		RegisterInterval: DefaultRegisterInterval,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterCheck:    DefaultRegisterCheck,
		Logger:           logger.DefaultLogger,
		Tracer:           tracer.DefaultTracer,
		Broker:           broker.DefaultBroker,
		Registry:         registry.DefaultRegistry,
		Transport:        transport.DefaultTransport,
		Address:          DefaultAddress,
		Name:             DefaultName,
		Version:          DefaultVersion,
		Id:               DefaultId,
		Namespace:        DefaultNamespace,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Name sets the server name option
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Namespace to register handlers in
func Namespace(n string) Option {
	return func(o *Options) {
		o.Namespace = n
	}
}

// Logger sets the logger option
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Id unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Advertise the address to advertise for discovery - host:port
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// Broker to use for pub/sub
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

// Codec to use to encode/decode requests for a given content type
func Codec(contentType string, c codec.Codec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Registry used for discovery
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Tracer mechanism for distributed tracking
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Auth mechanism for role based access control
func Auth(a auth.Auth) Option {
	return func(o *Options) {
		o.Auth = a
	}
}

// Transport mechanism for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Metadata associated with the server
func Metadata(md metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = metadata.Copy(md)
	}
}

// RegisterCheck run func before registry service
func RegisterCheck(fn func(context.Context) error) Option {
	return func(o *Options) {
		o.RegisterCheck = fn
	}
}

// RegisterTTL registers service with a TTL
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// RegisterInterval registers service with at interval
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

// TLSConfig specifies a *tls.Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		// set the internal tls
		o.TLSConfig = t

		// set the default transport if one is not
		// already set. Required for Init call below.

		// set the transport tls
		o.Transport.Init(
			transport.Secure(true),
			transport.TLSConfig(t),
		)
	}
}

// WithRouter sets the request router
func WithRouter(r Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// Wait tells the server to wait for requests to finish before exiting
// If `wg` is nil, server only wait for completion of rpc handler.
// For user need finer grained control, pass a concrete `wg` here, server will
// wait against it on stop.
func Wait(wg *sync.WaitGroup) Option {
	return func(o *Options) {
		if wg == nil {
			wg = new(sync.WaitGroup)
		}
		o.Wait = wg
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w HandlerWrapper) Option {
	return func(o *Options) {
		o.HdlrWrappers = append(o.HdlrWrappers, w)
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w SubscriberWrapper) Option {
	return func(o *Options) {
		o.SubWrappers = append(o.SubWrappers, w)
	}
}

// HandlerOption func
type HandlerOption func(*HandlerOptions)

// HandlerOptions struct
type HandlerOptions struct {
	Internal bool
	Metadata map[string]metadata.Metadata
	Context  context.Context
}

// NewHandlerOptions creates new HandlerOptions
func NewHandlerOptions(opts ...HandlerOption) HandlerOptions {
	options := HandlerOptions{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// SubscriberOption func
type SubscriberOption func(*SubscriberOptions)

// SubscriberOptions struct
type SubscriberOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck  bool
	Queue    string
	Internal bool
	Context  context.Context
}

// NewSubscriberOptions create new SubscriberOptions
func NewSubscriberOptions(opts ...SubscriberOption) SubscriberOptions {
	options := SubscriberOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// EndpointMetadata is a Handler option that allows metadata to be added to
// individual endpoints.
func EndpointMetadata(name string, md metadata.Metadata) HandlerOption {
	return func(o *HandlerOptions) {
		o.Metadata[name] = metadata.Copy(md)
	}
}

// InternalHandler options specifies that a handler is not advertised
// to the discovery system. In the future this may also limit request
// to the internal network or authorised user.
func InternalHandler(b bool) HandlerOption {
	return func(o *HandlerOptions) {
		o.Internal = b
	}
}

// InternalSubscriber options specifies that a subscriber is not advertised
// to the discovery system.
func InternalSubscriber(b bool) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Internal = b
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscriberOption {
	return func(o *SubscriberOptions) {
		o.AutoAck = false
	}
}

// SubscriberQueue sets the shared queue name distributed messages across subscribers
func SubscriberQueue(n string) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Queue = n
	}
}

// SubscriberContext set context options to allow broker SubscriberOption passed
func SubscriberContext(ctx context.Context) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Context = ctx
	}
}
