package server_test

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/server"
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

func TestNoopSub(t *testing.T) {
	ctx := context.Background()

	b := broker.NewBroker()

	if err := b.Init(); err != nil {
		t.Fatal(err)
	}

	if err := b.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	if err := logger.DefaultLogger.Init(logger.WithLevel(logger.ErrorLevel)); err != nil {
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

func TestHooks_Wrap(t *testing.T) {
	n := 5
	fn1 := func(next server.FuncSubHandler) server.FuncSubHandler {
		return func(ctx context.Context, msg server.Message) (err error) {
			n *= 2
			return next(ctx, msg)
		}
	}
	fn2 := func(next server.FuncSubHandler) server.FuncSubHandler {
		return func(ctx context.Context, msg server.Message) (err error) {
			n -= 10
			return next(ctx, msg)
		}
	}

	hs := &options.Hooks{}
	hs.Append(server.HookSubHandler(fn1), server.HookSubHandler(fn2))

	var fn = func(ctx context.Context, msg server.Message) error {
		return nil
	}

	hs.EachPrev(func(hook options.Hook) {
		if h, ok := hook.(server.HookSubHandler); ok {
			fn = h(fn)
		}
	})

	if err := fn(nil, nil); err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Fatalf("uncorrected hooks call")
	}
}
