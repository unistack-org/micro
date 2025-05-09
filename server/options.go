package server

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/register"
	msync "go.unistack.org/micro/v4/sync"
	"go.unistack.org/micro/v4/tracer"
	"go.unistack.org/micro/v4/util/id"
)

// Option func
type Option func(*Options)

// Options server struct
type Options struct {
	// Codecs map to handle content-type
	Codecs map[string]codec.Codec
	// Metadata holds the server metadata
	Metadata metadata.Metadata

	// ID holds the id of the server
	ID string
	// Namespace for te server
	Namespace string
	// Name holds the server name
	Name string
	// Address holds the server address
	Address string
	// Advertise holds the advertise address
	Advertise string
	// Version holds the server version
	Version string

	// Context holds the external options and can be used for server shutdown
	Context context.Context
	// Broker holds the server broker
	Broker broker.Broker
	// Register holds the register
	Register register.Register
	// Tracer holds the tracer
	Tracer tracer.Tracer
	// Logger holds the logger
	Logger logger.Logger
	// Meter holds the meter
	Meter meter.Meter
	// Listener may be passed if already created
	Listener net.Listener

	// TLSConfig specifies tls.Config for secure serving
	TLSConfig *tls.Config
	// Wait group
	Wait *msync.WaitGroup

	// RegisterCheck run before register server
	RegisterCheck func(context.Context) error

	// Hooks may contains hook actions that performs before/after server handler
	// or server subscriber handler
	Hooks options.Hooks

	// RegisterInterval holds he interval for re-register
	RegisterInterval time.Duration
	// RegisterTTL specifies TTL for register record
	RegisterTTL time.Duration
	// GracefulTimeout timeout for graceful stop server
	GracefulTimeout time.Duration
	// MaxConn limits number of connections
	MaxConn int
	// DeregisterAttempts holds the number of deregister attempts before error
	DeregisterAttempts int
	// RegisterAttempts holds the number of register attempts before error
	RegisterAttempts int
}

// NewOptions returns new options struct with default or passed values
func NewOptions(opts ...Option) Options {
	options := Options{
		Codecs:           make(map[string]codec.Codec),
		Context:          context.Background(),
		Metadata:         metadata.New(0),
		RegisterInterval: DefaultRegisterInterval,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterCheck:    DefaultRegisterCheck,
		Logger:           logger.DefaultLogger,
		Meter:            meter.DefaultMeter,
		Tracer:           tracer.DefaultTracer,
		Broker:           broker.DefaultBroker,
		Register:         register.DefaultRegister,
		Address:          DefaultAddress,
		Name:             DefaultName,
		Version:          DefaultVersion,
		ID:               id.MustNew(),
		Namespace:        DefaultNamespace,
		GracefulTimeout:  DefaultGracefulTimeout,
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

// Meter sets the meter option
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// ID unique server id
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
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

// Register used for discovery
func Register(r register.Register) Option {
	return func(o *Options) {
		o.Register = r
	}
}

// Tracer mechanism for distributed tracking
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Metadata associated with the server
func Metadata(md metadata.Metadata) Option {
	return func(o *Options) {
		o.Metadata = metadata.Copy(md)
	}
}

// RegisterCheck run func before register service
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
	}
}

/*
// WithRouter sets the request router
func WithRouter(r Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}
*/

// Wait tells the server to wait for requests to finish before exiting
// If `wg` is nil, server only wait for completion of rpc handler.
// For user need finer grained control, pass a concrete `wg` here, server will
// wait against it on stop.
func Wait(wg *sync.WaitGroup) Option {
	return func(o *Options) {
		if wg == nil {
			wg = new(sync.WaitGroup)
		}
		o.Wait = msync.WrapWaitGroup(wg)
	}
}

// MaxConn specifies maximum number of max simultaneous connections to server
func MaxConn(n int) Option {
	return func(o *Options) {
		o.MaxConn = n
	}
}

// Listener specifies the net.Listener to use instead of the default
func Listener(l net.Listener) Option {
	return func(o *Options) {
		o.Listener = l
	}
}

// HandlerOption func
type HandlerOption func(*HandlerOptions)

// GracefulTimeout duration
func GracefulTimeout(td time.Duration) Option {
	return func(o *Options) {
		o.GracefulTimeout = td
	}
}

// HandlerOptions struct
type HandlerOptions struct {
	// Context holds external options
	Context context.Context
	// Metadata for handler
	Metadata map[string]metadata.Metadata
}

// NewHandlerOptions creates new HandlerOptions
func NewHandlerOptions(opts ...HandlerOption) HandlerOptions {
	options := HandlerOptions{
		Context:  context.Background(),
		Metadata: make(map[string]metadata.Metadata),
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
	// Context holds the external options
	Context context.Context
	// Queue holds the subscription queue
	Queue string
	// BatchWait flag specifies max wait time for batch filling
	BatchWait time.Duration
	// BatchSize flag specifies max size of batch
	BatchSize int
	// AutoAck flag for auto ack messages after processing
	AutoAck bool
	// BodyOnly flag specifies that message without headers
	BodyOnly bool
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

// SubscriberGroup sets the shared group name distributed messages across subscribers
func SubscriberGroup(n string) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Queue = n
	}
}

// SubscriberBodyOnly says broker that message contains raw data with absence of micro broker.Message format
func SubscriberBodyOnly(b bool) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.BodyOnly = b
	}
}

// SubscriberContext set context options to allow broker SubscriberOption passed
func SubscriberContext(ctx context.Context) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.Context = ctx
	}
}

// SubscriberAck control auto ack processing for handler
func SubscriberAck(b bool) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.AutoAck = b
	}
}

// SubscriberBatchSize control batch filling size for handler
// Batch filling max waiting time controlled by SubscriberBatchWait
func SubscriberBatchSize(n int) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.BatchSize = n
	}
}

// SubscriberBatchWait control batch filling wait time for handler
func SubscriberBatchWait(td time.Duration) SubscriberOption {
	return func(o *SubscriberOptions) {
		o.BatchWait = td
	}
}

// Hooks sets hook runs before action
func Hooks(h ...options.Hook) Option {
	return func(o *Options) {
		o.Hooks = append(o.Hooks, h...)
	}
}
