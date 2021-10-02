package transport

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/tracer"
)

// Options struct holds the transport options
type Options struct {
	// Meter used for metrics
	Meter meter.Meter
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Codec used for marshal/unmarshal messages
	Codec codec.Codec
	// Logger used for logging
	Logger logger.Logger
	// Context holds external options
	Context context.Context
	// TLSConfig holds tls.TLSConfig options
	TLSConfig *tls.Config
	// Name holds the transport name
	Name string
	// Addrs holds the transport addrs
	Addrs []string
	// Timeout holds the timeout
	Timeout time.Duration
}

// NewOptions returns new options
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Meter:   meter.DefaultMeter,
		Tracer:  tracer.DefaultTracer,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// DialOptions struct
type DialOptions struct {
	// Context holds the external options
	Context context.Context
	// Timeout holds the timeout
	Timeout time.Duration
	// Stream flag
	Stream bool
}

// NewDialOptions returns new DialOptions
func NewDialOptions(opts ...DialOption) DialOptions {
	options := DialOptions{
		Timeout: DefaultDialTimeout,
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// ListenOptions struct
type ListenOptions struct {
	// TODO: add tls options when listening
	// Currently set in global options
	// Context holds the external options
	Context context.Context
	// TLSConfig holds the *tls.Config options
	TLSConfig *tls.Config
}

// NewListenOptions returns new ListenOptions
func NewListenOptions(opts ...ListenOption) ListenOptions {
	options := ListenOptions{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// Addrs to use for transport
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter sets the meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Context sets the context
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Codec sets the codec used for encoding where the transport
// does not support message headers
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// Timeout sets the timeout for Send/Recv execution
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// TLSConfig to be used for the transport.
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// WithStream indicates whether this is a streaming connection
func WithStream() DialOption {
	return func(o *DialOptions) {
		o.Stream = true
	}
}

// WithTimeout used when dialling the remote side
func WithTimeout(d time.Duration) DialOption {
	return func(o *DialOptions) {
		o.Timeout = d
	}
}

// Tracer to be used for tracing
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
