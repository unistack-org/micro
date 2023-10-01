// Package wrapper provides wrapper for Tracer
package wrapper // import "go.unistack.org/micro/v4/tracer/wrapper"

import (
	"context"
	"fmt"
	"strings"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/server"
	"go.unistack.org/micro/v4/tracer"
)

var DefaultHeadersExctract = []string{metadata.HeaderXRequestID}

func ExtractDefaultLabels(md metadata.Metadata) []interface{} {
	labels := make([]interface{}, 0, len(DefaultHeadersExctract))
	for _, k := range DefaultHeadersExctract {
		if v, ok := md.Get(k); ok {
			labels = append(labels, strings.ToLower(k), v)
		}
	}
	return labels
}

var (
	DefaultClientCallObserver = func(ctx context.Context, req client.Request, rsp interface{}, opts []options.Option, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		if err != nil {
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		sp.AddLabels(labels...)
	}

	DefaultClientStreamObserver = func(ctx context.Context, req client.Request, opts []options.Option, stream client.Stream, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		if err != nil {
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		sp.AddLabels(labels...)
	}

	DefaultServerHandlerObserver = func(ctx context.Context, req server.Request, rsp interface{}, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		if err != nil {
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		sp.AddLabels(labels...)
	}

	DefaultClientCallFuncObserver = func(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("%s.%s call", req.Service(), req.Method()))
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		if err != nil {
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		sp.AddLabels(labels...)
	}

	DefaultSkipEndpoints = []string{"Meter.Metrics", "Health.Live", "Health.Ready", "Health.Version"}
)

type tWrapper struct {
	client.Client
	serverHandler  server.HandlerFunc
	clientCallFunc client.CallFunc
	opts           Options
}

type (
	ClientCallObserver     func(context.Context, client.Request, interface{}, []options.Option, tracer.Span, error)
	ClientStreamObserver   func(context.Context, client.Request, []options.Option, client.Stream, tracer.Span, error)
	ClientCallFuncObserver func(context.Context, string, client.Request, interface{}, client.CallOptions, tracer.Span, error)
	ServerHandlerObserver  func(context.Context, server.Request, interface{}, tracer.Span, error)
)

// Options struct
type Options struct {
	// Tracer that used for tracing
	Tracer tracer.Tracer
	// ClientCallObservers funcs
	ClientCallObservers []ClientCallObserver
	// ClientStreamObservers funcs
	ClientStreamObservers []ClientStreamObserver
	// ClientCallFuncObservers funcs
	ClientCallFuncObservers []ClientCallFuncObserver
	// ServerHandlerObservers funcs
	ServerHandlerObservers []ServerHandlerObserver
	// SkipEndpoints
	SkipEndpoints []string
}

// Option func signature
type Option func(*Options)

// NewOptions create Options from Option slice
func NewOptions(opts ...Option) Options {
	options := Options{
		Tracer:                  tracer.DefaultTracer,
		ClientCallObservers:     []ClientCallObserver{DefaultClientCallObserver},
		ClientStreamObservers:   []ClientStreamObserver{DefaultClientStreamObserver},
		ClientCallFuncObservers: []ClientCallFuncObserver{DefaultClientCallFuncObserver},
		ServerHandlerObservers:  []ServerHandlerObserver{DefaultServerHandlerObserver},
		SkipEndpoints:           DefaultSkipEndpoints,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithTracer pass tracer
func WithTracer(t tracer.Tracer) Option {
	return func(o *Options) {
		o.Tracer = t
	}
}

// SkipEndponts
func SkipEndpoins(eps ...string) Option {
	return func(o *Options) {
		o.SkipEndpoints = append(o.SkipEndpoints, eps...)
	}
}

// WithClientCallObservers funcs
func WithClientCallObservers(ob ...ClientCallObserver) Option {
	return func(o *Options) {
		o.ClientCallObservers = ob
	}
}

// WithClientStreamObservers funcs
func WithClientStreamObservers(ob ...ClientStreamObserver) Option {
	return func(o *Options) {
		o.ClientStreamObservers = ob
	}
}

// WithClientCallFuncObservers funcs
func WithClientCallFuncObservers(ob ...ClientCallFuncObserver) Option {
	return func(o *Options) {
		o.ClientCallFuncObservers = ob
	}
}

// WithServerHandlerObservers funcs
func WithServerHandlerObservers(ob ...ServerHandlerObserver) Option {
	return func(o *Options) {
		o.ServerHandlerObservers = ob
	}
}

func (ot *tWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...options.Option) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range ot.opts.SkipEndpoints {
		if ep == endpoint {
			return ot.Client.Call(ctx, req, rsp, opts...)
		}
	}

	nctx, sp := ot.opts.Tracer.Start(ctx, fmt.Sprintf("%s.%s rpc-client", req.Service(), req.Method()),
		tracer.WithSpanKind(tracer.SpanKindClient),
		options.Labels(
			"rpc.service", req.Service(),
			"rpc.method", req.Method(),
			"rpc.flavor", "rpc",
			"rpc.call", "/"+req.Service()+"/"+req.Endpoint(),
			"rpc.call_type", "unary",
		),
	)
	defer sp.Finish()

	err := ot.Client.Call(nctx, req, rsp, opts...)

	for _, o := range ot.opts.ClientCallObservers {
		o(nctx, req, rsp, opts, sp, err)
	}

	return err
}

func (ot *tWrapper) Stream(ctx context.Context, req client.Request, opts ...options.Option) (client.Stream, error) {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range ot.opts.SkipEndpoints {
		if ep == endpoint {
			return ot.Client.Stream(ctx, req, opts...)
		}
	}

	nctx, sp := ot.opts.Tracer.Start(ctx, fmt.Sprintf("%s.%s rpc-client", req.Service(), req.Method()),
		tracer.WithSpanKind(tracer.SpanKindClient),
		options.Labels(
			"rpc.service", req.Service(),
			"rpc.method", req.Method(),
			"rpc.flavor", "rpc",
			"rpc.call", "/"+req.Service()+"/"+req.Endpoint(),
			"rpc.call_type", "stream",
		),
	)
	defer sp.Finish()

	stream, err := ot.Client.Stream(nctx, req, opts...)

	for _, o := range ot.opts.ClientStreamObservers {
		o(nctx, req, opts, stream, sp, err)
	}

	return stream, err
}

func (ot *tWrapper) ServerHandler(ctx context.Context, req server.Request, rsp interface{}) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Method())
	for _, ep := range ot.opts.SkipEndpoints {
		if ep == endpoint {
			return ot.serverHandler(ctx, req, rsp)
		}
	}

	callType := "unary"
	if req.Stream() {
		callType = "stream"
	}

	nctx, sp := ot.opts.Tracer.Start(ctx, fmt.Sprintf("%s.%s rpc-server", req.Service(), req.Method()),
		tracer.WithSpanKind(tracer.SpanKindServer),
		options.Labels(
			"rpc.service", req.Service(),
			"rpc.method", req.Method(),
			"rpc.flavor", "rpc",
			"rpc.call", "/"+req.Service()+"/"+req.Endpoint(),
			"rpc.call_type", callType,
		),
	)
	defer sp.Finish()

	err := ot.serverHandler(nctx, req, rsp)

	for _, o := range ot.opts.ServerHandlerObservers {
		o(nctx, req, rsp, sp, err)
	}

	return err
}

// NewClientWrapper accepts an open tracing Trace and returns a Client Wrapper
func NewClientWrapper(opts ...Option) client.Wrapper {
	return func(c client.Client) client.Client {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}
		return &tWrapper{opts: options, Client: c}
	}
}

// NewClientCallWrapper accepts an opentracing Tracer and returns a Call Wrapper
func NewClientCallWrapper(opts ...Option) client.CallWrapper {
	return func(h client.CallFunc) client.CallFunc {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}

		ot := &tWrapper{opts: options, clientCallFunc: h}
		return ot.ClientCallFunc
	}
}

func (ot *tWrapper) ClientCallFunc(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Method())
	for _, ep := range ot.opts.SkipEndpoints {
		if ep == endpoint {
			return ot.ClientCallFunc(ctx, addr, req, rsp, opts)
		}
	}

	nctx, sp := ot.opts.Tracer.Start(ctx, fmt.Sprintf("%s.%s rpc-client", req.Service(), req.Method()),
		tracer.WithSpanKind(tracer.SpanKindClient),
		options.Labels(
			"rpc.service", req.Service(),
			"rpc.method", req.Method(),
			"rpc.flavor", "rpc",
			"rpc.call", "/"+req.Service()+"/"+req.Endpoint(),
			"rpc.call_type", "unary",
		),
	)

	defer sp.Finish()

	err := ot.clientCallFunc(nctx, addr, req, rsp, opts)

	for _, o := range ot.opts.ClientCallFuncObservers {
		o(nctx, addr, req, rsp, opts, sp, err)
	}

	return err
}

// NewServerHandlerWrapper accepts an options and returns a Handler Wrapper
func NewServerHandlerWrapper(opts ...Option) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}

		ot := &tWrapper{opts: options, serverHandler: h}
		return ot.ServerHandler
	}
}
