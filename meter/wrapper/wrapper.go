package wrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/server"
)

var (
	ClientRequestDurationSeconds        = "client_request_duration_seconds"
	ClientRequestLatencyMicroseconds    = "client_request_latency_microseconds"
	ClientRequestTotal                  = "client_request_total"
	ServerRequestDurationSeconds        = "server_request_duration_seconds"
	ServerRequestLatencyMicroseconds    = "server_request_latency_microseconds"
	ServerRequestTotal                  = "server_request_total"
	PublishMessageDurationSeconds       = "publish_message_duration_seconds"
	PublishMessageLatencyMicroseconds   = "publish_message_latency_microseconds"
	PublishMessageTotal                 = "publish_message_total"
	SubscribeMessageDurationSeconds     = "subscribe_message_duration_seconds"
	SubscribeMessageLatencyMicroseconds = "subscribe_message_latency_microseconds"
	SubscribeMessageTotal               = "subscribe_message_total"

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
	ts := time.Now()
	err := w.callFunc(ctx, addr, req, rsp, opts)
	te := time.Since(ts)

	lopts := w.opts.lopts
	lopts = append(lopts, meter.Labels(labelEndpoint, endpoint))

	w.opts.Meter.Summary(ClientRequestLatencyMicroseconds, lopts...).Update(float64(te.Seconds()))
	w.opts.Meter.Histogram(ClientRequestDurationSeconds, lopts...).Update(float64(te.Seconds()))

	if err == nil {
		lopts = append(lopts, meter.Labels(labelStatus, labelSuccess))
	} else {
		lopts = append(lopts, meter.Labels(labelStatus, labelFailure))
	}
	w.opts.Meter.Counter(ClientRequestTotal, lopts...).Inc()

	return err
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range w.opts.SkipEndpoints {
		if ep == endpoint {
			return w.Client.Call(ctx, req, rsp, opts...)
		}
	}

	ts := time.Now()
	err := w.Client.Call(ctx, req, rsp, opts...)
	te := time.Since(ts)

	lopts := w.opts.lopts
	lopts = append(lopts, meter.Labels(labelEndpoint, endpoint))

	w.opts.Meter.Summary(ClientRequestLatencyMicroseconds, lopts...).Update(float64(te.Seconds()))
	w.opts.Meter.Histogram(ClientRequestDurationSeconds, lopts...).Update(float64(te.Seconds()))

	if err == nil {
		lopts = append(lopts, meter.Labels(labelStatus, labelSuccess))
	} else {
		lopts = append(lopts, meter.Labels(labelStatus, labelFailure))
	}
	w.opts.Meter.Counter(ClientRequestTotal, lopts...).Inc()

	return err
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	for _, ep := range w.opts.SkipEndpoints {
		if ep == endpoint {
			return w.Client.Stream(ctx, req, opts...)
		}
	}

	ts := time.Now()
	stream, err := w.Client.Stream(ctx, req, opts...)
	te := time.Since(ts)

	lopts := w.opts.lopts
	lopts = append(lopts, meter.Labels(labelEndpoint, endpoint))

	w.opts.Meter.Summary(ClientRequestLatencyMicroseconds, lopts...).Update(float64(te.Seconds()))
	w.opts.Meter.Histogram(ClientRequestDurationSeconds, lopts...).Update(float64(te.Seconds()))

	if err == nil {
		lopts = append(lopts, meter.Labels(labelStatus, labelSuccess))
	} else {
		lopts = append(lopts, meter.Labels(labelStatus, labelFailure))
	}
	w.opts.Meter.Counter(ClientRequestTotal, lopts...).Inc()

	return stream, err
}

func (w *wrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	endpoint := p.Topic()

	ts := time.Now()
	err := w.Client.Publish(ctx, p, opts...)
	te := time.Since(ts)

	lopts := w.opts.lopts
	lopts = append(lopts, meter.Labels(labelEndpoint, endpoint))

	w.opts.Meter.Summary(PublishMessageLatencyMicroseconds, lopts...).Update(float64(te.Seconds()))
	w.opts.Meter.Histogram(PublishMessageDurationSeconds, lopts...).Update(float64(te.Seconds()))

	if err == nil {
		lopts = append(lopts, meter.Labels(labelStatus, labelSuccess))
	} else {
		lopts = append(lopts, meter.Labels(labelStatus, labelFailure))
	}
	w.opts.Meter.Counter(PublishMessageTotal, lopts...).Inc()

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

		ts := time.Now()
		err := fn(ctx, req, rsp)
		te := time.Since(ts)

		lopts := w.opts.lopts
		lopts = append(lopts, meter.Labels(labelEndpoint, endpoint))

		w.opts.Meter.Summary(ServerRequestLatencyMicroseconds, lopts...).Update(float64(te.Seconds()))
		w.opts.Meter.Histogram(ServerRequestDurationSeconds, lopts...).Update(float64(te.Seconds()))

		if err == nil {
			lopts = append(lopts, meter.Labels(labelStatus, labelSuccess))
		} else {
			lopts = append(lopts, meter.Labels(labelStatus, labelFailure))
		}
		w.opts.Meter.Counter(ServerRequestTotal, lopts...).Inc()

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

		ts := time.Now()
		err := fn(ctx, msg)
		te := time.Since(ts)

		lopts := w.opts.lopts
		lopts = append(lopts, meter.Labels(labelEndpoint, endpoint))

		w.opts.Meter.Summary(SubscribeMessageLatencyMicroseconds, lopts...).Update(float64(te.Seconds()))
		w.opts.Meter.Histogram(SubscribeMessageDurationSeconds, lopts...).Update(float64(te.Seconds()))

		if err == nil {
			lopts = append(lopts, meter.Labels(labelStatus, labelSuccess))
		} else {
			lopts = append(lopts, meter.Labels(labelStatus, labelFailure))
		}
		w.opts.Meter.Counter(SubscribeMessageTotal, lopts...).Inc()

		return err
	}
}
