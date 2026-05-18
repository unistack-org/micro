package mock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.unistack.org/micro/v5/client"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/errors"
)

func TestMockClient_String(t *testing.T) {
	mc := NewClient()
	if mc.String() != "mock" {
		t.Errorf("expected 'mock', got %q", mc.String())
	}
}

func TestMockClient_Name(t *testing.T) {
	mc := NewClient(client.Name("test-mock"))
	if mc.Name() != "test-mock" {
		t.Errorf("expected 'test-mock', got %q", mc.Name())
	}
}

func TestMockClient_Options(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c))
	opts := mc.Options()
	if opts.Codecs == nil {
		t.Fatal("expected non-nil codecs")
	}
}

func TestMockClient_Init(t *testing.T) {
	mc := NewClient()
	if err := mc.Init(client.Name("updated")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mc.Name() != "updated" {
		t.Errorf("expected 'updated', got %q", mc.Name())
	}
}

func TestMockClient_Stream(t *testing.T) {
	mc := NewClient(client.ContentType("application/json"))
	req := mc.NewRequest("svc", "Method", nil)
	stream, err := mc.Stream(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stream != nil {
		t.Error("expected nil stream")
	}
}

func TestMockClient_NewRequest(t *testing.T) {
	mc := NewClient(client.ContentType("application/json"))
	req := mc.NewRequest("svc", "Method", "body")
	if req.Service() != "svc" {
		t.Errorf("expected svc, got %q", req.Service())
	}
	if req.Method() != "Method" {
		t.Errorf("expected Method, got %q", req.Method())
	}
	if req.Body() != "body" {
		t.Errorf("expected 'body', got %v", req.Body())
	}
	if req.Codec() != nil {
		t.Error("expected nil codec")
	}
	if req.Stream() {
		t.Error("expected non-streaming")
	}
}

func TestMockClient_NewRequestWithContentType(t *testing.T) {
	mc := NewClient(client.ContentType("application/json"))
	req := mc.NewRequest("svc", "Method", "body", client.RequestContentType("application/json"))
	if req.ContentType() != "application/json" {
		t.Errorf("expected application/json, got %q", req.ContentType())
	}
}

func TestMockClient_ExpectRequestAndCall_Success(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c), client.ContentType("application/json"))

	type Response struct {
		Value string
	}
	expected := Response{Value: "ok"}

	// Use an inner request as body so mock takes the client.Request branch (no body comparison)
	innerReq := mc.NewRequest("inner", "fn", nil)
	req := mc.NewRequest("svc", "Method", innerReq)
	mc.ExpectRequest(req).WillReturnResponse("application/json", expected)

	var rsp Response
	actualReq := mc.NewRequest("svc", "Method", innerReq)
	if err := mc.Call(context.Background(), actualReq, &rsp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rsp.Value != "ok" {
		t.Errorf("expected 'ok', got %q", rsp.Value)
	}
}

func TestMockClient_ExpectRequestAndCall_ReturnError(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c), client.ContentType("application/json"))

	req := mc.NewRequest("svc", "Method", client.Request(nil))
	mc.ExpectRequest(req).WillReturnError(errors.InternalServerError("test", "something failed"))

	actualReq := mc.NewRequest("svc", "Method", client.Request(nil))
	err := mc.Call(context.Background(), actualReq, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMockClient_ExpectRequestAndCall_WithDelay(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c), client.ContentType("application/json"))

	req := mc.NewRequest("svc", "Delayed", client.Request(nil))
	mc.ExpectRequest(req).WillDelayFor(10 * time.Millisecond).WillReturnError(nil)

	actualReq := mc.NewRequest("svc", "Delayed", client.Request(nil))
	start := time.Now()
	_ = mc.Call(context.Background(), actualReq, nil)
	if time.Since(start) < 5*time.Millisecond {
		t.Error("expected delay")
	}
}

func TestMockClient_ExpectationsWereMet_Unfulfilled(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c), client.ContentType("application/json"))

	req := mc.NewRequest("svc", "Method", client.Request(nil))
	mc.ExpectRequest(req)

	if err := mc.ExpectationsWereMet(); err == nil {
		t.Fatal("expected error for unfulfilled expectation")
	}
}

func TestMockClient_ExpectationsWereMet_Fulfilled(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c), client.ContentType("application/json"))

	req := mc.NewRequest("svc", "Method", client.Request(nil))
	mc.ExpectRequest(req)

	actualReq := mc.NewRequest("svc", "Method", client.Request(nil))
	_ = mc.Call(context.Background(), actualReq, nil)

	if err := mc.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMockClient_Call_NoCodec(t *testing.T) {
	mc := NewClient() // no codec registered
	req := mc.NewRequest("svc", "Method", nil)
	req.(*mockRequest).contentType = "application/json"
	err := mc.Call(context.Background(), req, nil)
	if err == nil {
		t.Fatal("expected error for missing codec")
	}
}

func TestMockClient_Call_ServiceNotFound(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c))
	req := mc.NewRequest("unknown_svc", "unknown_method", nil)
	err := mc.Call(context.Background(), req, nil)
	if err == nil {
		t.Fatal("expected error for service not found")
	}
}

func TestMockClient_Call_WithBytesBody(t *testing.T) {
	c := codec.NewCodec()
	mc := NewClient(client.Codec("application/json", c), client.ContentType("application/json"))

	type MyResp struct {
		Name string `json:"name"`
	}
	type MyReq struct {
		Name string `json:"name"`
	}

	actualBody := &MyReq{Name: "test"}
	// expected body as bytes (JSON), actual body as pointer — mock deserializes expected and compares
	expectedBodyBytes, _ := c.Marshal(actualBody)

	// expected request carries bytes; actual request carries the pointer
	expectedReq := mc.NewRequest("svc", "Method", expectedBodyBytes)
	mc.ExpectRequest(expectedReq).WillReturnResponse("application/json", MyResp{Name: "result"})

	var rsp MyResp
	actualReq := mc.NewRequest("svc", "Method", actualBody)
	if err := mc.Call(context.Background(), actualReq, &rsp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExpectedRequest_String(t *testing.T) {
	mc := NewClient()
	req := mc.NewRequest("svc", "Method", nil)
	e := mc.ExpectRequest(req)
	e.WillReturnError(fmt.Errorf("test error"))
	s := e.String()
	if s == "" {
		t.Error("expected non-empty String()")
	}
}
