package mock

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/errors"
	rutil "go.unistack.org/micro/v5/util/reflect"
)

var _ client.Client = (*MockClient)(nil)

type MockClient struct {
	opts     client.Options
	mu       sync.Mutex
	expected []expectation
}

func (c *MockClient) newCodec(ct string) (codec.Codec, error) {
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}

	if cc, ok := c.opts.Codecs[ct]; ok {
		return cc, nil
	}

	return nil, codec.ErrUnknownContentType
}

type expectation interface {
	fulfilled() bool
	Lock()
	Unlock()
	String() string
}

type commonExpectation struct {
	sync.Mutex
	triggered bool
	err       error
}

func (e *commonExpectation) fulfilled() bool {
	return e.triggered
}

type ExpectedRequest struct {
	commonExpectation
	delay      time.Duration
	rsp        any
	rspContent string
	req        client.Request
}

func (e *ExpectedRequest) WillDelayFor(duration time.Duration) *ExpectedRequest {
	e.delay = duration
	return e
}

func (e *ExpectedRequest) WillReturnError(err error) *ExpectedRequest {
	e.err = err
	return e
}

func (e *ExpectedRequest) WillReturnResponse(contentType string, rsp any) *ExpectedRequest {
	e.rspContent = contentType
	e.rsp = rsp
	return e
}

func (e *ExpectedRequest) String() string {
	msg := "ExpectedRequest => expecting client.Call request"
	if e.err != nil {
		msg += fmt.Sprintf(", which should return error: %s", e.err)
	}
	if e.rsp != nil {
		msg += fmt.Sprintf(", which should return rsp: %v", e.rsp)
	}
	return msg
}

func (c *MockClient) ExpectationsWereMet() error {
	for _, e := range c.expected {
		e.Lock()
		fulfilled := e.fulfilled()
		e.Unlock()

		if !fulfilled {
			return fmt.Errorf("there is a remaining expectation which was not matched: %s", e)
		}
	}

	return nil
}

func (c *MockClient) ExpectRequest(req client.Request) *ExpectedRequest {
	e := &ExpectedRequest{req: req}
	c.expected = append(c.expected, e)
	return e
}

func (c *MockClient) Call(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	options := client.NewCallOptions(opts...)
	ct := req.ContentType()
	if len(options.ContentType) > 0 {
		ct = options.ContentType
	}

	cf, err := c.newCodec(ct)
	if err != nil {
		return errors.BadRequest("go.micro.client", "%+v", err.Error())
	}

	for _, e := range c.expected {
		er, ok := e.(*ExpectedRequest)
		if !ok {
			continue
		}

		if er.delay > 0 {
			time.Sleep(er.delay)
		}

		if er.req.Service() != req.Service() ||
			er.req.Method() != req.Method() {
			continue
		}

		er.triggered = true

		if er.err != nil {
			return er.err
		}

		if er.req == nil {
			return errors.BadRequest("go.micro.client", "empty request passed")
		}

		src := er.req.Body()
		switch reqbody := er.req.Body().(type) {
		case []byte:
			src, err = rutil.Zero(req.Body())
			if err == nil {
				err = cf.Unmarshal(reqbody, src)
			}
			if err != nil {
				return errors.BadRequest("go.micro.client", "%+v", err.Error())
			}
		case client.Request:
			break
		default:
			return errors.BadRequest("go.micro.client", "unknown request passed: %v", reqbody)
		}

		if !reflect.DeepEqual(req.Body(), src) {
			return errors.BadRequest("go.micro.client", "unexpected request %v != %v", req.Body(), src)
		}

		if er.rsp == nil {
			return nil
		}

		rspCt := ct
		if len(er.rspContent) > 0 {
			rspCt = er.rspContent
		}
		cfRsp, err := c.newCodec(rspCt)
		if err != nil {
			return errors.BadRequest("go.micro.client", "%+v", err.Error())
		}

		switch rspbody := er.rsp.(type) {
		case []byte:
			if err = cfRsp.Unmarshal(rspbody, rsp); err != nil {
				return errors.BadRequest("go.micro.client", "%+v", err.Error())
			}
			return nil
		}

		v := reflect.ValueOf(rsp)

		if t := reflect.TypeOf(rsp); t.Kind() == reflect.Ptr {
			v = reflect.Indirect(v)
		}
		response := er.rsp
		if t := reflect.TypeOf(er.rsp); t.Kind() == reflect.Func {
			response = reflect.ValueOf(er.rsp).Call([]reflect.Value{})[0].Interface()
		}

		v.Set(reflect.ValueOf(response))

		return nil
	}

	return fmt.Errorf("can't find service %s", req.Method())
}

func (c *MockClient) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	return nil, nil
}

func (c *MockClient) Init(opts ...client.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}

	return nil
}

func (c *MockClient) String() string {
	return "mock"
}

func (c *MockClient) Name() string {
	return c.opts.Name
}

func (c *MockClient) Options() client.Options {
	return c.opts
}

func (c *MockClient) NewRequest(service, method string, req any, opts ...client.RequestOption) client.Request {
	return newMockRequest(service, method, req, c.opts.ContentType, opts...)
}

func NewClient(opts ...client.Option) *MockClient {
	options := client.NewOptions(opts...)
	return &MockClient{opts: options}
}
