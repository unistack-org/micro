package wrapper // import "go.unistack.org/micro/v3/meter/wrapper"

import (
	"context"
	"fmt"
	"time"

	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/server"
)

var (
	ClientRequestDurationSeconds     = "client_request_duration_seconds"
	ClientRequestLatencyMicroseconds = "client_request_latency_microseconds"
	ClientRequestTotal               = "client_request_total"
	ClientRequestInflight            = "client_request_inflight"

	ServerRequestDurationSeconds     = "server_request_duration_seconds"
	ServerRequestLatencyMicroseconds = "server_request_latency_microseconds"
	ServerRequestTotal               = "server_request_total"
	ServerRequestInflight            = "server_request_inflight"

	PublishMessageDurationSeconds     = "publish_message_duration_seconds"
	PublishMessageLatencyMicroseconds = "publish_message_latency_microseconds"
	PublishMessageTotal               = "publish_message_total"
	PublishMessageInflight            = "publish_message_inflight"

	SubscribeMessageDurationSeconds     = "subscribe_message_duration_seconds"
	SubscribeMessageLatencyMicroseconds = "subscribe_message_latency_microseconds"
	SubscribeMessageTotal               = "subscribe_message_total"
	SubscribeMessageInflight            = "subscribe_message_inflight"

	labelSuccess  = "success"
	labelFailure  = "failure"
	labelStatus   = "status"
	labelEndpoint = "endpoint"

	// DefaultSkipEndpoints contains list of endpoints that not evaluted by wrapper
	DefaultSkipEndpoints = []string{"Meter.Metrics"}
)

type Options struct {
	Meter         meter.Meter
	lopts         []meter.Option
	SkipEndpoints []string
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Meter:         meter.DefaultMeter,
		lopts:         make([]meter.Option, 0, 5),
		SkipEndpoints: DefaultSkipEndpoints,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func ServiceName(name string) Option {
	return func(o *Options) {
		o.lopts = append(o.lopts, meter.Labels("name", name))
	}
}

func ServiceVersion(version string) Option {
	return func(o *Options) {
		o.lopts = append(o.lopts, meter.Labels("version", version))
	}
}

func ServiceID(id string) Option {
	return func(o *Options) {
		o.lopts = append(o.lopts, meter.Labels("id", id))
	}
}

func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

func SkipEndoints(eps ...string) Option {
	return func(o *Options) {
		o.SkipEndpoints = append(o.SkipEndpoints, eps...)
	}
}

type wrapper struct {
	client.Client
	callFunc client.CallFunc
	opts     Options
}

func NewClientWrapper(opts ...Option) client.Wrapper {
	return func(c client.Client) client.Client {
		handler := &wrapper{
			opts:   NewOptions(opts...),
			Client: c,
		}
		return handler
	}
}

func NewCallWrapper(opts ...Option) client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		handler := &wrapper{
			opts:     NewOptions(opts...),
			callFunc: fn,
		}
		return handler.CallFunc
	}
}

func (w *wrapper) CallFunc(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range w.opts.SkipEndpoints {
		if ep == endpoint {
			return w.callFunc(ctx, addr, req, rsp, opts)
		}
	}

	labels := make([]string, 0, 4)
	labels = append(labels, labelEndpoint, endpoint)

	w.opts.Meter.Counter(ClientRequestInflight, labels...).Inc()
	ts := time.Now()
	err := w.callFunc(ctx, addr, req, rsp, opts)
	te := time.Since(ts)
	w.opts.Meter.Counter(ClientRequestInflight, labels...).Dec()

	w.opts.Meter.Summary(ClientRequestLatencyMicroseconds, labels...).Update(te.Seconds())
	w.opts.Meter.Histogram(ClientRequestDurationSeconds, labels...).Update(te.Seconds())

	if err == nil {
		labels = append(labels, labelStatus, labelSuccess)
	} else {
		labels = append(labels, labelStatus, labelFailure)
	}
	w.opts.Meter.Counter(ClientRequestTotal, labels...).Inc()

	return err
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range w.opts.SkipEndpoints {
		if ep == endpoint {
			return w.Client.Call(ctx, req, rsp, opts...)
		}
	}

	labels := make([]string, 0, 4)
	labels = append(labels, labelEndpoint, endpoint)

	w.opts.Meter.Counter(ClientRequestInflight, labels...).Inc()
	ts := time.Now()
	err := w.Client.Call(ctx, req, rsp, opts...)
	te := time.Since(ts)
	w.opts.Meter.Counter(ClientRequestInflight, labels...).Dec()

	w.opts.Meter.Summary(ClientRequestLatencyMicroseconds, labels...).Update(te.Seconds())
	w.opts.Meter.Histogram(ClientRequestDurationSeconds, labels...).Update(te.Seconds())

	if err == nil {
		labels = append(labels, labelStatus, labelSuccess)
	} else {
		labels = append(labels, labelStatus, labelFailure)
	}
	w.opts.Meter.Counter(ClientRequestTotal, labels...).Inc()

	return err
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range w.opts.SkipEndpoints {
		if ep == endpoint {
			return w.Client.Stream(ctx, req, opts...)
		}
	}

	labels := make([]string, 0, 4)
	labels = append(labels, labelEndpoint, endpoint)

	w.opts.Meter.Counter(ClientRequestInflight, labels...).Inc()
	ts := time.Now()
	stream, err := w.Client.Stream(ctx, req, opts...)
	te := time.Since(ts)
	w.opts.Meter.Counter(ClientRequestInflight, labels...).Dec()

	w.opts.Meter.Summary(ClientRequestLatencyMicroseconds, labels...).Update(te.Seconds())
	w.opts.Meter.Histogram(ClientRequestDurationSeconds, labels...).Update(te.Seconds())

	if err == nil {
		labels = append(labels, labelStatus, labelSuccess)
	} else {
		labels = append(labels, labelStatus, labelFailure)
	}
	w.opts.Meter.Counter(ClientRequestTotal, labels...).Inc()

	return stream, err
}

func (w *wrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	endpoint := p.Topic()

	labels := make([]string, 0, 4)
	labels = append(labels, labelEndpoint, endpoint)

	w.opts.Meter.Counter(PublishMessageInflight, labels...).Inc()
	ts := time.Now()
	err := w.Client.Publish(ctx, p, opts...)
	te := time.Since(ts)
	w.opts.Meter.Counter(PublishMessageInflight, labels...).Dec()

	w.opts.Meter.Summary(PublishMessageLatencyMicroseconds, labels...).Update(te.Seconds())
	w.opts.Meter.Histogram(PublishMessageDurationSeconds, labels...).Update(te.Seconds())

	if err == nil {
		labels = append(labels, labelStatus, labelSuccess)
	} else {
		labels = append(labels, labelStatus, labelFailure)
	}
	w.opts.Meter.Counter(PublishMessageTotal, labels...).Inc()

	return err
}

func NewHandlerWrapper(opts ...Option) server.HandlerWrapper {
	handler := &wrapper{
		opts: NewOptions(opts...),
	}
	return handler.HandlerFunc
}

func (w *wrapper) HandlerFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		endpoint := req.Endpoint()
		for _, ep := range w.opts.SkipEndpoints {
			if ep == endpoint {
				return fn(ctx, req, rsp)
			}
		}

		labels := make([]string, 0, 4)
		labels = append(labels, labelEndpoint, endpoint)

		w.opts.Meter.Counter(ServerRequestInflight, labels...).Inc()
		ts := time.Now()
		err := fn(ctx, req, rsp)
		te := time.Since(ts)
		w.opts.Meter.Counter(ServerRequestInflight, labels...).Dec()

		w.opts.Meter.Summary(ServerRequestLatencyMicroseconds, labels...).Update(te.Seconds())
		w.opts.Meter.Histogram(ServerRequestDurationSeconds, labels...).Update(te.Seconds())

		if err == nil {
			labels = append(labels, labelStatus, labelSuccess)
		} else {
			labels = append(labels, labelStatus, labelFailure)
		}
		w.opts.Meter.Counter(ServerRequestTotal, labels...).Inc()

		return err
	}
}

func NewSubscriberWrapper(opts ...Option) server.SubscriberWrapper {
	handler := &wrapper{
		opts: NewOptions(opts...),
	}
	return handler.SubscriberFunc
}

func (w *wrapper) SubscriberFunc(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) error {
		endpoint := msg.Topic()

		labels := make([]string, 0, 4)
		labels = append(labels, labelEndpoint, endpoint)

		w.opts.Meter.Counter(SubscribeMessageInflight, labels...).Inc()
		ts := time.Now()
		err := fn(ctx, msg)
		te := time.Since(ts)
		w.opts.Meter.Counter(SubscribeMessageInflight, labels...).Dec()

		w.opts.Meter.Summary(SubscribeMessageLatencyMicroseconds, labels...).Update(te.Seconds())
		w.opts.Meter.Histogram(SubscribeMessageDurationSeconds, labels...).Update(te.Seconds())

		if err == nil {
			labels = append(labels, labelStatus, labelSuccess)
		} else {
			labels = append(labels, labelStatus, labelFailure)
		}
		w.opts.Meter.Counter(SubscribeMessageTotal, labels...).Inc()

		return err
	}
}
