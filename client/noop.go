package client

import (
	"context"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/metadata"
)

// DefaultCodecs will be used to encode/decode data
var DefaultCodecs = map[string]codec.Codec{
	"application/octet-stream": codec.NewCodec(),
}

type noopClient struct {
	opts Options
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
	return &noopRequest{service: service, endpoint: endpoint}
}

func (n *noopClient) NewMessage(topic string, msg interface{}, opts ...MessageOption) Message {
	options := NewMessageOptions(append([]MessageOption{MessageContentType(n.opts.ContentType)}, opts...)...)
	return &noopMessage{topic: topic, payload: msg, opts: options}
}

func (n *noopClient) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	return &noopStream{}, nil
}

func (n *noopClient) BatchPublish(ctx context.Context, ps []Message, opts ...PublishOption) error {
	return n.publish(ctx, ps, opts...)
}

func (n *noopClient) Publish(ctx context.Context, p Message, opts ...PublishOption) error {
	return n.publish(ctx, []Message{p}, opts...)
}

func (n *noopClient) publish(ctx context.Context, ps []Message, opts ...PublishOption) error {
	options := NewPublishOptions(opts...)

	msgs := make([]*broker.Message, 0, len(ps))

	for _, p := range ps {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(0)
		}
		md[metadata.HeaderContentType] = p.ContentType()

		topic := p.Topic()

		// get the exchange
		if len(options.Exchange) > 0 {
			topic = options.Exchange
		}

		md[metadata.HeaderTopic] = topic

		var body []byte

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

		msgs = append(msgs, &broker.Message{Header: md, Body: body})
	}

	return n.opts.Broker.BatchPublish(ctx, msgs,
		broker.PublishContext(options.Context),
		broker.PublishBodyOnly(options.BodyOnly),
	)
}
