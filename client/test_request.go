package client

import (
	"github.com/unistack-org/micro/v3/codec"
)

type testRequest struct {
	opts        RequestOptions
	codec       codec.Codec
	body        interface{}
	method      string
	endpoint    string
	contentType string
	service     string
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
