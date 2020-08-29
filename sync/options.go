package sync

import (
	"time"

	"github.com/unistack-org/micro/v3/logger"
)

type Options struct {
	Nodes  []string
	Prefix string
	Logger logger.Logger
}

type Option func(o *Options)

type LeaderOptions struct{}

type LeaderOption func(o *LeaderOptions)

type LockOptions struct {
	TTL  time.Duration
	Wait time.Duration
}

type LockOption func(o *LockOptions)

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
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
