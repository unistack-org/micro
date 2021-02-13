// +build ignore

package micro

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/server"
)

// Function is a one time executing Service
type Function interface {
	// Inherits Service interface
	Service
	// Done signals to complete execution
	Done() error
	// Handle registers an RPC handler
	Handle(v interface{}) error
	// Subscribe registers a subscriber
	Subscribe(topic string, v interface{}) error
}

type function struct {
	cancel context.CancelFunc
	Service
}

// NewFunction returns a new Function for a one time executing Service
func NewFunction(opts ...Option) Function {
	return newFunction(opts...)
}

func fnHandlerWrapper(f Function) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			defer f.Done()
			return h(ctx, req, rsp)
		}
	}
}

func fnSubWrapper(f Function) server.SubscriberWrapper {
	return func(s server.SubscriberFunc) server.SubscriberFunc {
		return func(ctx context.Context, msg server.Message) error {
			defer f.Done()
			return s(ctx, msg)
		}
	}
}

func newFunction(opts ...Option) Function {
	ctx, cancel := context.WithCancel(context.Background())

	// force ttl/interval
	fopts := []Option{
		RegisterTTL(time.Minute),
		RegisterInterval(time.Second * 30),
	}

	// prepend to opts
	fopts = append(fopts, opts...)

	// make context the last thing
	fopts = append(fopts, Context(ctx))

	service := &service{opts: NewOptions(fopts...)}

	fn := &function{
		cancel:  cancel,
		Service: service,
	}

	service.Server().Init(
		// ensure the service waits for requests to finish
		server.Wait(nil),
		// wrap handlers and subscribers to finish execution
		server.WrapHandler(fnHandlerWrapper(fn)),
		server.WrapSubscriber(fnSubWrapper(fn)),
	)

	return fn
}

func (f *function) Done() error {
	f.cancel()
	return nil
}

func (f *function) Handle(v interface{}) error {
	return f.Service.Server().Handle(
		f.Service.Server().NewHandler(v),
	)
}

func (f *function) Subscribe(topic string, v interface{}) error {
	return f.Service.Server().Subscribe(
		f.Service.Server().NewSubscriber(topic, v),
	)
}
