package handler

import (
	"bytes"
	"context"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/meter"
)

// guard to fail early
var _ MeterServer = &Handler{}

type Handler struct {
	opts Options
}

type Option func(*Options)

type Options struct {
	Meter        meter.Meter
	Name         string
	MeterOptions []meter.Option
}

func Meter(m meter.Meter) Option {
	return func(o *Options) {
		o.Meter = m
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func MeterOptions(opts ...meter.Option) Option {
	return func(o *Options) {
		o.MeterOptions = append(o.MeterOptions, opts...)
	}
}

func NewOptions(opts ...Option) Options {
	options := Options{Meter: meter.DefaultMeter}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func NewHandler(opts ...Option) *Handler {
	options := NewOptions(opts...)
	return &Handler{opts: options}
}

func (h *Handler) Metrics(ctx context.Context, req *codec.Frame, rsp *codec.Frame) error {
	buf := bytes.NewBuffer(nil)
	if err := h.opts.Meter.Write(buf, h.opts.MeterOptions...); err != nil {
		return errors.InternalServerError(h.opts.Name, "%v", err)
	}

	rsp.Data = buf.Bytes()

	return nil
}
