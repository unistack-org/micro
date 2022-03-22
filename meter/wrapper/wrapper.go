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
	// ClientRequestDurationSeconds specifies meter metric name
	ClientRequestDurationSeconds = "client_request_duration_seconds"
	// ClientRequestLatencyMicroseconds specifies meter metric name
	ClientRequestLatencyMicroseconds = "client_request_latency_microseconds"
	// ClientRequestTotal specifies meter metric name
	ClientRequestTotal = "client_request_total"
	// ClientRequestInflight specifies meter metric name
	ClientRequestInflight = "client_request_inflight"
	// ServerRequestDurationSeconds specifies meter metric name
	ServerRequestDurationSeconds = "server_request_duration_seconds"
	// ServerRequestLatencyMicroseconds specifies meter metric name
	ServerRequestLatencyMicroseconds = "server_request_latency_microseconds"
	// ServerRequestTotal specifies meter metric name
	ServerRequestTotal = "server_request_total"
	// ServerRequestInflight specifies meter metric name
	ServerRequestInflight = "server_request_inflight"
	// PublishMessageDurationSeconds specifies meter metric name
	PublishMessageDurationSeconds = "publish_message_duration_seconds"
	// PublishMessageLatencyMicroseconds specifies meter metric name
	PublishMessageLatencyMicroseconds = "publish_message_latency_microseconds"
	// PublishMessageTotal specifies meter metric name
	PublishMessageTotal = "publish_message_total"
	// PublishMessageInflight specifies meter metric name
	PublishMessageInflight = "publish_message_inflight"
	// SubscribeMessageDurationSeconds specifies meter metric name
	SubscribeMessageDurationSeconds = "subscribe_message_duration_seconds"
	// SubscribeMessageLatencyMicroseconds specifies meter metric name
	SubscribeMessageLatencyMicroseconds = "subscribe_message_latency_microseconds"
	// SubscribeMessageTotal specifies meter metric name
	SubscribeMessageTotal = "subscribe_message_total"
	// SubscribeMessageInflight specifies meter metric name
	SubscribeMessageInflight = "subscribe_message_inflight"

	labelSuccess  = "success"
	labelFailure  = "failure"
	labelStatus   = "status"
	labelEndpoint = "endpoint"

	// DefaultSkipEndpoints contains list of endpoints that not evaluted by wrapper
	DefaultSkipEndpoints = []string{"Meter.Metrics"}
)

// Options struct
type Options struct {
	Meter         meter.Meter
	lopts         []meter.Option
	SkipEndpoints []string
}

// Option func signature
type Option func(*Options)

// NewOptions creates new Options struct
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

// ServiceName passes service name to meter label
func ServiceName(name string) Option {
	return func(o *Options) {
		o.lopts = append(o.lopts, meter.Labels("name", name))
	}
}

// ServiceVersion passes service version to meter label
func ServiceVersion(version string) Option {
	return func(o *Options) {
		o.lopts = append(o.lopts, meter.Labels("version", version))
	}
}

// ServiceID passes service id to meter label
func ServiceID(id string) Option {
	return func(o *Options) {
		o.lopts = append(o.lopts, meter.Labels("id", id))
	}
}

// Meter passes meter
func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

// SkipEndoints add endpoint to skip
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

// NewClientWrapper create new client wrapper
func NewClientWrapper(opts ...Option) client.Wrapper {
	return func(c client.Client) client.Client {
		handler := &wrapper{
			opts:   NewOptions(opts...),
			Client: c,
		}
		return handler
	}
}

// NewCallWrapper create new call wrapper
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

// NewHandlerWrapper create new server handler wrapper
func NewHandlerWrapper(opts ...Option) server.HandlerWrapper {
	handler := &wrapper{
		opts: NewOptions(opts...),
	}
	return handler.HandlerFunc
}

func (w *wrapper) HandlerFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		endpoint := req.Service() + "." + req.Endpoint()
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

// NewSubscriberWrapper create server subscribe wrapper
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
