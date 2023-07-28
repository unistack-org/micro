package store

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/tracer"
)

// Options contains configuration for the Store
type Options struct {
	// Meter used for metrics
	Meter meter.Meter
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Context holds external options
	Context context.Context
	// Codec used to marshal/unmarshal
	Codec codec.Codec
	// Logger used for logging
	Logger logger.Logger
	// TLSConfig holds tls.TLSConfig options
	TLSConfig *tls.Config
	// Name specifies store name
	Name string
	// Namespace of the records
	Namespace string
	// Separator used as key parts separator
	Separator string
	// Address contains store address
	Address []string
	// Timeout specifies timeout duration for all operations
	Timeout time.Duration
}

// NewOptions creates options struct
func NewOptions(opts ...options.Option) Options {
	options := Options{
		Logger:    logger.DefaultLogger,
		Context:   context.Background(),
		Codec:     codec.DefaultCodec,
		Tracer:    tracer.DefaultTracer,
		Meter:     meter.DefaultMeter,
		Separator: DefaultSeparator,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Separator the value used as key parts separator
func Separator(s string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, s, "Separator")
	}
}

// Timeout sets the timeout
func Timeout(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".Timeout")
	}
}

// ReadOptions configures an individual Read operation
type ReadOptions struct {
	// Context holds external options
	Context context.Context
	// Namespace holds namespace
	Namespace string
}

// NewReadOptions fills ReadOptions struct with opts slice
func NewReadOptions(opts ...options.Option) ReadOptions {
	options := ReadOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WriteOptions configures an individual Write operation
type WriteOptions struct {
	// Context holds external options
	Context context.Context
	// Metadata contains additional metadata
	Metadata metadata.Metadata
	// Namespace holds namespace
	Namespace string
	// TTL specifies key TTL
	TTL time.Duration
}

// NewWriteOptions fills WriteOptions struct with opts slice
func NewWriteOptions(opts ...options.Option) WriteOptions {
	options := WriteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WriteMetadata add metadata.Metadata
func WriteMetadata(md metadata.Metadata) options.Option {
	return func(src interface{}) error {
		return options.Set(src, metadata.Copy(md), ".Metadata")
	}
}

// WriteTTL is the time the record expires
func WriteTTL(td time.Duration) options.Option {
	return func(src interface{}) error {
		return options.Set(src, td, ".TTL")
	}
}

// DeleteOptions configures an individual Delete operation
type DeleteOptions struct {
	// Context holds external options
	Context context.Context
	// Namespace holds namespace
	Namespace string
}

// NewDeleteOptions fills DeleteOptions struct with opts slice
func NewDeleteOptions(opts ...options.Option) DeleteOptions {
	options := DeleteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ListOptions configures an individual List operation
type ListOptions struct {
	Context   context.Context
	Prefix    string
	Suffix    string
	Namespace string
	Limit     uint
	Offset    uint
}

// NewListOptions fills ListOptions struct with opts slice
func NewListOptions(opts ...options.Option) ListOptions {
	options := ListOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ListPrefix returns all keys that are prefixed with key
func ListPrefix(s string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, s, ".Prefix")
	}
}

// ListSuffix returns all keys that end with key
func ListSuffix(s string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, s, ".Prefix")
	}
}

// ListLimit limits the number of returned keys
func ListLimit(n uint) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".Limit")
	}
}

// ListOffset use with Limit for pagination
func ListOffset(n uint) options.Option {
	return func(src interface{}) error {
		return options.Set(src, n, ".Offset")
	}
}

// ExistsOptions holds options for Exists method
type ExistsOptions struct {
	// Context holds external options
	Context context.Context
	// Namespace contains namespace
	Namespace string
}

// NewExistsOptions helper for Exists method
func NewExistsOptions(opts ...options.Option) ExistsOptions {
	options := ExistsOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
