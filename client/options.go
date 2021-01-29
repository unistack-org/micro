package client

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/network/transport"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/router"
	"github.com/unistack-org/micro/v3/selector"
	"github.com/unistack-org/micro/v3/selector/random"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options holds client options
type Options struct {
	Name string
	// Used to select codec
	ContentType string
	// Proxy address to send requests via
	Proxy string

	// Plugged interfaces
	Broker    broker.Broker
	Codecs    map[string]codec.Codec
	Router    router.Router
	Selector  selector.Selector
	Transport transport.Transport
	Logger    logger.Logger
	Meter     meter.Meter
	// Lookup used for looking up routes
	Lookup LookupFunc

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration
	Tracer   tracer.Tracer
	// Wrapper that used client
	Wrappers []Wrapper

	// CallOptions that used by default
	CallOptions CallOptions

	// Context is used for non default options
	Context context.Context
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
	// Address of remote hosts
	Address []string
	// Backoff func
	Backoff BackoffFunc
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Number of Call attempts
	Retries int
	// Check if retriable func
	Retry RetryFunc
	// Request/Response timeout
	RequestTimeout time.Duration
	// Router to use for this call
	Router router.Router
	// Selector to use for the call
	Selector selector.Selector
	// SelectOptions to use when selecting a route
	SelectOptions []selector.SelectOption
	// Stream timeout for the stream
	StreamTimeout time.Duration
	// Use the auth token as the authorization header
	AuthToken bool
	// Network to lookup the route within
	Network string
	// Middleware for low level call func
	CallWrappers []CallWrapper
	// Context is uded for non default options
	Context context.Context
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
	// Exchange is the routing exchange for the message
	Exchange string
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// NewMessageOptions creates message options struct
func NewMessageOptions(opts ...MessageOption) MessageOptions {
	options := MessageOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// MessageOptions holds client message options
type MessageOptions struct {
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
	ContentType string
	Stream      bool
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// NewOptions creates new options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Context:     context.Background(),
		ContentType: "application/json",
		Codecs:      make(map[string]codec.Codec),
		CallOptions: CallOptions{
			Backoff:        DefaultBackoff,
			Retry:          DefaultRetry,
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
			DialTimeout:    transport.DefaultDialTimeout,
		},
		Lookup:   LookupRoute,
		PoolSize: DefaultPoolSize,
		PoolTTL:  DefaultPoolTTL,
		Selector: random.NewSelector(),
		Logger:   logger.DefaultLogger,
		Broker:   broker.DefaultBroker,
		Meter:    meter.DefaultMeter,
		Tracer:   tracer.DefaultTracer,
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

// Default content type of the client
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

// Transport to use for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Register sets the routers register
func Register(r register.Register) Option {
	return func(o *Options) {
		if o.Router != nil {
			o.Router.Init(router.Register(r))
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

// Adds a Wrapper to a list of options passed into the client
func Wrap(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

// Adds a Wrapper to the list of CallFunc wrappers
func WrapCall(cw ...CallWrapper) Option {
	return func(o *Options) {
		o.CallOptions.CallWrappers = append(o.CallOptions.CallWrappers, cw...)
	}
}

// Backoff is used to set the backoff function used
// when retrying Calls
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

// Retries sets the retry count when making the request.
// Should this be a Call Option?
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
// Should this be a Call Option?
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

// Transport dial timeout
func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.DialTimeout = d
	}
}

// Call Options

// WithExchange sets the exchange to route a message through
func WithExchange(e string) PublishOption {
	return func(o *PublishOptions) {
		o.Exchange = e
	}
}

// PublishContext sets the context in publish options
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// WithAddress sets the remote addresses to use rather than using service discovery
func WithAddress(a ...string) CallOption {
	return func(o *CallOptions) {
		o.Address = a
	}
}

// WithCallWrapper is a CallOption which adds to the existing CallFunc wrappers
func WithCallWrapper(cw ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.CallWrappers = append(o.CallWrappers, cw...)
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
func WithAuthToken() CallOption {
	return func(o *CallOptions) {
		o.AuthToken = true
	}
}

// WithNetwork is a CallOption which sets the network attribute
func WithNetwork(n string) CallOption {
	return func(o *CallOptions) {
		o.Network = n
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
func WithMessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.ContentType = ct
	}
}

// Request Options

// WithContentType specifies request content type
func WithContentType(ct string) RequestOption {
	return func(o *RequestOptions) {
		o.ContentType = ct
	}
}

// StreamingRequest specifies that request is streaming
func StreamingRequest() RequestOption {
	return func(o *RequestOptions) {
		o.Stream = true
	}
}
