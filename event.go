package micro

import (
	"context"

	"go.unistack.org/micro/v4/client"
)

// Event is used to publish messages to a topic
type Event interface {
	// Publish publishes a message to the event topic
	Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error
}

type event struct {
	c     client.Client
	topic string
}

// NewEvent creates a new event publisher
func NewEvent(topic string, c client.Client) Event {
	return &event{c, topic}
}

func (e *event) Publish(ctx context.Context, msg interface{}, opts ...client.PublishOption) error {
	return e.c.Publish(ctx, e.c.NewMessage(e.topic, msg), opts...)
}
