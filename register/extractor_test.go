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

	if endpoints[0].Request == "" {
		t.Fatal("Expected non nil Request")
	}

	if endpoints[0].Response == "" {
		t.Fatal("Expected non nil Request")
	}

	if endpoints[0].Request != "TestRequest" {
		t.Fatalf("Expected TestRequest got %s", endpoints[0].Request)
	}

	if endpoints[0].Response != "TestResponse" {
		t.Fatalf("Expected TestResponse got %s", endpoints[0].Response)
	}

	t.Logf("XXX %#+v\n", endpoints[0])
}
