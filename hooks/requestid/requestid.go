package requestid

import (
	"context"
	"net/textproto"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/server"
	"go.unistack.org/micro/v4/util/id"
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
		var ids []string

		for i := range imd.Get(DefaultMetadataKey) {
			if ids[i] != "" {
				xid = ids[i]
			}
		}

		for i := range omd.Get(DefaultMetadataKey) {
			if ids[i] != "" {
				xid = ids[i]
			}
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

func (w *hook) ServerHandler(next server.FuncHandler) server.FuncHandler {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		var err error
		if ctx, err = DefaultMetadataFunc(ctx); err != nil {
			return err
		}
		return next(ctx, req, rsp)
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
