package client

import (
	"context"

	raw "github.com/unistack-org/micro-codec-bytes"
	json "github.com/unistack-org/micro-codec-json"
	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/metadata"
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
	codec       codec.Writer
	stream      bool
}

// NewClient returns new noop client
func NewClient(opts ...Option) Client {
	return &noopClient{opts: NewOptions(opts...)}
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

func (n *noopRequest) Codec() codec.Writer {
	return n.codec
}

func (n *noopRequest) Stream() bool {
	return n.stream
}

type noopResponse struct {
	codec  codec.Reader
	header map[string]string
}

func (n *noopResponse) Codec() codec.Reader {
	return n.codec
}

func (n *noopResponse) Header() map[string]string {
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

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	md["Content-Type"] = p.ContentType()
	md["Micro-Topic"] = p.Topic()

	// passed in raw data
	if d, ok := p.Payload().(*raw.Frame); ok {
		body = d.Data
	} else {
		cf := n.opts.Broker.Options().Codec
		if cf == nil {
			cf = json.Marshaler{}
		}

		/*
			// use codec for payload
			cf, err := n.opts.Codecs[p.ContentType()]
			if err != nil {
				return errors.InternalServerError("go.micro.client", err.Error())
			}
		*/
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

	return nil
}
