package grpc_util

import (
	"context"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go.unistack.org/micro/v4/tracer"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

type gRPCContextKey struct{}

type gRPCContext struct {
	messagesReceived int64
	messagesSent     int64
}

type Options struct {
	Tracer tracer.Tracer
}

type Option func(*Options)

func Tracer(tr tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = tr
	}
}

// NewServerHandler creates a stats.Handler for gRPC server.
func NewServerHandler(opts ...Option) stats.Handler {
	options := Options{Tracer: tracer.DefaultTracer}
	for _, o := range opts {
		o(&options)
	}
	h := &serverHandler{
		opts: options,
	}
	return h
}

type serverHandler struct {
	opts Options
}

// TagRPC can attach some information to the given context.
func (h *serverHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	name, attrs := parseFullMethod(info.FullMethodName)
	attrs = append(attrs, "rpc.system", "grpc")
	ctx, _ = h.opts.Tracer.Start(
		ctx,
		name,
		tracer.WithSpanKind(tracer.SpanKindServer),
		tracer.WithSpanLabels(attrs...),
	)

	gctx := gRPCContext{}
	return context.WithValue(ctx, gRPCContextKey{}, &gctx)
}

// HandleRPC processes the RPC stats.
func (h *serverHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	handleRPC(ctx, rs)
}

// TagConn can attach some information to the given context.
func (h *serverHandler) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	if span, ok := tracer.SpanFromContext(ctx); ok {
		attrs := peerAttr(peerFromCtx(ctx))
		span.AddLabels(attrs...)
	}
	return ctx
}

// HandleConn processes the Conn stats.
func (h *serverHandler) HandleConn(_ context.Context, _ stats.ConnStats) {
}

type clientHandler struct {
	opts Options
}

// NewClientHandler creates a stats.Handler for gRPC client.
func NewClientHandler(opts ...Option) stats.Handler {
	options := Options{Tracer: tracer.DefaultTracer}
	for _, o := range opts {
		o(&options)
	}
	h := &clientHandler{
		opts: options,
	}
	return h
}

// TagRPC can attach some information to the given context.
func (h *clientHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	name, attrs := parseFullMethod(info.FullMethodName)
	attrs = append(attrs, "rpc.system", "grpc", "rpc.flavor", "grpc", "rpc.call", info.FullMethodName)
	ctx, _ = h.opts.Tracer.Start(
		ctx,
		name,
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels(attrs...),
	)

	gctx := gRPCContext{}

	return context.WithValue(ctx, gRPCContextKey{}, &gctx)
}

// HandleRPC processes the RPC stats.
func (h *clientHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	handleRPC(ctx, rs)
}

// TagConn can attach some information to the given context.
func (h *clientHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	// TODO
	if span, ok := tracer.SpanFromContext(ctx); ok {
		attrs := peerAttr(cti.RemoteAddr.String())
		span.AddLabels(attrs...)
	}
	return ctx
}

// HandleConn processes the Conn stats.
func (h *clientHandler) HandleConn(context.Context, stats.ConnStats) {
	// no-op
}

func handleRPC(ctx context.Context, rs stats.RPCStats) {
	span, ok := tracer.SpanFromContext(ctx)
	gctx, _ := ctx.Value(gRPCContextKey{}).(*gRPCContext)
	var messageID int64
	if rs.IsClient() {
		span.AddLabels("span.kind", "client")
	} else {
		span.AddLabels("span.kind", "server")
	}

	switch rs := rs.(type) {
	case *stats.Begin:
		if rs.IsClientStream || rs.IsServerStream {
			span.AddLabels("rpc.call_type", "stream")
		} else {
			span.AddLabels("rpc.call_type", "unary")
		}
		span.AddEvent("message",
			tracer.WithEventLabels(
				"message.begin_time", rs.BeginTime.Format(time.RFC3339),
			),
		)
	case *stats.InPayload:
		if gctx != nil {
			messageID = atomic.AddInt64(&gctx.messagesReceived, 1)
		}
		if ok {
			span.AddEvent("message",
				tracer.WithEventLabels(
					"message.recv_time", rs.RecvTime.Format(time.RFC3339),
					"message.type", "RECEIVED",
					"message.id", messageID,
					"message.compressed_size", rs.CompressedLength,
					"message.uncompressed_size", rs.Length,
				),
			)
		}
	case *stats.OutPayload:
		if gctx != nil {
			messageID = atomic.AddInt64(&gctx.messagesSent, 1)
		}
		if ok {
			span.AddEvent("message",
				tracer.WithEventLabels(
					"message.sent_time", rs.SentTime.Format(time.RFC3339),
					"message.type", "SENT",
					"message.id", messageID,
					"message.compressed_size", rs.CompressedLength,
					"message.uncompressed_size", rs.Length,
				),
			)
		}
	case *stats.End:
		if ok {
			span.AddEvent("message",
				tracer.WithEventLabels(
					"message.begin_time", rs.BeginTime.Format(time.RFC3339),
					"message.end_time", rs.EndTime.Format(time.RFC3339),
				),
			)
			if rs.Error != nil {
				s, _ := status.FromError(rs.Error)
				span.SetStatus(tracer.SpanStatusError, s.Message())
				span.AddLabels("rpc.grpc.status_code", s.Code())
			} else {
				span.AddLabels("rpc.grpc.status_code", grpc_codes.OK)
			}
			span.Finish()
		}
	default:
		return
	}
}

func parseFullMethod(fullMethod string) (string, []interface{}) {
	if !strings.HasPrefix(fullMethod, "/") {
		// Invalid format, does not follow `/package.service/method`.
		return fullMethod, nil
	}
	name := fullMethod[1:]
	pos := strings.LastIndex(name, "/")
	if pos < 0 {
		// Invalid format, does not follow `/package.service/method`.
		return name, nil
	}
	service, method := name[:pos], name[pos+1:]

	var attrs []interface{}
	if service != "" {
		attrs = append(attrs, "rpc.service", service)
	}
	if method != "" {
		attrs = append(attrs, "rpc.method", method)
	}
	return name, attrs
}

func peerAttr(addr string) []interface{} {
	host, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil
	}

	if host == "" {
		host = "127.0.0.1"
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		return nil
	}

	var attr []interface{}
	if ip := net.ParseIP(host); ip != nil {
		attr = []interface{}{
			"net.sock.peer.addr", host,
			"net.sock.peer.port", port,
		}
	} else {
		attr = []interface{}{
			"net.peer.name", host,
			"net.peer.port", port,
		}
	}

	return attr
}

func peerFromCtx(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	return p.Addr.String()
}
