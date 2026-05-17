package register

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/tracer"
)

func TestNoopName(t *testing.T) {
	r := NewRegister(Name("test"))
	if r.Name() != "test" {
		t.Fatalf("expected name 'test', got %q", r.Name())
	}
}

func TestNoopDefaultName(t *testing.T) {
	r := NewRegister()
	// noop returns "" by default since Name option is not set
	_ = r.Name()
}

func TestNoopInit(t *testing.T) {
	r := NewRegister()
	if err := r.Init(Name("updated")); err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
	if r.Name() != "updated" {
		t.Fatalf("expected name 'updated' after Init, got %q", r.Name())
	}
}

func TestNoopOptions(t *testing.T) {
	r := NewRegister(Name("opts"))
	opts := r.Options()
	if opts.Name != "opts" {
		t.Fatalf("expected opts.Name 'opts', got %q", opts.Name)
	}
}

func TestNoopConnect(t *testing.T) {
	r := NewRegister()
	if err := r.Connect(context.Background()); err != nil {
		t.Fatalf("Connect returned error: %v", err)
	}
}

func TestNoopDisconnect(t *testing.T) {
	r := NewRegister()
	if err := r.Disconnect(context.Background()); err != nil {
		t.Fatalf("Disconnect returned error: %v", err)
	}
}

func TestNoopRegister(t *testing.T) {
	r := NewRegister()
	svc := &Service{Name: "test-service", Version: "v1.0.0"}
	if err := r.Register(context.Background(), svc); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
}

func TestNoopRegisterWithOptions(t *testing.T) {
	r := NewRegister()
	svc := &Service{Name: "test-service", Version: "v1.0.0"}
	if err := r.Register(context.Background(), svc,
		RegisterTTL(time.Second),
		RegisterAttempts(3),
		RegisterNamespace("custom"),
		RegisterContext(context.Background()),
	); err != nil {
		t.Fatalf("Register with options returned error: %v", err)
	}
}

func TestNoopDeregister(t *testing.T) {
	r := NewRegister()
	svc := &Service{Name: "test-service", Version: "v1.0.0"}
	if err := r.Deregister(context.Background(), svc); err != nil {
		t.Fatalf("Deregister returned error: %v", err)
	}
}

func TestNoopDeregisterWithOptions(t *testing.T) {
	r := NewRegister()
	svc := &Service{Name: "test-service", Version: "v1.0.0"}
	if err := r.Deregister(context.Background(), svc,
		DeregisterAttempts(2),
		DeregisterNamespace("ns"),
		DeregisterContext(context.Background()),
	); err != nil {
		t.Fatalf("Deregister with options returned error: %v", err)
	}
}

func TestNoopLookupService(t *testing.T) {
	r := NewRegister()
	svcs, err := r.LookupService(context.Background(), "test-service")
	if err != nil {
		t.Fatalf("LookupService returned error: %v", err)
	}
	if svcs != nil {
		t.Fatalf("expected nil services, got %v", svcs)
	}
}

func TestNoopLookupServiceWithOptions(t *testing.T) {
	r := NewRegister()
	svcs, err := r.LookupService(context.Background(), "test-service",
		LookupNamespace("ns"),
		LookupContext(context.Background()),
	)
	if err != nil {
		t.Fatalf("LookupService with options returned error: %v", err)
	}
	if svcs != nil {
		t.Fatalf("expected nil services, got %v", svcs)
	}
}

func TestNoopListServices(t *testing.T) {
	r := NewRegister()
	svcs, err := r.ListServices(context.Background())
	if err != nil {
		t.Fatalf("ListServices returned error: %v", err)
	}
	if svcs != nil {
		t.Fatalf("expected nil services, got %v", svcs)
	}
}

func TestNoopListServicesWithOptions(t *testing.T) {
	r := NewRegister()
	svcs, err := r.ListServices(context.Background(),
		ListNamespace("ns"),
		ListContext(context.Background()),
	)
	if err != nil {
		t.Fatalf("ListServices with options returned error: %v", err)
	}
	if svcs != nil {
		t.Fatalf("expected nil services, got %v", svcs)
	}
}

func TestNoopWatch(t *testing.T) {
	r := NewRegister()
	w, err := r.Watch(context.Background())
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	if w == nil {
		t.Fatal("Watch returned nil watcher")
	}
}

func TestNoopWatchWithOptions(t *testing.T) {
	r := NewRegister()
	w, err := r.Watch(context.Background(),
		WatchService("test"),
		WatchNamespace("ns"),
		WatchContext(context.Background()),
	)
	if err != nil {
		t.Fatalf("Watch with options returned error: %v", err)
	}
	if w == nil {
		t.Fatal("Watch with options returned nil watcher")
	}
	w.Stop()
}

func TestNoopString(t *testing.T) {
	r := NewRegister()
	if r.String() != "noop" {
		t.Fatalf("expected 'noop', got %q", r.String())
	}
}

func TestNoopOptionsWithAddrs(t *testing.T) {
	r := NewRegister(Addrs("localhost:8500"))
	opts := r.Options()
	if len(opts.Addrs) == 0 || opts.Addrs[0] != "localhost:8500" {
		t.Fatalf("expected addrs to contain 'localhost:8500', got %v", opts.Addrs)
	}
}

func TestNoopOptionsWithTimeout(t *testing.T) {
	r := NewRegister(Timeout(5 * time.Second))
	opts := r.Options()
	if opts.Timeout != 5*time.Second {
		t.Fatalf("expected timeout 5s, got %v", opts.Timeout)
	}
}

func TestNoopMustContext(t *testing.T) {
	ctx := NewContext(context.Background(), NewRegister())
	r := MustContext(ctx)
	if r == nil {
		t.Fatal("MustContext returned nil")
	}
}

func TestNoopMustContextPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Fatal("expected panic from MustContext with empty context")
		}
	}()
	MustContext(context.Background())
}

func TestNewOptions(t *testing.T) {
	opts := NewOptions(
		Name("reg"),
		Addrs("127.0.0.1:8500"),
		Timeout(time.Second),
		Context(context.Background()),
	)
	if opts.Name != "reg" {
		t.Fatalf("expected name 'reg', got %q", opts.Name)
	}
}

func TestNewRegisterOptions(t *testing.T) {
	opts := NewRegisterOptions(
		RegisterTTL(5*time.Second),
		RegisterAttempts(2),
		RegisterNamespace("ns"),
		RegisterContext(context.Background()),
	)
	if opts.Namespace != "ns" {
		t.Fatalf("expected namespace 'ns', got %q", opts.Namespace)
	}
	if opts.Attempts != 2 {
		t.Fatalf("expected attempts 2, got %d", opts.Attempts)
	}
}

func TestNewWatchOptions(t *testing.T) {
	opts := NewWatchOptions(
		WatchService("svc"),
		WatchNamespace("ns"),
		WatchContext(context.Background()),
	)
	if opts.Service != "svc" {
		t.Fatalf("expected service 'svc', got %q", opts.Service)
	}
	if opts.Namespace != "ns" {
		t.Fatalf("expected namespace 'ns', got %q", opts.Namespace)
	}
}

func TestNewDeregisterOptions(t *testing.T) {
	opts := NewDeregisterOptions(
		DeregisterAttempts(3),
		DeregisterNamespace("ns"),
		DeregisterContext(context.Background()),
	)
	if opts.Namespace != "ns" {
		t.Fatalf("expected namespace 'ns', got %q", opts.Namespace)
	}
	if opts.Attempts != 3 {
		t.Fatalf("expected attempts 3, got %d", opts.Attempts)
	}
}

func TestNewLookupOptions(t *testing.T) {
	opts := NewLookupOptions(
		LookupNamespace("ns"),
		LookupContext(context.Background()),
	)
	if opts.Namespace != "ns" {
		t.Fatalf("expected namespace 'ns', got %q", opts.Namespace)
	}
}

func TestNewListOptions(t *testing.T) {
	opts := NewListOptions(
		ListNamespace("ns"),
		ListContext(context.Background()),
	)
	if opts.Namespace != "ns" {
		t.Fatalf("expected namespace 'ns', got %q", opts.Namespace)
	}
}

func TestNoopOptionsWithLogger(t *testing.T) {
	l := logger.DefaultLogger
	r := NewRegister(Logger(l))
	opts := r.Options()
	if opts.Logger == nil {
		t.Fatal("expected logger to be set")
	}
}

func TestNoopOptionsWithMeter(t *testing.T) {
	m := meter.DefaultMeter
	r := NewRegister(Meter(m))
	opts := r.Options()
	if opts.Meter == nil {
		t.Fatal("expected meter to be set")
	}
}

func TestNoopOptionsWithTracer(t *testing.T) {
	tr := tracer.DefaultTracer
	r := NewRegister(Tracer(tr))
	opts := r.Options()
	if opts.Tracer == nil {
		t.Fatal("expected tracer to be set")
	}
}

func TestNoopOptionsWithContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{ key string }{"k"}, "v")
	r := NewRegister(Context(ctx))
	opts := r.Options()
	if opts.Context == nil {
		t.Fatal("expected context to be set")
	}
}

func TestNoopOptionsWithTLSConfig(t *testing.T) {
	tlsCfg := &tls.Config{InsecureSkipVerify: true} // nolint: gosec
	r := NewRegister(TLSConfig(tlsCfg))
	opts := r.Options()
	if opts.TLSConfig == nil {
		t.Fatal("expected TLSConfig to be set")
	}
}

func TestNoopOptionsWithCodec(t *testing.T) {
	c := codec.NewCodec()
	r := NewRegister(Codec(c))
	opts := r.Options()
	if opts.Codec == nil {
		t.Fatal("expected codec to be set")
	}
}
