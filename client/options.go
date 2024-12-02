package client

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/options"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/selector"
	"go.unistack.org/micro/v3/selector/random"
	"go.unistack.org/micro/v3/tracer"
)

// Options holds client options
type Options struct {
	// Selector used to select needed address
	Selector selector.Selector
	// Logger used to log messages
	Logger logger.Logger
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Broker used to publish messages
	Broker broker.Broker
	// Meter used for metrics
	Meter meter.Meter
	// Context is used for external options
	Context context.Context
	// Router used to get route
	Router router.Router
	// TLSConfig specifies tls.Config for secure connection
	TLSConfig *tls.Config
	// Codecs map
	Codecs map[string]codec.Codec
	// Lookup func used to get destination addr
	Lookup LookupFunc
	// Proxy is used for proxy requests
	Proxy string
	// ContentType is used to select codec
	ContentType string
	// Name is the client name
	Name string
	// Wrappers contains wrappers
	Wrappers []Wrapper
	// CallOptions contains default CallOptions
	CallOptions CallOptions
	// PoolSize connection pool size
	PoolSize int
	// PoolTTL connection pool ttl
	PoolTTL time.Duration
	// ContextDialer used to connect
	ContextDialer func(context.Context, string) (net.Conn, error)
	// Hooks can be run before broker Publish/BatchPublish and
	// Subscribe/BatchSubscribe methods
	Hooks options.Hooks
}

// NewCallOptions creates new call options struct
func NewCallOptions(opts ...CallOption) CallOptions {
	options := CallOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// CallOptions holds client call options
type CallOptions struct {
	// Selector selects addr
	Selector selector.Selector
	// Context used for deadline
	Context context.Context
	// Router used for route
	Router router.Router
	// Retry func used for retries
	Retry RetryFunc
	// Backoff func used for backoff when retry
	Backoff BackoffFunc
	// Network name
	Network string
	// Content-Type
	ContentType string
	// AuthToken string
	AuthToken string
	// Address specifies static addr list
	Address []string
	// SelectOptions selector options
	SelectOptions []selector.SelectOption
	// StreamTimeout stream timeout
	StreamTimeout time.Duration
	// RequestTimeout request timeout
	RequestTimeout time.Duration
	// RequestMetadata holds additional metadata for call
	RequestMetadata metadata.Metadata
	// ResponseMetadata holds additional metadata from call
	ResponseMetadata *metadata.Metadata
	// DialTimeout dial timeout
	DialTimeout time.Duration
	// Retries specifies retries num
	Retries int
	// ContextDialer used to connect
	ContextDialer func(context.Context, string) (net.Conn, error)
}

// ContextDialer pass ContextDialer to client
func ContextDialer(fn func(context.Context, string) (net.Conn, error)) Option {
	return func(o *Options) {
		o.ContextDialer = fn
	}
}

// Context pass context to client
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// NewPublishOptions create new PublishOptions struct from option
func NewPublishOptions(opts ...PublishOption) PublishOptions {
	options := PublishOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// PublishOptions holds publish options
type PublishOptions struct {
	// Context used for external options
	Context context.Context
	// Exchange topic exchange name
	Exchange string
	// BodyOnly will publish only message body
	BodyOnly bool
}

// NewMessageOptions creates message options struct
func NewMessageOptions(opts ...MessageOption) MessageOptions {
	options := MessageOptions{Metadata: metadata.New(1)}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// MessageOptions holds client message options
type MessageOptions struct {
	// Metadata additional metadata
	Metadata metadata.Metadata
	// ContentType specify content-type of message
	// deprecated
	ContentType string
}

// NewRequestOptions creates new RequestOptions struct
func NewRequestOptions(opts ...RequestOption) RequestOptions {
	options := RequestOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// RequestOptions holds client request options
type RequestOptions struct {
	// Context used for external options
	Context context.Context
	// ContentType specify content-type of message
	ContentType string
	// Stream flag
	Stream bool
}

// NewOptions creates new options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Context:     context.Background(),
		ContentType: DefaultContentType,
		Codecs:      DefaultCodecs,
		CallOptions: CallOptions{
			Context:        context.Background(),
			Backoff:        DefaultBackoff,
			Retry:          DefaultRetry,
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
		},
		Lookup:   LookupRoute,
		PoolSize: DefaultPoolSize,
		PoolTTL:  DefaultPoolTTL,
		Selector: random.NewSelector(),
		Logger:   logger.DefaultLogger,
		Broker:   broker.DefaultBroker,
		Meter:    meter.DefaultMeter,
		Tracer:   tracer.DefaultTracer,
		Router:   router.DefaultRouter,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Broker to be used for pub/sub
func Broker(b broker.Broker) Option {
	return func(o *Options) {
		o.Broker = b
	}
}

// Tracer to be used for tracing
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Logger to be used for log mesages
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter to be used for metrics
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c codec.Codec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

// ContentType used by default if not specified
func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

// Proxy sets the proxy address
func Proxy(addr string) Option {
	return func(o *Options) {
		o.Proxy = addr
	}
}

// PoolSize sets the connection pool size
func PoolSize(d int) Option {
	return func(o *Options) {
		o.PoolSize = d
	}
}

// PoolTTL sets the connection pool ttl
func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}

// Register sets the routers register
func Register(r register.Register) Option {
	return func(o *Options) {
		if o.Router != nil {
			_ = o.Router.Init(router.Register(r))
		}
	}
}

// Router is used to lookup routes for a service
func Router(r router.Router) Option {
	return func(o *Options) {
		o.Router = r
	}
}

// Selector is used to select a route
func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

// Backoff is used to set the backoff function used when retrying Calls
func Backoff(fn BackoffFunc) Option {
	return func(o *Options) {
		o.CallOptions.Backoff = fn
	}
}

// Name sets the client name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Lookup sets the lookup function to use for resolving service names
func Lookup(l LookupFunc) Option {
	return func(o *Options) {
		o.Lookup = l
	}
}

// TLSConfig specifies a *tls.Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		// set the internal tls
		o.TLSConfig = t
	}
}

// Retries sets the retry count when making the request.
func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

// Retry sets the retry function to be used when re-trying.
func Retry(fn RetryFunc) Option {
	return func(o *Options) {
		o.CallOptions.Retry = fn
	}
}

// RequestTimeout is the request timeout.
func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.RequestTimeout = d
	}
}

// StreamTimeout sets the stream timeout
func StreamTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.StreamTimeout = d
	}
}

// DialTimeout sets the dial timeout
func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.DialTimeout = d
	}
}

// WithExchange sets the exchange to route a message through
// Deprecated
func WithExchange(e string) PublishOption {
	return func(o *PublishOptions) {
		o.Exchange = e
	}
}

// PublishExchange sets the exchange to route a message through
func PublishExchange(e string) PublishOption {
	return func(o *PublishOptions) {
		o.Exchange = e
	}
}

// WithBodyOnly publish only message body
// DERECATED
func WithBodyOnly(b bool) PublishOption {
	return func(o *PublishOptions) {
		o.BodyOnly = b
	}
}

// PublishBodyOnly publish only message body
func PublishBodyOnly(b bool) PublishOption {
	return func(o *PublishOptions) {
		o.BodyOnly = b
	}
}

// PublishContext sets the context in publish options
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// WithContextDialer pass ContextDialer to client call
func WithContextDialer(fn func(context.Context, string) (net.Conn, error)) CallOption {
	return func(o *CallOptions) {
		o.ContextDialer = fn
	}
}

// WithContentType specifies call content type
func WithContentType(ct string) CallOption {
	return func(o *CallOptions) {
		o.ContentType = ct
	}
}

// WithAddress sets the remote addresses to use rather than using service discovery
func WithAddress(a ...string) CallOption {
	return func(o *CallOptions) {
		o.Address = a
	}
}

// WithBackoff is a CallOption which overrides that which
// set in Options.CallOptions
func WithBackoff(fn BackoffFunc) CallOption {
	return func(o *CallOptions) {
		o.Backoff = fn
	}
}

// WithRetry is a CallOption which overrides that which
// set in Options.CallOptions
func WithRetry(fn RetryFunc) CallOption {
	return func(o *CallOptions) {
		o.Retry = fn
	}
}

// WithRetries is a CallOption which overrides that which
// set in Options.CallOptions
func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

// WithResponseMetadata is a CallOption which adds metadata.Metadata to Options.CallOptions
func WithResponseMetadata(md *metadata.Metadata) CallOption {
	return func(o *CallOptions) {
		o.ResponseMetadata = md
	}
}

// WithRequestMetadata is a CallOption which adds metadata.Metadata to Options.CallOptions
func WithRequestMetadata(md metadata.Metadata) CallOption {
	return func(o *CallOptions) {
		o.RequestMetadata = md
	}
}

// WithRequestTimeout is a CallOption which overrides that which
// set in Options.CallOptions
func WithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
	}
}

// WithStreamTimeout sets the stream timeout
func WithStreamTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.StreamTimeout = d
	}
}

// WithDialTimeout is a CallOption which overrides that which
// set in Options.CallOptions
func WithDialTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.DialTimeout = d
	}
}

// WithAuthToken is a CallOption which overrides the
// authorization header with the services own auth token
func WithAuthToken(t string) CallOption {
	return func(o *CallOptions) {
		o.AuthToken = t
	}
}

// WithRouter sets the router to use for this call
func WithRouter(r router.Router) CallOption {
	return func(o *CallOptions) {
		o.Router = r
	}
}

// WithSelector sets the selector to use for this call
func WithSelector(s selector.Selector) CallOption {
	return func(o *CallOptions) {
		o.Selector = s
	}
}

// WithSelectOptions sets the options to pass to the selector for this call
func WithSelectOptions(sops ...selector.SelectOption) CallOption {
	return func(o *CallOptions) {
		o.SelectOptions = sops
	}
}

// WithMessageContentType sets the message content type
// Deprecated
func WithMessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.Metadata.Set(metadata.HeaderContentType, ct)
		o.ContentType = ct
	}
}

// MessageContentType sets the message content type
func MessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.Metadata.Set(metadata.HeaderContentType, ct)
		o.ContentType = ct
	}
}

// MessageMetadata sets the message metadata
func MessageMetadata(k, v string) MessageOption {
	return func(o *MessageOptions) {
		o.Metadata.Set(k, v)
	}
}

// StreamingRequest specifies that request is streaming
func StreamingRequest(b bool) RequestOption {
	return func(o *RequestOptions) {
		o.Stream = b
	}
}

// RequestContentType specifies request content type
func RequestContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		o.ContentType = ct
	}
}

// Hooks sets hook runs before action
func Hooks(h ...options.Hook) Option {
	return func(o *Options) {
		o.Hooks = append(o.Hooks, h...)
	}
}
