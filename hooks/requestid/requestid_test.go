package requestid

import (
	"context"
	"slices"
	"testing"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/server"
)

func TestDefaultMetadataFunc(t *testing.T) {
	ctx := context.TODO()

	nctx, err := DefaultMetadataFunc(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	imd, ok := metadata.FromIncomingContext(nctx)
	if !ok {
		t.Fatalf("md missing in incoming context")
	}
	omd, ok := metadata.FromOutgoingContext(nctx)
	if !ok {
		t.Fatalf("md missing in outgoing context")
	}

	iv := imd.Get(DefaultMetadataKey)
	ov := omd.Get(DefaultMetadataKey)

	if !slices.Equal(iv, ov) {
		t.Fatalf("missing metadata key value %v != %v", iv, ov)
	}
}

// mockServerRequest implements server.Request for testing.
type mockServerRequest struct{}

func (r *mockServerRequest) Header() metadata.Metadata { return metadata.New(0) }
func (r *mockServerRequest) Body() any                 { return nil }
func (r *mockServerRequest) Method() string            { return "test" }

// mockClientRequest implements client.Request for testing.
type mockClientRequest struct{}

func (r *mockClientRequest) Service() string      { return "test" }
func (r *mockClientRequest) Method() string       { return "test" }
func (r *mockClientRequest) Endpoint() string     { return "test" }
func (r *mockClientRequest) ContentType() string  { return "application/json" }
func (r *mockClientRequest) Body() any            { return nil }
func (r *mockClientRequest) Codec() codec.Codec   { return nil }
func (r *mockClientRequest) Stream() bool         { return false }

func TestNewHook(t *testing.T) {
	h := NewHook()
	if h == nil {
		t.Fatal("expected non-nil hook")
	}
}

func TestServerHandler(t *testing.T) {
	h := NewHook()
	called := false
	next := func(ctx context.Context, req server.Request, rsp any) error {
		called = true
		return nil
	}
	wrapped := h.ServerHandler(next)
	err := wrapped(context.Background(), &mockServerRequest{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("next handler was not called")
	}
}

func TestServerHandlerPropagatesRequestID(t *testing.T) {
	h := NewHook()
	var capturedCtx context.Context
	next := func(ctx context.Context, req server.Request, rsp any) error {
		capturedCtx = ctx
		return nil
	}
	wrapped := h.ServerHandler(next)
	err := wrapped(context.Background(), &mockServerRequest{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	xid, ok := capturedCtx.Value(XRequestIDKey{}).(string)
	if !ok || xid == "" {
		t.Fatal("request-id not set in context by ServerHandler")
	}
}

func TestClientCall(t *testing.T) {
	h := NewHook()
	called := false
	next := func(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
		called = true
		return nil
	}
	wrapped := h.ClientCall(next)
	err := wrapped(context.Background(), &mockClientRequest{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("next was not called")
	}
}

func TestClientCallWithExistingRequestID(t *testing.T) {
	h := NewHook()
	const existingID = "my-request-id"
	ctx := context.WithValue(context.Background(), XRequestIDKey{}, existingID)

	var capturedCtx context.Context
	next := func(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
		capturedCtx = ctx
		return nil
	}
	wrapped := h.ClientCall(next)
	err := wrapped(ctx, &mockClientRequest{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	xid, _ := capturedCtx.Value(XRequestIDKey{}).(string)
	if xid != existingID {
		t.Fatalf("expected request-id %q, got %q", existingID, xid)
	}
}

func TestClientStream(t *testing.T) {
	h := NewHook()
	called := false
	next := func(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
		called = true
		return nil, nil
	}
	wrapped := h.ClientStream(next)
	_, err := wrapped(context.Background(), &mockClientRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("next was not called")
	}
}

func TestDefaultMetadataFuncWithExistingIncoming(t *testing.T) {
	imd := metadata.New(1)
	imd.Set(DefaultMetadataKey, "incoming-id")
	ctx := metadata.NewIncomingContext(context.Background(), imd)

	nctx, err := DefaultMetadataFunc(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nimd, ok := metadata.FromIncomingContext(nctx)
	if !ok {
		t.Fatal("incoming metadata missing")
	}
	vals := nimd.Get(DefaultMetadataKey)
	if len(vals) == 0 || vals[0] != "incoming-id" {
		t.Fatalf("expected incoming-id, got %v", vals)
	}
}

func TestDefaultMetadataFuncWithExistingOutgoing(t *testing.T) {
	omd := metadata.New(1)
	omd.Set(DefaultMetadataKey, "outgoing-id")
	ctx := metadata.NewOutgoingContext(context.Background(), omd)

	nctx, err := DefaultMetadataFunc(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nomd, ok := metadata.FromOutgoingContext(nctx)
	if !ok {
		t.Fatal("outgoing metadata missing")
	}
	vals := nomd.Get(DefaultMetadataKey)
	if len(vals) == 0 || vals[0] != "outgoing-id" {
		t.Fatalf("expected outgoing-id, got %v", vals)
	}
}
