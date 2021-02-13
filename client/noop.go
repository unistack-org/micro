package client

import (
	"context"

	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/metadata"
)

var (
	// DefaultCodecs will be used to encode/decode data
	DefaultCodecs = map[string]codec.Codec{
		//"application/json":         cjson.NewCodec,
		//"application/json-rpc":     cjsonrpc.NewCodec,
		//"application/protobuf":     cproto.NewCodec,
		//"application/proto-rpc":    cprotorpc.NewCodec,
		"application/octet-stream": codec.NewCodec(),
	}
)

type noopClient struct {
	opts Options
}

type noopMessage struct {
	topic   string
	payload interface{}
	opts    MessageOptions
}

type noopRequest struct {
	service     string
	method      string
	endpoint    string
	contentType string
	body        interface{}
	codec       codec.Codec
	stream      bool
}

// NewClient returns new noop client
func NewClient(opts ...Option) Client {
	nc := &noopClient{opts: NewOptions(opts...)}
	// wrap in reverse

	c := Client(nc)

	for i := len(nc.opts.Wrappers); i > 0; i-- {
		c = nc.opts.Wrappers[i-1](c)
	}

	return c
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

type noopStream struct{}

func (n *noopStream) Context() context.Context {
	return context.Background()
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

func (n *noopStream) Error() error {
	return nil
}

func (n *noopStream) Close() error {
	return nil
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
	return nil
}

func (n *noopClient) Options() Options {
	return n.opts
}

func (n *noopClient) String() string {
	return "noop"
}

func (n *noopClient) Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
	return nil
}

func (n *noopClient) NewRequest(service, endpoint string, req interface{}, opts ...RequestOption) Request {
	return &noopRequest{}
}

func (n *noopClient) NewMessage(topic string, msg interface{}, opts ...MessageOption) Message {
	options := NewMessageOptions(opts...)
	return &noopMessage{topic: topic, payload: msg, opts: options}
}

func (n *noopClient) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	return &noopStream{}, nil
}

func (n *noopClient) Publish(ctx context.Context, p Message, opts ...PublishOption) error {
	var body []byte

	options := NewPublishOptions(opts...)

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(0)
	}
	md["Content-Type"] = p.ContentType()
	md["Micro-Topic"] = p.Topic()

	// passed in raw data
	if d, ok := p.Payload().(*codec.Frame); ok {
		body = d.Data
	} else {
		// use codec for payload
		cf, err := n.newCodec(p.ContentType())
		if err != nil {
			return errors.InternalServerError("go.micro.client", err.Error())
		}

		// set the body
		b, err := cf.Marshal(p.Payload())
		if err != nil {
			return errors.InternalServerError("go.micro.client", err.Error())
		}
		body = b
	}

	topic := p.Topic()

	// get the exchange
	if len(options.Exchange) > 0 {
		topic = options.Exchange
	}

	return n.opts.Broker.Publish(ctx, topic, &broker.Message{
		Header: md,
		Body:   body,
	}, broker.PublishContext(options.Context))
}
