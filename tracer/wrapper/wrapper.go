// Package wrapper provides wrapper for Tracer
package wrapper // import "go.unistack.org/micro/v3/tracer/wrapper"

import (
	"context"
	"fmt"
	"strings"

	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/server"
	"go.unistack.org/micro/v3/tracer"
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
	DefaultClientCallObserver = func(ctx context.Context, req client.Request, rsp interface{}, opts []client.CallOption, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		if err != nil {
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		sp.AddLabels(labels...)
	}

	DefaultClientStreamObserver = func(ctx context.Context, req client.Request, opts []client.CallOption, stream client.Stream, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		if err != nil {
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		sp.AddLabels(labels...)
	}

	DefaultClientPublishObserver = func(ctx context.Context, msg client.Message, opts []client.PublishOption, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		labels = append(labels, ExtractDefaultLabels(msg.Metadata())...)
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

	DefaultServerSubscriberObserver = func(ctx context.Context, msg server.Message, sp tracer.Span, err error) {
		var labels []interface{}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			labels = append(labels, ExtractDefaultLabels(md)...)
		}
		labels = append(labels, ExtractDefaultLabels(msg.Header())...)

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
	serverHandler    server.HandlerFunc
	serverSubscriber server.SubscriberFunc
	clientCallFunc   client.CallFunc
	opts             Options
}

type (
	ClientCallObserver       func(context.Context, client.Request, interface{}, []client.CallOption, tracer.Span, error)
	ClientStreamObserver     func(context.Context, client.Request, []client.CallOption, client.Stream, tracer.Span, error)
	ClientPublishObserver    func(context.Context, client.Message, []client.PublishOption, tracer.Span, error)
	ClientCallFuncObserver   func(context.Context, string, client.Request, interface{}, client.CallOptions, tracer.Span, error)
	ServerHandlerObserver    func(context.Context, server.Request, interface{}, tracer.Span, error)
	ServerSubscriberObserver func(context.Context, server.Message, tracer.Span, error)
)

// Options struct
type Options struct {
	// Tracer that used for tracing
	Tracer tracer.Tracer
	// ClientCallObservers funcs
	ClientCallObservers []ClientCallObserver
	// ClientStreamObservers funcs
	ClientStreamObservers []ClientStreamObserver
	// ClientPublishObservers funcs
	ClientPublishObservers []ClientPublishObserver
	// ClientCallFuncObservers funcs
	ClientCallFuncObservers []ClientCallFuncObserver
	// ServerHandlerObservers funcs
	ServerHandlerObservers []ServerHandlerObserver
	// ServerSubscriberObservers funcs
	ServerSubscriberObservers []ServerSubscriberObserver
	// SkipEndpoints
	SkipEndpoints []string
}

// Option func signature
type Option func(*Options)

// NewOptions create Options from Option slice
func NewOptions(opts ...Option) Options {
	options := Options{
		Tracer:                    tracer.DefaultTracer,
		ClientCallObservers:       []ClientCallObserver{DefaultClientCallObserver},
		ClientStreamObservers:     []ClientStreamObserver{DefaultClientStreamObserver},
		ClientPublishObservers:    []ClientPublishObserver{DefaultClientPublishObserver},
		ClientCallFuncObservers:   []ClientCallFuncObserver{DefaultClientCallFuncObserver},
		ServerHandlerObservers:    []ServerHandlerObserver{DefaultServerHandlerObserver},
		ServerSubscriberObservers: []ServerSubscriberObserver{DefaultServerSubscriberObserver},
		SkipEndpoints:             DefaultSkipEndpoints,
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

// WithClientPublishObservers funcs
func WithClientPublishObservers(ob ...ClientPublishObserver) Option {
	return func(o *Options) {
		o.ClientPublishObservers = ob
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

// WithServerSubscriberObservers funcs
func WithServerSubscriberObservers(ob ...ServerSubscriberObserver) Option {
	return func(o *Options) {
		o.ServerSubscriberObservers = ob
	}
}

func (ot *tWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range ot.opts.SkipEndpoints {
		if ep == endpoint {
			return ot.Client.Call(ctx, req, rsp, opts...)
		}
	}

	nctx, sp := ot.opts.Tracer.Start(ctx, fmt.Sprintf("%s.%s rpc-client", req.Service(), req.Method()),
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels(
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

func (ot *tWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range ot.opts.SkipEndpoints {
		if ep == endpoint {
			return ot.Client.Stream(ctx, req, opts...)
		}
	}

	nctx, sp := ot.opts.Tracer.Start(ctx, fmt.Sprintf("%s.%s rpc-client", req.Service(), req.Method()),
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels(
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

func (ot *tWrapper) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	nctx, sp := ot.opts.Tracer.Start(ctx, msg.Topic()+" publish", tracer.WithSpanKind(tracer.SpanKindProducer))
	defer sp.Finish()
	sp.AddLabels("messaging.destination.name", msg.Topic())
	sp.AddLabels("messaging.operation", "publish")
	err := ot.Client.Publish(nctx, msg, opts...)

	for _, o := range ot.opts.ClientPublishObservers {
		o(nctx, msg, opts, sp, err)
	}

	return err
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
		tracer.WithSpanLabels(
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

func (ot *tWrapper) ServerSubscriber(ctx context.Context, msg server.Message) error {
	nctx, sp := ot.opts.Tracer.Start(ctx, msg.Topic()+" process", tracer.WithSpanKind(tracer.SpanKindConsumer))
	defer sp.Finish()
	sp.AddLabels("messaging.operation", "process")
	sp.AddLabels("messaging.source.name", msg.Topic())
	err := ot.serverSubscriber(nctx, msg)

	for _, o := range ot.opts.ServerSubscriberObservers {
		o(nctx, msg, sp, err)
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
		tracer.WithSpanLabels(
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

// NewServerSubscriberWrapper accepts an opentracing Tracer and returns a Subscriber Wrapper
func NewServerSubscriberWrapper(opts ...Option) server.SubscriberWrapper {
	return func(h server.SubscriberFunc) server.SubscriberFunc {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}

		ot := &tWrapper{opts: options, serverSubscriber: h}
		return ot.ServerSubscriber
	}
}
