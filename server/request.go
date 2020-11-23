package server

import (
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/metadata"
)

type rpcMessage struct {
	topic       string
	contentType string
	payload     interface{}
	header      metadata.Metadata
	body        []byte
	codec       codec.Codec
}

func (r *rpcMessage) ContentType() string {
	return r.contentType
}

func (r *rpcMessage) Topic() string {
	return r.topic
}

func (r *rpcMessage) Payload() interface{} {
	return r.payload
}

func (r *rpcMessage) Header() metadata.Metadata {
	return r.header
}

func (r *rpcMessage) Body() []byte {
	return r.body
}

func (r *rpcMessage) Codec() codec.Codec {
	return r.codec
}
