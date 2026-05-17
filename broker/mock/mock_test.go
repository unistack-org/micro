package mock_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/broker/mock"
	"go.unistack.org/micro/v5/metadata"
)

func TestNewMockMessage(t *testing.T) {
	ctx := context.Background()
	hdr := metadata.Metadata{"key": []string{"val"}}
	body := []byte(`{"id":1}`)

	msg := mock.NewMockMessage(ctx, "orders", hdr, body, nil)

	if msg.Topic() != "orders" {
		t.Fatalf("want topic %q, got %q", "orders", msg.Topic())
	}
	if string(msg.Body()) != string(body) {
		t.Fatalf("want body %q, got %q", body, msg.Body())
	}
	if msg.Context() != ctx {
		t.Fatal("context mismatch")
	}
	if v, ok := msg.Header()["key"]; !ok || len(v) == 0 || v[0] != "val" {
		t.Fatal("header mismatch")
	}
	if err := msg.Ack(); err != nil {
		t.Fatalf("Ack: %v", err)
	}
	if err := msg.Error(); err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func TestConnect(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestConnect_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	if err := b.Connect(ctx); err == nil {
		t.Fatal("expected error for unexpected Connect, got nil")
	}
}

func TestConnect_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("conn failed")
	b.ExpectConnect().WillReturnError(want)

	if err := b.Connect(ctx); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestDisconnect(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectDisconnect()

	_ = b.Connect(ctx)
	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestDisconnect_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	if err := b.Disconnect(ctx); err == nil {
		t.Fatal("expected error for unexpected Disconnect, got nil")
	}
}

func TestDisconnect_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	want := fmt.Errorf("disconnect failed")
	b.ExpectDisconnect().WillReturnError(want)

	_ = b.Connect(ctx)
	if err := b.Disconnect(ctx); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestExpectationsWereMet_Unfulfilled(t *testing.T) {
	b := mock.NewMockBroker()
	b.ExpectConnect()

	if err := b.ExpectationsWereMet(); err == nil {
		t.Fatal("expected error for unfulfilled expectation, got nil")
	}
}

func TestExpectationsWereMet_Empty(t *testing.T) {
	b := mock.NewMockBroker()
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestPublish(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectPublish("orders")

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestPublish_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err == nil {
		t.Fatal("expected error for unexpected Publish, got nil")
	}
}

func TestPublish_WrongTopic(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectPublish("invoices")

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err == nil {
		t.Fatal("expected error for wrong topic, got nil")
	}
}

func TestPublish_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("publish failed")
	b.ExpectConnect()
	b.ExpectPublish("orders").WillReturnError(want)

	_ = b.Connect(ctx)
	if err := b.Publish(ctx, "orders"); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestPublish_NotConnected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectPublish("orders")

	if err := b.Publish(ctx, "orders"); err != broker.ErrNotConnected {
		t.Fatalf("want ErrNotConnected, got %v", err)
	}
}

func TestPublish_Delay(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectPublish("orders").WillDelayFor(50 * time.Millisecond)

	_ = b.Connect(ctx)
	start := time.Now()
	_ = b.Publish(ctx, "orders")
	if time.Since(start) < 50*time.Millisecond {
		t.Fatal("expected delay of at least 50ms")
	}
}

func TestSubscribe(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)
	sub, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if sub == nil {
		t.Fatal("expected non-nil subscriber")
	}
	if sub.Topic() != "orders" {
		t.Fatalf("want topic %q, got %q", "orders", sub.Topic())
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestSubscribe_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()

	_ = b.Connect(ctx)
	_, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err == nil {
		t.Fatal("expected error for unexpected Subscribe, got nil")
	}
}

func TestSubscribe_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("subscribe failed")
	b.ExpectConnect()
	b.ExpectSubscribe("orders").WillReturnError(want)

	_ = b.Connect(ctx)
	_, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestSubscribe_NotConnected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectSubscribe("orders")

	_, err := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })
	if err != broker.ErrNotConnected {
		t.Fatalf("want ErrNotConnected, got %v", err)
	}
}

func TestInjectMessage_SingleHandler(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)

	received := make([]broker.Message, 0)
	_, _ = b.Subscribe(ctx, "orders", func(m broker.Message) error {
		received = append(received, m)
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{"x": []string{"y"}}, []byte(`{}`), nil)
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("want 1 message, got %d", len(received))
	}
	if string(received[0].Body()) != `{}` {
		t.Fatalf("unexpected body: %s", received[0].Body())
	}
}

func TestInjectMessage_BatchHandler(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)

	var batchReceived []broker.Message
	_, _ = b.Subscribe(ctx, "orders", func(msgs []broker.Message) error {
		batchReceived = append(batchReceived, msgs...)
		return nil
	})

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":2}`), nil)
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage: %v", err)
	}
	if len(batchReceived) != 1 {
		t.Fatalf("want 1 message in batch, got %d", len(batchReceived))
	}
}

func TestInjectMessage_NoHandlers(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	msg := mock.NewMockMessage(ctx, "orders", nil, []byte(`{}`), nil)
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage with no handlers should return nil, got: %v", err)
	}
}

func TestInjectMessage_HandlerError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)
	want := fmt.Errorf("handler error")
	_, _ = b.Subscribe(ctx, "orders", func(broker.Message) error { return want })

	msg := mock.NewMockMessage(ctx, "orders", nil, nil, nil)
	if err := b.InjectMessage(ctx, "orders", msg); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestUnsubscribe(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders")

	_ = b.Connect(ctx)
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unsubscribe: %v", err)
	}
	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}
}

func TestUnsubscribe_Unexpected(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")

	_ = b.Connect(ctx)
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })

	if err := sub.Unsubscribe(ctx); err == nil {
		t.Fatal("expected error for unexpected Unsubscribe, got nil")
	}
}

func TestUnsubscribe_ReturnsError(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	want := fmt.Errorf("unsubscribe failed")
	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders").WillReturnError(want)

	_ = b.Connect(ctx)
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error { return nil })

	if err := sub.Unsubscribe(ctx); err != want {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func TestUnsubscribe_RemovesHandler(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()
	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectUnsubscribe("orders")

	_ = b.Connect(ctx)

	called := false
	sub, _ := b.Subscribe(ctx, "orders", func(broker.Message) error {
		called = true
		return nil
	})
	_ = sub.Unsubscribe(ctx)

	msg := mock.NewMockMessage(ctx, "orders", nil, nil, nil)
	_ = b.InjectMessage(ctx, "orders", msg)

	if called {
		t.Fatal("handler should not be called after Unsubscribe")
	}
}

func TestFullLifecycle(t *testing.T) {
	ctx := context.Background()
	b := mock.NewMockBroker()

	b.ExpectConnect()
	b.ExpectSubscribe("orders")
	b.ExpectPublish("results")
	b.ExpectUnsubscribe("orders")
	b.ExpectDisconnect()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	processed := make([]string, 0)
	sub, err := b.Subscribe(ctx, "orders", func(m broker.Message) error {
		processed = append(processed, string(m.Body()))
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	msg := mock.NewMockMessage(ctx, "orders", metadata.Metadata{}, []byte(`{"id":42}`), nil)
	if err := b.InjectMessage(ctx, "orders", msg); err != nil {
		t.Fatalf("InjectMessage: %v", err)
	}

	if err := b.Publish(ctx, "results"); err != nil {
		t.Fatalf("Publish: %v", err)
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatalf("Unsubscribe: %v", err)
	}

	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect: %v", err)
	}

	if err := b.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet: %v", err)
	}

	if len(processed) != 1 || processed[0] != `{"id":42}` {
		t.Fatalf("unexpected processed messages: %v", processed)
	}
}
