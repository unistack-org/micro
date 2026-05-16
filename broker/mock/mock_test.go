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

// keep time import used
var _ = time.Millisecond
var _ broker.Message
