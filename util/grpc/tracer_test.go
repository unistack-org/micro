package grpc_util

import (
	"context"
	"net"
	"testing"
	"time"

	"go.unistack.org/micro/v5/tracer"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

func TestTracer_Option(t *testing.T) {
	tr := tracer.DefaultTracer
	opts := Options{}
	Tracer(tr)(&opts)
	if opts.Tracer != tr {
		t.Error("expected tracer to be set")
	}
}

func TestNewServerHandler(t *testing.T) {
	h := NewServerHandler()
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestNewClientHandler(t *testing.T) {
	h := NewClientHandler()
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestServerHandler_TagRPC(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Service/Method"})
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestServerHandler_TagRPC_InvalidFormat(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "noSlash"})
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestServerHandler_TagRPC_NoMethod(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Service/"})
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestServerHandler_HandleRPC_Begin(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.Begin{
		BeginTime:      time.Now(),
		IsClientStream: false,
		IsServerStream: false,
	})
}

func TestServerHandler_HandleRPC_BeginStream(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Stream"})
	h.HandleRPC(ctx, &stats.Begin{
		BeginTime:      time.Now(),
		IsClientStream: true,
		IsServerStream: false,
	})
}

func TestServerHandler_HandleRPC_InPayload(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.InPayload{
		RecvTime:         time.Now(),
		Length:           100,
		CompressedLength: 90,
	})
}

func TestServerHandler_HandleRPC_OutPayload(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.OutPayload{
		SentTime:         time.Now(),
		Length:           200,
		CompressedLength: 180,
	})
}

func TestServerHandler_HandleRPC_End_NoError(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.End{
		BeginTime: time.Now(),
		EndTime:   time.Now(),
		Error:     nil,
	})
}

func TestServerHandler_HandleRPC_End_WithError(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	st := status.New(grpc_codes.Internal, "internal error")
	h.HandleRPC(ctx, &stats.End{
		BeginTime: time.Now(),
		EndTime:   time.Now(),
		Error:     st.Err(),
	})
}

func TestServerHandler_HandleRPC_Unknown(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	// unknown stats type — use InHeader which implements RPCStats but is not handled
	h.HandleRPC(ctx, &stats.InHeader{})
}

func TestServerHandler_TagConn(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	ctx2 := h.TagConn(ctx, &stats.ConnTagInfo{})
	if ctx2 == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestServerHandler_TagConn_WithPeer(t *testing.T) {
	h := NewServerHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
	p := &peer.Peer{Addr: addr}
	ctx = peer.NewContext(ctx, p)
	ctx2 := h.TagConn(ctx, &stats.ConnTagInfo{})
	if ctx2 == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestServerHandler_HandleConn(t *testing.T) {
	h := NewServerHandler()
	// no-op
	h.HandleConn(context.Background(), &stats.ConnEnd{})
}

func TestClientHandler_TagRPC(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Service/Method"})
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestClientHandler_HandleRPC_Begin(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.Begin{BeginTime: time.Now()})
}

func TestClientHandler_HandleRPC_InPayload(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.InPayload{RecvTime: time.Now(), Length: 50, CompressedLength: 40})
}

func TestClientHandler_HandleRPC_OutPayload(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.OutPayload{SentTime: time.Now(), Length: 80, CompressedLength: 60})
}

func TestClientHandler_HandleRPC_End(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.End{BeginTime: time.Now(), EndTime: time.Now()})
}

func TestClientHandler_TagConn_WithRemoteAddr(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	addr := &net.TCPAddr{IP: net.ParseIP("10.0.0.1"), Port: 9000}
	ctx2 := h.TagConn(ctx, &stats.ConnTagInfo{RemoteAddr: addr})
	if ctx2 == nil {
		t.Fatal("expected non-nil context")
	}
}

func TestClientHandler_HandleConn(t *testing.T) {
	h := NewClientHandler()
	h.HandleConn(context.Background(), &stats.ConnEnd{})
}

func TestParseFullMethod_Valid(t *testing.T) {
	name, attrs := parseFullMethod("/my.package.Service/MyMethod")
	if name != "my.package.Service/MyMethod" {
		t.Errorf("unexpected name: %q", name)
	}
	if len(attrs) == 0 {
		t.Error("expected attrs")
	}
}

func TestParseFullMethod_NoLeadingSlash(t *testing.T) {
	name, attrs := parseFullMethod("noSlash")
	if name != "noSlash" {
		t.Errorf("unexpected name: %q", name)
	}
	if len(attrs) != 0 {
		t.Errorf("expected no attrs, got %v", attrs)
	}
}

func TestParseFullMethod_NoInternalSlash(t *testing.T) {
	name, attrs := parseFullMethod("/onlyservice")
	if name != "onlyservice" {
		t.Errorf("unexpected name: %q", name)
	}
	if len(attrs) != 0 {
		t.Errorf("expected no attrs, got %v", attrs)
	}
}

func TestPeerAttr_ValidIP(t *testing.T) {
	attrs := peerAttr("127.0.0.1:8080")
	if len(attrs) == 0 {
		t.Error("expected non-empty attrs for valid IP")
	}
}

func TestPeerAttr_Hostname(t *testing.T) {
	attrs := peerAttr("myhost:9090")
	if len(attrs) == 0 {
		t.Error("expected non-empty attrs for hostname")
	}
}

func TestPeerAttr_InvalidAddr(t *testing.T) {
	attrs := peerAttr("invalid")
	if len(attrs) != 0 {
		t.Errorf("expected empty attrs for invalid addr, got %v", attrs)
	}
}

func TestPeerAttr_InvalidPort(t *testing.T) {
	attrs := peerAttr("host:notaport")
	if len(attrs) != 0 {
		t.Errorf("expected empty attrs for invalid port, got %v", attrs)
	}
}

func TestPeerFromCtx_NoPeer(t *testing.T) {
	addr := peerFromCtx(context.Background())
	if addr != "" {
		t.Errorf("expected empty addr, got %q", addr)
	}
}

func TestPeerFromCtx_WithPeer(t *testing.T) {
	addr := &net.TCPAddr{IP: net.ParseIP("192.168.1.1"), Port: 1234}
	p := &peer.Peer{Addr: addr}
	ctx := peer.NewContext(context.Background(), p)
	result := peerFromCtx(ctx)
	if result == "" {
		t.Error("expected non-empty addr")
	}
}

func TestHandleRPC_ClientSide(t *testing.T) {
	h := NewClientHandler()
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	// InPayload on client side (IsClient() == true)
	h.HandleRPC(ctx, &stats.InPayload{RecvTime: time.Now()})
}

func TestHandleRPC_InHeader(t *testing.T) {
	h := NewServerHandler()
	// TagRPC first to establish span, then call HandleRPC with unhandled type
	ctx := h.TagRPC(context.Background(), &stats.RPCTagInfo{FullMethodName: "/pkg.Svc/Method"})
	h.HandleRPC(ctx, &stats.InHeader{})
}
