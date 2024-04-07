// Package broker is an interface used for asynchronous messaging
package broker // import "go.unistack.org/micro/v4/broker"

import (
	"context"
	"errors"
	"time"

	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
)

// DefaultBroker default memory broker
var DefaultBroker Broker = NewBroker()

var (
	// ErrNotConnected returns when broker used but not connected yet
	ErrNotConnected = errors.New("broker not connected")
	// ErrDisconnected returns when broker disconnected
	ErrDisconnected = errors.New("broker disconnected")
	// ErrInvalidMessage returns when message has nvalid format
	ErrInvalidMessage = errors.New("broker message has invalid format")
	// DefaultGracefulTimeout
	DefaultGracefulTimeout = 5 * time.Second
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	// Name returns broker instance name
	Name() string
	// Init initilize broker
	Init(opts ...options.Option) error
	// Options returns broker options
	Options() Options
	// Address return configured address
	Address() string
	// Connect connects to broker
	Connect(ctx context.Context) error
	// Disconnect disconnect from broker
	Disconnect(ctx context.Context) error
	// Publish message, msg can be single broker.Message or []broker.Message
	Publish(ctx context.Context, msg interface{}, opts ...options.Option) error
	// Subscribe subscribes to topic message via handler
	Subscribe(ctx context.Context, topic string, handler interface{}, opts ...options.Option) (Subscriber, error)
	// String type of broker
	String() string
}

// Message is given to a subscription handler for processing
type Message interface {
	// Context for the message
	Context() context.Context
	// Topic
	Topic() string
	// Header returns message headers
	Header() metadata.Metadata
	// Body returns broker message may be []byte slice or some go struct
	Body() interface{}
	// Ack acknowledge message
	Ack() error
	// Error returns message error (like decoding errors or some other)
	// In this case Body contains raw []byte from broker
	Error() error
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

// MessageHandler func signature for single message processing
type MessageHandler func(Message) error

// MessagesHandler func signature for batch message processing
type MessagesHandler func([]Message) error
