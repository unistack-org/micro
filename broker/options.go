package broker

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/tracer"
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
	// ErrorHandler used when broker have error while processing message
	ErrorHandler interface{}
	// Name holds the broker name
	Name string
	// Address holds the broker address
	Address []string
}

// NewOptions create new Options
func NewOptions(opts ...options.Option) Options {
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

// PublishOptions struct
type PublishOptions struct {
	// Context holds external options
	Context context.Context
	// Message metadata usually passed as message headers
	Metadata metadata.Metadata
	// Content-Type of message for marshal
	ContentType string
	// Topic destination
	Topic string
	// BodyOnly flag says the message contains raw body bytes
	BodyOnly bool
}

// NewPublishOptions creates PublishOptions struct
func NewPublishOptions(opts ...options.Option) PublishOptions {
	options := PublishOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// PublishTopic pass topic for messages
func PublishTopic(t string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, t, ".Topic")
	}
}

// SubscribeOptions struct
type SubscribeOptions struct {
	// Context holds external options
	Context context.Context
	// ErrorHandler used when broker have error while processing message
	ErrorHandler interface{}
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

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors
func ErrorHandler(h interface{}) options.Option {
	return func(src interface{}) error {
		return options.Set(src, h, ".ErrorHandler")
	}
}

// NewSubscribeOptions creates new SubscribeOptions
func NewSubscribeOptions(opts ...options.Option) SubscribeOptions {
	options := SubscribeOptions{
		AutoAck: true,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// SubscribeAutoAck contol auto acking of messages
// after they have been handled.
func SubscribeAutoAck(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".AutoAck")
	}
}

// BodyOnly transfer only body without
func BodyOnly(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".BodyOnly")
	}
}

// SubscribeBatchSize specifies max batch size
func SubscribeBatchSize(n int) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".BatchSize")
	}
}

// SubscribeBatchWait specifies max batch wait time
func SubscribeBatchWait(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".BatchWait")
	}
}

// SubscribeQueueGroup sets the shared queue name distributed messages across subscribers
func SubscribeQueueGroup(n string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".QueueGroup")
	}
}
