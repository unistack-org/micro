package client

import (
	"context"

	"github.com/unistack-org/micro/v3/codec"
)

type noopClient struct {
	opts Options
}

type noopMessage struct {
	topic       string
	payload     interface{}
	contentType string
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
	return n.contentType
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
	return &noopMessage{}
}

func (n *noopClient) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	return &noopStream{}, nil
}

func (n *noopClient) Publish(ctx context.Context, msg Message, opts ...PublishOption) error {
	return nil
}

func newClient(opts ...Option) Client {
	options := NewOptions()
	for _, o := range opts {
		o(&options)
	}
	return &noopClient{opts: options}
}
