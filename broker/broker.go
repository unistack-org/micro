// Package broker is an interface used for asynchronous messaging
package broker // import "go.unistack.org/micro/v4/broker"

import (
	"context"
	"errors"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
)

// DefaultBroker default memory broker
var DefaultBroker Broker // =  NewBroker()

var (
	// ErrNotConnected returns when broker used but not connected yet
	ErrNotConnected = errors.New("broker not connected")
	// ErrDisconnected returns when broker disconnected
	ErrDisconnected = errors.New("broker disconnected")
	// ErrInvalidMessage returns when message has nvalid format
	ErrInvalidMessage = errors.New("broker message has invalid format")
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
	// NewMessage creates new broker message
	NewMessage(endpoint string, req interface{}, opts ...MessageOption) Message
	// Publish message to broker topic
	Publish(ctx context.Context, msg interface{}, opts ...PublishOption) error
	// Subscribe subscribes to topic message via handler
	Subscribe(ctx context.Context, topic string, handler interface{}, opts ...SubscribeOption) (Subscriber, error)
	// String type of broker
	String() string
}

// Message is given to a subscription handler for processing
type Message interface {
	// Context for the message
	Context() context.Context
	// Topic returns event topic
	Topic() string
	// Body returns broker message
	Body() interface{}
	// Ack acknowledge message
	Ack() error
	// Error returns message error (like decoding errors or some other)
	// In this case Body contains raw []byte from broker
	Error() error
}

// RawMessage is used to transfer data
type RawMessage struct {
	// Header contains message metadata
	Header metadata.Metadata
	// Body contains message body
	Body codec.RawMessage
}

// NewMessage create broker message with topic filled
func NewRawMessage(topic string) *RawMessage {
	m := &RawMessage{Header: metadata.New(2)}
	m.Header.Set(metadata.HeaderTopic, topic)
	return m
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
