package broker

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/tracer"
)

var (
	// PublishMessageDurationSeconds specifies meter metric name
	PublishMessageDurationSeconds = "publish_message_duration_seconds"
	// PublishMessageLatencyMicroseconds specifies meter metric name
	PublishMessageLatencyMicroseconds = "publish_message_latency_microseconds"
	// PublishMessageTotal specifies meter metric name
	PublishMessageTotal = "publish_message_total"
	// PublishMessageInflight specifies meter metric name
	PublishMessageInflight = "publish_message_inflight"
	// SubscribeMessageDurationSeconds specifies meter metric name
	SubscribeMessageDurationSeconds = "subscribe_message_duration_seconds"
	// SubscribeMessageLatencyMicroseconds specifies meter metric name
	SubscribeMessageLatencyMicroseconds = "subscribe_message_latency_microseconds"
	// SubscribeMessageTotal specifies meter metric name
	SubscribeMessageTotal = "subscribe_message_total"
	// SubscribeMessageInflight specifies meter metric name
	SubscribeMessageInflight = "subscribe_message_inflight"
)

// Options struct
type Options struct {
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Register can be used for clustering
	Register register.Register
	// Codecs holds the codec for marshal/unmarshal
	Codecs map[string]codec.Codec
	// Logger used for logging
	Logger logger.Logger
	// Meter used for metrics
	Meter meter.Meter
	// Context holds external options
	Context context.Context
	// TLSConfig holds tls.TLSConfig options
	TLSConfig *tls.Config
	// ErrorHandler used when broker can't unmarshal incoming message
	ErrorHandler func(Message)
	// Name holds the broker name
	Name string
	// Addrs holds the broker address
	Addrs []string
}

// Option func
type Option func(*Options)

// NewOptions create new Options
func NewOptions(opts ...Option) Options {
	options := Options{
		Register: register.DefaultRegister,
		Logger:   logger.DefaultLogger,
		Context:  context.Background(),
		Meter:    meter.DefaultMeter,
		Codecs:   make(map[string]codec.Codec),
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

// MessageOption func
type MessageOption func(*MessageOptions)

// MessageOptions struct
type MessageOptions struct {
	Metadata    metadata.Metadata
	ContentType string
}

// MessageMetadata pass additional message metadata
func MessageMetadata(md metadata.Metadata) MessageOption {
	return func(o *MessageOptions) {
		o.Metadata = md
	}
}

// MessageContentType pass ContentType for message data
func MessageContentType(ct string) MessageOption {
	return func(o *MessageOptions) {
		o.ContentType = ct
	}
}

// PublishOption func
type PublishOption func(*PublishOptions)

// PublishOptions struct
type PublishOptions struct {
	// Context holds external options
	Context context.Context
	// BodyOnly flag says the message contains raw body bytes
	BodyOnly bool
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
	// Context holds external options
	Context context.Context
	// ErrorHandler used when broker can't unmarshal incoming message
	ErrorHandler func(Message)
	// QueueGroup holds consumer group
	QueueGroup string
	// AutoAck flag specifies auto ack of incoming message when no error happens
	AutoAck bool
	// BodyOnly flag specifies that message contains only body bytes without header
	BodyOnly bool
	// BatchSize flag specifies max batch size
	BatchSize int
	// BatchWait flag specifies max wait time for batch filling
	BatchWait time.Duration
}

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

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding used where
// a broker does not support headers
// Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c codec.Codec) Option {
	return func(o *Options) {
		o.Codecs[contentType] = c
	}
}

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func ErrorHandler(h func(Message)) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

// SubscribeErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func SubscribeErrorHandler(h func(Message)) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.ErrorHandler = h
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

// SubscribeOption func signature
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

// SubscribeQueueGroup sets the shared queue name distributed messages across subscribers
func SubscribeQueueGroup(n string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.QueueGroup = n
	}
}

// SubscribeAutoAck control auto ack processing for handler
func SubscribeAuthAck(b bool) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = b
	}
}
