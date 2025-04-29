package metadata

import (
	"context"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/server"
)

type wrapper struct {
	keys []string

	client.Client
}

func NewClientWrapper(keys ...string) client.Wrapper {
	return func(c client.Client) client.Client {
		handler := &wrapper{
			Client: c,
			keys:   keys,
		}
		return handler
	}
}

func NewClientCallWrapper(keys ...string) client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		return func(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
			if keys == nil {
				return fn(ctx, addr, req, rsp, opts)
			}
			if imd, iok := metadata.FromIncomingContext(ctx); iok && imd != nil {
				omd, ook := metadata.FromOutgoingContext(ctx)
				if !ook || omd == nil {
					omd = metadata.New(len(imd))
				}
				for _, k := range keys {
					if v := imd.Get(k); v != nil {
						omd.Add(k, v...)
					}
				}
				if !ook {
					ctx = metadata.NewOutgoingContext(ctx, omd)
				}
			}
			return fn(ctx, addr, req, rsp, opts)
		}
	}
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if w.keys == nil {
		return w.Client.Call(ctx, req, rsp, opts...)
	}
	if imd, iok := metadata.FromIncomingContext(ctx); iok && imd != nil {
		omd, ook := metadata.FromOutgoingContext(ctx)
		if !ook || omd == nil {
			omd = metadata.New(len(imd))
		}
		for _, k := range w.keys {
			if v := imd.Get(k); v != nil {
				omd.Add(k, v...)
			}
		}
		if !ook {
			ctx = metadata.NewOutgoingContext(ctx, omd)
		}
	}
	return w.Client.Call(ctx, req, rsp, opts...)
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	if w.keys == nil {
		return w.Client.Stream(ctx, req, opts...)
	}
	if imd, iok := metadata.FromIncomingContext(ctx); iok && imd != nil {
		omd, ook := metadata.FromOutgoingContext(ctx)
		if !ook || omd == nil {
			omd = metadata.New(len(imd))
		}
		for _, k := range w.keys {
			if v := imd.Get(k); v != nil {
				omd.Add(k, v...)
			}
		}
		if !ook {
			ctx = metadata.NewOutgoingContext(ctx, omd)
		}
	}
	return w.Client.Stream(ctx, req, opts...)
}

func NewServerHandlerWrapper(keys ...string) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			if keys == nil {
				return fn(ctx, req, rsp)
			}
			if imd, iok := metadata.FromIncomingContext(ctx); iok && imd != nil {
				omd, ook := metadata.FromOutgoingContext(ctx)
				if !ook || omd == nil {
					omd = metadata.New(len(imd))
				}
				for _, k := range keys {
					if v := imd.Get(k); v != nil {
						omd.Add(k, v...)
					}
				}
				if !ook {
					ctx = metadata.NewOutgoingContext(ctx, omd)
				}
			}
			return fn(ctx, req, rsp)
		}
	}
}
