package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/errors"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/options"
	"go.unistack.org/micro/v5/selector"
	"go.unistack.org/micro/v5/semconv"
	"go.unistack.org/micro/v5/tracer"
)

const statusSuccess = "200"

type noopClient struct {
	funcCall   FuncCall
	funcStream FuncStream
	opts       Options
}

type noopRequest struct {
	body        any
	codec       codec.Codec
	service     string
	method      string
	endpoint    string
	contentType string
	stream      bool
}

// NewClient returns new noop client
func NewClient(opts ...Option) Client {
	options := NewOptions(opts...)
	if len(options.Name) == 0 {
		options.Name = "noop"
	}
	n := &noopClient{opts: options}

	n.funcCall = n.fnCall
	n.funcStream = n.fnStream

	return n
}

func (n *noopClient) Name() string {
	return n.opts.Name
}

func (n *noopRequest) Service() string {
	return n.service
}

func (n *noopRequest) Method() string {
	return n.method
}

func (n *noopRequest) Endpoint() string {
	return n.endpoint
}

func (n *noopRequest) ContentType() string {
	return n.contentType
}

func (n *noopRequest) Body() any {
	return n.body
}

func (n *noopRequest) Codec() codec.Codec {
	return n.codec
}

func (n *noopRequest) Stream() bool {
	return n.stream
}

type noopResponse struct {
	codec  codec.Codec
	header metadata.Metadata
}

func (n *noopResponse) Codec() codec.Codec {
	return n.codec
}

func (n *noopResponse) Header() metadata.Metadata {
	return n.header
}

func (n *noopResponse) Read() ([]byte, error) {
	return nil, nil
}

type noopStream struct {
	err error
	ctx context.Context
}

func (n *noopStream) Context() context.Context {
	return n.ctx
}

func (n *noopStream) Request() Request {
	return &noopRequest{}
}

func (n *noopStream) Response() Response {
	return &noopResponse{}
}

func (n *noopStream) Send(any) error {
	return nil
}

func (n *noopStream) Recv(any) error {
	return nil
}

func (n *noopStream) SendMsg(any) error {
	return nil
}

func (n *noopStream) RecvMsg(any) error {
	return nil
}

func (n *noopStream) Error() error {
	return n.err
}

func (n *noopStream) Close() error {
	if sp, ok := tracer.SpanFromContext(n.ctx); ok && sp != nil {
		if n.err != nil {
			sp.SetStatus(tracer.SpanStatusError, n.err.Error())
		}
		sp.Finish()
	}
	return n.err
}

func (n *noopStream) CloseSend() error {
	return n.err
}

func (n *noopClient) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	n.funcCall = n.fnCall
	n.funcStream = n.fnStream

	n.opts.Hooks.EachPrev(func(hook options.Hook) {
		switch h := hook.(type) {
		case HookCall:
			n.funcCall = h(n.funcCall)
		case HookStream:
			n.funcStream = h(n.funcStream)
		}
	})

	return nil
}

func (n *noopClient) Options() Options {
	return n.opts
}

func (n *noopClient) String() string {
	return "noop"
}

func (n *noopClient) Call(ctx context.Context, req Request, rsp any, opts ...CallOption) error {
	endpoint := req.Endpoint()
	ts := time.Now()
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", endpoint).Inc()
	var sp tracer.Span
	ctx, sp = n.opts.Tracer.Start(ctx, "rpc-client",
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels("endpoint", endpoint),
	)
	err := n.funcCall(ctx, req, rsp, opts...)
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", endpoint).Dec()
	te := time.Since(ts)
	n.opts.Meter.Summary(semconv.ClientRequestLatencyMicroseconds, "endpoint", endpoint).Update(te.Seconds())
	n.opts.Meter.Histogram(semconv.ClientRequestDurationSeconds, "endpoint", endpoint).Update(te.Seconds())

	if me := errors.FromError(err); me == nil {
		sp.Finish()
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", endpoint, "status", "success", "code", statusSuccess).Inc()
	} else {
		sp.SetStatus(tracer.SpanStatusError, err.Error())
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", endpoint, "status", "failure", "code", strconv.Itoa(int(me.Code))).Inc()
	}

	return err
}

func (n *noopClient) fnCall(ctx context.Context, req Request, rsp any, opts ...CallOption) error {
	// make a copy of call opts
	callOpts := n.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	// check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		// no deadline so we create a new one
		ctx, cancel = context.WithTimeout(ctx, callOpts.RequestTimeout)
		defer cancel()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		opt := WithRequestTimeout(time.Until(d))
		opt(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return errors.New("go.micro.client", fmt.Sprintf("%v", ctx.Err()), 408)
	default:
	}

	// make copy of call method
	hcall := func(ctx context.Context, addr string, req Request, rsp any, opts CallOptions) error {
		return nil
	}

	// use the router passed as a call option, or fallback to the rpc clients router
	if callOpts.Router == nil {
		callOpts.Router = n.opts.Router
	}

	if callOpts.Selector == nil {
		callOpts.Selector = n.opts.Selector
	}

	// inject proxy address
	// TODO: don't even bother using Lookup/Select in this case
	if len(n.opts.Proxy) > 0 {
		callOpts.Address = []string{n.opts.Proxy}
	}

	var next selector.Next

	// return errors.New("go.micro.client", "request timeout", 408)
	call := func(i int) error {
		// call backoff first. Someone may want an initial start delay
		t, err := callOpts.Backoff(ctx, req, i)
		if err != nil {
			return errors.InternalServerError("go.micro.client", "%s", err)
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		if next == nil {
			var routes []string
			// lookup the route to send the reques to
			// TODO apply any filtering here
			routes, err = n.opts.Lookup(ctx, req, callOpts)
			if err != nil {
				return errors.InternalServerError("go.micro.client", "%s", err)
			}

			// balance the list of nodes
			next, err = callOpts.Selector.Select(routes)
			if err != nil {
				return err
			}
		}

		node := next()

		// make the call
		err = hcall(ctx, node, req, rsp, callOpts)
		// record the result of the call to inform future routing decisions
		if verr := n.opts.Selector.Record(node, err); verr != nil {
			return verr
		}

		// try and transform the error to a go-micro error
		if verr, ok := err.(*errors.Error); ok {
			return verr
		}

		return err
	}

	var gerr error

	for i := 0; i <= callOpts.Retries; i++ {
		err := call(i)
		// if the call succeeded lets bail early
		if err == nil {
			return nil
		}

		retry, rerr := callOpts.Retry(ctx, req, i, err)
		if rerr != nil {
			return rerr
		}

		if !retry {
			return err
		}

		gerr = err
	}

	return gerr
}

func (n *noopClient) NewRequest(service, endpoint string, body any, _ ...RequestOption) Request {
	return &noopRequest{service: service, endpoint: endpoint, body: body, contentType: n.opts.ContentType}
}

func (n *noopClient) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	endpoint := req.Endpoint()
	ts := time.Now()
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", endpoint).Inc()
	var sp tracer.Span
	ctx, sp = n.opts.Tracer.Start(ctx, "rpc-client",
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels("endpoint", endpoint),
	)
	stream, err := n.funcStream(ctx, req, opts...)
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", endpoint).Dec()
	te := time.Since(ts)
	n.opts.Meter.Summary(semconv.ClientRequestLatencyMicroseconds, "endpoint", endpoint).Update(te.Seconds())
	n.opts.Meter.Histogram(semconv.ClientRequestDurationSeconds, "endpoint", endpoint).Update(te.Seconds())

	if me := errors.FromError(err); me == nil {
		sp.Finish()
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", endpoint, "status", "success", "code", statusSuccess).Inc()
	} else {
		sp.SetStatus(tracer.SpanStatusError, err.Error())
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", endpoint, "status", "failure", "code", strconv.Itoa(int(me.Code))).Inc()
	}

	return stream, err
}

func (n *noopClient) fnStream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	var err error

	// make a copy of call opts
	callOpts := n.opts.CallOptions
	for _, o := range opts {
		o(&callOpts)
	}

	// check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok && callOpts.StreamTimeout > time.Duration(0) {
		var cancel context.CancelFunc
		// no deadline so we create a new one
		ctx, cancel = context.WithTimeout(ctx, callOpts.StreamTimeout)
		defer cancel()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		o := WithStreamTimeout(time.Until(d))
		o(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return nil, errors.New("go.micro.client", fmt.Sprintf("%v", ctx.Err()), 408)
	default:
	}

	/*
		// make copy of call method
		hstream := h.stream
		// wrap the call in reverse
		for i := len(callOpts.CallWrappers); i > 0; i-- {
			hstream = callOpts.CallWrappers[i-1](hstream)
		}
	*/

	// use the router passed as a call option, or fallback to the rpc clients router
	if callOpts.Router == nil {
		callOpts.Router = n.opts.Router
	}

	if callOpts.Selector == nil {
		callOpts.Selector = n.opts.Selector
	}

	// inject proxy address
	// TODO: don't even bother using Lookup/Select in this case
	if len(n.opts.Proxy) > 0 {
		callOpts.Address = []string{n.opts.Proxy}
	}

	var next selector.Next

	call := func(i int) (Stream, error) {
		// call backoff first. Someone may want an initial start delay
		t, cerr := callOpts.Backoff(ctx, req, i)
		if cerr != nil {
			return nil, errors.InternalServerError("go.micro.client", "%s", cerr)
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		if next == nil {
			var routes []string
			// lookup the route to send the reques to
			// TODO apply any filtering here
			routes, err = n.opts.Lookup(ctx, req, callOpts)
			if err != nil {
				return nil, errors.InternalServerError("go.micro.client", "%s", err)
			}

			// balance the list of nodes
			next, err = callOpts.Selector.Select(routes)
			if err != nil {
				return nil, err
			}
		}

		node := next()

		stream, cerr := n.stream(ctx, node, req, callOpts)

		// record the result of the call to inform future routing decisions
		if verr := n.opts.Selector.Record(node, cerr); verr != nil {
			return nil, verr
		}

		// try and transform the error to a go-micro error
		if verr, ok := cerr.(*errors.Error); ok {
			return nil, verr
		}

		return stream, cerr
	}

	var grr error

	for i := 0; i <= callOpts.Retries; i++ {
		s, cerr := call(i)
		// if the call succeeded lets bail early
		if cerr == nil {
			return s, nil
		}

		retry, rerr := callOpts.Retry(ctx, req, i, cerr)
		if rerr != nil {
			return nil, rerr
		}

		if !retry {
			return nil, cerr
		}

		grr = cerr
	}

	return nil, grr
}

func (n *noopClient) stream(ctx context.Context, _ string, _ Request, _ CallOptions) (Stream, error) {
	return &noopStream{ctx: ctx}, nil
}
