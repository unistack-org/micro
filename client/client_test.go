package client

import (
	"context"
	"net"
	"testing"
	"time"

	"crypto/tls"

	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
)

func TestNoopClient_Name(t *testing.T) {
	c := NewClient()
	if name := c.Name(); name != "" {
		t.Errorf("expected '', got %q", name)
	}
}

func TestNoopClient_Options(t *testing.T) {
	c := NewClient()
	opts := c.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
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

func TestNoopClient_NewRequest(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("service", "endpoint", "body")
	if req == nil {
		t.Error("expected non-nil request")
	}
	if req.Service() != "service" {
		t.Errorf("expected service 'service', got %q", req.Service())
	}
	if req.Endpoint() != "endpoint" {
		t.Errorf("expected endpoint 'endpoint', got %q", req.Endpoint())
	}
	if req.Method() != "" {
		t.Errorf("expected empty method, got %q", req.Method())
	}
}

func TestNoopClient_Call(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	req := c.NewRequest("service", "endpoint", "body")
	rsp := &struct{}{}
	// Use static address to skip service discovery
	err := c.Call(ctx, req, rsp, WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopClient_Stream(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	req := c.NewRequest("service", "endpoint", "body")
	// Use static address to skip service discovery
	stream, err := c.Stream(ctx, req, WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Error("expected non-nil stream")
	}
}

func TestMustContext(t *testing.T) {
	c := NewClient()
	ctx := NewContext(context.Background(), c)
	retrieved := MustContext(ctx)
	if retrieved != c {
		t.Error("expected client from context")
	}
}

func TestMustContextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing client")
		}
	}()
	MustContext(context.Background())
}

func TestNoopRequest_ContentType(t *testing.T) {
	req := &noopRequest{contentType: "application/json"}
	if req.ContentType() != "application/json" {
		t.Errorf("expected application/json, got %q", req.ContentType())
	}
}

func TestNoopRequest_Body(t *testing.T) {
	body := "test body"
	req := &noopRequest{body: body}
	if req.Body() != body {
		t.Error("expected body match")
	}
}

func TestNoopRequest_Codec(t *testing.T) {
	req := &noopRequest{}
	if req.Codec() != nil {
		t.Error("expected nil codec")
	}
}

func TestNoopRequest_Stream(t *testing.T) {
	req := &noopRequest{stream: true}
	if !req.Stream() {
		t.Error("expected stream true")
	}
}

func TestNoopResponse_Codec(t *testing.T) {
	resp := &noopResponse{}
	if resp.Codec() != nil {
		t.Error("expected nil codec")
	}
}

func TestNoopResponse_Header(t *testing.T) {
	resp := &noopResponse{header: metadata.Metadata{"key": []string{"value"}}}
	if resp.Header()["key"][0] != "value" {
		t.Error("expected header value")
	}
}

func TestNoopResponse_Read(t *testing.T) {
	resp := &noopResponse{}
	data, err := resp.Read()
	if err != nil || data != nil {
		t.Error("expected nil data and no error")
	}
}

func TestNoopStream_Context(t *testing.T) {
	ctx := context.Background()
	stream := &noopStream{ctx: ctx}
	if stream.Context() != ctx {
		t.Error("expected stream context")
	}
}

func TestNoopStream_Request(t *testing.T) {
	stream := &noopStream{}
	req := stream.Request()
	if req == nil {
		t.Error("expected non-nil request")
	}
}

func TestNoopStream_Response(t *testing.T) {
	stream := &noopStream{}
	resp := stream.Response()
	if resp == nil {
		t.Error("expected non-nil response")
	}
}

func TestNoopStream_Send(t *testing.T) {
	stream := &noopStream{}
	if err := stream.Send("test"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopStream_Recv(t *testing.T) {
	stream := &noopStream{}
	if err := stream.Recv("test"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopStream_Error(t *testing.T) {
	stream := &noopStream{err: context.Canceled}
	if stream.Error() != context.Canceled {
		t.Errorf("expected canceled error, got %v", stream.Error())
	}
}

func TestNoopStream_Close(t *testing.T) {
	stream := &noopStream{}
	if err := stream.Close(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNameOption(t *testing.T) {
	c := NewClient(Name("test-client"))
	if c.Name() != "test-client" {
		t.Errorf("expected test-client, got %q", c.Name())
	}
}

func TestInitWithOptions(t *testing.T) {
	c := NewClient()
	err := c.Init(Name("updated"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.Name() != "updated" {
		t.Errorf("expected updated, got %q", c.Name())
	}
}

func TestContextOption(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	c := NewClient(Context(ctx))
	if c.Options().Context.Value("key") != "value" {
		t.Error("expected context value")
	}
}

func TestContextDialerOption(t *testing.T) {
	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		return nil, nil
	}
	c := NewClient(ContextDialer(dialer))
	if c.Options().ContextDialer == nil {
		t.Error("expected context dialer")
	}
}

func TestTLSConfigOption(t *testing.T) {
	tlsConfig := &tls.Config{}
	c := NewClient(TLSConfig(tlsConfig))
	if c.Options().TLSConfig != tlsConfig {
		t.Error("expected tls config")
	}
}

func TestRetriesOption(t *testing.T) {
	c := NewClient(Retries(5))
	if c.Options().CallOptions.Retries != 5 {
		t.Error("expected retries 5")
	}
}

func TestRequestTimeoutOption(t *testing.T) {
	timeout := time.Second * 10
	c := NewClient(RequestTimeout(timeout))
	if c.Options().CallOptions.RequestTimeout != timeout {
		t.Error("expected request timeout")
	}
}

func TestDialTimeoutOption(t *testing.T) {
	timeout := time.Second * 5
	c := NewClient(DialTimeout(timeout))
	if c.Options().CallOptions.DialTimeout != timeout {
		t.Error("expected dial timeout")
	}
}

func TestHooksOption(t *testing.T) {
	// Test that Hooks option is accepted
	hook := func(next options.Hook) options.Hook {
		return next
	}
	c := NewClient(Hooks(hook))
	if len(c.Options().Hooks) == 0 {
		t.Error("expected hook to be added")
	}
}

func TestLookupOption(t *testing.T) {
	lookup := func(ctx context.Context, req Request, opts CallOptions) ([]string, error) {
		return []string{"addr"}, nil
	}
	c := NewClient(Lookup(lookup))
	if c.Options().Lookup == nil {
		t.Error("expected lookup")
	}
}

func TestContentTypeOption(t *testing.T) {
	c := NewClient(ContentType("application/json"))
	if c.Options().ContentType != "application/json" {
		t.Error("expected content type")
	}
}

func TestProxyOption(t *testing.T) {
	c := NewClient(Proxy("proxy:8080"))
	if c.Options().Proxy != "proxy:8080" {
		t.Error("expected proxy")
	}
}

func TestPoolSizeOption(t *testing.T) {
	c := NewClient(PoolSize(200))
	if c.Options().PoolSize != 200 {
		t.Error("expected pool size 200")
	}
}

func TestPoolTTLOption(t *testing.T) {
	ttl := time.Minute * 5
	c := NewClient(PoolTTL(ttl))
	if c.Options().PoolTTL != ttl {
		t.Error("expected pool ttl")
	}
}

func TestNewRequestWithOptions(t *testing.T) {
	req := NewClient().NewRequest("svc", "ep", "body")
	if req.Stream() {
		t.Error("expected non-streaming request")
	}
	if req.ContentType() != "" {
		t.Error("expected empty content type")
	}
}

func TestCallOptionWithAddress(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "ep", "body")
	err := c.Call(context.Background(), req, &struct{}{}, WithAddress("addr1", "addr2"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCallOptionWithRetries(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "ep", "body")
	err := c.Call(context.Background(), req, &struct{}{}, WithRetries(3), WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStreamOptionWithAddress(t *testing.T) {
	c := NewClient()
	req := c.NewRequest("svc", "ep", "body")
	stream, err := c.Stream(context.Background(), req, WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Error("expected non-nil stream")
	}
}

func TestWithRequestTimeoutCallOption(t *testing.T) {
	timeout := time.Second * 5
	req := NewClient().NewRequest("svc", "ep", "body")
	err := NewClient().Call(context.Background(), req, &struct{}{}, WithRequestTimeout(timeout), WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWithStreamTimeoutCallOption(t *testing.T) {
	timeout := time.Second * 10
	req := NewClient().NewRequest("svc", "ep", "body")
	stream, err := NewClient().Stream(context.Background(), req, WithStreamTimeout(timeout), WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Error("expected non-nil stream")
	}
}

func TestWithAuthTokenCallOption(t *testing.T) {
	req := NewClient().NewRequest("svc", "ep", "body")
	err := NewClient().Call(context.Background(), req, &struct{}{}, WithAuthToken("token"), WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWithRetriesCallOption(t *testing.T) {
	req := NewClient().NewRequest("svc", "ep", "body")
	err := NewClient().Call(context.Background(), req, &struct{}{}, WithRetries(3), WithAddress("test"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWrapperTypes(t *testing.T) {
	// Test Wrapper type
	var w Wrapper
	w = func(c Client) Client { return c }
	if w == nil {
		t.Error("expected wrapper")
	}
	// Test CallWrapper type
	var cw CallWrapper
	cw = func(cf CallFunc) CallFunc { return cf }
	if cw == nil {
		t.Error("expected call wrapper")
	}
	// Test StreamWrapper type
	var sw StreamWrapper
	sw = func(s Stream) Stream { return s }
	if sw == nil {
		t.Error("expected stream wrapper")
	}
}

func TestNewClientCallOptions(t *testing.T) {
	c := NewClient()
	// Add a call option that sets address
	cco := NewClientCallOptions(c, WithAddress("test"))
	req := c.NewRequest("svc", "ep", "body")
	rsp := &struct{}{}
	// Call should use the preset address
	err := cco.Call(context.Background(), req, rsp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Test Stream
	stream, err := cco.Stream(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if stream == nil {
		t.Error("expected non-nil stream")
	}
}


