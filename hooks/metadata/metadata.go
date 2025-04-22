package metadata

import (
	"context"

	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/server"
)

var DefaultMetadataKeys = []string{"x-request-id"}

type hook struct {
	keys []string
}

func NewHook(keys ...string) *hook {
	return &hook{keys: keys}
}

func metadataCopy(ctx context.Context, keys []string) context.Context {
	if keys == nil {
		return ctx
	}
	if imd, iok := metadata.FromIncomingContext(ctx); iok && imd != nil {
		omd, ook := metadata.FromOutgoingContext(ctx)
		if !ook || omd == nil {
			omd = metadata.New(len(keys))
		}
		for _, k := range keys {
			if v, ok := imd.Get(k); ok && v != "" {
				omd.Set(k, v)
			}
		}
		if !ook {
			ctx = metadata.NewOutgoingContext(ctx, omd)
		}
	}
	return ctx
}

func (w *hook) ClientCall(next client.FuncCall) client.FuncCall {
	return func(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
		return next(metadataCopy(ctx, w.keys), req, rsp, opts...)
	}
}

func (w *hook) ClientStream(next client.FuncStream) client.FuncStream {
	return func(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
		return next(metadataCopy(ctx, w.keys), req, opts...)
	}
}

func (w *hook) ClientPublish(next client.FuncPublish) client.FuncPublish {
	return func(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
		return next(metadataCopy(ctx, w.keys), msg, opts...)
	}
}

func (w *hook) ClientBatchPublish(next client.FuncBatchPublish) client.FuncBatchPublish {
	return func(ctx context.Context, msgs []client.Message, opts ...client.PublishOption) error {
		return next(metadataCopy(ctx, w.keys), msgs, opts...)
	}
}

func (w *hook) ServerHandler(next server.FuncHandler) server.FuncHandler {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		return next(metadataCopy(ctx, w.keys), req, rsp)
	}
}

func (w *hook) ServerSubscriber(next server.FuncSubHandler) server.FuncSubHandler {
	return func(ctx context.Context, msg server.Message) error {
		return next(metadataCopy(ctx, w.keys), msg)
	}
}
