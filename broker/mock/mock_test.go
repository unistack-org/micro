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

	msg := mock.NewMockMessage(ctx, "orders", hdr, body)

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

var _ broker.Message
