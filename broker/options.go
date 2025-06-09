package broker

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/sync"
	"go.unistack.org/micro/v4/tracer"
)

// Options struct
type Options struct {
	// Name holds the broker name
	Name string

	// Tracer used for tracing
	Tracer tracer.Tracer
	// Register can be used for clustering
	Register register.Register
	// Codecs holds the codecs for marshal/unmarshal based on content-type
	Codecs map[string]codec.Codec
	// Logger used for logging
	Logger logger.Logger
	// Meter used for metrics
	Meter meter.Meter
	// Context holds external options
	Context context.Context

	// Wait waits for a collection of goroutines to finish
	Wait *sync.WaitGroup
	// TLSConfig holds tls.TLSConfig options
	TLSConfig *tls.Config

	// Addrs holds the broker address
	Addrs []string
	// Hooks can be run before broker Publish/BatchPublish and
	// Subscribe/BatchSubscribe methods
	Hooks options.Hooks

	// GracefulTimeout contains time to wait to finish in flight requests
	GracefulTimeout time.Duration

	// ContentType will be used if no content-type set when creating message
	ContentType string
}

// NewOptions create new Options
func NewOptions(opts ...Option) Options {
	options := Options{
		Register:        register.DefaultRegister,
		Logger:          logger.DefaultLogger,
		Context:         context.Background(),
		Meter:           meter.DefaultMeter,
		Codecs:          make(map[string]codec.Codec),
		Tracer:          tracer.DefaultTracer,
		GracefulTimeout: DefaultGracefulTimeout,
		ContentType:     DefaultContentType,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// DefaultContentType is the default content-type if not specified
var DefaultContentType = ""

// Context sets the context option
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// ContentType used by default if not specified
func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

// MessageOptions struct
type MessageOptions struct {
	// ContentType for message body
	ContentType string
	// BodyOnly flag says the message contains raw body bytes and don't need
	// codec Marshal method
	BodyOnly bool
	// Context holds custom options
	Context context.Context
}

// NewMessageOptions creates MessageOptions struct
func NewMessageOptions(opts ...MessageOption) MessageOptions {
	options := MessageOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// SubscribeOptions struct
type SubscribeOptions struct {
	// Context holds external options
	Context context.Context
	// Group holds consumer group
	Group string
	// AutoAck flag specifies auto ack of incoming message when no error happens
	AutoAck bool
	// BodyOnly flag specifies that message contains only body bytes without header
	BodyOnly bool
	// BatchSize flag specifies max batch size
	BatchSize int
	// BatchWait flag specifies max wait time for batch filling
	BatchWait time.Duration
}

// Option func
type Option func(*Options)

// MessageOption func
type MessageOption func(*MessageOptions)

// MessageContentType sets message content-type that used to Marshal
func MessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.ContentType = ct
	}
}

// MessageBodyOnly publish only body of the message
func MessageBodyOnly(b bool) MessageOption {
	return func(o *MessageOptions) {
		o.BodyOnly = b
	}
}

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding messages
func Codec(ct string, c codec.Codec) Option {
	return func(o *Options) {
		o.Codecs[ct] = c
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

// Hooks sets hook runs before action
func Hooks(h ...options.Hook) Option {
	return func(o *Options) {
		o.Hooks = append(o.Hooks, h...)
	}
}

// SubscribeContext set context
func SubscribeContext(ctx context.Context) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Context = ctx
	}
}

// SubscribeAutoAck contol auto acking of messages
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

// SubscribeBatchSize specifies max batch size
func SubscribeBatchSize(n int) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.BatchSize = n
	}
}

// SubscribeBatchWait specifies max batch wait time
func SubscribeBatchWait(td time.Duration) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.BatchWait = td
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
