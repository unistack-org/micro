// Package server is an interface for a micro server
package server // import "go.unistack.org/micro/v4/server"

import (
	"context"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
)

// DefaultServer default server
var DefaultServer = NewServer()

var (
	// DefaultAddress will be used if no address passed, use secure localhost
	DefaultAddress = "127.0.0.1:0"
	// DefaultName will be used if no name passed
	DefaultName = "server"
	// DefaultVersion will be used if no version passed
	DefaultVersion = "latest"
	// DefaultRegisterCheck holds func that run before register server
	DefaultRegisterCheck = func(context.Context) error { return nil }
	// DefaultRegisterInterval holds interval for register
	DefaultRegisterInterval = time.Second * 30
	// DefaultRegisterTTL holds register record ttl, must be multiple of DefaultRegisterInterval
	DefaultRegisterTTL = time.Second * 90
	// DefaultMaxMsgSize holds default max msg size
	DefaultMaxMsgSize = 1024 * 1024 * 4 // 4Mb
	// DefaultMaxMsgRecvSize holds default max recv size
	DefaultMaxMsgRecvSize = 1024 * 1024 * 4 // 4Mb
	// DefaultMaxMsgSendSize holds default max send size
	DefaultMaxMsgSendSize = 1024 * 1024 * 4 // 4Mb
	// DefaultGracefulTimeout default time for graceful stop
	DefaultGracefulTimeout = 5 * time.Second
)

// Server is a simple micro server abstraction
type Server interface {
	// Name returns server name
	Name() string
	// Initialise options
	Init(...options.Option) error
	// Retrieve the options
	Options() Options
	// Create and register new handler
	Handle(h interface{}, opts ...options.Option) error
	// Start the server
	Start() error
	// Stop the server
	Stop() error
	// Server implementation
	String() string
}

// Request is a synchronous request interface
type Request interface {
	// Service name requested
	Service() string
	// The action requested
	Method() string
	// Endpoint name requested
	Endpoint() string
	// Content type provided
	ContentType() string
	// Header of the request
	Header() metadata.Metadata
	// Body is the initial decoded value
	Body() interface{}
	// Read the undecoded request body
	Read() ([]byte, error)
	// The encoded message stream
	Codec() codec.Codec
	// Indicates whether its a stream
	Stream() bool
}

// Response is the response writer for unencoded messages
type Response interface {
	// Encoded writer
	Codec() codec.Codec
	// Write the header
	WriteHeader(md metadata.Metadata)
	// write a response directly to the client
	Write([]byte) error
}

// Stream represents a stream established with a client.
// A stream can be bidirectional which is indicated by the request.
// The last error will be left in Error().
// EOF indicates end of the stream.
type Stream interface {
	// Context for the stream
	Context() context.Context
	// Request returns request
	Request() Request
	// Send will encode and send a request
	Send(msg interface{}) error
	// Recv will decode and read a response
	Recv(msg interface{}) error
	// SendMsg will encode and send a request
	SendMsg(msg interface{}) error
	// RecvMsg will decode and read a response
	RecvMsg(msg interface{}) error
	// Error returns stream error
	Error() error
	// Close closes the stream
	Close() error
}
