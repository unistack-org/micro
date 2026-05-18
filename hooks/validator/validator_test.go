package validator

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/errors"
	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/server"
)

type validatableBody struct {
	body string
}

func (v *validatableBody) Validate() error {
	if v.body == "" {
		return errors.BadRequest("test", "missing body")
	}
	return nil
}

type validRequest struct {
	body string
}

func (v *validRequest) Body() any      { return &validatableBody{body: v.body} }
func (v *validRequest) Service() string { return "test" }
func (v *validRequest) Method() string  { return "test" }
func (v *validRequest) Endpoint() string { return "test" }
func (v *validRequest) ContentType() string { return "test" }
func (v *validRequest) Stream() bool    { return false }
func (v *validRequest) Codec() codec.Codec { return nil }

type validResponse struct {
	ok bool
}

func (v *validResponse) Validate() error {
	if !v.ok {
		return errors.BadRequest("test", "invalid")
	}
	return nil
}

func TestClientCallValidRequest(t *testing.T) {
	h := NewHook()
	next := func(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
		return nil
	}
	wrapped := h.ClientCall(next)
	req := &validRequest{body: "x"}
	rsp := &validResponse{ok: true}
	err := wrapped(context.Background(), req, rsp)
	if err != nil {
		t.Errorf("unexpected: %v", err)
	}
}

func TestClientCallInvalidRequest(t *testing.T) {
	h := NewHook()
	next := func(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
		return nil
	}
	wrapped := h.ClientCall(next)
	req := &validRequest{body: ""}
	rsp := &validResponse{ok: true}
	err := wrapped(context.Background(), req, rsp)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClientStreamInvalid(t *testing.T) {
	h := NewHook()
	next := func(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
		return nil, nil
	}
	wrapped := h.ClientStream(next)
	req := &validRequest{body: ""}
	_, err := wrapped(context.Background(), req)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClientValidateResponseEnabled(t *testing.T) {
	h := NewHook(ClientValidateResponse(true))
	next := func(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
		return nil
	}
	wrapped := h.ClientCall(next)
	req := &validRequest{body: "x"}
	rsp := &validResponse{ok: false}
	err := wrapped(context.Background(), req, rsp)
	if err == nil {
		t.Error("expected error for invalid response")
	}
}

func TestClientValidateResponseDisabled(t *testing.T) {
	h := NewHook(ClientValidateResponse(false))
	next := func(ctx context.Context, req client.Request, rsp any, opts ...client.CallOption) error {
		return nil
	}
	wrapped := h.ClientCall(next)
	req := &validRequest{body: "x"}
	rsp := &validResponse{ok: false}
	err := wrapped(context.Background(), req, rsp)
	if err != nil {
		t.Errorf("unexpected: %v", err)
	}
}

// mockServerRequest implements server.Request for testing.
type mockServerRequest struct {
	body any
}

func (r *mockServerRequest) Header() metadata.Metadata { return metadata.New(0) }
func (r *mockServerRequest) Body() any                 { return r.body }
func (r *mockServerRequest) Method() string            { return "test" }

func TestServerHandlerValidRequest(t *testing.T) {
	h := NewHook()
	called := false
	next := func(ctx context.Context, req server.Request, rsp any) error {
		called = true
		return nil
	}
	wrapped := h.ServerHandler(next)
	req := &mockServerRequest{body: &validatableBody{body: "ok"}}
	err := wrapped(context.Background(), req, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("next was not called")
	}
}

func TestServerHandlerInvalidRequest(t *testing.T) {
	h := NewHook()
	next := func(ctx context.Context, req server.Request, rsp any) error {
		return nil
	}
	wrapped := h.ServerHandler(next)
	req := &mockServerRequest{body: &validatableBody{body: ""}}
	err := wrapped(context.Background(), req, nil)
	if err == nil {
		t.Error("expected error for invalid request body")
	}
}

func TestServerHandlerNoValidator(t *testing.T) {
	h := NewHook()
	called := false
	next := func(ctx context.Context, req server.Request, rsp any) error {
		called = true
		return nil
	}
	wrapped := h.ServerHandler(next)
	// body does not implement validator
	req := &mockServerRequest{body: "plain string body"}
	err := wrapped(context.Background(), req, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("next was not called")
	}
}

func TestServerHandlerValidateResponse(t *testing.T) {
	// ServerValidateResponse option sets opts.ServerValidateResponse on the hook directly.
	h := NewHook()
	h.opts.ServerValidateResponse = true
	next := func(ctx context.Context, req server.Request, rsp any) error {
		return nil
	}
	wrapped := h.ServerHandler(next)
	req := &mockServerRequest{body: nil}
	rsp := &validResponse{ok: false}
	err := wrapped(context.Background(), req, rsp)
	if err == nil {
		t.Error("expected error for invalid response")
	}
}

func TestServerHandlerValidateResponseValid(t *testing.T) {
	h := NewHook()
	h.opts.ServerValidateResponse = true
	next := func(ctx context.Context, req server.Request, rsp any) error {
		return nil
	}
	wrapped := h.ServerHandler(next)
	req := &mockServerRequest{body: nil}
	rsp := &validResponse{ok: true}
	err := wrapped(context.Background(), req, rsp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClientStreamValid(t *testing.T) {
	h := NewHook()
	called := false
	next := func(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
		called = true
		return nil, nil
	}
	wrapped := h.ClientStream(next)
	req := &validRequest{body: "ok"}
	_, err := wrapped(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !called {
		t.Error("next was not called")
	}
}

func TestServerValidateResponseOption(t *testing.T) {
	opts := NewOptions(ServerValidateResponse(true))
	if !opts.ClientValidateResponse {
		t.Error("ServerValidateResponse option should set ClientValidateResponse")
	}
}

func TestClientReqErrorFnOption(t *testing.T) {
	customFn := func(req client.Request, rsp any, err error) error {
		return fmt.Errorf("custom: %w", err)
	}
	opts := NewOptions(ClientReqErrorFn(customFn))
	if opts.ClientErrorFn == nil {
		t.Error("ClientReqErrorFn option should set ClientErrorFn")
	}
}

func TestServerErrorFnOption(t *testing.T) {
	customFn := func(req server.Request, rsp any, err error) error {
		return fmt.Errorf("custom: %w", err)
	}
	opts := NewOptions(ServerErrorFn(customFn))
	if opts.ServerErrorFn == nil {
		t.Error("ServerErrorFn option should set ServerErrorFn")
	}
}