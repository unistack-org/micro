package store

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options contains configuration for the Store
type Options struct {
	// Name specifies store name
	Name string
	// Nodes contains the addresses or other connection information of the backing storage.
	// For example, an etcd implementation would contain the nodes of the cluster.
	// A SQL implementation could contain one or more connection strings.
	Nodes []string
	// Database allows multiple isolated stores to be kept in one backend, if supported.
	Database string
	// Table is analag for a table in database backends or a key prefix in KV backends
	Table string
	// Codec that used for marshal/unmarshal value
	Codec codec.Codec
	// Logger the logger
	Logger logger.Logger
	// Meter the meter
	Meter meter.Meter
	// Tracer the tacer
	Tracer tracer.Tracer
	// TLSConfig specifies tls.Config for secure
	TLSConfig *tls.Config
	// Context should contain all implementation specific options
	Context context.Context
}

// NewOptions creates options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Context: context.Background(),
		Codec:   codec.DefaultCodec,
		Tracer:  tracer.DefaultTracer,
		Meter:   meter.DefaultMeter,
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

// Name the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Tracer sets the tracer
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Nodes contains the addresses or other connection information of the backing storage.
// For example, an etcd implementation would contain the nodes of the cluster.
// A SQL implementation could contain one or more connection strings.
func Nodes(a ...string) Option {
	return func(o *Options) {
		o.Nodes = a
	}
}

// Database allows multiple isolated stores to be kept in one backend, if supported.
func Database(db string) Option {
	return func(o *Options) {
		o.Database = db
	}
}

// Table is analag for a table in database backends or a key prefix in KV backends
func Table(t string) Option {
	return func(o *Options) {
		o.Table = t
	}
}

// NewReadOptions fills ReadOptions struct with opts slice
func NewReadOptions(opts ...ReadOption) ReadOptions {
	options := ReadOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ReadOptions configures an individual Read operation
type ReadOptions struct {
	Database  string
	Table     string
	Namespace string
	Context   context.Context
}

// ReadOption sets values in ReadOptions
type ReadOption func(r *ReadOptions)

// ReadFrom the database and table
func ReadFrom(database, table string) ReadOption {
	return func(r *ReadOptions) {
		r.Database = database
		r.Table = table
	}
}

// NewWriteOptions fills WriteOptions struct with opts slice
func NewWriteOptions(opts ...WriteOption) WriteOptions {
	options := WriteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// WriteOptions configures an individual Write operation
type WriteOptions struct {
	Database  string
	Table     string
	TTL       time.Duration
	Metadata  metadata.Metadata
	Namespace string
	Context   context.Context
}

// WriteOption sets values in WriteOptions
type WriteOption func(w *WriteOptions)

// WriteTo the database and table
func WriteTo(database, table string) WriteOption {
	return func(w *WriteOptions) {
		w.Database = database
		w.Table = table
	}
}

// WriteTTL is the time the record expires
func WriteTTL(d time.Duration) WriteOption {
	return func(w *WriteOptions) {
		w.TTL = d
	}
}

// WriteMetadata add metadata.Metadata
func WriteMetadata(md metadata.Metadata) WriteOption {
	return func(w *WriteOptions) {
		w.Metadata = metadata.Copy(md)
	}
}

// NewDeleteOptions fills DeleteOptions struct with opts slice
func NewDeleteOptions(opts ...DeleteOption) DeleteOptions {
	options := DeleteOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// DeleteOptions configures an individual Delete operation
type DeleteOptions struct {
	Database  string
	Table     string
	Namespace string
	Context   context.Context
}

// DeleteOption sets values in DeleteOptions
type DeleteOption func(d *DeleteOptions)

// DeleteFrom the database and table
func DeleteFrom(database, table string) DeleteOption {
	return func(d *DeleteOptions) {
		d.Database = database
		d.Table = table
	}
}

// NewListOptions fills ListOptions struct with opts slice
func NewListOptions(opts ...ListOption) ListOptions {
	options := ListOptions{}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// ListOptions configures an individual List operation
type ListOptions struct {
	// List from the following
	Database, Table string
	// Prefix returns all keys that are prefixed with key
	Prefix string
	// Suffix returns all keys that end with key
	Suffix string
	// Limit limits the number of returned keys
	Limit uint
	// Offset when combined with Limit supports pagination
	Offset    uint
	Namespace string
	Context   context.Context
}

// ListOption sets values in ListOptions
type ListOption func(l *ListOptions)

// ListFrom the database and table
func ListFrom(database, table string) ListOption {
	return func(l *ListOptions) {
		l.Database = database
		l.Table = table
	}
}

// ListPrefix returns all keys that are prefixed with key
func ListPrefix(p string) ListOption {
	return func(l *ListOptions) {
		l.Prefix = p
	}
}

// ListSuffix returns all keys that end with key
func ListSuffix(s string) ListOption {
	return func(l *ListOptions) {
		l.Suffix = s
	}
}

// ListLimit limits the number of returned keys to l
func ListLimit(l uint) ListOption {
	return func(lo *ListOptions) {
		lo.Limit = l
	}
}

// ListOffset starts returning responses from o. Use in conjunction with Limit for pagination.
func ListOffset(o uint) ListOption {
	return func(l *ListOptions) {
		l.Offset = o
	}
}

// ExistsOption specifies Exists call options
type ExistsOption func(*ExistsOptions)

// ExistsOptions holds options for Exists method
type ExistsOptions struct {
	Namespace string
	Context   context.Context
}

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
