package tracer

import (
	"context"

	"go.unistack.org/micro/v4/logger"
)

type SpanStatus int

const (
	// SpanStatusUnset is the default status code.
	SpanStatusUnset SpanStatus = 0

	// SpanStatusError indicates the operation contains an error.
	SpanStatusError SpanStatus = 1

	// SpanStatusOK indicates operation has been validated by an Application developers
	// or Operator to have completed successfully, or contain no error.
	SpanStatusOK SpanStatus = 2
)

func (s SpanStatus) String() string {
	switch s {
	case SpanStatusUnset:
		return "Unset"
	case SpanStatusError:
		return "Error"
	case SpanStatusOK:
		return "OK"
	default:
		return "Unset"
	}
}

type SpanKind int

const (
	// SpanKindUnspecified is an unspecified SpanKind and is not a valid
	// SpanKind. SpanKindUnspecified should be replaced with SpanKindInternal
	// if it is received.
	SpanKindUnspecified SpanKind = 0
	// SpanKindInternal is a SpanKind for a Span that represents an internal
	// operation within an application.
	SpanKindInternal SpanKind = 1
	// SpanKindServer is a SpanKind for a Span that represents the operation
	// of handling a request from a client.
	SpanKindServer SpanKind = 2
	// SpanKindClient is a SpanKind for a Span that represents the operation
	// of client making a request to a server.
	SpanKindClient SpanKind = 3
	// SpanKindProducer is a SpanKind for a Span that represents the operation
	// of a producer sending a message to a message broker. Unlike
	// SpanKindClient and SpanKindServer, there is often no direct
	// relationship between this kind of Span and a SpanKindConsumer kind. A
	// SpanKindProducer Span will end once the message is accepted by the
	// message broker which might not overlap with the processing of that
	// message.
	SpanKindProducer SpanKind = 4
	// SpanKindConsumer is a SpanKind for a Span that represents the operation
	// of a consumer receiving a message from a message broker. Like
	// SpanKindProducer Spans, there is often no direct relationship between
	// this Span and the Span that produced the message.
	SpanKindConsumer SpanKind = 5
)

func (sk SpanKind) String() string {
	switch sk {
	case SpanKindInternal:
		return "internal"
	case SpanKindServer:
		return "server"
	case SpanKindClient:
		return "client"
	case SpanKindProducer:
		return "producer"
	case SpanKindConsumer:
		return "consumer"
	default:
		return "unspecified"
	}
}

// SpanOptions contains span option
type SpanOptions struct {
	Labels []interface{}
	Kind   SpanKind
}

// SpanOption func signature
type SpanOption func(o *SpanOptions)

// EventOptions contains event options
type EventOptions struct{}

// EventOption func signature
type EventOption func(o *EventOptions)

func WithSpanLabels(labels ...interface{}) SpanOption {
	return func(o *SpanOptions) {
		o.Labels = labels
	}
}

func WithSpanKind(k SpanKind) SpanOption {
	return func(o *SpanOptions) {
		o.Kind = k
	}
}

// Options struct
type Options struct {
	// Context used to store custome tracer options
	Context context.Context
	// Logger used for logging
	Logger logger.Logger
	// Name of the tracer
	Name string
}

// Option func signature
type Option func(o *Options)

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// NewSpanOptions returns default SpanOptions
func NewSpanOptions(opts ...SpanOption) SpanOptions {
	options := SpanOptions{
		Kind: SpanKindInternal,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// NewOptions returns default options
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger: logger.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

// Name sets the name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}
