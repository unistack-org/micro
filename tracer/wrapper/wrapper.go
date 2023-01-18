// Package wrapper provides wrapper for Tracer
package wrapper // import "go.unistack.org/micro/v3/tracer/wrapper"

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/server"
	"go.unistack.org/micro/v3/tracer"
)

var (
	DefaultClientCallObserver = func(ctx context.Context, req client.Request, rsp interface{}, opts []client.CallOption, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("Call %s.%s", req.Service(), req.Method()))
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = make([]interface{}, 0, len(md)+1)
			for k, v := range md {
				labels = append(labels, k, v)
			}
		}
		if err != nil {
			labels = append(labels, "error", err.Error())
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		labels = append(labels, "kind", sp.Kind())
		sp.SetLabels(labels...)
	}

	DefaultClientStreamObserver = func(ctx context.Context, req client.Request, opts []client.CallOption, stream client.Stream, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("Stream %s.%s", req.Service(), req.Method()))
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = make([]interface{}, 0, len(md))
			for k, v := range md {
				labels = append(labels, k, v)
			}
		}
		if err != nil {
			labels = append(labels, "error", err.Error())
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		labels = append(labels, "kind", sp.Kind())
		sp.SetLabels(labels...)
	}

	DefaultClientPublishObserver = func(ctx context.Context, msg client.Message, opts []client.PublishOption, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("Publish %s", msg.Topic()))
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = make([]interface{}, 0, len(md))
			for k, v := range md {
				labels = append(labels, k, v)
			}
		}
		if err != nil {
			labels = append(labels, "error", err.Error())
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		labels = append(labels, "kind", sp.Kind())
		sp.SetLabels(labels...)
	}

	DefaultServerHandlerObserver = func(ctx context.Context, req server.Request, rsp interface{}, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("Handler %s.%s", req.Service(), req.Method()))
		var labels []interface{}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			labels = make([]interface{}, 0, len(md))
			for k, v := range md {
				labels = append(labels, k, v)
			}
		}
		if err != nil {
			labels = append(labels, "error", err.Error())
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		labels = append(labels, "kind", sp.Kind())
		sp.SetLabels(labels...)
	}

	DefaultServerSubscriberObserver = func(ctx context.Context, msg server.Message, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("Subscriber %s", msg.Topic()))
		var labels []interface{}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			labels = make([]interface{}, 0, len(md))
			for k, v := range md {
				labels = append(labels, k, v)
			}
		}
		if err != nil {
			labels = append(labels, "error", err.Error())
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		labels = append(labels, "kind", sp.Kind())
		sp.SetLabels(labels...)
	}

	DefaultClientCallFuncObserver = func(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions, sp tracer.Span, err error) {
		sp.SetName(fmt.Sprintf("Call %s.%s", req.Service(), req.Method()))
		var labels []interface{}
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			labels = make([]interface{}, 0, len(md))
			for k, v := range md {
				labels = append(labels, k, v)
			}
		}
		if err != nil {
			labels = append(labels, "error", err.Error())
			sp.SetStatus(tracer.SpanStatusError, err.Error())
		}
		labels = append(labels, "kind", sp.Kind())
		sp.SetLabels(labels...)
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

	sp, ok := tracer.SpanFromContext(ctx)
	if !ok {
		ctx, sp = ot.opts.Tracer.Start(ctx, "", tracer.WithSpanKind(tracer.SpanKindClient))
	}
	defer sp.Finish()

	err := ot.Client.Call(ctx, req, rsp, opts...)

	for _, o := range ot.opts.ClientCallObservers {
		o(ctx, req, rsp, opts, sp, err)
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

	sp, ok := tracer.SpanFromContext(ctx)
	if !ok {
		ctx, sp = ot.opts.Tracer.Start(ctx, "", tracer.WithSpanKind(tracer.SpanKindClient))
	}
	defer sp.Finish()

	stream, err := ot.Client.Stream(ctx, req, opts...)

	for _, o := range ot.opts.ClientStreamObservers {
		o(ctx, req, opts, stream, sp, err)
	}

	return stream, err
}

func (ot *tWrapper) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	sp, ok := tracer.SpanFromContext(ctx)
	if !ok {
		ctx, sp = ot.opts.Tracer.Start(ctx, "", tracer.WithSpanKind(tracer.SpanKindProducer))
	}
	defer sp.Finish()

	err := ot.Client.Publish(ctx, msg, opts...)

	for _, o := range ot.opts.ClientPublishObservers {
		o(ctx, msg, opts, sp, err)
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

	sp, ok := tracer.SpanFromContext(ctx)
	if !ok {
		ctx, sp = ot.opts.Tracer.Start(ctx, "", tracer.WithSpanKind(tracer.SpanKindServer))
	}
	defer sp.Finish()

	err := ot.serverHandler(ctx, req, rsp)

	for _, o := range ot.opts.ServerHandlerObservers {
		o(ctx, req, rsp, sp, err)
	}

	return err
}

func (ot *tWrapper) ServerSubscriber(ctx context.Context, msg server.Message) error {
	sp, ok := tracer.SpanFromContext(ctx)
	if !ok {
		ctx, sp = ot.opts.Tracer.Start(ctx, "", tracer.WithSpanKind(tracer.SpanKindConsumer))
	}
	defer sp.Finish()

	err := ot.serverSubscriber(ctx, msg)

	for _, o := range ot.opts.ServerSubscriberObservers {
		o(ctx, msg, sp, err)
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

	sp, ok := tracer.SpanFromContext(ctx)
	if !ok {
		ctx, sp = ot.opts.Tracer.Start(ctx, "", tracer.WithSpanKind(tracer.SpanKindClient))
	}
	defer sp.Finish()

	err := ot.clientCallFunc(ctx, addr, req, rsp, opts)

	for _, o := range ot.opts.ClientCallFuncObservers {
		o(ctx, addr, req, rsp, opts, sp, err)
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
