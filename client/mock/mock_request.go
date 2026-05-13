package mock

import (
	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/codec"
)

type mockRequest struct {
	service     string
	method      string
	contentType string
	request     any
	opts        client.RequestOptions
}

func newMockRequest(service, method string, request any, contentType string, opts ...client.RequestOption) client.Request {
	options := client.NewRequestOptions(opts...)
	if len(options.ContentType) == 0 {
		options.ContentType = contentType
	}

	return &mockRequest{
		service:     service,
		method:      method,
		request:     request,
		contentType: options.ContentType,
		opts:        options,
	}
}

func (r *mockRequest) ContentType() string {
	return r.contentType
}

func (r *mockRequest) Service() string {
	return r.service
}

func (r *mockRequest) Method() string {
	return r.method
}

func (r *mockRequest) Endpoint() string {
	return r.method
}

func (r *mockRequest) Codec() codec.Codec {
	return nil
}

func (r *mockRequest) Body() any {
	return r.request
}

func (r *mockRequest) Stream() bool {
	return r.opts.Stream
}