package store

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
)

// Options contains configuration for the Store
type Options struct {
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
	// Context should contain all implementation specific options, using context.WithValue.
	Context context.Context
}

// NewOptions creates options struct
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Option sets values in Options
type Option func(o *Options)

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

// ReadOptions configures an individual Read operation
type ReadOptions struct {
	Database string
	Table    string
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

// WriteOptions configures an individual Write operation
// If Expiry and TTL are set TTL takes precedence
type WriteOptions struct {
	Database string
	Table    string
	// Expiry is the time the record expires
	Expiry time.Time
	// TTL is the time until the record expires
	TTL      time.Duration
	Metadata metadata.Metadata
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

// WriteExpiry is the time the record expires
func WriteExpiry(t time.Time) WriteOption {
	return func(w *WriteOptions) {
		w.Expiry = t
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

// DeleteOptions configures an individual Delete operation
type DeleteOptions struct {
	Database, Table string
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
	Offset uint
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
