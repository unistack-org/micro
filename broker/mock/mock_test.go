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

// imports used by later tasks — declared here to avoid repeated additions
var _ = fmt.Sprintf
var _ = time.Millisecond
var _ broker.Message
