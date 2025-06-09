// Package broker is an interface used for asynchronous messaging
package broker

import (
	"context"
	"errors"
	"time"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
)

// DefaultBroker default memory broker
var DefaultBroker Broker = NewBroker()

var (
	// ErrNotConnected returns when broker used but not connected yet
	ErrNotConnected = errors.New("broker not connected")
	// ErrDisconnected returns when broker disconnected
	ErrDisconnected = errors.New("broker disconnected")
	// ErrInvalidMessage returns when invalid Message passed
	ErrInvalidMessage = errors.New("invalid message")
	// ErrInvalidHandler returns when subscriber passed to Subscribe
	ErrInvalidHandler = errors.New("invalid handler, ony func(Message) error and func([]Message) error supported")
	// DefaultGracefulTimeout
	DefaultGracefulTimeout = 5 * time.Second
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	// Name returns broker instance name
	Name() string
	// Init initilize broker
	Init(opts ...Option) error
	// Options returns broker options
	Options() Options
	// Address return configured address
	Address() string
	// Connect connects to broker
	Connect(ctx context.Context) error
	// Disconnect disconnect from broker
	Disconnect(ctx context.Context) error
	// NewMessage create new broker message to publish.
	NewMessage(ctx context.Context, hdr metadata.Metadata, body interface{}, opts ...MessageOption) (Message, error)
	// Publish message to broker topic
	Publish(ctx context.Context, topic string, messages ...Message) error
	// Subscribe subscribes to topic message via handler
	Subscribe(ctx context.Context, topic string, handler interface{}, opts ...SubscribeOption) (Subscriber, error)
	// String type of broker
	String() string
	// Live returns broker liveness
	Live() bool
	// Ready returns broker readiness
	Ready() bool
	// Health returns broker health
	Health() bool
}

type (
	FuncPublish   func(ctx context.Context, topic string, messages ...Message) error
	HookPublish   func(next FuncPublish) FuncPublish
	FuncSubscribe func(ctx context.Context, topic string, handler interface{}, opts ...SubscribeOption) (Subscriber, error)
	HookSubscribe func(next FuncSubscribe) FuncSubscribe
)

// Message is given to a subscription handler for processing
type Message interface {
	// Context for the message.
	Context() context.Context
	// Topic returns message destination topic.
	Topic() string
	// Header returns message headers.
	Header() metadata.Metadata
	// Body returns broker message []byte slice
	Body() []byte
	// Unmarshal try to decode message body to dst.
	// This is helper method that uses codec.Unmarshal.
	Unmarshal(dst interface{}, opts ...codec.Option) error
	// Ack acknowledge message if supported.
	Ack() error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	// Options returns subscriber options
	Options() SubscribeOptions
	// Topic returns topic for subscription
	Topic() string
	// Unsubscribe from topic
	Unsubscribe(ctx context.Context) error
}
