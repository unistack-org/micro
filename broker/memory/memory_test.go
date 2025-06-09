package broker

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/atomic"
	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
)

type hldr struct {
	c atomic.Int64
}

func (h *hldr) Handler(m broker.Message) error {
	h.c.Add(1)
	return nil
}

func TestMemoryBroker(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()

	if err := b.Init(); err != nil {
		t.Fatalf("Unexpected init error %v", err)
	}

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}

	topic := "test"
	count := int64(10)

	h := &hldr{}

	sub, err := b.Subscribe(ctx, topic, h.Handler)
	if err != nil {
		t.Fatalf("Unexpected error subscribing %v", err)
	}

	for i := int64(0); i < count; i++ {
		message, err := b.NewMessage(ctx,
			metadata.Pairs(
				"foo", "bar",
				"id", fmt.Sprintf("%d", i),
			),
			[]byte(`"hello world"`),
			broker.MessageContentType("application/octet-stream"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := b.Publish(ctx, topic, message); err != nil {
			t.Fatalf("Unexpected error publishing %d err: %v", i, err)
		}
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unexpected error unsubscribing from %s: %v", topic, err)
	}

	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Unexpected connect error %v", err)
	}

	if h.c.Load() != count {
		t.Fatal("invalid messages count received")
	}
}
