package server

import (
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
)

type rpcMessage struct {
	payload     interface{}
	codec       codec.Codec
	header      metadata.Metadata
	topic       string
	contentType string
	body        []byte
}

func (r *rpcMessage) ContentType() string {
	return r.contentType
}

func (r *rpcMessage) Topic() string {
	return r.topic
}

func (r *rpcMessage) Body() interface{} {
	return r.payload
}

func (r *rpcMessage) Header() metadata.Metadata {
	return r.header
}

func (r *rpcMessage) Codec() codec.Codec {
	return r.codec
}
