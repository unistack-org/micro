package store

import (
	"context"
	"crypto/tls"
	"time"

	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/tracer"
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
	// Addrs contains store address
	Addrs []string
	// Wrappers store wrapper that called before actual functions
	// Wrappers []Wrapper
	// Timeout specifies timeout duration for all operations
	Timeout time.Duration
}

// NewOptions creates options struct
func NewOptions(opts ...Option) Options {
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

// Option sets values in Options
type Option func(o *Options)

// TLSConfig specifies a *tls.Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// Context pass context to store
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Codec sets the codec
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
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

// Name the name of the store
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Separator the value used as key parts separator
func Separator(s string) Option {
	return func(o *Options) {
		o.Separator = s
	}
}

// Namespace sets namespace of the store
func Namespace(ns string) Option {
	return func(o *Options) {
		o.Namespace = ns
	}
}

// Tracer sets the tracer
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Timeout sets the timeout
func Timeout(td time.Duration) Option {
	return func(o *Options) {
		o.Timeout = td
	}
}

// Addrs contains the addresses or other connection information of the backing storage.
// For example, an etcd implementation would contain the nodes of the cluster.
// A SQL implementation could contain one or more connection strings.
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// ReadOptions configures an individual Read operation
type ReadOptions struct {
	// Context holds external options
	Context context.Context
	// Namespace holds namespace
	Namespace string
	// Name holds mnemonic name
	Name string
	// Timeout specifies max timeout for operation
	Timeout time.Duration
}

// NewReadOptions fills ReadOptions struct with opts slice
func NewReadOptions(opts ...ReadOption) ReadOptions {
	options := ReadOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ReadOption sets values in ReadOptions
type ReadOption func(r *ReadOptions)

// ReadTimeout pass timeout to ReadOptions
func ReadTimeout(td time.Duration) ReadOption {
	return func(o *ReadOptions) {
		o.Timeout = td
	}
}

// ReadName pass name to ReadOptions
func ReadName(name string) ReadOption {
	return func(o *ReadOptions) {
		o.Name = name
	}
}

// ReadContext pass context.Context to ReadOptions
func ReadContext(ctx context.Context) ReadOption {
	return func(o *ReadOptions) {
		o.Context = ctx
	}
}

// ReadNamespace pass namespace to ReadOptions
func ReadNamespace(ns string) ReadOption {
	return func(o *ReadOptions) {
		o.Namespace = ns
	}
}

// WriteOptions configures an individual Write operation
type WriteOptions struct {
	// Context holds external options
	Context context.Context
	// Metadata contains additional metadata
	Metadata metadata.Metadata
	// Namespace holds namespace
	Namespace string
	// Name holds mnemonic name
	Name string
	// Timeout specifies max timeout for operation
	Timeout time.Duration
	// TTL specifies key TTL
	TTL time.Duration
}

// NewWriteOptions fills WriteOptions struct with opts slice
func NewWriteOptions(opts ...WriteOption) WriteOptions {
	options := WriteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WriteOption sets values in WriteOptions
type WriteOption func(w *WriteOptions)

// WriteContext pass context.Context to wirte options
func WriteContext(ctx context.Context) WriteOption {
	return func(o *WriteOptions) {
		o.Context = ctx
	}
}

// WriteMetadata add metadata.Metadata
func WriteMetadata(md metadata.Metadata) WriteOption {
	return func(o *WriteOptions) {
		o.Metadata = metadata.Copy(md)
	}
}

// WriteTTL is the time the record expires
func WriteTTL(d time.Duration) WriteOption {
	return func(o *WriteOptions) {
		o.TTL = d
	}
}

// WriteNamespace pass namespace to write options
func WriteNamespace(ns string) WriteOption {
	return func(o *WriteOptions) {
		o.Namespace = ns
	}
}

// WriteName pass name to WriteOptions
func WriteName(name string) WriteOption {
	return func(o *WriteOptions) {
		o.Name = name
	}
}

// WriteTimeout pass timeout to WriteOptions
func WriteTimeout(td time.Duration) WriteOption {
	return func(o *WriteOptions) {
		o.Timeout = td
	}
}

// DeleteOptions configures an individual Delete operation
type DeleteOptions struct {
	// Context holds external options
	Context context.Context
	// Namespace holds namespace
	Namespace string
	// Name holds mnemonic name
	Name string
	// Timeout specifies max timeout for operation
	Timeout time.Duration
}

// NewDeleteOptions fills DeleteOptions struct with opts slice
func NewDeleteOptions(opts ...DeleteOption) DeleteOptions {
	options := DeleteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// DeleteOption sets values in DeleteOptions
type DeleteOption func(d *DeleteOptions)

// DeleteContext pass context.Context to delete options
func DeleteContext(ctx context.Context) DeleteOption {
	return func(o *DeleteOptions) {
		o.Context = ctx
	}
}

// DeleteNamespace pass namespace to delete options
func DeleteNamespace(ns string) DeleteOption {
	return func(o *DeleteOptions) {
		o.Namespace = ns
	}
}

// DeleteName pass name to DeleteOptions
func DeleteName(name string) DeleteOption {
	return func(o *DeleteOptions) {
		o.Name = name
	}
}

// DeleteTimeout pass timeout to DeleteOptions
func DeleteTimeout(td time.Duration) DeleteOption {
	return func(o *DeleteOptions) {
		o.Timeout = td
	}
}

// ListOptions configures an individual List operation
type ListOptions struct {
	Context   context.Context
	Prefix    string
	Suffix    string
	Namespace string
	// Name holds mnemonic name
	Name   string
	Limit  uint
	Offset uint
	// Timeout specifies max timeout for operation
	Timeout time.Duration
}

// NewListOptions fills ListOptions struct with opts slice
func NewListOptions(opts ...ListOption) ListOptions {
	options := ListOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ListOption sets values in ListOptions
type ListOption func(l *ListOptions)

// ListContext pass context.Context to list options
func ListContext(ctx context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = ctx
	}
}

// ListPrefix returns all keys that are prefixed with key
func ListPrefix(s string) ListOption {
	return func(o *ListOptions) {
		o.Prefix = s
	}
}

// ListSuffix returns all keys that end with key
func ListSuffix(s string) ListOption {
	return func(o *ListOptions) {
		o.Suffix = s
	}
}

// ListLimit limits the number of returned keys
func ListLimit(n uint) ListOption {
	return func(o *ListOptions) {
		o.Limit = n
	}
}

// ListOffset use with Limit for pagination
func ListOffset(n uint) ListOption {
	return func(o *ListOptions) {
		o.Offset = n
	}
}

// ListNamespace pass namespace to list options
func ListNamespace(ns string) ListOption {
	return func(o *ListOptions) {
		o.Namespace = ns
	}
}

// ListTimeout pass timeout to ListOptions
func ListTimeout(td time.Duration) ListOption {
	return func(o *ListOptions) {
		o.Timeout = td
	}
}

// ExistsOptions holds options for Exists method
type ExistsOptions struct {
	// Context holds external options
	Context context.Context
	// Namespace contains namespace
	Namespace string
	// Name holds mnemonic name
	Name string
	// Timeout specifies max timeout for operation
	Timeout time.Duration
}

// ExistsOption specifies Exists call options
type ExistsOption func(*ExistsOptions)

// NewExistsOptions helper for Exists method
func NewExistsOptions(opts ...ExistsOption) ExistsOptions {
	options := ExistsOptions{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ExistsContext pass context.Context to exist options
func ExistsContext(ctx context.Context) ExistsOption {
	return func(o *ExistsOptions) {
		o.Context = ctx
	}
}

// ExistsNamespace pass namespace to exist options
func ExistsNamespace(ns string) ExistsOption {
	return func(o *ExistsOptions) {
		o.Namespace = ns
	}
}

// ExistsName pass name to exist options
func ExistsName(name string) ExistsOption {
	return func(o *ExistsOptions) {
		o.Name = name
	}
}

// ExistsTimeout timeout to ListOptions
func ExistsTimeout(td time.Duration) ExistsOption {
	return func(o *ExistsOptions) {
		o.Timeout = td
	}
}

/*
// WrapStore adds a store Wrapper to a list of options passed into the store
func WrapStore(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}
*/
