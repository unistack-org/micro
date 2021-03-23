package pb

import (
	"bytes"
	"context"

	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/meter"
)

var (
	// guard to fail early
	_ MeterServer = &handler{}
)

type Empty struct{}

type handler struct {
	meter meter.Meter
	opts  []meter.Option
}

func NewHandler(meter meter.Meter, opts ...meter.Option) *handler {
	return &handler{meter: meter, opts: opts}
}

func (h *handler) Metrics(ctx context.Context, req *Empty, rsp *codec.Frame) error {
	buf := bytes.NewBuffer(nil)
	if err := h.meter.Write(buf, h.opts...); err != nil {
		return err
	}

	rsp.Data = buf.Bytes()

	return nil
}
