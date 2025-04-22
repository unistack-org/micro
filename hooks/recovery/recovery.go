package recovery

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/server"
)

func NewOptions(opts ...Option) Options {
	options := Options{
		ServerHandlerFn:    DefaultServerHandlerFn,
		ServerSubscriberFn: DefaultServerSubscriberFn,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

type Options struct {
	ServerHandlerFn    func(context.Context, server.Request, interface{}, error) error
	ServerSubscriberFn func(context.Context, server.Message, error) error
}

type Option func(*Options)

func ServerHandlerFunc(fn func(context.Context, server.Request, interface{}, error) error) Option {
	return func(o *Options) {
		o.ServerHandlerFn = fn
	}
}

func ServerSubscriberFunc(fn func(context.Context, server.Message, error) error) Option {
	return func(o *Options) {
		o.ServerSubscriberFn = fn
	}
}

var (
	DefaultServerHandlerFn = func(ctx context.Context, req server.Request, rsp interface{}, err error) error {
		return errors.BadRequest("", "%v", err)
	}
	DefaultServerSubscriberFn = func(ctx context.Context, req server.Message, err error) error {
		return errors.BadRequest("", "%v", err)
	}
)

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

func (w *hook) ServerSubscriber(next server.FuncSubHandler) server.FuncSubHandler {
	return func(ctx context.Context, msg server.Message) (err error) {
		defer func() {
			r := recover()
			switch verr := r.(type) {
			case nil:
				return
			case error:
				err = w.opts.ServerSubscriberFn(ctx, msg, verr)
			default:
				err = w.opts.ServerSubscriberFn(ctx, msg, fmt.Errorf("%v", r))
			}
		}()
		err = next(ctx, msg)
		return err
	}
}
