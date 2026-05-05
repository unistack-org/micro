package broker_test

import (
	"context"
	"testing"

	broker "go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	brokermemory "go.unistack.org/micro/v4/broker/memory"
)

type testCodec struct{}

func (c *testCodec) Marshal(v any, opts ...codec.Option) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	return nil, nil
}

func (c *testCodec) Unmarshal(b []byte, v any, opts ...codec.Option) error {
	return nil
}

func (c *testCodec) ContentType() string {
	return "application/octet-stream"
}

func (c *testCodec) String() string {
	return c.ContentType()
}

func TestNoopBroker_Name(t *testing.T) {
	b := broker.NewBroker(broker.Name("noop"))
	if name := b.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopBroker_Options(t *testing.T) {
	b := broker.NewBroker()
	opts := b.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopBroker_Address(t *testing.T) {
	b := broker.NewBroker(broker.Addrs(":0"))
	if addr := b.Address(); addr != ":0" {
		t.Errorf("expected ':0', got %q", addr)
	}
}

func TestNoopBroker_String(t *testing.T) {
	b := broker.NewBroker()
	if str := b.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopBroker_Init(t *testing.T) {
	b := broker.NewBroker()
	if err := b.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopBroker_ConnectDisconnect(t *testing.T) {
	b := broker.NewBroker()
	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := b.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
}

func TestNoopBroker_PublishSubscribe(t *testing.T) {
	c := &testCodec{}
	b := broker.NewBroker(
		broker.ContentType(c.ContentType()),
		broker.Codec(c.ContentType(), c),
	)
	ctx := context.Background()

	msg, err := b.NewMessage(ctx, metadata.Metadata{}, []byte("body"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := b.Publish(ctx, "topic", msg); err != nil {
		t.Errorf("unexpected error on Publish: %v", err)
	}

	sub, err := b.Subscribe(ctx, "topic", func(msg broker.Message) error {
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error on Subscribe: %v", err)
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Errorf("unexpected error on Unsubscribe: %v", err)
	}
}

func TestNoopBroker_LiveReadyHealth(t *testing.T) {
	b := broker.NewBroker()
	if !b.Live() {
		t.Error("expected broker to be live")
	}
	if !b.Ready() {
		t.Error("expected broker to be ready")
	}
	if !b.Health() {
		t.Error("expected broker to be healthy")
	}
}

func TestMemoryBroker_PublishSubscribe(t *testing.T) {
	c := &testCodec{}
	b := brokermemory.NewBroker(
		broker.Codec(c.ContentType(), c),
	)
	if err := b.Init(); err != nil {
		t.Fatalf("unexpected error on Init: %v", err)
	}

	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("unexpected error on Connect: %v", err)
	}
	defer b.Disconnect(ctx)

	received := make(chan bool, 1)
	_, err := b.Subscribe(ctx, "test-topic", func(msg broker.Message) error {
		received <- true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error on Subscribe: %v", err)
	}

	msg, err := b.NewMessage(ctx, metadata.Metadata{}, []byte("test"), broker.MessageContentType(c.ContentType()))
	if err != nil {
		t.Fatalf("unexpected error on NewMessage: %v", err)
	}
	if err := b.Publish(ctx, "test-topic", msg); err != nil {
		t.Errorf("unexpected error on Publish: %v", err)
	}

	select {
	case <-received:
		// OK
	case <-ctx.Done():
		t.Error("timeout waiting for message")
	}
}
