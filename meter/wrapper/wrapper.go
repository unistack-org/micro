package wrapper // import "go.unistack.org/micro/v4/meter/wrapper"

import (
	"context"
	"fmt"
	"time"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/server"
)

var (
	labelSuccess  = "success"
	labelFailure  = "failure"
	labelStatus   = "status"
	labelEndpoint = "endpoint"

	// DefaultSkipEndpoints contains list of endpoints that not evaluted by wrapper
	DefaultSkipEndpoints = []string{"Meter.Metrics", "Health.Live", "Health.Ready", "Health.Version"}
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
	return w.Client.Publish(ctx, p, opts...)
}

// NewServerHandlerWrapper create new server handler wrapper
func NewServerHandlerWrapper(opts ...Option) server.HandlerWrapper {
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
