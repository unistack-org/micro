package client

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/router"
	"go.unistack.org/micro/v4/selector"
	"go.unistack.org/micro/v4/selector/random"
	"go.unistack.org/micro/v4/tracer"
)

// Options holds client options
type Options struct {
	// Selector used to select needed address
	Selector selector.Selector
	// Logger used to log messages
	Logger logger.Logger
	// Tracer used for tracing
	Tracer tracer.Tracer
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
	// CallOptions contains default CallOptions
	CallOptions CallOptions
	// PoolSize connection pool size
	PoolSize int
	// PoolTTL connection pool ttl
	PoolTTL time.Duration
	// ContextDialer used to connect
	ContextDialer func(context.Context, string) (net.Conn, error)
	// Hooks may contains Client func wrapper
	Hooks options.Hooks
}

// NewCallOptions creates new call options struct
func NewCallOptions(opts ...options.Option) CallOptions {
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
	// CallWrappers call wrappers
	CallWrappers []CallWrapper
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
func ContextDialer(fn func(context.Context, string) (net.Conn, error)) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".ContextDialer")
	}
}

// NewRequestOptions creates new RequestOptions struct
func NewRequestOptions(opts ...options.Option) RequestOptions {
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
func NewOptions(opts ...options.Option) Options {
	options := Options{
		Context:     context.Background(),
		ContentType: DefaultContentType,
		Codecs:      make(map[string]codec.Codec),
		CallOptions: CallOptions{
			Context:        context.Background(),
			Backoff:        DefaultBackoff,
			Retry:          DefaultRetry,
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
			DialTimeout:    DefaultDialTimeout,
		},
		Lookup:   LookupRoute,
		PoolSize: DefaultPoolSize,
		PoolTTL:  DefaultPoolTTL,
		Selector: random.NewSelector(),
		Logger:   logger.DefaultLogger,
		Meter:    meter.DefaultMeter,
		Tracer:   tracer.DefaultTracer,
		Router:   router.DefaultRouter,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Proxy sets the proxy address
func Proxy(addr string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, addr, ".Proxy")
	}
}

// PoolSize sets the connection pool size
func PoolSize(d int) options.Option {
	return func(src interface{}) error {
		return options.Set(src, d, ".PoolSize")
	}
}

// PoolTTL sets the connection pool ttl
func PoolTTL(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".PoolTTL")
	}
}

// Selector is used to select a route
func Selector(s selector.Selector) options.Option {
	return func(src interface{}) error {
		return options.Set(src, s, ".Selector")
	}
}

// Backoff is used to set the backoff function used when retrying Calls
func Backoff(fn BackoffFunc) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".Backoff")
	}
}

// Lookup sets the lookup function to use for resolving service names
func Lookup(fn LookupFunc) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".Lookup")
	}
}

// WithCallWrapper sets the retry function to be used when re-trying.
func WithCallWrapper(fn CallWrapper) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".CallWrappers")
	}
}

// Retries sets the retry count when making the request.
func Retries(n int) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".Retries")
	}
}

// Retry sets the retry function to be used when re-trying.
func Retry(fn RetryFunc) options.Option {
	return func(src interface{}) error {
		return options.Set(src, fn, ".Retry")
	}
}

// RequestTimeout is the request timeout.
func RequestTimeout(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".RequestTimeout")
	}
}

// StreamTimeout sets the stream timeout
func StreamTimeout(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".StreamTimeout")
	}
}

// DialTimeout sets the dial timeout
func DialTimeout(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".DialTimeout")
	}
}

// WithResponseMetadata is a CallOption which adds metadata.Metadata to Options.CallOptions
func ResponseMetadata(md *metadata.Metadata) options.Option {
	return func(src interface{}) error {
		return options.Set(src, md, ".ResponseMetadata")
	}
}

// WithRequestMetadata is a CallOption which adds metadata.Metadata to Options.CallOptions
func RequestMetadata(md metadata.Metadata) options.Option {
	return func(src interface{}) error {
		return options.Set(src, metadata.Copy(md), ".RequestMetadata")
	}
}

// AuthToken is a CallOption which overrides the
// authorization header with the services own auth token
func AuthToken(t string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, t, ".AuthToken")
	}
}

// Network is a CallOption which sets the network attribute
func Network(n string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".Network")
	}
}

/*
// WithSelectOptions sets the options to pass to the selector for this call
func WithSelectOptions(sops ...selector.SelectOption) options.Option {
	return func(o *CallOptions) {
		o.SelectOptions = sops
	}
}
*/

// StreamingRequest specifies that request is streaming
func StreamingRequest(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".Stream")
	}
}
