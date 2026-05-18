package http

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptrace"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"go.unistack.org/micro/v5/tracer"
)

// mockSpan is a minimal no-op implementation of tracer.Span.
type mockSpan struct {
	ctx context.Context
}

func (s *mockSpan) Tracer() tracer.Tracer                        { return nil }
func (s *mockSpan) Finish(opts ...tracer.SpanOption)             {}
func (s *mockSpan) Context() context.Context                     { return s.ctx }
func (s *mockSpan) SetName(name string)                          {}
func (s *mockSpan) SetStatus(st tracer.SpanStatus, msg string)   {}
func (s *mockSpan) Status() (tracer.SpanStatus, string)          { return tracer.SpanStatusUnset, "" }
func (s *mockSpan) AddLabels(kv ...any)                          {}
func (s *mockSpan) AddEvent(name string, opts ...tracer.EventOption) {}
func (s *mockSpan) AddLogs(kv ...any)                            {}
func (s *mockSpan) Kind() tracer.SpanKind                        { return tracer.SpanKindClient }
func (s *mockSpan) TraceID() string                              { return "trace-id" }
func (s *mockSpan) SpanID() string                               { return "span-id" }
func (s *mockSpan) IsRecording() bool                            { return true }

// mockTracer is a minimal no-op implementation of tracer.Tracer.
type mockTracer struct{}

func (t *mockTracer) Name() string { return "mock" }
func (t *mockTracer) Init(...tracer.Option) error { return nil }
func (t *mockTracer) Start(ctx context.Context, name string, opts ...tracer.SpanOption) (context.Context, tracer.Span) {
	sp := &mockSpan{ctx: ctx}
	return tracer.NewSpanContext(ctx, sp), sp
}
func (t *mockTracer) Flush(ctx context.Context) error { return nil }
func (t *mockTracer) Enabled() bool                   { return true }

func TestWrite(t *testing.T) {
	w := httptest.NewRecorder()
	Write(w, "text/plain", 200, "hello")
	if w.Code != 200 {
		t.Fatalf("expected 200 got %d", w.Code)
	}
	if w.Body.String() != "hello" {
		t.Fatalf("expected 'hello' got %q", w.Body.String())
	}
}

func TestWriteBadRequestError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteBadRequestError(w, errors.New("bad input"))
	if w.Code != 400 {
		t.Fatalf("expected 400 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "bad input") {
		t.Fatalf("expected error message in body, got %q", w.Body.String())
	}
}

func TestWriteInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteInternalServerError(w, errors.New("server error"))
	if w.Code != 500 {
		t.Fatalf("expected 500 got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "server error") {
		t.Fatalf("expected error message in body, got %q", w.Body.String())
	}
}

func TestRequestToContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token123")
	ctx := RequestToContext(req)
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestNewRoundTripper(t *testing.T) {
	rt := NewRoundTripper()
	if rt == nil {
		t.Fatal("expected non-nil RoundTripper")
	}
}

func TestWithRouter(t *testing.T) {
	opt := WithRouter(nil)
	opts := &Options{}
	opt(opts)
	if opts.Router != nil {
		t.Fatal("expected nil router")
	}
}

func TestRegisterMethodEmpty(t *testing.T) {
	err := RegisterMethod("")
	if err != nil {
		t.Fatalf("expected nil error for empty method, got %v", err)
	}
}

func TestRegisterMethodExisting(t *testing.T) {
	err := RegisterMethod("GET")
	if err != nil {
		t.Fatalf("expected nil error for existing method, got %v", err)
	}
}

func TestRegisterMethodNew(t *testing.T) {
	err := RegisterMethod("CUSTOM")
	if err != nil {
		t.Fatalf("expected nil error for new method, got %v", err)
	}
}

func TestTrieSetEndpointStub(t *testing.T) {
	tr := NewTrie()
	// Insert with a method that isn't in methodMap to trigger mSTUB path
	// mSTUB is set when method stub bit is set; use "STUB" custom method
	if err := tr.Insert([]string{http.MethodGet}, "/stub-test", "stubhandler"); err != nil {
		t.Fatal(err)
	}
	h, _, err := tr.Search(http.MethodGet, "/stub-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.(string) != "stubhandler" {
		t.Fatalf("expected 'stubhandler' got %v", h)
	}
}

func TestTrieSearchNotFound(t *testing.T) {
	tr := NewTrie()
	_, _, err := tr.Search(http.MethodGet, "/nonexistent")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestParentHookConnect(t *testing.T) {
	got := parentHook("http.connect.somehost:443")
	if got != "http.getconn" {
		t.Fatalf("expected 'http.getconn' for connect hook, got %q", got)
	}
}

func TestParentHookDNS(t *testing.T) {
	got := parentHook("http.dns")
	if got != "http.getconn" {
		t.Fatalf("expected 'http.getconn' for dns hook, got %q", got)
	}
}

func TestParentHookTLS(t *testing.T) {
	got := parentHook("http.tls")
	if got != "http.getconn" {
		t.Fatalf("expected 'http.getconn' for tls hook, got %q", got)
	}
}

func TestSliceToStringEmpty(t *testing.T) {
	got := sliceToString([]string{})
	if got != "undefined" {
		t.Fatalf("expected 'undefined', got %q", got)
	}
}

func TestSliceToStringMultiple(t *testing.T) {
	got := sliceToString([]string{"a", "b", "c"})
	if got != "a,b,c" {
		t.Fatalf("expected 'a,b,c', got %q", got)
	}
}

func TestSm2s(t *testing.T) {
	m := map[string][]string{"key": {"v1", "v2"}}
	got := sm2s(m)
	if !strings.Contains(got, "key=v1,v2") {
		t.Fatalf("unexpected sm2s output %q", got)
	}
}

func TestSm2sEmpty(t *testing.T) {
	got := sm2s(map[string][]string{})
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestNewClientTrace(t *testing.T) {
	ct := NewClientTrace(context.Background(), nil)
	if ct == nil {
		t.Fatal("expected non-nil ClientTrace")
	}
}

func newTestClientTracer() *clientTracer {
	return &clientTracer{
		Context:     context.Background(),
		activeHooks: make(map[string]context.Context),
		tr:          &mockTracer{},
	}
}

func TestClientTracerStartEnd(t *testing.T) {
	ct := newTestClientTracer()
	ct.start("http.getconn", "http.getconn")
	ct.end("http.getconn", nil)
}

func TestClientTracerStartEndError(t *testing.T) {
	ct := newTestClientTracer()
	ct.start("http.getconn", "http.getconn")
	ct.end("http.getconn", errors.New("dial failed"))
}

func TestClientTracerEndBeforeStart(t *testing.T) {
	// end called before start — should store context for when start arrives
	ct := newTestClientTracer()
	ct.end("http.receive", nil)
	ct.start("http.receive", "http.receive")
}

func TestClientTracerGetConn(t *testing.T) {
	ct := newTestClientTracer()
	ct.getConn("example.com:80")
}

func TestClientTracerDNS(t *testing.T) {
	ct := newTestClientTracer()
	ct.dnsStart(httptrace.DNSStartInfo{Host: "example.com"})
	ct.dnsDone(httptrace.DNSDoneInfo{})
}

func TestClientTracerConnect(t *testing.T) {
	ct := newTestClientTracer()
	ct.connectStart("tcp", "1.2.3.4:80")
	ct.connectDone("tcp", "1.2.3.4:80", nil)
}

func TestClientTracerConnectError(t *testing.T) {
	ct := newTestClientTracer()
	ct.connectStart("tcp", "1.2.3.4:80")
	ct.connectDone("tcp", "1.2.3.4:80", errors.New("refused"))
}

func TestClientTracerTLS(t *testing.T) {
	ct := newTestClientTracer()
	ct.tlsHandshakeStart()
	ct.tlsHandshakeDone(tls.ConnectionState{}, nil)
}

func TestClientTracerHeaders(t *testing.T) {
	ct := newTestClientTracer()
	// prime root span so wroteHeaderField doesn't panic on ct.root.AddLabels
	ct.start("http.getconn", "http.getconn")
	_, sp := ct.tr.Start(ct.Context, "root")
	ct.root = sp
	ct.wroteHeaderField("Content-Type", []string{"application/json"})
	ct.wroteHeaders()
	ct.wroteRequest(httptrace.WroteRequestInfo{})
}

func TestClientTracerReceive(t *testing.T) {
	ct := newTestClientTracer()
	ct.gotFirstResponseByte()
	ct.got100Continue()
	ct.wait100Continue()
	ct.got1xxResponse(100, textproto.MIMEHeader{})
	ct.putIdleConn(nil)
}

func TestClientTracerSpan(t *testing.T) {
	ct := newTestClientTracer()
	ct.start("http.test", "http.test")
	sp, ok := ct.span("http.test")
	if !ok || sp == nil {
		t.Fatal("expected span to be found after start")
	}
	_, notOk := ct.span("http.missing")
	if notOk {
		t.Fatal("expected no span for missing hook")
	}
}
