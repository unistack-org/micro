// Package broker is an interface used for asynchronous messaging
package broker

import (
	"context"
	"errors"

	"github.com/unistack-org/micro/v3/metadata"
)

// DefaultBroker default broker
var DefaultBroker Broker = NewBroker()

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Name() string
	Init(...Option) error
	Options() Options
	Address() string
	Connect(context.Context) error
	Disconnect(context.Context) error
	Publish(context.Context, string, *Message, ...PublishOption) error
	Subscribe(context.Context, string, Handler, ...SubscribeOption) (Subscriber, error)
	String() string
}

// Handler is used to process messages via a subscription of a topic.
type Handler func(Event) error

// Event is given to a subscription handler for processing
type Event interface {
	Topic() string
	Message() *Message
	Ack() error
	Error() error
}

// RawMessage is a raw encoded JSON value.
// It implements Marshaler and Unmarshaler and can be used to delay decoding or precompute a encoding.
type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return nil, nil
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
	Header metadata.Metadata // contains message metadata
	Body   RawMessage        // contains message body
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe(context.Context) error
}
