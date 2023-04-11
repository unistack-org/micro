// Package wrapper provides wrapper for Logger
package wrapper

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/server"
)

var (
	// DefaultClientCallObserver called by wrapper in client Call
	DefaultClientCallObserver = func(ctx context.Context, req client.Request, rsp interface{}, opts []client.CallOption, err error) []string {
		labels := []string{"service", req.Service(), "endpoint", req.Endpoint()}
		if err != nil {
			labels = append(labels, "error", err.Error())
		}
		return labels
	}

	// DefaultClientStreamObserver called by wrapper in client Stream
	DefaultClientStreamObserver = func(ctx context.Context, req client.Request, opts []client.CallOption, stream client.Stream, err error) []string {
		labels := []string{"service", req.Service(), "endpoint", req.Endpoint()}
		if err != nil {
			labels = append(labels, "error", err.Error())
		}
		return labels
	}

	// DefaultClientPublishObserver called by wrapper in client Publish
	DefaultClientPublishObserver = func(ctx context.Context, msg client.Message, opts []client.PublishOption, err error) []string {
		labels := []string{"endpoint", msg.Topic()}
		if err != nil {
			labels = append(labels, "error", err.Error())
		}
		return labels
	}

	// DefaultServerHandlerObserver called by wrapper in server Handler
	DefaultServerHandlerObserver = func(ctx context.Context, req server.Request, rsp interface{}, err error) []string {
		labels := []string{"service", req.Service(), "endpoint", req.Endpoint()}
		if err != nil {
			labels = append(labels, "error", err.Error())
		}
		return labels
	}

	// DefaultServerSubscriberObserver called by wrapper in server Subscriber
	DefaultServerSubscriberObserver = func(ctx context.Context, msg server.Message, err error) []string {
		labels := []string{"endpoint", msg.Topic()}
		if err != nil {
			labels = append(labels, "error", err.Error())
		}
		return labels
	}

	// DefaultClientCallFuncObserver called by wrapper in client CallFunc
	DefaultClientCallFuncObserver = func(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions, err error) []string {
		labels := []string{"service", req.Service(), "endpoint", req.Endpoint()}
		if err != nil {
			labels = append(labels, "error", err.Error())
		}
		return labels
	}

	// DefaultSkipEndpoints wrapper not called for this endpoints
	DefaultSkipEndpoints = []string{"Meter.Metrics", "Health.Live", "Health.Ready", "Health.Version"}
)

type lWrapper struct {
	client.Client
	serverHandler    server.HandlerFunc
	serverSubscriber server.SubscriberFunc
	clientCallFunc   client.CallFunc
	opts             Options
}

type (
	// ClientCallObserver func signature
	ClientCallObserver func(context.Context, client.Request, interface{}, []client.CallOption, error) []string
	// ClientStreamObserver func signature
	ClientStreamObserver func(context.Context, client.Request, []client.CallOption, client.Stream, error) []string
	// ClientPublishObserver func signature
	ClientPublishObserver func(context.Context, client.Message, []client.PublishOption, error) []string
	// ClientCallFuncObserver func signature
	ClientCallFuncObserver func(context.Context, string, client.Request, interface{}, client.CallOptions, error) []string
	// ServerHandlerObserver func signature
	ServerHandlerObserver func(context.Context, server.Request, interface{}, error) []string
	// ServerSubscriberObserver func signature
	ServerSubscriberObserver func(context.Context, server.Message, error) []string
)

// Options struct for wrapper
type Options struct {
	// Logger that used for log
	Logger logger.Logger
	// ServerHandlerObservers funcs
	ServerHandlerObservers []ServerHandlerObserver
	// ServerSubscriberObservers funcs
	ServerSubscriberObservers []ServerSubscriberObserver
	// ClientCallObservers funcs
	ClientCallObservers []ClientCallObserver
	// ClientStreamObservers funcs
	ClientStreamObservers []ClientStreamObserver
	// ClientPublishObservers funcs
	ClientPublishObservers []ClientPublishObserver
	// ClientCallFuncObservers funcs
	ClientCallFuncObservers []ClientCallFuncObserver
	// SkipEndpoints
	SkipEndpoints []string
	// Level for logger
	Level logger.Level
	// Enabled flag
	Enabled bool
}

// Option func signature
type Option func(*Options)

// NewOptions creates Options from Option slice
func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:                    logger.DefaultLogger,
		Level:                     logger.TraceLevel,
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

// WithEnabled enable/diable flag
func WithEnabled(b bool) Option {
	return func(o *Options) {
		o.Enabled = b
	}
}

// WithLevel log level
func WithLevel(l logger.Level) Option {
	return func(o *Options) {
		o.Level = l
	}
}

// WithLogger logger
func WithLogger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
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

// SkipEndpoins
func SkipEndpoints(eps ...string) Option {
	return func(o *Options) {
		o.SkipEndpoints = append(o.SkipEndpoints, eps...)
	}
}

func (l *lWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	err := l.Client.Call(ctx, req, rsp, opts...)

	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range l.opts.SkipEndpoints {
		if ep == endpoint {
			return err
		}
	}

	if !l.opts.Enabled {
		return err
	}

	var labels []string
	for _, o := range l.opts.ClientCallObservers {
		labels = append(labels, o(ctx, req, rsp, opts, err)...)
	}
	l.opts.Logger.Fields(labels).Log(ctx, l.opts.Level)

	return err
}

func (l *lWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	stream, err := l.Client.Stream(ctx, req, opts...)

	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range l.opts.SkipEndpoints {
		if ep == endpoint {
			return stream, err
		}
	}

	if !l.opts.Enabled {
		return stream, err
	}

	var labels []string
	for _, o := range l.opts.ClientStreamObservers {
		labels = append(labels, o(ctx, req, opts, stream, err)...)
	}
	l.opts.Logger.Fields(labels).Log(ctx, l.opts.Level)

	return stream, err
}

func (l *lWrapper) Publish(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
	err := l.Client.Publish(ctx, msg, opts...)

	endpoint := msg.Topic()
	for _, ep := range l.opts.SkipEndpoints {
		if ep == endpoint {
			return err
		}
	}

	if !l.opts.Enabled {
		return err
	}

	var labels []string
	for _, o := range l.opts.ClientPublishObservers {
		labels = append(labels, o(ctx, msg, opts, err)...)
	}
	l.opts.Logger.Fields(labels).Log(ctx, l.opts.Level)

	return err
}

func (l *lWrapper) ServerHandler(ctx context.Context, req server.Request, rsp interface{}) error {
	err := l.serverHandler(ctx, req, rsp)

	endpoint := req.Endpoint()
	for _, ep := range l.opts.SkipEndpoints {
		if ep == endpoint {
			return err
		}
	}

	if !l.opts.Enabled {
		return err
	}

	var labels []string
	for _, o := range l.opts.ServerHandlerObservers {
		labels = append(labels, o(ctx, req, rsp, err)...)
	}
	l.opts.Logger.Fields(labels).Log(ctx, l.opts.Level)

	return err
}

func (l *lWrapper) ServerSubscriber(ctx context.Context, msg server.Message) error {
	err := l.serverSubscriber(ctx, msg)

	endpoint := msg.Topic()
	for _, ep := range l.opts.SkipEndpoints {
		if ep == endpoint {
			return err
		}
	}

	if !l.opts.Enabled {
		return err
	}

	var labels []string
	for _, o := range l.opts.ServerSubscriberObservers {
		labels = append(labels, o(ctx, msg, err)...)
	}
	l.opts.Logger.Fields(labels).Log(ctx, l.opts.Level)

	return err
}

// NewClientWrapper accepts an open options and returns a Client Wrapper
func NewClientWrapper(opts ...Option) client.Wrapper {
	return func(c client.Client) client.Client {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}
		return &lWrapper{opts: options, Client: c}
	}
}

// NewClientCallWrapper accepts an options and returns a Call Wrapper
func NewClientCallWrapper(opts ...Option) client.CallWrapper {
	return func(h client.CallFunc) client.CallFunc {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}

		l := &lWrapper{opts: options, clientCallFunc: h}
		return l.ClientCallFunc
	}
}

func (l *lWrapper) ClientCallFunc(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
	err := l.clientCallFunc(ctx, addr, req, rsp, opts)

	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range l.opts.SkipEndpoints {
		if ep == endpoint {
			return err
		}
	}

	if !l.opts.Enabled {
		return err
	}

	var labels []string
	for _, o := range l.opts.ClientCallFuncObservers {
		labels = append(labels, o(ctx, addr, req, rsp, opts, err)...)
	}
	l.opts.Logger.Fields(labels).Log(ctx, l.opts.Level)

	return err
}

// NewServerHandlerWrapper accepts an options and returns a Handler Wrapper
func NewServerHandlerWrapper(opts ...Option) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}

		l := &lWrapper{opts: options, serverHandler: h}
		return l.ServerHandler
	}
}

// NewServerSubscriberWrapper accepts an options and returns a Subscriber Wrapper
func NewServerSubscriberWrapper(opts ...Option) server.SubscriberWrapper {
	return func(h server.SubscriberFunc) server.SubscriberFunc {
		options := NewOptions()
		for _, o := range opts {
			o(&options)
		}

		l := &lWrapper{opts: options, serverSubscriber: h}
		return l.ServerSubscriber
	}
}
