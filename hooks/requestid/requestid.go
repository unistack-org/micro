package requestid

import (
	"context"
	"net/textproto"
	
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/metadata"
	"go.unistack.org/micro/v3/server"
	"go.unistack.org/micro/v3/util/id"
)

type XRequestIDKey struct{}

// DefaultMetadataKey contains metadata key
var DefaultMetadataKey = textproto.CanonicalMIMEHeaderKey("x-request-id")

// DefaultMetadataFunc wil be used if user not provide own func to fill metadata
var DefaultMetadataFunc = func(ctx context.Context) (context.Context, error) {
	var xid string

	cid, cok := ctx.Value(XRequestIDKey{}).(string)
	if cok && cid != "" {
		xid = cid
	}

	imd, iok := metadata.FromIncomingContext(ctx)
	if !iok || imd == nil {
		imd = metadata.New(1)
		ctx = metadata.NewIncomingContext(ctx, imd)
	}

	omd, ook := metadata.FromOutgoingContext(ctx)
	if !ook || omd == nil {
		omd = metadata.New(1)
		ctx = metadata.NewOutgoingContext(ctx, omd)
	}

	if xid == "" {
		var id string
		if id, iok = imd.Get(DefaultMetadataKey); iok && id != "" {
			xid = id
		}
		if id, ook = omd.Get(DefaultMetadataKey); ook && id != "" {
			xid = id
		}
	}

	if xid == "" {
		var err error
		xid, err = id.New()
		if err != nil {
			return ctx, err
		}
	}

	if !cok {
		ctx = context.WithValue(ctx, XRequestIDKey{}, xid)
	}

	if !iok {
		imd.Set(DefaultMetadataKey, xid)
	}

	if !ook {
		omd.Set(DefaultMetadataKey, xid)
	}

	return ctx, nil
}

type hook struct{}

func NewHook() *hook {
	return &hook{}
}

func (w *hook) ServerSubscriber(next server.FuncSubHandler) server.FuncSubHandler {
	return func(ctx context.Context, msg server.Message) error {
		var err error
		if xid, ok := msg.Header()[DefaultMetadataKey]; ok {
			ctx = context.WithValue(ctx, XRequestIDKey{}, xid)
		}
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return err
		}
		return next(ctx, msg)
	}
}

func (w *hook) ServerHandler(next server.FuncHandler) server.FuncHandler {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		var err error
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return err
		}
		return next(ctx, req, rsp)
	}
}

func (w *hook) ClientBatchPublish(next client.FuncBatchPublish) client.FuncBatchPublish {
	return func(ctx context.Context, msgs []client.Message, opts ...client.PublishOption) error {
		var err error
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return err
		}
		return next(ctx, msgs, opts...)
	}
}

func (w *hook) ClientPublish(next client.FuncPublish) client.FuncPublish {
	return func(ctx context.Context, msg client.Message, opts ...client.PublishOption) error {
		var err error
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return err
		}
		return next(ctx, msg, opts...)
	}
}

func (w *hook) ClientCall(next client.FuncCall) client.FuncCall {
	return func(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
		var err error
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return err
		}
		return next(ctx, req, rsp, opts...)
	}
}

func (w *hook) ClientStream(next client.FuncStream) client.FuncStream {
	return func(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
		var err error
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return nil, err
		}
		return next(ctx, req, opts...)
	}
}
