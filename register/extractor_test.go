package register

import (
	"context"
	"reflect"
	"testing"
)

type TestHandler struct{}

type TestRequest struct{}

type TestResponse struct{}

func (t *TestHandler) Test(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	return nil
}

func TestExtractEndpoint(t *testing.T) {
	handler := &TestHandler{}
	typ := reflect.TypeOf(handler)

	var endpoints []*Endpoint

	for m := 0; m < typ.NumMethod(); m++ {
		if e := ExtractEndpoint(typ.Method(m)); e != nil {
			endpoints = append(endpoints, e)
		}
	}

	if i := len(endpoints); i != 1 {
		t.Fatalf("Expected 1 endpoint, have %d", i)
	}

	if endpoints[0].Name != "Test" {
		t.Fatalf("Expected handler Test, got %s", endpoints[0].Name)
	}

	if endpoints[0].Request == nil {
		t.Fatal("Expected non nil Request")
	}

	if endpoints[0].Response == nil {
		t.Fatal("Expected non nil Request")
	}

	if endpoints[0].Request.Name != "TestRequest" {
		t.Fatalf("Expected TestRequest got %s", endpoints[0].Request.Name)
	}

	if endpoints[0].Response.Name != "TestResponse" {
		t.Fatalf("Expected TestResponse got %s", endpoints[0].Response.Name)
	}

	if endpoints[0].Request.Type != "TestRequest" {
		t.Fatalf("Expected TestRequest type got %s", endpoints[0].Request.Type)
	}

	if endpoints[0].Response.Type != "TestResponse" {
		t.Fatalf("Expected TestResponse type got %s", endpoints[0].Response.Type)
	}

}
