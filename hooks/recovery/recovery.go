package recovery

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v4/errors"
	"go.unistack.org/micro/v4/server"
)

func NewOptions(opts ...Option) Options {
	options := Options{
		ServerHandlerFn: DefaultServerHandlerFn,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

type Options struct {
	ServerHandlerFn func(context.Context, server.Request, interface{}, error) error
}

type Option func(*Options)

func ServerHandlerFunc(fn func(context.Context, server.Request, interface{}, error) error) Option {
	return func(o *Options) {
		o.ServerHandlerFn = fn
	}
}

var DefaultServerHandlerFn = func(ctx context.Context, req server.Request, rsp interface{}, err error) error {
	return errors.BadRequest("", "%v", err)
}

var Hook = NewHook()

type hook struct {
	opts Options
}

func NewHook(opts ...Option) *hook {
	return &hook{opts: NewOptions(opts...)}
}

func (w *hook) ServerHandler(next server.FuncHandler) server.FuncHandler {
	return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
		defer func() {
			r := recover()
			switch verr := r.(type) {
			case nil:
				return
			case error:
				err = w.opts.ServerHandlerFn(ctx, req, rsp, verr)
			default:
				err = w.opts.ServerHandlerFn(ctx, req, rsp, fmt.Errorf("%v", r))
			}
		}()
		err = next(ctx, req, rsp)
		return err
	}
}
