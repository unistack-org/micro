// Package server is an interface for a micro server
package server // import "go.unistack.org/micro/v3/server"

import (
	"context"
	"time"

	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/util/id"
)

// DefaultServer default server
var DefaultServer Server = NewServer()

var (
	// DefaultAddress will be used if no address passed, use secure localhost
	DefaultAddress = "127.0.0.1:0"
	// DefaultName will be used if no name passed
	DefaultName = "server"
	// DefaultVersion will be used if no version passed
	DefaultVersion = "latest"
	// DefaultID will be used if no id passed
	DefaultID = id.Must()
	// DefaultRegisterCheck holds func that run before register server
	DefaultRegisterCheck = func(context.Context) error { return nil }
	// DefaultRegisterInterval holds interval for register
	DefaultRegisterInterval = time.Second * 30
	// DefaultRegisterTTL holds register record ttl, must be multiple of DefaultRegisterInterval
	DefaultRegisterTTL = time.Second * 90
	// DefaultNamespace will be used if no namespace passed
	DefaultNamespace = "micro"
	// DefaultMaxMsgSize holds default max msg ssize
	DefaultMaxMsgSize = 1024 * 1024 * 4 // 4Mb
	// DefaultMaxMsgRecvSize holds default max recv size
	DefaultMaxMsgRecvSize = 1024 * 1024 * 4 // 4Mb
	// DefaultMaxMsgSendSize holds default max send size
	DefaultMaxMsgSendSize = 1024 * 1024 * 4 // 4Mb
)

// Server is a simple micro server abstraction
type Server interface {
	// Name returns server name
	Name() string
	// Initialise options
	Init(...Option) error
	// Retrieve the options
	Options() Options
	// Register a handler
	Handle(h Handler) error
	// Create a new handler
	NewHandler(h interface{}, opts ...HandlerOption) Handler
	// Create a new subscriber
	NewSubscriber(topic string, h interface{}, opts ...SubscriberOption) Subscriber
	// Register a subscriber
	Subscribe(s Subscriber) error
	// Start the server
	Start() error
	// Stop the server
	Stop() error
	// Server implementation
	String() string
}

/*
// Router handle serving messages
type Router interface {
	// ProcessMessage processes a message
	ProcessMessage(ctx context.Context, msg Message) error
	// ServeRequest processes a request to completion
	ServeRequest(ctx context.Context, req Request, rsp Response) error
}
*/

// Message is an async message interface
type Message interface {
	// Topic of the message
	Topic() string
	// The decoded payload value
	Payload() interface{}
	// The content type of the payload
	ContentType() string
	// The raw headers of the message
	Header() metadata.Metadata
	// The raw body of the message
	Body() []byte
	// Codec used to decode the message
	Codec() codec.Codec
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
	Context() context.Context
	Request() Request
	Send(msg interface{}) error
	Recv(msg interface{}) error
	Error() error
	Close() error
}

// Handler interface represents a request handler. It's generated
// by passing any type of public concrete object with endpoints into server.NewHandler.
// Most will pass in a struct.
//
// Example:
//
//      type Greeter struct {}
//
//      func (g *Greeter) Hello(context, request, response) error {
//              return nil
//      }
//
type Handler interface {
	Name() string
	Handler() interface{}
	Endpoints() []*register.Endpoint
	Options() HandlerOptions
}

// Subscriber interface represents a subscription to a given topic using
// a specific subscriber function or object with endpoints. It mirrors
// the handler in its behaviour.
type Subscriber interface {
	Topic() string
	Subscriber() interface{}
	Endpoints() []*register.Endpoint
	Options() SubscriberOptions
}
