package server

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/unistack-org/micro/v3/broker"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/registry"
)

const (
	subSig = "func(context.Context, interface{}) error"
)

var (
	// Precompute the reflect type for error. Can't use error directly
	// because Typeof takes an empty interface value. This is annoying.
	typeOfError = reflect.TypeOf((*error)(nil)).Elem()
)

type handler struct {
	method  reflect.Value
	reqType reflect.Type
	ctxType reflect.Type
}

type subscriber struct {
	topic      string
	rcvr       reflect.Value
	typ        reflect.Type
	subscriber interface{}
	handlers   []*handler
	endpoints  []*registry.Endpoint
	opts       SubscriberOptions
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// ValidateSubscriber func
func ValidateSubscriber(sub Subscriber) error {
	typ := reflect.TypeOf(sub.Subscriber())
	var argType reflect.Type

	if typ.Kind() == reflect.Func {
		name := "Func"
		switch typ.NumIn() {
		case 2:
			argType = typ.In(1)
		default:
			return fmt.Errorf("subscriber %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), subSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("subscriber %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("subscriber %v has wrong number of outs: %v require signature %s",
				name, typ.NumOut(), subSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("subscriber %v returns %v not error", name, returnType.String())
		}
	} else {
		hdlr := reflect.ValueOf(sub.Subscriber())
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)

			switch method.Type.NumIn() {
			case 3:
				argType = method.Type.In(2)
			default:
				return fmt.Errorf("subscriber %v.%v takes wrong number of args: %v required signature %s",
					name, method.Name, method.Type.NumIn(), subSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"subscriber %v.%v has wrong number of outs: %v require signature %s",
					name, method.Name, method.Type.NumOut(), subSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("subscriber %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}

func newSubscriber(topic string, sub interface{}, opts ...SubscriberOption) Subscriber {
	var endpoints []*registry.Endpoint
	var handlers []*handler

	options := NewSubscriberOptions(opts...)

	if typ := reflect.TypeOf(sub); typ.Kind() == reflect.Func {
		h := &handler{
			method: reflect.ValueOf(sub),
		}

		switch typ.NumIn() {
		case 1:
			h.reqType = typ.In(0)
		case 2:
			h.ctxType = typ.In(0)
			h.reqType = typ.In(1)
		}

		handlers = append(handlers, h)
		ep := &registry.Endpoint{
			Name:     "Func",
			Request:  registry.ExtractSubValue(typ),
			Metadata: metadata.New(2),
		}
		ep.Metadata.Set("topic", topic)
		ep.Metadata.Set("subscriber", "true")
		endpoints = append(endpoints, ep)
	} else {
		hdlr := reflect.ValueOf(sub)
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			h := &handler{
				method: method.Func,
			}

			switch method.Type.NumIn() {
			case 2:
				h.reqType = method.Type.In(1)
			case 3:
				h.ctxType = method.Type.In(1)
				h.reqType = method.Type.In(2)
			}

			handlers = append(handlers, h)
			ep := &registry.Endpoint{
				Name:     name + "." + method.Name,
				Request:  registry.ExtractSubValue(method.Type),
				Metadata: metadata.New(2),
			}
			ep.Metadata.Set("topic", topic)
			ep.Metadata.Set("subscriber", "true")
			endpoints = append(endpoints, ep)
		}
	}

	return &subscriber{
		rcvr:       reflect.ValueOf(sub),
		typ:        reflect.TypeOf(sub),
		topic:      topic,
		subscriber: sub,
		handlers:   handlers,
		endpoints:  endpoints,
		opts:       options,
	}
}

func (n *noopServer) createSubHandler(sb *subscriber, opts Options) broker.Handler {
	return func(p broker.Event) (err error) {
		defer func() {
			if r := recover(); r != nil {
				n.RLock()
				config := n.opts
				n.RUnlock()
				if config.Logger.V(logger.ErrorLevel) {
					config.Logger.Error("panic recovered: ", r)
					config.Logger.Error(string(debug.Stack()))
				}
				err = errors.InternalServerError(n.opts.Name+".subscriber", "panic recovered: %v", r)
			}
		}()

		msg := p.Message()
		// if we don't have headers, create empty map
		if msg.Header == nil {
			msg.Header = make(map[string]string)
		}

		ct := msg.Header["Content-Type"]
		if len(ct) == 0 {
			msg.Header["Content-Type"] = defaultContentType
			ct = defaultContentType
		}
		cf, err := n.newCodec(ct)
		if err != nil {
			return err
		}

		hdr := make(map[string]string, len(msg.Header))
		for k, v := range msg.Header {
			hdr[k] = v
		}
		delete(hdr, "Content-Type")
		ctx := metadata.NewContext(sb.opts.Context, hdr)

		results := make(chan error, len(sb.handlers))

		for i := 0; i < len(sb.handlers); i++ {
			handler := sb.handlers[i]

			var isVal bool
			var req reflect.Value

			if handler.reqType.Kind() == reflect.Ptr {
				req = reflect.New(handler.reqType.Elem())
			} else {
				req = reflect.New(handler.reqType)
				isVal = true
			}
			if isVal {
				req = req.Elem()
			}

			if err = cf.ReadBody(bytes.NewBuffer(msg.Body), req.Interface()); err != nil {
				return err
			}

			fn := func(ctx context.Context, msg Message) error {
				var vals []reflect.Value
				if sb.typ.Kind() != reflect.Func {
					vals = append(vals, sb.rcvr)
				}
				if handler.ctxType != nil {
					vals = append(vals, reflect.ValueOf(ctx))
				}

				vals = append(vals, reflect.ValueOf(msg.Payload()))

				returnValues := handler.method.Call(vals)
				if rerr := returnValues[0].Interface(); rerr != nil {
					return rerr.(error)
				}
				return nil
			}

			for i := len(opts.SubWrappers); i > 0; i-- {
				fn = opts.SubWrappers[i-1](fn)
			}

			if n.wg != nil {
				n.wg.Add(1)
			}
			go func() {
				if n.wg != nil {
					defer n.wg.Done()
				}
				err := fn(ctx, &rpcMessage{
					topic:       sb.topic,
					contentType: ct,
					payload:     req.Interface(),
					header:      msg.Header,
					body:        msg.Body,
				})
				results <- err
			}()
		}
		var errors []string
		for i := 0; i < len(sb.handlers); i++ {
			if rerr := <-results; rerr != nil {
				errors = append(errors, rerr.Error())
			}
		}
		if len(errors) > 0 {
			err = fmt.Errorf("subscriber error: %s", strings.Join(errors, "\n"))
		}

		return err
	}
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Subscriber() interface{} {
	return s.subscriber
}

func (s *subscriber) Endpoints() []*registry.Endpoint {
	return s.endpoints
}

func (s *subscriber) Options() SubscriberOptions {
	return s.opts
}
