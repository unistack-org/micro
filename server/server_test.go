package server

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"testing"
	"time"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/options"
	"go.unistack.org/micro/v5/register"
	"go.unistack.org/micro/v5/tracer"
)

type contextKey string

const testContextKey contextKey = "key"

type testCodec struct{}

func (c *testCodec) Marshal(any, ...codec.Option) ([]byte, error) { return nil, nil }
func (c *testCodec) Unmarshal([]byte, any, ...codec.Option) error { return nil }
func (c *testCodec) String() string                              { return "test" }
func (c *testCodec) ContentType() string                         { return "test" }

func TestNoopServer_Name(t *testing.T) {
	s := NewServer()
	if name := s.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopServer_Options(t *testing.T) {
	s := NewServer()
	opts := s.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopServer_Handle(t *testing.T) {
	s := NewServer()
	h := s.NewHandler(func(ctx context.Context, req Request, rsp any) error {
		return nil
	})
	if err := s.Handle(h); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopServer_String(t *testing.T) {
	s := NewServer()
	if str := s.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopServer_Live(t *testing.T) {
	s := NewServer()
	if live := s.Live(); !live {
		t.Error("expected server to be live")
	}
}

func TestNoopServer_Ready(t *testing.T) {
	s := NewServer()
	if ready := s.Ready(); !ready {
		t.Error("expected server to be ready")
	}
}

func TestNoopServer_Health(t *testing.T) {
	s := NewServer()
	if health := s.Health(); !health {
		t.Error("expected server to be healthy")
	}
}

func TestNoopServer_Init(t *testing.T) {
	s := NewServer()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopServer_StartStop(t *testing.T) {
	s := NewServer()
	if err := s.Start(); err != nil {
		t.Errorf("unexpected error on Start: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Errorf("unexpected error on Stop: %v", err)
	}
}

func TestNoopServer_NewHandler(t *testing.T) {
	s := NewServer()
	handler := s.NewHandler(func(ctx context.Context, req Request, rsp any) error {
		return nil
	})
	if handler == nil {
		t.Error("expected non-nil handler")
	}
}

func TestMustContext(t *testing.T) {
	s := NewServer()
	ctx := NewContext(context.Background(), s)
	retrieved := MustContext(ctx)
	if retrieved != s {
		t.Error("expected server from MustContext")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustContext with nil context")
		}
	}()
	MustContext(context.TODO())
}

func TestSetHandlerOption(t *testing.T) {
	type key struct{}
	o := SetHandlerOption(key{}, "test")
	opts := &HandlerOptions{}
	o(opts)
	if v, ok := opts.Context.Value(key{}).(string); !ok || v != "test" {
		t.Error("SetHandlerOption not working")
	}
}

func TestRPCHandlerMethods(t *testing.T) {
	type testHandler struct{}
	handler := newRPCHandler(&testHandler{})
	if handler.Name() == "" {
		t.Error("expected non-empty handler name")
	}
	if handler.Handler() == nil {
		t.Error("expected non-nil handler function")
	}
	opts := handler.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in handler options")
	}
}

func TestOptionFunctions(t *testing.T) {
	s := NewServer(Name("test"), Namespace("test-ns"), Version("1.0"))
	opts := s.Options()
	if opts.Name != "test" {
		t.Errorf("expected name 'test', got %q", opts.Name)
	}
	if opts.Namespace != "test-ns" {
		t.Errorf("expected namespace 'test-ns', got %q", opts.Namespace)
	}
	if opts.Version != "1.0" {
		t.Errorf("expected version '1.0', got %q", opts.Version)
	}
}

func TestResponseMetadata(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(1)
	md.Set("key", "value")
	newCtx := ResponseMetadata(ctx, &md)
	if newCtx == nil {
		t.Error("expected non-nil context from ResponseMetadata")
	}
}

func TestSetResponseMetadata(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(1)
	md.Set("key", "value")
	newCtx := ResponseMetadata(ctx, &metadata.Metadata{})
	err := SetResponseMetadata(newCtx, md)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInitWithOptions(t *testing.T) {
	s := NewServer()
	err := s.Init(Name("test-init"), Namespace("init-ns"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	opts := s.Options()
	if opts.Name != "test-init" {
		t.Errorf("expected name 'test-init', got %q", opts.Name)
	}
	if opts.Namespace != "init-ns" {
		t.Errorf("expected namespace 'init-ns', got %q", opts.Namespace)
	}
}

func TestSetResponseMetadataError(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(1)
	md.Set("key", "value")
	err := SetResponseMetadata(ctx, md)
	if err == nil {
		t.Error("expected error for missing metadata in context")
	}
}

func TestNoopServer_RegisterDeregister(t *testing.T) {
	s := NewServer()
	ns := s.(*noopServer)
	if err := ns.Register(); err != nil {
		t.Errorf("unexpected error on Register: %v", err)
	}
	if err := ns.Deregister(); err != nil {
		t.Errorf("unexpected error on Deregister: %v", err)
	}
}

func TestMoreOptionFunctions(t *testing.T) {
	s := NewServer(
		Logger(logger.DefaultLogger),
		Meter(meter.DefaultMeter),
		ID("test-id"),
		Address("127.0.0.1:0"),
	)
	opts := s.Options()
	if opts.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got %q", opts.ID)
	}
	if opts.Address != "127.0.0.1:0" {
		t.Errorf("expected address '127.0.0.1:0', got %q", opts.Address)
	}
}

func TestAllOptions(t *testing.T) {
	s := NewServer(
		Advertise("127.0.0.1:8080"),
		Broker(broker.DefaultBroker),
		Context(context.WithValue(context.Background(), testContextKey, "value")),
		Register(register.DefaultRegister),
		Tracer(tracer.DefaultTracer),
		Metadata(metadata.New(1)),
		RegisterCheck(func(ctx context.Context) error { return nil }),
		RegisterTTL(30*time.Second),
		RegisterInterval(15*time.Second),
		TLSConfig(&tls.Config{}),
		MaxConn(100),
		GracefulTimeout(5*time.Second),
		Hooks(options.Hooks{}),
	)
	opts := s.Options()
	if opts.Advertise != "127.0.0.1:8080" {
		t.Errorf("expected advertise '127.0.0.1:8080', got %q", opts.Advertise)
	}
	if opts.MaxConn != 100 {
		t.Errorf("expected max conn 100, got %d", opts.MaxConn)
	}
	if opts.GracefulTimeout != 5*time.Second {
		t.Errorf("expected graceful timeout 5s, got %v", opts.GracefulTimeout)
	}
}

func TestNewRegisterService(t *testing.T) {
	s := NewServer()
	svc, err := NewRegisterService(s)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestCodecOption(t *testing.T) {
	c := &testCodec{}
	s := NewServer(Codec("test", c))
	opts := s.Options()
	if opts.Codecs == nil || opts.Codecs["test"] == nil {
		t.Error("expected codec to be set")
	}
}

func TestStartStopNoop(t *testing.T) {
	s := NewServer()
	if err := s.Start(); err != nil {
		t.Errorf("unexpected error on Start: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Errorf("unexpected error on Stop: %v", err)
	}
}

func TestWaitOption(t *testing.T) {
	var wg sync.WaitGroup
	s := NewServer(Wait(&wg))
	opts := s.Options()
	if opts.Wait == nil {
		t.Error("expected wait group to be set")
	}
}

type testListener struct{}

func (t *testListener) Accept() (net.Conn, error) { return nil, nil }
func (t *testListener) Close() error               { return nil }
func (t *testListener) Addr() net.Addr            { return nil }

func TestListenerOption(t *testing.T) {
	l := &testListener{}
	s := NewServer(Listener(l))
	opts := s.Options()
	if opts.Listener != l {
		t.Error("expected listener to be set")
	}
}

func TestEndpointMetadataOption(t *testing.T) {
	md := metadata.New(1)
	md.Set("key", "value")
	o := EndpointMetadata(md)
	opts := &HandlerOptions{}
	o(opts)
	if opts.Metadata == nil {
		t.Error("expected metadata to be set")
	}
}

func TestNewHandlerOptionsFunc(t *testing.T) {
	opts := NewHandlerOptions()
	if opts.Context == nil {
		t.Error("expected non-nil context")
	}
}

func TestNoopServer_NameWithOption(t *testing.T) {
	s := NewServer(Name("custom-server"))
	if name := s.Name(); name != "custom-server" {
		t.Errorf("expected 'custom-server', got %q", name)
	}
}

func TestNoopServer_InitMultiple(t *testing.T) {
	s := NewServer()
	_ = s.Init()
	_ = s.Init()
	opts := s.Options()
	if opts.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestNoopServer_HandleMultiple(t *testing.T) {
	s := NewServer()
	h1 := s.NewHandler(func(ctx context.Context, req Request, rsp any) error { return nil })
	h2 := s.NewHandler(func(ctx context.Context, req Request, rsp any) error { return nil })
	_ = s.Handle(h1)
	_ = s.Handle(h2)
}

func TestRegistryOptions(t *testing.T) {
	s := NewServer(RegisterTTL(time.Minute), RegisterInterval(time.Second))
	opts := s.Options()
	if opts.RegisterTTL != time.Minute {
		t.Errorf("expected 1m, got %v", opts.RegisterTTL)
	}
}

func TestServerWithMetadata(t *testing.T) {
	md := metadata.New(1)
	md.Set("key", "value")
	s := NewServer(Metadata(md))
	opts := s.Options()
	if opts.Metadata == nil {
		t.Error("expected metadata to be set")
	}
}

func TestServerWithBrokerOptions(t *testing.T) {
	s := NewServer(Broker(broker.DefaultBroker))
	opts := s.Options()
	if opts.Broker == nil {
		t.Error("expected broker to be set")
	}
}
