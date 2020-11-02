// Package transport is an interface for synchronous connection based communication
package transport

import (
	"time"
)

var (
	// DefaultTransport is the global default transport
	DefaultTransport Transport = NewTransport()
)

// Transport is an interface which is used for communication between
// services. It uses connection based socket send/recv semantics and
// has various implementations; http, grpc, quic.
type Transport interface {
	Init(...Option) error
	Options() Options
	Dial(addr string, opts ...DialOption) (Client, error)
	Listen(addr string, opts ...ListenOption) (Listener, error)
	String() string
}

// Message is used to transfer data
type Message struct {
	Header map[string]string
	Body   []byte
}

// Socket bastraction interface
type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
	Local() string
	Remote() string
}

// Client is the socket owner
type Client interface {
	Socket
}

// Listener is the interface for stream oriented messaging
type Listener interface {
	Addr() string
	Close() error
	Accept(func(Socket)) error
}

// Option is the option signature
type Option func(*Options)

// DialOption is the option signature
type DialOption func(*DialOptions)

// ListenOption is the option signature
type ListenOption func(*ListenOptions)

var (
	// Default dial timeout
	DefaultDialTimeout = time.Second * 5
)
