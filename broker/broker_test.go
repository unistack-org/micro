package broker_test

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	broker "go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/tracer"
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

func TestNewContext(t *testing.T) {
	b := broker.NewBroker()
	ctx := context.Background()
	newCtx := broker.NewContext(ctx, b)
	got, ok := broker.FromContext(newCtx)
	if !ok {
		t.Error("expected broker to be in context")
	}
	if got != b {
		t.Error("expected same broker instance")
	}
}

func TestFromContext(t *testing.T) {
	t.Run("nil context", func(t *testing.T) {
		_, ok := broker.FromContext(nil)
		if ok {
			t.Error("expected false for nil context")
		}
	})
	t.Run("empty context", func(t *testing.T) {
		_, ok := broker.FromContext(context.Background())
		if ok {
			t.Error("expected false for empty context")
		}
	})
	t.Run("with broker", func(t *testing.T) {
		b := broker.NewBroker()
		ctx := broker.NewContext(context.Background(), b)
		got, ok := broker.FromContext(ctx)
		if !ok {
			t.Error("expected ok for context with broker")
		}
		if got != b {
			t.Error("expected same broker")
		}
	})
}

func TestMustContext(t *testing.T) {
	t.Run("valid context", func(t *testing.T) {
		b := broker.NewBroker()
		ctx := broker.NewContext(context.Background(), b)
		got := broker.MustContext(ctx)
		if got != b {
			t.Error("expected same broker")
		}
	})
	t.Run("invalid context", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid context")
			}
		}()
		broker.MustContext(context.Background())
	})
}

func TestSetMessageOption(t *testing.T) {
	key := "test-key"
	val := "test-val"
	opt := broker.SetMessageOption(key, val)
	opts := broker.NewMessageOptions(opt)
	v := opts.Context.Value(key)
	if v != val {
		t.Errorf("expected %v, got %v", val, v)
	}
}

func TestSetSubscribeOption(t *testing.T) {
	key := "sub-key"
	val := "sub-val"
	opt := broker.SetSubscribeOption(key, val)
	opts := broker.NewSubscribeOptions(opt)
	v := opts.Context.Value(key)
	if v != val {
		t.Errorf("expected %v, got %v", val, v)
	}
}

func TestSetOption(t *testing.T) {
	key := "opt-key"
	val := "opt-val"
	opt := broker.SetOption(key, val)
	b := broker.NewBroker(opt)
	v := b.Options().Context.Value(key)
	if v != val {
		t.Errorf("expected %v, got %v", val, v)
	}
}

func TestOptions_Context(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "val")
	b := broker.NewBroker(broker.Context(ctx))
	if b.Options().Context.Value("key") != "val" {
		t.Error("Context option not applied")
	}
}

func TestOptions_GracefulTimeout(t *testing.T) {
	timeout := 5 * time.Second
	b := broker.NewBroker(broker.GracefulTimeout(timeout))
	if b.Options().GracefulTimeout != timeout {
		t.Errorf("expected %v, got %v", timeout, b.Options().GracefulTimeout)
	}
}

func TestOptions_ContentType(t *testing.T) {
	ct := "application/json"
	b := broker.NewBroker(broker.ContentType(ct))
	if b.Options().ContentType != ct {
		t.Errorf("expected %s, got %s", ct, b.Options().ContentType)
	}
}

func TestOptions_ErrorHandler(t *testing.T) {
	handler := func(msg broker.Message) error { return nil }
	b := broker.NewBroker(broker.ErrorHandler(handler))
	if b.Options().ErrorHandler == nil {
		t.Error("ErrorHandler not set")
	}
}

func TestOptions_Codec(t *testing.T) {
	c := &testCodec{}
	b := broker.NewBroker(broker.Codec("application/octet-stream", c))
	if _, ok := b.Options().Codecs["application/octet-stream"]; !ok {
		t.Error("Codec not added to Codecs map")
	}
}

func TestOptions_Register(t *testing.T) {
	r := register.NewRegister()
	b := broker.NewBroker(broker.Register(r))
	if b.Options().Register != r {
		t.Error("Register option not applied")
	}
}

func TestOptions_TLSConfig(t *testing.T) {
	tlsCfg := &tls.Config{}
	b := broker.NewBroker(broker.TLSConfig(tlsCfg))
	if b.Options().TLSConfig != tlsCfg {
		t.Error("TLSConfig option not applied")
	}
}

func TestOptions_Logger(t *testing.T) {
	l := logger.NewLogger()
	b := broker.NewBroker(broker.Logger(l))
	if b.Options().Logger != l {
		t.Error("Logger option not applied")
	}
}

func TestOptions_Tracer(t *testing.T) {
	tr := tracer.NewTracer()
	b := broker.NewBroker(broker.Tracer(tr))
	if b.Options().Tracer != tr {
		t.Error("Tracer option not applied")
	}
}

func TestOptions_Meter(t *testing.T) {
	m := meter.NewMeter()
	b := broker.NewBroker(broker.Meter(m))
	if b.Options().Meter != m {
		t.Error("Meter option not applied")
	}
}

func TestOptions_Name(t *testing.T) {
	name := "test-broker"
	b := broker.NewBroker(broker.Name(name))
	if b.Name() != name {
		t.Errorf("expected %s, got %s", name, b.Name())
	}
}

func TestOptions_Addrs(t *testing.T) {
	addrs := []string{"localhost:8080", "localhost:8081"}
	b := broker.NewBroker(broker.Addrs(addrs...))
	if len(b.Options().Addrs) != len(addrs) {
		t.Errorf("expected %d addrs, got %d", len(addrs), len(b.Options().Addrs))
	}
}

func TestSubscribeOptions_SubscribeContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), "sub-key", "sub-val")
	opts := broker.NewSubscribeOptions(broker.SubscribeContext(ctx))
	if opts.Context.Value("sub-key") != "sub-val" {
		t.Error("SubscribeContext option not applied")
	}
}

func TestSubscribeOptions_SubscribeAutoAck(t *testing.T) {
	opts := broker.NewSubscribeOptions(broker.SubscribeAutoAck(false))
	if opts.AutoAck {
		t.Error("expected AutoAck to be false")
	}
}

func TestSubscribeOptions_SubscribeBodyOnly(t *testing.T) {
	opts := broker.NewSubscribeOptions(broker.SubscribeBodyOnly(true))
	if !opts.BodyOnly {
		t.Error("expected BodyOnly to be true")
	}
}

func TestSubscribeOptions_SubscribeBatchSize(t *testing.T) {
	opts := broker.NewSubscribeOptions(broker.SubscribeBatchSize(10))
	if opts.BatchSize != 10 {
		t.Errorf("expected 10, got %d", opts.BatchSize)
	}
}

func TestSubscribeOptions_SubscribeBatchWait(t *testing.T) {
	wait := 2 * time.Second
	opts := broker.NewSubscribeOptions(broker.SubscribeBatchWait(wait))
	if opts.BatchWait != wait {
		t.Errorf("expected %v, got %v", wait, opts.BatchWait)
	}
}

func TestSubscribeOptions_SubscribeGroup(t *testing.T) {
	opts := broker.NewSubscribeOptions(broker.SubscribeGroup("test-group"))
	if opts.Group != "test-group" {
		t.Errorf("expected test-group, got %s", opts.Group)
	}
}

func TestMessageOptions_MessageContentType(t *testing.T) {
	ct := "application/json"
	opts := broker.NewMessageOptions(broker.MessageContentType(ct))
	if opts.ContentType != ct {
		t.Errorf("expected %s, got %s", ct, opts.ContentType)
	}
}

func TestMessageOptions_MessageBodyOnly(t *testing.T) {
	opts := broker.NewMessageOptions(broker.MessageBodyOnly(true))
	if !opts.BodyOnly {
		t.Error("expected BodyOnly to be true")
	}
}

func TestNewSubscribeOptions(t *testing.T) {
	opts := broker.NewSubscribeOptions(
		broker.SubscribeGroup("group1"),
		broker.SubscribeAutoAck(false),
		broker.SubscribeBodyOnly(true),
		broker.SubscribeBatchSize(5),
		broker.SubscribeBatchWait(1*time.Second),
	)
	if opts.Group != "group1" {
		t.Errorf("Group: expected group1, got %s", opts.Group)
	}
	if opts.AutoAck {
		t.Error("AutoAck: expected false")
	}
	if !opts.BodyOnly {
		t.Error("BodyOnly: expected true")
	}
	if opts.BatchSize != 5 {
		t.Errorf("BatchSize: expected 5, got %d", opts.BatchSize)
	}
	if opts.BatchWait != 1*time.Second {
		t.Errorf("BatchWait: expected 1s, got %v", opts.BatchWait)
	}
}

func TestNewMessageOptions(t *testing.T) {
	opts := broker.NewMessageOptions(
		broker.MessageContentType("application/json"),
		broker.MessageBodyOnly(true),
	)
	if opts.ContentType != "application/json" {
		t.Errorf("ContentType: expected application/json, got %s", opts.ContentType)
	}
	if !opts.BodyOnly {
		t.Error("BodyOnly: expected true")
	}
}

func TestIsValidHandler(t *testing.T) {
	t.Run("valid func(Message) error", func(t *testing.T) {
		err := broker.IsValidHandler(func(broker.Message) error { return nil })
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
	t.Run("valid func([]Message) error", func(t *testing.T) {
		err := broker.IsValidHandler(func([]broker.Message) error { return nil })
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
	t.Run("invalid handler", func(t *testing.T) {
		err := broker.IsValidHandler("invalid")
		if err != broker.ErrInvalidHandler {
			t.Errorf("expected ErrInvalidHandler, got %v", err)
		}
	})
}

func TestNoopSubscriber_Topic(t *testing.T) {
	b := broker.NewBroker()
	ctx := context.Background()
	sub, err := b.Subscribe(ctx, "test-topic", func(msg broker.Message) error { return nil })
	if err != nil {
		t.Fatalf("Subscribe error: %v", err)
	}
	if sub.Topic() != "test-topic" {
		t.Errorf("expected test-topic, got %s", sub.Topic())
	}
}

func TestNoopMessage_Methods(t *testing.T) {
	c := &testCodec{}
	b := broker.NewBroker(broker.Codec(c.ContentType(), c))
	ctx := context.Background()
	hdr := metadata.Pairs("key", "val")
	body := []byte("test-body")
	msg, err := b.NewMessage(ctx, hdr, body, broker.MessageContentType(c.ContentType()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
