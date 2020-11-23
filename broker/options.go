package broker

import (
	"context"
	"crypto/tls"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/registry"
)

// Options struct
type Options struct {
	Addrs  []string
	Secure bool
	Codec  codec.Codec

	// Logger
	Logger logger.Logger
	// Handler executed when errors occur processing messages
	ErrorHandler Handler

	TLSConfig *tls.Config
	// Registry used for clustering
	Registry registry.Registry
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// NewOptions create new Options
func NewOptions(opts ...Option) Options {
	options := Options{
		Registry: registry.DefaultRegistry,
		Logger:   logger.DefaultLogger,
		Context:  context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Context sets the context option
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// PublishOptions struct
type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// NewPublishOptions creates PublishOptions struct
func NewPublishOptions(opts ...PublishOption) PublishOptions {
	options := PublishOptions{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// SubscribeOptions struct
type SubscribeOptions struct {
	// AutoAck ack messages if handler returns nil err
	AutoAck bool

	// Handler executed when errors occur processing messages
	ErrorHandler Handler

	// Subscribers with the same group name
	// will create a shared subscription where each
	// receives a subset of messages.
	Group string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option func
type Option func(*Options)

// PublishOption func
type PublishOption func(*PublishOptions)

// PublishContext sets the context
func PublishContext(ctx context.Context) PublishOption {
	return func(o *PublishOptions) {
		o.Context = ctx
	}
}

// SubscribeOption func
type SubscribeOption func(*SubscribeOptions)

// NewSubscribeOptions creates new SubscribeOptions
func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	options := SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding used where
// a broker does not support headers
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// DisableAutoAck disables auto ack
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// SubscribeAutoAck will disable auto acking of messages
// after they have been handled.
func SubscribeAutoAck(b bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = b
	}
}

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

// SubscribeErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func SubscribeErrorHandler(h Handler) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.ErrorHandler = h
	}
}

// Queue sets the subscribers sueue
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Group = name
	}
}

// SubscribeGroup sets the name of the queue to share messages on
func SubscribeGroup(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Group = name
	}
}

// Registry sets registry option
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Secure communication with the broker
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig sets the TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
