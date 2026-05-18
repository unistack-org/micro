package client

import (
	"context"
	"net"
	"testing"
	"time"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/logger"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/meter"
	"go.unistack.org/micro/v5/options"
	"go.unistack.org/micro/v5/register"
	"go.unistack.org/micro/v5/router"
	"go.unistack.org/micro/v5/selector/random"
	"go.unistack.org/micro/v5/tracer"
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

// --- Options / CallOptions setters ---

func TestOptionsPoolSize(t *testing.T) {
	c := NewClient(PoolSize(10))
	if c.Options().PoolSize != 10 {
		t.Error("expected PoolSize 10")
	}
}

func TestOptionsPoolTTL(t *testing.T) {
	c := NewClient(PoolTTL(5 * time.Second))
	if c.Options().PoolTTL != 5*time.Second {
		t.Error("expected PoolTTL 5s")
	}
}

func TestOptionsRegister(t *testing.T) {
	// Register sets the router's register; just ensure it doesn't panic
	c := NewClient(Register(register.DefaultRegister))
	if c == nil {
		t.Error("expected non-nil client")
	}
}

func TestOptionsBackoff(t *testing.T) {
	fn := func(ctx context.Context, req Request, attempts int) (time.Duration, error) {
		return 0, nil
	}
	c := NewClient(Backoff(fn))
	if c.Options().CallOptions.Backoff == nil {
		t.Error("expected non-nil Backoff")
	}
}

func TestOptionsLookup(t *testing.T) {
	fn := func(ctx context.Context, req Request, opts CallOptions) ([]string, error) {
		return []string{"addr:8080"}, nil
	}
	c := NewClient(Lookup(fn))
	if c.Options().Lookup == nil {
		t.Error("expected non-nil Lookup")
	}
}

func TestOptionsRetries(t *testing.T) {
	c := NewClient(Retries(3))
	if c.Options().CallOptions.Retries != 3 {
		t.Error("expected Retries 3")
	}
}

func TestOptionsRetry(t *testing.T) {
	fn := func(ctx context.Context, req Request, retries int, err error) (bool, error) {
		return false, nil
	}
	c := NewClient(Retry(fn))
	if c.Options().CallOptions.Retry == nil {
		t.Error("expected non-nil Retry")
	}
}

func TestOptionsRequestTimeout(t *testing.T) {
	c := NewClient(RequestTimeout(2 * time.Second))
	if c.Options().CallOptions.RequestTimeout != 2*time.Second {
		t.Error("expected RequestTimeout 2s")
	}
}

func TestOptionsStreamTimeout(t *testing.T) {
	c := NewClient(StreamTimeout(3 * time.Second))
	if c.Options().CallOptions.StreamTimeout != 3*time.Second {
		t.Error("expected StreamTimeout 3s")
	}
}

func TestOptionsDialTimeout(t *testing.T) {
	c := NewClient(DialTimeout(1 * time.Second))
	if c.Options().CallOptions.DialTimeout != 1*time.Second {
		t.Error("expected DialTimeout 1s")
	}
}

func TestOptionsContextDialer(t *testing.T) {
	fn := func(ctx context.Context, addr string) (net.Conn, error) { return nil, nil }
	c := NewClient(ContextDialer(fn))
	if c.Options().ContextDialer == nil {
		t.Error("expected non-nil ContextDialer")
	}
}

func TestOptionsContext(t *testing.T) {
	type key struct{}
	ctx := context.WithValue(context.Background(), key{}, "val")
	c := NewClient(Context(ctx))
	if c.Options().Context == nil {
		t.Error("expected non-nil Context")
	}
}

func TestOptionsHooks(t *testing.T) {
	var called bool
	hook := HookCall(func(next FuncCall) FuncCall {
		return func(ctx context.Context, req Request, rsp any, opts ...CallOption) error {
			called = true
			return next(ctx, req, rsp, opts...)
		}
	})
	// Hooks are wired during Init; pass the hook and then call Init to activate it.
	c := NewClient(Hooks(hook))
	if err := c.Init(); err != nil {
		t.Fatalf("unexpected Init error: %v", err)
	}
	req := c.NewRequest("svc", "ep", nil)
	_ = c.Call(context.Background(), req, nil, WithAddress("127.0.0.1:8080"))
	if !called {
		t.Error("expected hook to be called")
	}
}

// --- CallOption setters via NewCallOptions ---

func TestNewCallOptions(t *testing.T) {
	md := metadata.New(2)
	md.Set("k", "v")
	var respMd metadata.Metadata
	copts := NewCallOptions(
		WithContentType("application/json"),
		WithAddress("addr1:8080", "addr2:8080"),
		WithRetries(2),
		WithRequestTimeout(1*time.Second),
		WithStreamTimeout(2*time.Second),
		WithDialTimeout(3*time.Second),
		WithAuthToken("token123"),
		WithRequestMetadata(md),
		WithResponseMetadata(&respMd),
		WithRouter(router.DefaultRouter),
		WithSelector(random.NewSelector()),
	)
	if copts.ContentType != "application/json" {
		t.Errorf("expected application/json, got %q", copts.ContentType)
	}
	if len(copts.Address) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(copts.Address))
	}
	if copts.Retries != 2 {
		t.Error("expected Retries 2")
	}
	if copts.RequestTimeout != 1*time.Second {
		t.Error("expected RequestTimeout 1s")
	}
	if copts.StreamTimeout != 2*time.Second {
		t.Error("expected StreamTimeout 2s")
	}
	if copts.DialTimeout != 3*time.Second {
		t.Error("expected DialTimeout 3s")
	}
	if copts.AuthToken != "token123" {
		t.Errorf("expected token123, got %q", copts.AuthToken)
	}
	if copts.ResponseMetadata == nil {
		t.Error("expected non-nil ResponseMetadata")
	}
	if copts.Router == nil {
		t.Error("expected non-nil Router")
	}
	if copts.Selector == nil {
		t.Error("expected non-nil Selector")
	}
}

func TestWithContextDialerCallOption(t *testing.T) {
	fn := func(ctx context.Context, addr string) (net.Conn, error) { return nil, nil }
	copts := NewCallOptions(WithContextDialer(fn))
	if copts.ContextDialer == nil {
		t.Error("expected non-nil ContextDialer")
	}
}

func TestWithBackoffCallOption(t *testing.T) {
	fn := func(ctx context.Context, req Request, attempts int) (time.Duration, error) { return 0, nil }
	copts := NewCallOptions(WithBackoff(fn))
	if copts.Backoff == nil {
		t.Error("expected non-nil Backoff")
	}
}

func TestWithRetryCallOption(t *testing.T) {
	fn := func(ctx context.Context, req Request, retries int, err error) (bool, error) { return false, nil }
	copts := NewCallOptions(WithRetry(fn))
	if copts.Retry == nil {
		t.Error("expected non-nil Retry")
	}
}

func TestWithSelectOptionsCallOption(t *testing.T) {
	copts := NewCallOptions(WithSelectOptions())
	_ = copts // just ensure no panic; SelectOptions may be empty but set
}

// --- NewRequestOptions ---

func TestNewRequestOptions(t *testing.T) {
	ropts := NewRequestOptions(
		StreamingRequest(true),
		RequestContentType("application/proto"),
	)
	if !ropts.Stream {
		t.Error("expected Stream true")
	}
	if ropts.ContentType != "application/proto" {
		t.Errorf("expected application/proto, got %q", ropts.ContentType)
	}
}

// --- noopRequest methods ---

func TestNoopRequestFields(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("my-svc", "MyMethod", "body")
	if req.Service() != "my-svc" {
		t.Errorf("expected my-svc, got %q", req.Service())
	}
	if req.Endpoint() != "MyMethod" {
		t.Errorf("expected MyMethod, got %q", req.Endpoint())
	}
	if req.Body() == nil {
		t.Error("expected non-nil body")
	}
	if req.Stream() {
		t.Error("expected Stream false")
	}
}

func TestNoopRequestMethod(t *testing.T) {
	nr := &noopRequest{method: "TestMethod"}
	if nr.Method() != "TestMethod" {
		t.Errorf("expected TestMethod, got %q", nr.Method())
	}
}

func TestNoopRequestCodec(t *testing.T) {
	nr := &noopRequest{}
	if nr.Codec() != nil {
		t.Error("expected nil codec")
	}
}

// --- noopResponse methods ---

func TestNoopResponse(t *testing.T) {
	nr := &noopResponse{header: metadata.New(1)}
	if nr.Codec() != nil {
		t.Error("expected nil codec")
	}
	hdr := nr.Header()
	if hdr == nil {
		t.Error("expected non-nil header")
	}
	data, err := nr.Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil data")
	}
}

// --- noopStream methods ---

func TestNoopStream(t *testing.T) {
	ctx := context.Background()
	ns := &noopStream{ctx: ctx}

	if ns.Context() != ctx {
		t.Error("expected same context")
	}
	if ns.Request() == nil {
		t.Error("expected non-nil request")
	}
	if ns.Response() == nil {
		t.Error("expected non-nil response")
	}
	if err := ns.Send("msg"); err != nil {
		t.Errorf("unexpected Send error: %v", err)
	}
	if err := ns.Recv("msg"); err != nil {
		t.Errorf("unexpected Recv error: %v", err)
	}
	if err := ns.SendMsg("msg"); err != nil {
		t.Errorf("unexpected SendMsg error: %v", err)
	}
	if err := ns.RecvMsg("msg"); err != nil {
		t.Errorf("unexpected RecvMsg error: %v", err)
	}
	if ns.Error() != nil {
		t.Error("expected nil error")
	}
	if err := ns.CloseSend(); err != nil {
		t.Errorf("unexpected CloseSend error: %v", err)
	}
	if err := ns.Close(); err != nil {
		t.Errorf("unexpected Close error: %v", err)
	}
}

// --- noopClient.Call and noopClient.Stream ---

func TestNoopClient_Call(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "Endpoint", nil)
	if err := c.Call(context.Background(), req, nil, WithAddress("127.0.0.1:8080")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopClient_CallWithAddress(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "Endpoint", nil)
	if err := c.Call(context.Background(), req, nil, WithAddress("127.0.0.1:8080")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopClient_Stream(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "StreamEndpoint", nil)
	stream, err := c.Stream(context.Background(), req, WithAddress("127.0.0.1:8080"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Fatal("expected non-nil stream")
	}
	_ = stream.Close()
}

func TestNoopClient_StreamWithTimeout(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "StreamEndpoint", nil)
	stream, err := c.Stream(context.Background(), req,
		WithAddress("127.0.0.1:8080"),
		WithStreamTimeout(5*time.Second),
	)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream != nil {
		_ = stream.Close()
	}
}

// --- MustContext panic ---

func TestMustContext_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustContext with empty context")
		}
	}()
	MustContext(context.Background())
}

func TestMustContext_OK(t *testing.T) {
	c := NewClient()
	ctx := NewContext(context.Background(), c)
	got := MustContext(ctx)
	if got != c {
		t.Error("expected same client from MustContext")
	}
}

// --- LookupRoute ---

func TestLookupRoute_WithAddress(t *testing.T) {
	req := &testRequest{service: "svc"}
	opts := CallOptions{Address: []string{"addr:8080"}}
	addrs, err := LookupRoute(context.Background(), req, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addrs) != 1 || addrs[0] != "addr:8080" {
		t.Errorf("expected [addr:8080], got %v", addrs)
	}
}

func TestLookupRoute_NilRouter(t *testing.T) {
	req := &testRequest{service: "svc"}
	opts := CallOptions{}
	_, err := LookupRoute(context.Background(), req, opts)
	if err == nil {
		t.Error("expected error when Router is nil")
	}
}

func TestLookupRoute_WithRouter(t *testing.T) {
	req := &testRequest{service: "svc"}
	opts := CallOptions{Router: router.NewRouter()}
	// dns router does a DNS lookup; it may return routes or an error
	// we just verify LookupRoute doesn't panic
	_, _ = LookupRoute(context.Background(), req, opts)
}

// --- routesByMetric sort interface ---

func TestRoutesByMetric(t *testing.T) {
	routes := routesByMetric{
		{Metric: 3, Address: "c"},
		{Metric: 1, Address: "a"},
		{Metric: 2, Address: "b"},
	}
	if routes.Len() != 3 {
		t.Errorf("expected Len 3, got %d", routes.Len())
	}
	if !routes.Less(1, 0) {
		t.Error("expected r[1] < r[0] by metric")
	}
	routes.Swap(0, 1)
	if routes[0].Address != "a" {
		t.Errorf("expected 'a' after swap, got %q", routes[0].Address)
	}
}

// --- NewClientCallOptions ---

func TestNewClientCallOptions_Call(t *testing.T) {
	base := NewClient()
	wrapped := NewClientCallOptions(base, WithAddress("127.0.0.1:8080"))
	req := wrapped.NewRequest("svc", "ep", nil)
	if err := wrapped.Call(context.Background(), req, nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewClientCallOptions_Stream(t *testing.T) {
	base := NewClient()
	wrapped := NewClientCallOptions(base, WithAddress("127.0.0.1:8080"))
	req := wrapped.NewRequest("svc", "ep", nil)
	stream, err := wrapped.Stream(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream != nil {
		_ = stream.Close()
	}
}

// --- testRequest methods (test_request.go) ---

func TestTestRequest_Fields(t *testing.T) {
	tr := &testRequest{
		service:     "svc",
		method:      "method",
		endpoint:    "ep",
		contentType: "application/json",
		body:        "body",
		opts:        RequestOptions{Stream: true},
	}
	if tr.Service() != "svc" {
		t.Errorf("expected svc, got %q", tr.Service())
	}
	if tr.Method() != "method" {
		t.Errorf("expected method, got %q", tr.Method())
	}
	if tr.Endpoint() != "ep" {
		t.Errorf("expected ep, got %q", tr.Endpoint())
	}
	if tr.ContentType() != "application/json" {
		t.Errorf("expected application/json, got %q", tr.ContentType())
	}
	if tr.Body() != "body" {
		t.Errorf("expected body, got %v", tr.Body())
	}
	if tr.Codec() != nil {
		t.Error("expected nil codec")
	}
	if !tr.Stream() {
		t.Error("expected Stream true")
	}
}

// --- options.Hook integration with Init ---

func TestNoopClient_InitWithHook(t *testing.T) {
	var called bool
	hook := HookCall(func(next FuncCall) FuncCall {
		return func(ctx context.Context, req Request, rsp any, opts ...CallOption) error {
			called = true
			return next(ctx, req, rsp, opts...)
		}
	})
	c := NewClient()
	if err := c.Init(Hooks(hook)); err != nil {
		t.Fatalf("unexpected Init error: %v", err)
	}
	req := c.NewRequest("svc", "ep", nil)
	_ = c.Call(context.Background(), req, nil, WithAddress("127.0.0.1:8080"))
	if !called {
		t.Error("expected hook to be called after Init")
	}
}

func TestNoopClient_InitWithStreamHook(t *testing.T) {
	var called bool
	hook := HookStream(func(next FuncStream) FuncStream {
		return func(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
			called = true
			return next(ctx, req, opts...)
		}
	})
	c := NewClient()
	if err := c.Init(Hooks(hook)); err != nil {
		t.Fatalf("unexpected Init error: %v", err)
	}
	req := c.NewRequest("svc", "ep", nil)
	stream, err := c.Stream(context.Background(), req, WithAddress("127.0.0.1:8080"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected stream hook to be called")
	}
	if stream != nil {
		_ = stream.Close()
	}
}

// --- Codec option ---

func TestOptionsCodec(t *testing.T) {
	c := NewClient(Codec("application/json", nil))
	if _, ok := c.Options().Codecs["application/json"]; !ok {
		t.Error("expected codec entry for application/json")
	}
}

// Verify _ usage to satisfy imports
var _ options.Hook = HookCall(nil)