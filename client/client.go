// Package client is an interface for an RPC client
package client // import "go.unistack.org/micro/v3/client"

import (
	"context"
	"time"

	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/metadata"
)

var (
	// DefaultClient is the global default client
	DefaultClient Client = NewClient()
	// DefaultContentType is the default content-type if not specified
	DefaultContentType = "application/json"
	// DefaultBackoff is the default backoff function for retries
	DefaultBackoff = exponentialBackoff
	// DefaultRetry is the default check-for-retry function for retries
	DefaultRetry = RetryNever
	// DefaultRetries is the default number of times a request is tried
	DefaultRetries = 0
	// DefaultRequestTimeout is the default request timeout
	DefaultRequestTimeout = time.Second * 5
	// DefaultPoolSize sets the connection pool size
	DefaultPoolSize = 100
	// DefaultPoolTTL sets the connection pool ttl
	DefaultPoolTTL = time.Minute
)

// Client is the interface used to make requests to services.
// It supports Request/Response via Transport and Publishing via the Broker.
// It also supports bidirectional streaming of requests.
type Client interface {
	Name() string
	Init(opts ...Option) error
	Options() Options
	NewMessage(topic string, msg interface{}, opts ...MessageOption) Message
	NewRequest(service string, endpoint string, req interface{}, opts ...RequestOption) Request
	Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error
	Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error)
	Publish(ctx context.Context, msg Message, opts ...PublishOption) error
	BatchPublish(ctx context.Context, msg []Message, opts ...PublishOption) error
	String() string
}

// Message is the interface for publishing asynchronously
type Message interface {
	Topic() string
	Payload() interface{}
	ContentType() string
	Metadata() metadata.Metadata
}

// Request is the interface for a synchronous request used by Call or Stream
type Request interface {
	// The service to call
	Service() string
	// The action to take
	Method() string
	// The endpoint to invoke
	Endpoint() string
	// The content type
	ContentType() string
	// The unencoded request body
	Body() interface{}
	// Write to the encoded request writer. This is nil before a call is made
	Codec() codec.Codec
	// indicates whether the request will be a streaming one rather than unary
	Stream() bool
}

// Response is the response received from a service
type Response interface {
	// Read the response
	Codec() codec.Codec
	// read the header
	Header() metadata.Metadata
	// Read the undecoded response
	Read() ([]byte, error)
}

// Stream is the interface for a bidirectional synchronous stream
type Stream interface {
	// Context for the stream
	Context() context.Context
	// The request made
	Request() Request
	// The response read
	Response() Response
	// Send will encode and send a request
	Send(msg interface{}) error
	// Recv will decode and read a response
	Recv(msg interface{}) error
	// SendMsg will encode and send a request
	SendMsg(msg interface{}) error
	// RecvMsg will decode and read a response
	RecvMsg(msg interface{}) error
	// Error returns the stream error
	Error() error
	// Close closes the stream
	Close() error
	// CloseSend closes the send direction of the stream
	CloseSend() error
}

// Option used by the Client
type Option func(*Options)

// CallOption used by Call or Stream
type CallOption func(*CallOptions)

// PublishOption used by Publish
type PublishOption func(*PublishOptions)

// MessageOption used by NewMessage
type MessageOption func(*MessageOptions)

// RequestOption used by NewRequest
type RequestOption func(*RequestOptions)
