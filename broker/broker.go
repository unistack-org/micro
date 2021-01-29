// Package broker is an interface used for asynchronous messaging
package broker

import (
	"context"

	"github.com/unistack-org/micro/v3/metadata"
)

var (
	// DefaultBroker default broker
	DefaultBroker Broker = NewBroker()
)

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

// Message is used to transfer data
type Message struct {
	Header metadata.Metadata // contains message metadata
	Body   []byte            // contains message body
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe(context.Context) error
}
