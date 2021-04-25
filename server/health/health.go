package health

import (
	"context"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/server"
)

var (
	// guard to fail early
	_ HealthServer = &handler{}
)

type handler struct {
	server server.Server
	opts   Options
}

type CheckFunc func(context.Context) error

type Option func(*Options)

type Options struct {
	Version     string
	Name        string
	LiveChecks  []CheckFunc
	ReadyChecks []CheckFunc
}

func LiveChecks(fns ...CheckFunc) Option {
	return func(o *Options) {
		o.LiveChecks = append(o.LiveChecks, fns...)
	}
}

func ReadyChecks(fns ...CheckFunc) Option {
	return func(o *Options) {
		o.ReadyChecks = append(o.ReadyChecks, fns...)
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func NewHandler(opts ...Option) *handler {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &handler{opts: options}
}

func (h *handler) Live(ctx context.Context, req *codec.Frame, rsp *codec.Frame) error {
	var err error
	for _, fn := range h.opts.LiveChecks {
		if err = fn(ctx); err != nil {
			return errors.ServiceUnavailable(h.opts.Name, "%v", err)
		}
	}
	return nil
}

func (h *handler) Ready(ctx context.Context, req *codec.Frame, rsp *codec.Frame) error {
	var err error
	for _, fn := range h.opts.ReadyChecks {
		if err = fn(ctx); err != nil {
			return errors.ServiceUnavailable(h.opts.Name, "%v", err)
		}
	}
	return nil
}

func (h *handler) Version(ctx context.Context, req *codec.Frame, rsp *codec.Frame) error {
	rsp.Data = []byte(h.opts.Version)
	return nil
}
