package broker

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v4/metadata"
)

func TestMemoryBatchBroker(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}

	topic := "test"
	count := 10

	fn := func(evts []Message) error {
		var err error
		for _, evt := range evts {
			if err = evt.Ack(); err != nil {
				return err
			}
		}
		return nil
	}

	sub, err := b.Subscribe(ctx, topic, fn)
	if err != nil {
		t.Fatalf("Unexpected error subscribing %v", err)
	}

	msgs := make([]Message, 0, count)
	for i := 0; i < count; i++ {
		message := &memoryMessage{
			header: map[string]string{
				metadata.HeaderTopic: topic,
				"foo":                "bar",
				"id":                 fmt.Sprintf("%d", i),
			},
			body: []byte(`"hello world"`),
		}
		msgs = append(msgs, message)
	}

	if err := b.Publish(ctx, msgs); err != nil {
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

	fn := func(p Message) error {
		return p.Ack()
	}

	sub, err := b.Subscribe(ctx, topic, fn)
	if err != nil {
		t.Fatalf("Unexpected error subscribing %v", err)
	}

	msgs := make([]Message, 0, count)
	for i := 0; i < count; i++ {
		message := &memoryMessage{
			header: map[string]string{
				metadata.HeaderTopic: topic,
				"foo":                "bar",
				"id":                 fmt.Sprintf("%d", i),
			},
			body: []byte(`"hello world"`),
		}
		msgs = append(msgs, message)
	}

	if err := b.Publish(ctx, msgs); err != nil {
		t.Fatalf("Unexpected error publishing %v", err)
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unexpected error unsubscribing from %s: %v", topic, err)
	}

	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}
}
