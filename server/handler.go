package server

import (
	"reflect"

	"github.com/unistack-org/micro/v3/register"
)

type rpcHandler struct {
	name      string
	handler   interface{}
	endpoints []*register.Endpoint
	opts      HandlerOptions
}

func newRpcHandler(handler interface{}, opts ...HandlerOption) Handler {
	options := NewHandlerOptions(opts...)

	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	var endpoints []*register.Endpoint

	for m := 0; m < typ.NumMethod(); m++ {
		if e := register.ExtractEndpoint(typ.Method(m)); e != nil {
			e.Name = name + "." + e.Name

			for k, v := range options.Metadata[e.Name] {
				e.Metadata[k] = v
			}

			endpoints = append(endpoints, e)
		}
	}

	return &rpcHandler{
		name:      name,
		handler:   handler,
		endpoints: endpoints,
		opts:      options,
	}
}

func (r *rpcHandler) Name() string {
	return r.name
}

func (r *rpcHandler) Handler() interface{} {
	return r.handler
}

func (r *rpcHandler) Endpoints() []*register.Endpoint {
	return r.endpoints
}

func (r *rpcHandler) Options() HandlerOptions {
	return r.opts
}
