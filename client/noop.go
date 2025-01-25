package client

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/errors"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
	"go.unistack.org/micro/v4/selector"
	"go.unistack.org/micro/v4/semconv"
	"go.unistack.org/micro/v4/tracer"
)

// DefaultCodecs will be used to encode/decode data
var DefaultCodecs = map[string]codec.Codec{
	"application/octet-stream": codec.NewCodec(),
}

type noopClient struct {
	funcPublish      FuncPublish
	funcBatchPublish FuncBatchPublish
	funcCall         FuncCall
	funcStream       FuncStream
	opts             Options
}

type noopMessage struct {
	topic   string
	payload interface{}
	opts    MessageOptions
}

type noopRequest struct {
	body        interface{}
	codec       codec.Codec
	service     string
	method      string
	endpoint    string
	contentType string
	stream      bool
}

// NewClient returns new noop client
func NewClient(opts ...Option) Client {
	n := &noopClient{opts: NewOptions(opts...)}

	n.funcCall = n.fnCall
	n.funcStream = n.fnStream
	n.funcPublish = n.fnPublish
	n.funcBatchPublish = n.fnBatchPublish

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

func (n *noopRequest) Body() interface{} {
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

func (n *noopStream) Send(interface{}) error {
	return nil
}

func (n *noopStream) Recv(interface{}) error {
	return nil
}

func (n *noopStream) SendMsg(interface{}) error {
	return nil
}

func (n *noopStream) RecvMsg(interface{}) error {
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

func (n *noopMessage) Topic() string {
	return n.topic
}

func (n *noopMessage) Payload() interface{} {
	return n.payload
}

func (n *noopMessage) ContentType() string {
	return n.opts.ContentType
}

func (n *noopMessage) Metadata() metadata.Metadata {
	return n.opts.Metadata
}

func (n *noopClient) newCodec(contentType string) (codec.Codec, error) {
	if cf, ok := n.opts.Codecs[contentType]; ok {
		return cf, nil
	}
	if cf, ok := DefaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, codec.ErrUnknownContentType
}

func (n *noopClient) Init(opts ...Option) error {
	for _, o := range opts {
		o(&n.opts)
	}

	n.funcCall = n.fnCall
	n.funcStream = n.fnStream
	n.funcPublish = n.fnPublish
	n.funcBatchPublish = n.fnBatchPublish

	n.opts.Hooks.EachPrev(func(hook options.Hook) {
		switch h := hook.(type) {
		case HookCall:
			n.funcCall = h(n.funcCall)
		case HookStream:
			n.funcStream = h(n.funcStream)
		case HookPublish:
			n.funcPublish = h(n.funcPublish)
		case HookBatchPublish:
			n.funcBatchPublish = h(n.funcBatchPublish)
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

func (n *noopClient) Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
	ts := time.Now()
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", req.Endpoint()).Inc()
	var sp tracer.Span
	ctx, sp = n.opts.Tracer.Start(ctx, "rpc-client",
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels("endpoint", req.Endpoint()),
	)
	err := n.funcCall(ctx, req, rsp, opts...)
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", req.Endpoint()).Dec()
	te := time.Since(ts)
	n.opts.Meter.Summary(semconv.ClientRequestLatencyMicroseconds, "endpoint", req.Endpoint()).Update(te.Seconds())
	n.opts.Meter.Histogram(semconv.ClientRequestDurationSeconds, "endpoint", req.Endpoint()).Update(te.Seconds())

	if me := errors.FromError(err); me == nil {
		sp.Finish()
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", req.Endpoint(), "status", "success", "code", strconv.Itoa(int(200))).Inc()
	} else {
		sp.SetStatus(tracer.SpanStatusError, err.Error())
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", req.Endpoint(), "status", "failure", "code", strconv.Itoa(int(me.Code))).Inc()
	}

	return err
}

func (n *noopClient) fnCall(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
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
	hcall := func(ctx context.Context, addr string, req Request, rsp interface{}, opts CallOptions) error {
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

	ch := make(chan error, callOpts.Retries)
	var gerr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func() {
			ch <- call(i)
		}()

		select {
		case <-ctx.Done():
			return errors.New("go.micro.client", fmt.Sprintf("%v", ctx.Err()), 408)
		case err := <-ch:
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
	}

	return gerr
}

func (n *noopClient) NewRequest(service, endpoint string, _ interface{}, _ ...RequestOption) Request {
	return &noopRequest{service: service, endpoint: endpoint}
}

func (n *noopClient) NewMessage(topic string, msg interface{}, opts ...MessageOption) Message {
	options := NewMessageOptions(append([]MessageOption{MessageContentType(n.opts.ContentType)}, opts...)...)
	return &noopMessage{topic: topic, payload: msg, opts: options}
}

func (n *noopClient) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	ts := time.Now()
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", req.Endpoint()).Inc()
	var sp tracer.Span
	ctx, sp = n.opts.Tracer.Start(ctx, "rpc-client",
		tracer.WithSpanKind(tracer.SpanKindClient),
		tracer.WithSpanLabels("endpoint", req.Endpoint()),
	)
	stream, err := n.funcStream(ctx, req, opts...)
	n.opts.Meter.Counter(semconv.ClientRequestInflight, "endpoint", req.Endpoint()).Dec()
	te := time.Since(ts)
	n.opts.Meter.Summary(semconv.ClientRequestLatencyMicroseconds, "endpoint", req.Endpoint()).Update(te.Seconds())
	n.opts.Meter.Histogram(semconv.ClientRequestDurationSeconds, "endpoint", req.Endpoint()).Update(te.Seconds())

	if me := errors.FromError(err); me == nil {
		sp.Finish()
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", req.Endpoint(), "status", "success", "code", strconv.Itoa(int(200))).Inc()
	} else {
		sp.SetStatus(tracer.SpanStatusError, err.Error())
		n.opts.Meter.Counter(semconv.ClientRequestTotal, "endpoint", req.Endpoint(), "status", "failure", "code", strconv.Itoa(int(me.Code))).Inc()
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

	type response struct {
		stream Stream
		err    error
	}

	ch := make(chan response, callOpts.Retries)
	var grr error

	for i := 0; i <= callOpts.Retries; i++ {
		go func() {
			s, cerr := call(i)
			ch <- response{s, cerr}
		}()

		select {
		case <-ctx.Done():
			return nil, errors.New("go.micro.client", fmt.Sprintf("%v", ctx.Err()), 408)
		case rsp := <-ch:
			// if the call succeeded lets bail early
			if rsp.err == nil {
				return rsp.stream, nil
			}

			retry, rerr := callOpts.Retry(ctx, req, i, err)
			if rerr != nil {
				return nil, rerr
			}

			if !retry {
				return nil, rsp.err
			}

			grr = rsp.err
		}
	}

	return nil, grr
}

func (n *noopClient) stream(ctx context.Context, _ string, _ Request, _ CallOptions) (Stream, error) {
	return &noopStream{ctx: ctx}, nil
}

func (n *noopClient) BatchPublish(ctx context.Context, ps []Message, opts ...PublishOption) error {
	return n.funcBatchPublish(ctx, ps, opts...)
}

func (n *noopClient) fnBatchPublish(ctx context.Context, ps []Message, opts ...PublishOption) error {
	return n.publish(ctx, ps, opts...)
}

func (n *noopClient) Publish(ctx context.Context, p Message, opts ...PublishOption) error {
	return n.funcPublish(ctx, p, opts...)
}

func (n *noopClient) fnPublish(ctx context.Context, p Message, opts ...PublishOption) error {
	return n.publish(ctx, []Message{p}, opts...)
}

func (n *noopClient) publish(ctx context.Context, ps []Message, opts ...PublishOption) error {
	options := NewPublishOptions(opts...)

	msgs := make([]*broker.Message, 0, len(ps))

	// get proxy
	exchange := ""
	if v, ok := os.LookupEnv("MICRO_PROXY"); ok {
		exchange = v
	}
	// get the exchange
	if len(options.Exchange) > 0 {
		exchange = options.Exchange
	}

	omd, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		omd = metadata.New(0)
	}

	for _, p := range ps {
		md := metadata.Copy(omd)
		topic := p.Topic()
		if len(exchange) > 0 {
			topic = exchange
		}
		md[metadata.HeaderTopic] = topic
		iter := p.Metadata().Iterator()
		var k, v string
		for iter.Next(&k, &v) {
			md.Set(k, v)
		}

		md[metadata.HeaderContentType] = p.ContentType()

		var body []byte

		// passed in raw data
		if d, ok := p.Payload().(*codec.Frame); ok {
			body = d.Data
		} else {
			// use codec for payload
			cf, err := n.newCodec(p.ContentType())
			if err != nil {
				return errors.InternalServerError("go.micro.client", "%s", err)
			}

			// set the body
			b, err := cf.Marshal(p.Payload())
			if err != nil {
				return errors.InternalServerError("go.micro.client", "%s", err)
			}
			body = b
		}

		msgs = append(msgs, &broker.Message{Header: md, Body: body})
	}

	if len(msgs) == 1 {
		return n.opts.Broker.Publish(ctx, msgs[0].Header[metadata.HeaderTopic], msgs[0],
			broker.PublishContext(options.Context),
			broker.PublishBodyOnly(options.BodyOnly),
		)
	}

	return n.opts.Broker.BatchPublish(ctx, msgs,
		broker.PublishContext(options.Context),
		broker.PublishBodyOnly(options.BodyOnly),
	)
}
