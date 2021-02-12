package broker

import (
	"context"
	"fmt"
	"testing"
)

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

	for i := 0; i < count; i++ {
		message := &Message{
			Header: map[string]string{
				"foo": "bar",
				"id":  fmt.Sprintf("%d", i),
			},
			Body: []byte(`hello world`),
		}

		if err := b.Publish(ctx, topic, message); err != nil {
			t.Fatalf("Unexpected error publishing %d", i)
		}
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unexpected error unsubscribing from %s: %v", topic, err)
	}

	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}
}
