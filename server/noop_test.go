package server_test

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/server"
)

type TestHandler struct {
	t *testing.T
}

type TestMessage struct {
	Name string
}

func (h *TestHandler) SingleSubHandler(ctx context.Context, msg *codec.Frame) error {
	// fmt.Printf("msg %s\n", msg.Data)
	return nil
}

func (h *TestHandler) BatchSubHandler(ctxs []context.Context, msgs []*codec.Frame) error {
	if len(msgs) != 8 {
		h.t.Fatal("invalid number of messages received")
	}
	for idx := 0; idx < len(msgs); idx++ {
		md, _ := metadata.FromIncomingContext(ctxs[idx])
		_ = md
		//	fmt.Printf("msg md %v\n", md)
	}
	return nil
}

func TestNoopSub(t *testing.T) {
	ctx := context.Background()

	b := broker.NewBroker()

	if err := b.Init(); err != nil {
		t.Fatal(err)
	}

	if err := b.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	s := server.NewServer(
		server.Broker(b),
		server.Codec("application/octet-stream", codec.NewCodec()),
	)
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}

	c := client.NewClient(
		client.Broker(b),
		client.Codec("application/octet-stream", codec.NewCodec()),
		client.ContentType("application/octet-stream"),
	)
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	h := &TestHandler{t: t}

	if err := s.Subscribe(s.NewSubscriber("single_topic", h.SingleSubHandler,
		server.SubscriberQueue("queue"),
	)); err != nil {
		t.Fatal(err)
	}

	if err := s.Subscribe(s.NewSubscriber("batch_topic", h.BatchSubHandler,
		server.SubscriberQueue("queue"),
		server.SubscriberBatch(true),
	)); err != nil {
		t.Fatal(err)
	}

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}

	msgs := make([]client.Message, 0, 8)
	for i := 0; i < 8; i++ {
		msgs = append(msgs, c.NewMessage("batch_topic", &codec.Frame{Data: []byte(fmt.Sprintf(`{"name": "test_name %d"}`, i))}))
	}

	if err := c.BatchPublish(ctx, msgs); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()
}
