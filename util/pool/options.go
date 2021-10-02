package pool

import (
	"time"

	"go.unistack.org/micro/v3/network/transport"
)

// Options struct
type Options struct {
	Transport transport.Transport
	TTL       time.Duration
	Size      int
}

// Option func signature
type Option func(*Options)

// Size sets the size
func Size(i int) Option {
	return func(o *Options) {
		o.Size = i
	}
}

// Transport sets the transport
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// TTL specifies ttl
func TTL(t time.Duration) Option {
	return func(o *Options) {
		o.TTL = t
	}
}
