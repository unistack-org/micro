package broker

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"
	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/metadata"
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

	for i := range count {
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

func TestMemoryBroker_ConnectDisconnect(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()

	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	if err := b.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect error: %v", err)
	}
	// Test double connect
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Double connect error: %v", err)
	}
	// Test disconnect when not connected
	b2 := NewBroker()
	if err := b2.Disconnect(ctx); err != nil {
		t.Fatalf("Disconnect when not connected error: %v", err)
	}
}

func TestMemoryBroker_String(t *testing.T) {
	b := NewBroker()
	if str := b.String(); str != "memory" {
		t.Errorf("expected 'memory', got %q", str)
	}
}

func TestMemoryBroker_Name(t *testing.T) {
	b := NewBroker(broker.Name("test-memory"))
	if name := b.Name(); name != "test-memory" {
		t.Errorf("expected 'test-memory', got %q", name)
	}
}

func TestMemoryBroker_LiveReadyHealth(t *testing.T) {
	b := NewBroker()
	if !b.Live() {
		t.Error("expected Live to be true")
	}
	if !b.Ready() {
		t.Error("expected Ready to be true")
	}
	if !b.Health() {
		t.Error("expected Health to be true")
	}
}

func TestMemoryBroker_Address(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	addr := b.Address()
	if addr == "" {
		t.Error("expected non-empty address")
	}
}

func TestMemoryBroker_MessageMethods(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()

	hdr := metadata.Pairs("key", "val")
	body := []byte("test-body")
	msg, err := b.NewMessage(ctx, hdr, body, broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}

	if string(msg.Body()) != "test-body" {
		t.Errorf("Body: expected test-body, got %s", string(msg.Body()))
	}
	gotHdr := msg.Header()
	if len(gotHdr["key"]) == 0 || gotHdr["key"][0] != "val" {
		t.Errorf("Header: expected val, got %v", gotHdr["key"])
	}
	if msg.Context() != ctx {
		t.Error("Context: expected original context")
	}
	if msg.Topic() != "" {
		t.Errorf("Topic: expected empty string, got %s", msg.Topic())
	}
	if err := msg.Ack(); err != nil {
		t.Errorf("Ack: unexpected error %v", err)
	}
	if msg.Error() != nil {
		t.Errorf("Error: expected nil, got %v", msg.Error())
	}
	// Test Unmarshal
	var dst []byte
	if err := msg.Unmarshal(&dst); err != nil {
		t.Errorf("Unmarshal: unexpected error %v", err)
	}
}

func TestMemoryBroker_ConcurrentSubscribe(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()

	var wg sync.WaitGroup
	subCount := 5
	msgCount := 10
	received := make([]atomic.Int64, subCount)

	for i := range subCount {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := b.Subscribe(ctx, "concurrent-topic", func(msg broker.Message) error {
				received[idx].Add(1)
				return nil
			})
			if err != nil {
				t.Errorf("Subscribe error: %v", err)
			}
		}(i)
	}
	wg.Wait()

	for range msgCount {
		msg, err := b.NewMessage(ctx, metadata.Pairs(), []byte("test"), broker.MessageContentType("application/octet-stream"))
		if err != nil {
			t.Fatalf("NewMessage error: %v", err)
		}
		if err := b.Publish(ctx, "concurrent-topic", msg); err != nil {
			t.Errorf("Publish error: %v", err)
		}
	}

	time.Sleep(100 * time.Millisecond)

	for i := range subCount {
		if received[i].Load() != int64(msgCount) {
			t.Errorf("sub %d: expected %d, got %d", i, msgCount, received[i].Load())
		}
	}
}

func TestMemoryBroker_Unsubscribe(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()

	sub, err := b.Subscribe(ctx, "unsub-topic", func(msg broker.Message) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}

	if err := sub.Unsubscribe(ctx); err != nil {
		t.Errorf("Unsubscribe error: %v", err)
	}

	// Publish after unsubscribe should not be received
	msg, err := b.NewMessage(ctx, metadata.Pairs(), []byte("test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "unsub-topic", msg); err != nil {
		t.Errorf("Publish error: %v", err)
	}
}

func TestMemoryBroker_WithOptions(t *testing.T) {
	ctx := context.Background()
	b := NewBroker(
		broker.Name("test"),
		broker.Addrs("localhost:0"),
		broker.ContentType("application/octet-stream"),
	)
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()

	if b.Name() != "test" {
		t.Errorf("Name: expected test, got %s", b.Name())
	}
	if len(b.Options().Addrs) != 1 {
		t.Errorf("Addrs: expected 1 addr, got %d", len(b.Options().Addrs))
	}
}

func TestMemoryBroker_Init(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	// Test double init
	if err := b.Init(); err != nil {
		t.Fatalf("First init error: %v", err)
	}
	if err := b.Init(); err != nil {
		t.Fatalf("Second init error: %v", err)
	}
	// Test init with hooks
	hookCalled := false
	if err := b.Init(broker.Hooks(broker.HookPublish(func(next broker.FuncPublish) broker.FuncPublish {
		return func(ctx context.Context, topic string, msg ...broker.Message) error {
			hookCalled = true
			return next(ctx, topic, msg...)
		}
	}))); err != nil {
		t.Fatalf("Init with hook error: %v", err)
	}
	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()
	// Publish to trigger hook
	msg, err := b.NewMessage(ctx, metadata.Pairs(), []byte("test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "init-test", msg); err != nil {
		t.Errorf("Publish error: %v", err)
	}
	if !hookCalled {
		t.Error("Hook not called")
	}
}

func TestMemoryBroker_PublishErrors(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	// Test publish when not connected
	msg, err := b.NewMessage(ctx, metadata.Pairs(), []byte("test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "test", msg); err != broker.ErrNotConnected {
		t.Errorf("Expected ErrNotConnected, got %v", err)
	}
	// Connect and test invalid handler
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()
	// Subscribe with invalid handler
	_, err = b.Subscribe(ctx, "invalid-handler", "invalid")
	if err != broker.ErrInvalidHandler {
		t.Errorf("Expected ErrInvalidHandler, got %v", err)
	}
}

func TestMemoryBroker_ConnectContextCancel(t *testing.T) {
	b := NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := b.Connect(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestMemoryBroker_DisconnectContextCancel(t *testing.T) {
	b := NewBroker()
	ctx := context.Background()
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := b.Disconnect(ctx)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestMemoryBroker_newCodecErrors(t *testing.T) {
	b := NewBroker()
	br, ok := b.(*Broker)
	if !ok {
		t.Fatal("failed to cast to *Broker")
	}
	_, err := br.newCodec("unknown-content-type")
	if err != codec.ErrUnknownContentType {
		t.Errorf("Expected ErrUnknownContentType, got %v", err)
	}
}

func TestMemoryBroker_PublishEdgeCases(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()

	// Test publish with no subscribers
	msg, err := b.NewMessage(ctx, metadata.Pairs(), []byte("test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "no-sub", msg); err != nil {
		t.Errorf("Publish with no subs error: %v", err)
	}

	// Test publish with batch handler
	batchReceived := make(chan []broker.Message, 1)
	_, err = b.Subscribe(ctx, "batch-topic", func(msgs []broker.Message) error {
		batchReceived <- msgs
		return nil
	})
	if err != nil {
		t.Fatalf("Subscribe batch error: %v", err)
	}
	msg, err = b.NewMessage(ctx, metadata.Pairs(), []byte("batch-test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "batch-topic", msg); err != nil {
		t.Errorf("Publish batch error: %v", err)
	}
	select {
	case msgs := <-batchReceived:
		if len(msgs) != 1 {
			t.Errorf("Expected 1 message, got %d", len(msgs))
		}
	case <-time.After(time.Second):
		t.Error("Timeout waiting for batch message")
	}

	// Test publish with handler returning error
	handlerErr := fmt.Errorf("test error")
	_, err = b.Subscribe(ctx, "error-topic", func(msg broker.Message) error {
		return handlerErr
	})
	if err != nil {
		t.Fatalf("Subscribe error error: %v", err)
	}
	msg, err = b.NewMessage(ctx, metadata.Pairs(), []byte("error-test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "error-topic", msg); err != nil {
		t.Errorf("Publish error topic error: %v", err)
	}

	// Test publish with AutoAck false
	ackReceived := make(chan bool, 1)
	_, err = b.Subscribe(ctx, "autoack-topic", func(msg broker.Message) error {
		ackReceived <- true
		return nil
	}, broker.SubscribeAutoAck(false))
	if err != nil {
		t.Fatalf("Subscribe autoack error: %v", err)
	}
	msg, err = b.NewMessage(ctx, metadata.Pairs(), []byte("autoack-test"), broker.MessageContentType("application/octet-stream"))
	if err != nil {
		t.Fatalf("NewMessage error: %v", err)
	}
	if err := b.Publish(ctx, "autoack-topic", msg); err != nil {
		t.Errorf("Publish autoack error: %v", err)
	}
	select {
	case <-ackReceived:
		// Message received, but Ack not called
	case <-time.After(time.Second):
		t.Error("Timeout waiting for autoack message")
	}
}

func TestSubscriberMethods(t *testing.T) {
	b := NewBroker(broker.Codec("application/octet-stream", codec.NewCodec()))
	ctx := context.Background()
	if err := b.Init(); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	if err := b.Connect(ctx); err != nil {
		t.Fatalf("Connect error: %v", err)
	}
	defer func() { _ = b.Disconnect(ctx) }()

	sub, err := b.Subscribe(ctx, "sub-methods", func(msg broker.Message) error {
		return nil
	}, broker.SubscribeGroup("test-group"))
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}

	if sub.Topic() != "sub-methods" {
		t.Errorf("Topic: expected sub-methods, got %s", sub.Topic())
	}
	if sub.Options().Group != "test-group" {
		t.Errorf("Group: expected test-group, got %s", sub.Options().Group)
	}
}
