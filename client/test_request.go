package client

import (
	"go.unistack.org/micro/v4/codec"
)

type testRequest struct {
	codec       codec.Codec
	body        interface{}
	service     string
	method      string
	endpoint    string
	contentType string
	opts        RequestOptions
}

func (r *testRequest) ContentType() string {
	return r.contentType
}

func (r *testRequest) Service() string {
	return r.service
}

func (r *testRequest) Method() string {
	return r.method
}

func (r *testRequest) Endpoint() string {
	return r.endpoint
}

func (r *testRequest) Body() interface{} {
	return r.body
}

func (r *testRequest) Codec() codec.Codec {
	return r.codec
}

func (r *testRequest) Stream() bool {
	return r.opts.Stream
}
