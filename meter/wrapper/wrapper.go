// +build ignore

package wrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/server"
)

type Options struct {
	Meter   meter.Meter
	Name    string
	Version string
	ID      string
}

type Option func(*Options)

func ServiceName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func ServiceVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func ServiceID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

type wrapper struct {
	options  Options
	callFunc client.CallFunc
	client.Client
}

func NewClientWrapper(opts ...Option) client.Wrapper {
	return func(c client.Client) client.Client {
		handler := &wrapper{
			labels: labels,
			Client: c,
		}

		return handler
	}
}

func NewCallWrapper(opts ...Option) client.CallWrapper {
	labels := getLabels(opts...)

	return func(fn client.CallFunc) client.CallFunc {
		handler := &wrapper{
			labels:   labels,
			callFunc: fn,
		}

		return handler.CallFunc
	}
}

func (w *wrapper) CallFunc(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	wlabels := append(w.labels, fmt.Sprintf(`%sendpoint="%s"`, DefaultLabelPrefix, endpoint))

	timeCounterSummary := metrics.GetOrCreateSummary(getName("client_request_latency_microseconds", wlabels))
	timeCounterHistogram := metrics.GetOrCreateSummary(getName("client_request_duration_seconds", wlabels))

	ts := time.Now()
	err := w.callFunc(ctx, addr, req, rsp, opts)
	te := time.Since(ts)

	timeCounterSummary.Update(float64(te.Seconds()))
	timeCounterHistogram.Update(te.Seconds())
	if err == nil {
		metrics.GetOrCreateCounter(getName("client_request_total", append(wlabels, fmt.Sprintf(`%sstatus="success"`, DefaultLabelPrefix)))).Inc()
	} else {
		metrics.GetOrCreateCounter(getName("client_request_total", append(wlabels, fmt.Sprintf(`%sstatus="failure"`, DefaultLabelPrefix)))).Inc()
	}

	return err
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	wlabels := append(w.labels, fmt.Sprintf(`%sendpoint="%s"`, DefaultLabelPrefix, endpoint))

	timeCounterSummary := metrics.GetOrCreateSummary(getName("client_request_latency_microseconds", wlabels))
	timeCounterHistogram := metrics.GetOrCreateSummary(getName("client_request_duration_seconds", wlabels))

	ts := time.Now()
	err := w.Client.Call(ctx, req, rsp, opts...)
	te := time.Since(ts)

	timeCounterSummary.Update(float64(te.Seconds()))
	timeCounterHistogram.Update(te.Seconds())
	if err == nil {
		metrics.GetOrCreateCounter(getName("client_request_total", append(wlabels, fmt.Sprintf(`%sstatus="success"`, DefaultLabelPrefix)))).Inc()
	} else {
		metrics.GetOrCreateCounter(getName("client_request_total", append(wlabels, fmt.Sprintf(`%sstatus="failure"`, DefaultLabelPrefix)))).Inc()
	}

	return err
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	wlabels := append(w.labels, fmt.Sprintf(`%sendpoint="%s"`, DefaultLabelPrefix, endpoint))

	timeCounterSummary := metrics.GetOrCreateSummary(getName("client_request_latency_microseconds", wlabels))
	timeCounterHistogram := metrics.GetOrCreateSummary(getName("client_request_duration_seconds", wlabels))

	ts := time.Now()
	stream, err := w.Client.Stream(ctx, req, opts...)
	te := time.Since(ts)

	timeCounterSummary.Update(float64(te.Seconds()))
	timeCounterHistogram.Update(te.Seconds())
	if err == nil {
		metrics.GetOrCreateCounter(getName("client_request_total", append(wlabels, fmt.Sprintf(`%sstatus="success"`, DefaultLabelPrefix)))).Inc()
	} else {
		metrics.GetOrCreateCounter(getName("client_request_total", append(wlabels, fmt.Sprintf(`%sstatus="failure"`, DefaultLabelPrefix)))).Inc()
	}

	return stream, err
}

func (w *wrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	endpoint := p.Topic()
	wlabels := append(w.labels, fmt.Sprintf(`%sendpoint="%s"`, DefaultLabelPrefix, endpoint))

	timeCounterSummary := metrics.GetOrCreateSummary(getName("publish_message_latency_microseconds", wlabels))
	timeCounterHistogram := metrics.GetOrCreateSummary(getName("publish_message_duration_seconds", wlabels))

	ts := time.Now()
	err := w.Client.Publish(ctx, p, opts...)
	te := time.Since(ts)

	timeCounterSummary.Update(float64(te.Seconds()))
	timeCounterHistogram.Update(te.Seconds())
	if err == nil {
		metrics.GetOrCreateCounter(getName("publish_message_total", append(wlabels, fmt.Sprintf(`%sstatus="success"`, DefaultLabelPrefix)))).Inc()
	} else {
		metrics.GetOrCreateCounter(getName("publish_message_total", append(wlabels, fmt.Sprintf(`%sstatus="failure"`, DefaultLabelPrefix)))).Inc()
	}

	return err
}

func NewHandlerWrapper(opts ...Option) server.HandlerWrapper {
	labels := getLabels(opts...)

	handler := &wrapper{
		labels: labels,
	}

	return handler.HandlerFunc
}

func (w *wrapper) HandlerFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		endpoint := req.Endpoint()
		wlabels := append(w.labels, fmt.Sprintf(`%sendpoint="%s"`, DefaultLabelPrefix, endpoint))

		timeCounterSummary := metrics.GetOrCreateSummary(getName("server_request_latency_microseconds", wlabels))
		timeCounterHistogram := metrics.GetOrCreateSummary(getName("server_request_duration_seconds", wlabels))

		ts := time.Now()
		err := fn(ctx, req, rsp)
		te := time.Since(ts)

		timeCounterSummary.Update(float64(te.Seconds()))
		timeCounterHistogram.Update(te.Seconds())
		if err == nil {
			metrics.GetOrCreateCounter(getName("server_request_total", append(wlabels, fmt.Sprintf(`%sstatus="success"`, DefaultLabelPrefix)))).Inc()
		} else {
			metrics.GetOrCreateCounter(getName("server_request_total", append(wlabels, fmt.Sprintf(`%sstatus="failure"`, DefaultLabelPrefix)))).Inc()
		}

		return err
	}
}

func NewSubscriberWrapper(opts ...Option) server.SubscriberWrapper {
	labels := getLabels(opts...)

	handler := &wrapper{
		labels: labels,
	}

	return handler.SubscriberFunc
}

func (w *wrapper) SubscriberFunc(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) error {
		endpoint := msg.Topic()
		wlabels := append(w.labels, fmt.Sprintf(`%sendpoint="%s"`, DefaultLabelPrefix, endpoint))

		timeCounterSummary := metrics.GetOrCreateSummary(getName("subscribe_message_latency_microseconds", wlabels))
		timeCounterHistogram := metrics.GetOrCreateSummary(getName("subscribe_message_duration_seconds", wlabels))

		ts := time.Now()
		err := fn(ctx, msg)
		te := time.Since(ts)

		timeCounterSummary.Update(float64(te.Seconds()))
		timeCounterHistogram.Update(te.Seconds())
		if err == nil {
			metrics.GetOrCreateCounter(getName("subscribe_message_total", append(wlabels, fmt.Sprintf(`%sstatus="success"`, DefaultLabelPrefix)))).Inc()
		} else {
			metrics.GetOrCreateCounter(getName("subscribe_message_total", append(wlabels, fmt.Sprintf(`%sstatus="failure"`, DefaultLabelPrefix)))).Inc()
		}

		return err
	}
}
