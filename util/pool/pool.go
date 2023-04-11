// Package pool is a connection pool
package pool // import "go.unistack.org/micro/v4/util/pool"

import (
	"context"
	"time"

	"go.unistack.org/micro/v4/network/transport"
)

// Pool is an interface for connection pooling
type Pool interface {
	// Close the pool
	Close() error
	// Get a connection
	Get(ctx context.Context, addr string, opts ...transport.DialOption) (Conn, error)
	// Release the connection
	Release(c Conn, status error) error
}

// Conn conn pool interface
type Conn interface {
	// unique id of connection
	ID() string
	// time it was created
	Created() time.Time
	// embedded connection
	transport.Client
}

// NewPool creates new connection pool
func NewPool(opts ...Option) Pool {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return newPool(options)
}
