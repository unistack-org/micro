package wrapper

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/meter"
	"github.com/unistack-org/micro/v3/server"
)

// Wrapper provides a HandlerFunc for meter.Meter implementations
type Wrapper struct {
	meter meter.Meter
}

// New returns a *Wrapper configured with the given meter.Meter
func New(meter meter.Meter) *Wrapper {
	return &Wrapper{
		meter: meter,
	}
}

// HandlerFunc instruments handlers registered to a service:
func (w *Wrapper) HandlerFunc(handlerFunction server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		// Build some tags to describe the call:
		tags := metadata.New(2)
		tags.Set("method", req.Method())

		// Start the clock:
		callTime := time.Now()

		// Run the handlerFunction:
		err := handlerFunction(ctx, req, rsp)

		// Add a result tag:
		if err != nil {
			tags["status"] = "failure"
		} else {
			tags["status"] = "success"
		}

		// Instrument the result (if the DefaultClient has been configured):
		w.meter.Summary("service.handler", tags).Update(float64(time.Since(callTime).Seconds()))

		return err
	}
}
