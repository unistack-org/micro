// Package broker is an interface used for asynchronous messaging
package broker

var (
	DefaultBroker Broker = newBroker()
)

// Broker is an interface used for asynchronous messaging.
type Broker interface {
	Init(...Option) error
	Options() Options
	Address() string
	Connect() error
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opts ...SubscribeOption) (Subscriber, error)
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
	Header map[string]string
	Body   []byte
	Error  error
}

// Subscriber is a convenience return type for the Subscribe method
type Subscriber interface {
	Options() SubscribeOptions
	Topic() string
	Unsubscribe() error
}
