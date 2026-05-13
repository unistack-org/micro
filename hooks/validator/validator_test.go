package validator

import (
	"context"
	"testing"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/errors"
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