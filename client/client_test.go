package client

import (
	"testing"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/router"
	"go.unistack.org/micro/v4/selector/random"
	"go.unistack.org/micro/v4/tracer"
)

func TestNoopClient_Name(t *testing.T) {
	c := NewClient()
	if name := c.Name(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopClient_String(t *testing.T) {
	c := NewClient()
	if str := c.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopClient_Init(t *testing.T) {
	c := NewClient()
	if err := c.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopClient_Options(t *testing.T) {
	c := NewClient()
	opts := c.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context")
	}
	if opts.Codecs == nil {
		t.Error("expected non-nil codecs map")
	}
	if opts.Selector == nil {
		t.Error("expected non-nil selector")
	}
}

func TestNoopClient_NewRequest(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "ep", "body")
	if req == nil {
		t.Fatal("expected non-nil request")
	}
}

func TestNoopClient_NewRequestWithInterface(t *testing.T) {
	c := NewClient()
	type MyRequest struct {
		Msg string
	}
	req := c.NewRequest("svc", "ep", MyRequest{Msg: "hello"})
	if req == nil {
		t.Fatal("expected non-nil request")
	}
}

func TestNoopClient_NewRequestWithNil(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "ep", nil)
	if req == nil {
		t.Fatal("expected non-nil request")
	}
}

func TestOptionsName(t *testing.T) {
	c := NewClient(Name("test-client"))
	opts := c.Options()
	if opts.Name != "test-client" {
		t.Errorf("expected 'test-client', got %q", opts.Name)
	}
}

func TestOptionsContentType(t *testing.T) {
	c := NewClient(ContentType("application/json"))
	opts := c.Options()
	if opts.ContentType != "application/json" {
		t.Errorf("expected application/json, got %q", opts.ContentType)
	}
}

func TestOptionsProxy(t *testing.T) {
	c := NewClient(Proxy("http://proxy:8080"))
	opts := c.Options()
	if opts.Proxy != "http://proxy:8080" {
		t.Errorf("expected proxy URL, got %q", opts.Proxy)
	}
}

func TestOptionsLogger(t *testing.T) {
	c := NewClient(Logger(logger.DefaultLogger))
	opts := c.Options()
	if opts.Logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestOptionsTracer(t *testing.T) {
	c := NewClient(Tracer(tracer.DefaultTracer))
	opts := c.Options()
	if opts.Tracer == nil {
		t.Error("expected non-nil tracer")
	}
}

func TestOptionsBroker(t *testing.T) {
	c := NewClient(Broker(broker.DefaultBroker))
	opts := c.Options()
	if opts.Broker == nil {
		t.Error("expected non-nil broker")
	}
}

func TestOptionsMeter(t *testing.T) {
	c := NewClient(Meter(meter.DefaultMeter))
	opts := c.Options()
	if opts.Meter == nil {
		t.Error("expected non-nil meter")
	}
}

func TestOptionsSelector(t *testing.T) {
	c := NewClient(Selector(random.NewSelector()))
	opts := c.Options()
	if opts.Selector == nil {
		t.Error("expected non-nil selector")
	}
}

func TestOptionsRouter(t *testing.T) {
	c := NewClient(Router(router.DefaultRouter))
	opts := c.Options()
	if opts.Router == nil {
		t.Error("expected non-nil router")
	}
}

func TestOptionsTLSConfig(t *testing.T) {
	c := NewClient(TLSConfig(nil))
	opts := c.Options()
	if opts.TLSConfig != nil {
		t.Error("expected nil TLSConfig")
	}
}