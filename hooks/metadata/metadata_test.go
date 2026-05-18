package metadata

import (
	"context"
	"testing"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/server"
)

// --- fake client used in Call/Stream tests ---

type fakeClient struct {
	client.Client
	lastCallCtx   context.Context
	lastStreamCtx context.Context
}

func (f *fakeClient) Call(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
	f.lastCallCtx = ctx
	return nil
}

func (f *fakeClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	f.lastStreamCtx = ctx
	return nil, nil
}

func (f *fakeClient) NewRequest(service, endpoint string, body any, opts ...client.RequestOption) client.Request {
	return client.NewClient().NewRequest(service, endpoint, body, opts...)
}

// --- helper ---

type noopServerRequest struct{}

func (r *noopServerRequest) Header() metadata.Metadata { return metadata.New(0) }
func (r *noopServerRequest) Body() any                 { return nil }
func (r *noopServerRequest) Method() string            { return "test" }

// =============================================================================
// NewClientWrapper — Call path
// =============================================================================

func TestClientWrapperCall_NilKeys(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper()(fc)
	req := fc.NewRequest("svc", "Method", nil)
	ctx := context.Background()
	if err := wrapped.Call(ctx, req, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.lastCallCtx != ctx {
		t.Error("context should be unchanged when no keys")
	}
}

func TestClientWrapperCall_WithKey_NoIncoming(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper("key1")(fc)
	req := fc.NewRequest("svc", "Method", nil)
	ctx := context.Background()
	if err := wrapped.Call(ctx, req, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No incoming metadata — no-op, context unchanged
	if fc.lastCallCtx != ctx {
		t.Error("context should be unchanged when no incoming metadata")
	}
}

func TestClientWrapperCall_WithKey_Propagates(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper("key1")(fc)
	req := fc.NewRequest("svc", "Method", nil)

	imd := metadata.New(1)
	imd.Set("key1", "val1")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	if err := wrapped.Call(ctx, req, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	omd, ok := metadata.FromOutgoingContext(fc.lastCallCtx)
	if !ok || omd == nil {
		t.Fatal("expected outgoing metadata to be set on the forwarded context")
	}
	if vals := omd.Get("key1"); len(vals) == 0 || vals[0] != "val1" {
		t.Errorf("expected key1=val1, got %v", vals)
	}
}

func TestClientWrapperCall_WithKey_ExistingOutgoing(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper("key1")(fc)
	req := fc.NewRequest("svc", "Method", nil)

	imd := metadata.New(1)
	imd.Set("key1", "val2")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	omd := metadata.New(1)
	omd.Set("other", "yes")
	ctx = metadata.NewOutgoingContext(ctx, omd)

	if err := wrapped.Call(ctx, req, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// NewClientWrapper — Stream path
// =============================================================================

func TestClientWrapperStream_NilKeys(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper()(fc)
	req := fc.NewRequest("svc", "Method", nil)
	ctx := context.Background()
	_, err := wrapped.Stream(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.lastStreamCtx != ctx {
		t.Error("context should be unchanged when no keys")
	}
}

func TestClientWrapperStream_WithKey_NoIncoming(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper("key1")(fc)
	req := fc.NewRequest("svc", "Method", nil)
	ctx := context.Background()
	_, err := wrapped.Stream(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.lastStreamCtx != ctx {
		t.Error("context should be unchanged when no incoming metadata")
	}
}

func TestClientWrapperStream_WithKey_Propagates(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper("key1")(fc)
	req := fc.NewRequest("svc", "Method", nil)

	imd := metadata.New(1)
	imd.Set("key1", "sval1")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	_, err := wrapped.Stream(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	omd, ok := metadata.FromOutgoingContext(fc.lastStreamCtx)
	if !ok || omd == nil {
		t.Fatal("expected outgoing metadata to be set on the forwarded context")
	}
	if vals := omd.Get("key1"); len(vals) == 0 || vals[0] != "sval1" {
		t.Errorf("expected key1=sval1, got %v", vals)
	}
}

func TestClientWrapperStream_WithKey_ExistingOutgoing(t *testing.T) {
	fc := &fakeClient{}
	wrapped := NewClientWrapper("key1")(fc)
	req := fc.NewRequest("svc", "Method", nil)

	imd := metadata.New(1)
	imd.Set("key1", "sval2")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	omd := metadata.New(1)
	omd.Set("existing", "yes")
	ctx = metadata.NewOutgoingContext(ctx, omd)

	_, err := wrapped.Stream(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =============================================================================
// NewClientCallWrapper
// =============================================================================

func TestClientCallWrapper_NilKeys(t *testing.T) {
	called := false
	inner := func(ctx context.Context, addr string, req client.Request, rsp any, opts client.CallOptions) error {
		called = true
		return nil
	}
	fn := NewClientCallWrapper()(inner)
	c := client.NewClient()
	req := c.NewRequest("svc", "Method", nil)
	if err := fn(context.Background(), "addr", req, nil, client.CallOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner function to be called")
	}
}

func TestClientCallWrapper_WithKey_NoIncoming(t *testing.T) {
	called := false
	inner := func(ctx context.Context, addr string, req client.Request, rsp any, opts client.CallOptions) error {
		called = true
		return nil
	}
	fn := NewClientCallWrapper("x-request-id")(inner)
	c := client.NewClient()
	req := c.NewRequest("svc", "Method", nil)
	if err := fn(context.Background(), "addr", req, nil, client.CallOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner function to be called")
	}
}

func TestClientCallWrapper_WithKey_Propagates(t *testing.T) {
	var capturedCtx context.Context
	inner := func(ctx context.Context, addr string, req client.Request, rsp any, opts client.CallOptions) error {
		capturedCtx = ctx
		return nil
	}
	fn := NewClientCallWrapper("x-request-id")(inner)

	imd := metadata.New(1)
	imd.Set("x-request-id", "req-123")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	c := client.NewClient()
	req := c.NewRequest("svc", "Method", nil)
	if err := fn(ctx, "addr", req, nil, client.CallOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedCtx == nil {
		t.Fatal("captured ctx is nil")
	}
	omd, ok := metadata.FromOutgoingContext(capturedCtx)
	if !ok || omd == nil {
		t.Fatal("expected outgoing metadata on forwarded context")
	}
	if vals := omd.Get("x-request-id"); len(vals) == 0 || vals[0] != "req-123" {
		t.Errorf("expected x-request-id=req-123, got %v", vals)
	}
}

func TestClientCallWrapper_WithKey_ExistingOutgoing(t *testing.T) {
	var capturedCtx context.Context
	inner := func(ctx context.Context, addr string, req client.Request, rsp any, opts client.CallOptions) error {
		capturedCtx = ctx
		return nil
	}
	fn := NewClientCallWrapper("x-request-id")(inner)

	imd := metadata.New(1)
	imd.Set("x-request-id", "req-456")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	omd := metadata.New(1)
	omd.Set("existing", "yes")
	ctx = metadata.NewOutgoingContext(ctx, omd)

	c := client.NewClient()
	req := c.NewRequest("svc", "Method", nil)
	if err := fn(ctx, "addr", req, nil, client.CallOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = capturedCtx
}

// =============================================================================
// NewServerHandlerWrapper
// =============================================================================

func TestServerHandlerWrapper_NilKeys(t *testing.T) {
	called := false
	inner := func(ctx context.Context, req server.Request, rsp any) error {
		called = true
		return nil
	}
	fn := NewServerHandlerWrapper()(inner)
	if err := fn(context.Background(), &noopServerRequest{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner handler to be called")
	}
}

func TestServerHandlerWrapper_WithKey_NoIncoming(t *testing.T) {
	called := false
	inner := func(ctx context.Context, req server.Request, rsp any) error {
		called = true
		return nil
	}
	fn := NewServerHandlerWrapper("trace-id")(inner)
	if err := fn(context.Background(), &noopServerRequest{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected inner handler to be called")
	}
}

func TestServerHandlerWrapper_WithKey_Propagates(t *testing.T) {
	var capturedCtx context.Context
	inner := func(ctx context.Context, req server.Request, rsp any) error {
		capturedCtx = ctx
		return nil
	}
	fn := NewServerHandlerWrapper("trace-id")(inner)

	imd := metadata.New(1)
	imd.Set("trace-id", "trace-789")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	if err := fn(ctx, &noopServerRequest{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedCtx == nil {
		t.Fatal("captured ctx is nil")
	}
	omd, ok := metadata.FromOutgoingContext(capturedCtx)
	if !ok || omd == nil {
		t.Fatal("expected outgoing metadata to be set on the forwarded context")
	}
	if vals := omd.Get("trace-id"); len(vals) == 0 || vals[0] != "trace-789" {
		t.Errorf("expected trace-id=trace-789, got %v", vals)
	}
}

func TestServerHandlerWrapper_WithKey_ExistingOutgoing(t *testing.T) {
	var capturedCtx context.Context
	inner := func(ctx context.Context, req server.Request, rsp any) error {
		capturedCtx = ctx
		return nil
	}
	fn := NewServerHandlerWrapper("trace-id")(inner)

	imd := metadata.New(1)
	imd.Set("trace-id", "trace-abc")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	omd := metadata.New(1)
	omd.Set("existing", "yes")
	ctx = metadata.NewOutgoingContext(ctx, omd)

	if err := fn(ctx, &noopServerRequest{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = capturedCtx
}
