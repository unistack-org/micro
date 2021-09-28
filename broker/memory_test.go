package broker

import (
	"context"
	"fmt"
	"testing"

	"github.com/unistack-org/micro/v3/metadata"
)

func TestMemoryBatchBroker(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}

	topic := "test"
	count := 10

	fn := func(evts Events) error {
		return evts.Ack()
	}

	sub, err := b.BatchSubscribe(ctx, topic, fn)
	if err != nil {
		t.Fatalf("Unexpected error subscribing %v", err)
	}

	msgs := make([]*Message, 0, 0)
	for i := 0; i < count; i++ {
		message := &Message{
			Header: map[string]string{
				metadata.HeaderTopic: topic,
				"foo":                "bar",
				"id":                 fmt.Sprintf("%d", i),
			},
			Body: []byte(`"hello world"`),
		}
		msgs = append(msgs, message)
	}

	if err := b.BatchPublish(ctx, msgs); err != nil {
		t.Fatalf("Unexpected error publishing %v", err)
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unexpected error unsubscribing from %s: %v", topic, err)
	}

	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}
}

func TestMemoryBroker(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}

	topic := "test"
	count := 10

	fn := func(p Event) error {
		return nil
	}

	sub, err := b.Subscribe(ctx, topic, fn)
	if err != nil {
		t.Fatalf("Unexpected error subscribing %v", err)
	}

	msgs := make([]*Message, 0, 0)
	for i := 0; i < count; i++ {
		message := &Message{
			Header: map[string]string{
				metadata.HeaderTopic: topic,
				"foo":                "bar",
				"id":                 fmt.Sprintf("%d", i),
			},
			Body: []byte(`"hello world"`),
		}
		msgs = append(msgs, message)

		if err := b.Publish(ctx, topic, message); err != nil {
			t.Fatalf("Unexpected error publishing %d err: %v", i, err)
		}
	}

	if err := b.BatchPublish(ctx, msgs); err != nil {
		t.Fatalf("Unexpected error publishing %v", err)
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unexpected error unsubscribing from %s: %v", topic, err)
	}

	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}
}
