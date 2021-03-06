package sync

import (
	"time"

	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/tracer"
)

// Options holds the sync options
type Options struct {
	// Logger used for logging
	Logger logger.Logger
	// Tracer used for tracing
	Tracer tracer.Tracer
	// Meter used for merics
	Meter meter.Meter
	// Prefix holds prefix?
	Prefix string
	// Nodes holds the nodes
	// TODO: change to Addrs ?
	Nodes []string
}

// Option func signature
type Option func(o *Options)

// NewOptions returns options that filled by opts
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger: logger.DefaultLogger,
		Meter:  meter.DefaultMeter,
		Tracer: tracer.DefaultTracer,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// LeaderOptions holds the leader options
type LeaderOptions struct{}

// LeaderOption func signature
type LeaderOption func(o *LeaderOptions)

// LockOptions holds the lock options
type LockOptions struct {
	TTL  time.Duration
	Wait time.Duration
}

// LockOption func signature
type LockOption func(o *LockOptions)

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Meter sets the logger
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// Tracer sets the tracer
func Tracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// Nodes sets the addresses to use
func Nodes(a ...string) Option {
	return func(o *Options) {
		o.Nodes = a
	}
}

// Prefix sets a prefix to any lock ids used
func Prefix(p string) Option {
	return func(o *Options) {
		o.Prefix = p
	}
}

// LockTTL sets the lock ttl
func LockTTL(t time.Duration) LockOption {
	return func(o *LockOptions) {
		o.TTL = t
	}
}

// LockWait sets the wait time
func LockWait(t time.Duration) LockOption {
	return func(o *LockOptions) {
		o.Wait = t
	}
}
