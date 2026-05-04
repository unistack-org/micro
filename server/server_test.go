package server

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v4/metadata"
)

func TestNoopServer_Name(t *testing.T) {
	s := NewServer(Name("noop"))
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

func TestServerOptionFunctions(t *testing.T) {
	s := NewServer(Namespace("test-ns"))
	if s.Options().Namespace != "test-ns" {
		t.Errorf("Namespace option not applied, got %q", s.Options().Namespace)
	}

	s = NewServer(ID("test-id"))
	if s.Options().ID != "test-id" {
		t.Errorf("ID option not applied, got %q", s.Options().ID)
	}

	s = NewServer(Version("1.0.0"))
	if s.Options().Version != "1.0.0" {
		t.Errorf("Version option not applied, got %q", s.Options().Version)
	}

	s = NewServer(Address("localhost:8080"))
	if s.Options().Address != "localhost:8080" {
		t.Errorf("Address option not applied, got %q", s.Options().Address)
	}

	s = NewServer(Advertise("localhost:9090"))
	if s.Options().Advertise != "localhost:9090" {
		t.Errorf("Advertise option not applied, got %q", s.Options().Advertise)
	}
}

func TestMustContext(t *testing.T) {
	s := NewServer()
	ctx := NewContext(context.Background(), s)
	retrieved := MustContext(ctx)
	if retrieved != s {
		t.Error("MustContext did not retrieve server")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustContext with empty context")
		}
	}()
	MustContext(context.Background())
}

func TestSetHandlerOption(t *testing.T) {
	type key struct{}
	opt := SetHandlerOption(key{}, "test")
	opts := &HandlerOptions{}
	opt(opts)
	if v, ok := opts.Context.Value(key{}).(string); !ok || v != "test" {
		t.Error("SetHandlerOption not working")
	}
}

func TestRPCHandler(t *testing.T) {
	type Greeter struct{}
	h := newRPCHandler(&Greeter{})
	if h.Name() != "Greeter" {
		t.Errorf("expected handler name 'Greeter', got %q", h.Name())
	}
	if h.Handler() == nil {
		t.Error("expected non-nil handler")
	}
	if h.Options().Context == nil {
		t.Error("expected non-nil context in handler options")
	}
}

func TestResponseMetadata(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(1)
	md.Append("key", "value")
	newCtx := ResponseMetadata(ctx, &md)
	if newCtx == nil {
		t.Error("expected non-nil context from ResponseMetadata")
	}
}

func TestSetResponseMetadata(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(1)
	md.Append("key", "value")
	newCtx := ResponseMetadata(ctx, &md)
	err := SetResponseMetadata(newCtx, md)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopServer_RegisterDeregister(t *testing.T) {
	s := NewServer()
	ns, ok := s.(*noopServer)
	if !ok {
		t.Fatal("failed to cast to *noopServer")
	}
	if err := ns.Start(); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	if err := ns.Register(); err != nil {
		t.Errorf("Register error: %v", err)
	}
	if err := ns.Deregister(); err != nil {
		t.Errorf("Deregister error: %v", err)
	}
	if err := ns.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}
}

func TestNewHandlerOptions(t *testing.T) {
	opts := NewHandlerOptions()
	if opts.Context == nil {
		t.Error("expected non-nil context")
	}
}

func TestSetResponseMetadataError(t *testing.T) {
	ctx := context.Background()
	md := metadata.New(1)
	md.Append("key", "value")
	err := SetResponseMetadata(ctx, md)
	if err == nil {
		t.Error("expected error from SetResponseMetadata with no metadata in context")
	}
}

func TestRemainingOptions(t *testing.T) {
	// Logger
	_ = NewServer(Logger(nil))

	// Meter
	_ = NewServer(Meter(nil))

	// Broker
	_ = NewServer(Broker(nil))

	// Codec
	_ = NewServer(Codec("application/json", nil))

	// Context
	ctx := context.Background()
	s := NewServer(Context(ctx))
	if s.Options().Context != ctx {
		t.Error("Context option not applied")
	}

	// Register
	_ = NewServer(Register(nil))

	// Tracer
	_ = NewServer(Tracer(nil))

	// Metadata
	md := metadata.New(1)
	s = NewServer(Metadata(md))
	if s.Options().Metadata.Len() != md.Len() {
		t.Error("Metadata option not applied")
	}

	// RegisterCheck
	fn := func(context.Context) error { return nil }
	s = NewServer(RegisterCheck(fn))
	if s.Options().RegisterCheck == nil {
		t.Error("RegisterCheck option not applied")
	}

	// RegisterTTL
	s = NewServer(RegisterTTL(10 * time.Second))
	if s.Options().RegisterTTL != 10*time.Second {
		t.Error("RegisterTTL option not applied")
	}

	// RegisterInterval
	s = NewServer(RegisterInterval(5 * time.Second))
	if s.Options().RegisterInterval != 5*time.Second {
		t.Error("RegisterInterval option not applied")
	}

	// TLSConfig
	_ = NewServer(TLSConfig(nil))

	// Wait
	_ = NewServer(Wait(nil))

	// MaxConn
	s = NewServer(MaxConn(10))
	if s.Options().MaxConn != 10 {
		t.Error("MaxConn option not applied")
	}

	// Listener
	_ = NewServer(Listener(nil))

	// GracefulTimeout
	s = NewServer(GracefulTimeout(3 * time.Second))
	if s.Options().GracefulTimeout != 3*time.Second {
		t.Error("GracefulTimeout option not applied")
	}

	// Hooks
	_ = NewServer(Hooks(nil))

	// EndpointMetadata (HandlerOption)
	emd := metadata.New(1)
	emd.Append("k", "v")
	h := NewServer().NewHandler(struct{}{}, EndpointMetadata(emd))
	if len(h.Options().Metadata) == 0 {
		t.Error("EndpointMetadata not applied")
	}
}

func TestNoopServer_InitWithOptions(t *testing.T) {
	s := NewServer()
	if err := s.Init(Name("test")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if s.Options().Name != "test" {
		t.Errorf("Init did not apply options, got %q", s.Options().Name)
	}
}

func TestNoopServer_StartAlreadyStarted(t *testing.T) {
	s := NewServer()
	if err := s.Start(); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	if err := s.Start(); err != nil {
		t.Errorf("unexpected error on second start: %v", err)
	}
	if err := s.Stop(); err != nil {
		t.Fatalf("Stop error: %v", err)
	}
}

func TestNoopServer_StopWithoutStart(t *testing.T) {
	s := NewServer()
	if err := s.Stop(); err != nil {
		t.Errorf("unexpected error stopping non-started server: %v", err)
	}
}
