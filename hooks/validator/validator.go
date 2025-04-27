package validator

import (
	"context"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/errors"
	"go.unistack.org/micro/v4/server"
)

var (
	DefaultClientErrorFunc = func(req client.Request, rsp interface{}, err error) error {
		if rsp != nil {
			return errors.BadGateway(req.Service(), "%v", err)
		}
		return errors.BadRequest(req.Service(), "%v", err)
	}

	DefaultServerErrorFunc = func(req server.Request, rsp interface{}, err error) error {
		if rsp != nil {
			return errors.BadGateway(req.Service(), "%v", err)
		}
		return errors.BadRequest(req.Service(), "%v", err)
	}
)

type (
	ClientErrorFunc func(client.Request, interface{}, error) error
	ServerErrorFunc func(server.Request, interface{}, error) error
)

// Options struct holds wrapper options
type Options struct {
	ClientErrorFn          ClientErrorFunc
	ServerErrorFn          ServerErrorFunc
	ClientValidateResponse bool
	ServerValidateResponse bool
}

// Option func signature
type Option func(*Options)

func ClientValidateResponse(b bool) Option {
	return func(o *Options) {
		o.ClientValidateResponse = b
	}
}

func ServerValidateResponse(b bool) Option {
	return func(o *Options) {
		o.ClientValidateResponse = b
	}
}

func ClientReqErrorFn(fn ClientErrorFunc) Option {
	return func(o *Options) {
		o.ClientErrorFn = fn
	}
}

func ServerErrorFn(fn ServerErrorFunc) Option {
	return func(o *Options) {
		o.ServerErrorFn = fn
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{
		ClientErrorFn: DefaultClientErrorFunc,
		ServerErrorFn: DefaultServerErrorFunc,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func NewHook(opts ...Option) *hook {
	return &hook{opts: NewOptions(opts...)}
}

type validator interface {
	Validate() error
}

type hook struct {
	opts Options
}

func (w *hook) ClientCall(next client.FuncCall) client.FuncCall {
	return func(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
		if v, ok := req.Body().(validator); ok {
			if err := v.Validate(); err != nil {
				return w.opts.ClientErrorFn(req, nil, err)
			}
		}
		err := next(ctx, req, rsp, opts...)
		if v, ok := rsp.(validator); ok && w.opts.ClientValidateResponse {
			if verr := v.Validate(); verr != nil {
				return w.opts.ClientErrorFn(req, rsp, verr)
			}
		}
		return err
	}
}

func (w *hook) ClientStream(next client.FuncStream) client.FuncStream {
	return func(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
		if v, ok := req.Body().(validator); ok {
			if err := v.Validate(); err != nil {
				return nil, w.opts.ClientErrorFn(req, nil, err)
			}
		}
		return next(ctx, req, opts...)
	}
}

func (w *hook) ServerHandler(next server.FuncHandler) server.FuncHandler {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if v, ok := req.Body().(validator); ok {
			if err := v.Validate(); err != nil {
				return w.opts.ServerErrorFn(req, nil, err)
			}
		}
		err := next(ctx, req, rsp)
		if v, ok := rsp.(validator); ok && w.opts.ServerValidateResponse {
			if verr := v.Validate(); verr != nil {
				return w.opts.ServerErrorFn(req, rsp, verr)
			}
		}
		return err
	}
}
