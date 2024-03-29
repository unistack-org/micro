// Package broker is an interface used for asynchronous messaging
package broker // import "go.unistack.org/micro/v3/broker"

import (
	"context"
	"errors"

	"go.unistack.org/micro/v3/metadata"
)

// DefaultBroker default memory broker
var DefaultBroker = NewBroker()

var (
	// ErrNotConnected returns when broker used but not connected yet
	ErrNotConnected = errors.New("broker not connected")
	// ErrDisconnected returns when broker disconnected
	ErrDisconnected = errors.New("broker disconnected")
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
	// Publish message to broker topic
	Publish(ctx context.Context, topic string, msg *Message, opts ...PublishOption) error
	// Subscribe subscribes to topic message via handler
	Subscribe(ctx context.Context, topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
	// BatchPublish messages to broker with multiple topics
	BatchPublish(ctx context.Context, msgs []*Message, opts ...PublishOption) error
	// BatchSubscribe subscribes to topic messages via handler
	BatchSubscribe(ctx context.Context, topic string, h BatchHandler, opts ...SubscribeOption) (Subscriber, error)
	// String type of broker
	String() string
}

// Handler is used to process messages via a subscription of a topic.
type Handler func(Event) error

// Events contains multiple events
type Events []Event

// Ack try to ack all events and return
func (evs Events) Ack() error {
	var err error
	for _, ev := range evs {
		if err = ev.Ack(); err != nil {
			return err
		}
	}
	return nil
}

// SetError sets error on event
func (evs Events) SetError(err error) {
	for _, ev := range evs {
		ev.SetError(err)
	}
}

// BatchHandler is used to process messages in batches via a subscription of a topic.
type BatchHandler func(Events) error

// Event is given to a subscription handler for processing
type Event interface {
	// Topic returns event topic
	Topic() string
	// Message returns broker message
	Message() *Message
	// Ack acknowledge message
	Ack() error
	// Error returns message error (like decoding errors or some other)
	Error() error
	// SetError set event processing error
	SetError(err error)
}

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can be used to delay decoding or precompute a encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return *m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("RawMessage UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

// Message is used to transfer data
type Message struct {
	// Header contains message metadata
	Header metadata.Metadata
	// Body contains message body
	Body RawMessage
}

// NewMessage create broker message with topic filled
func NewMessage(topic string) *Message {
	m := &Message{Header: metadata.New(2)}
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
