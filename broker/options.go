package broker

import (
	"context"
	"crypto/tls"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options struct
type Options struct {
	Name string
	// Addrs useed by broker
	Addrs []string
	// ErrorHandler executed when errors occur processing messages
	ErrorHandler Handler
	// Codec used to marshal/unmarshal messages
	Codec codec.Codec
	// Logger the used logger
	Logger logger.Logger
	// Meter the used for metrics
	Meter meter.Meter
	// Tracer used for trace
	Tracer tracer.Tracer
	// TLSConfig for secure communication
	TLSConfig *tls.Config
	// Register used for clustering
	Register register.Register
	// Context is used for non default options
	Context context.Context
}

// NewOptions create new Options
func NewOptions(opts ...Option) Options {
	options := Options{
		Register: register.DefaultRegister,
		Logger:   logger.DefaultLogger,
		Context:  context.Background(),
		Meter:    meter.DefaultMeter,
		Codec:    codec.DefaultCodec,
		Tracer:   tracer.DefaultTracer,
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
	// BodyOnly says that only body of the message must be published
	BodyOnly bool
	// Context for non default options
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

	// ErrorHandler executed when errors occur processing messages
	ErrorHandler Handler

	// Group for subscriber, Subscribers with the same group name
	// will create a shared subscription where each
	// receives a subset of messages.
	Group string

	// BodyOnly says that consumed only body of the message
	BodyOnly bool

	// Context is used for non default options
	Context context.Context
}

// Option func
type Option func(*Options)

// PublishOption func
type PublishOption func(*PublishOptions)

// PublishBodyOnly publish only body of the message
func PublishBodyOnly(b bool) PublishOption {
	return func(o *PublishOptions) {
		o.BodyOnly = b
	}
}

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

// SubscribeBodyOnly consumes only body of the message
func SubscribeBodyOnly(b bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.BodyOnly = b
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

// Queue sets the subscribers queue
// Deprecated
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

// Register sets register option
func Register(r register.Register) Option {
	return func(o *Options) {
		o.Register = r
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

// Tracer to be used for tracing
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Meter sets the meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}
