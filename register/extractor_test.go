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

func TestExtractValueNil(t *testing.T) {
	if got := ExtractValue(nil, 0); got != "" {
		t.Fatalf("expected empty string for nil type, got %q", got)
	}
}

func TestExtractValueDepthLimit(t *testing.T) {
	typ := reflect.TypeOf(TestRequest{})
	if got := ExtractValue(typ, 3); got != "" {
		t.Fatalf("expected empty string at depth 3, got %q", got)
	}
}

func TestExtractValuePointer(t *testing.T) {
	typ := reflect.TypeOf(&TestRequest{})
	got := ExtractValue(typ, 0)
	if got != "TestRequest" {
		t.Fatalf("expected 'TestRequest', got %q", got)
	}
}

func TestExtractValueStruct(t *testing.T) {
	typ := reflect.TypeOf(TestRequest{})
	got := ExtractValue(typ, 0)
	if got != "TestRequest" {
		t.Fatalf("expected 'TestRequest', got %q", got)
	}
}

func TestExtractValueSlice(t *testing.T) {
	typ := reflect.TypeOf([]TestRequest{})
	got := ExtractValue(typ, 0)
	if got != "" {
		t.Fatalf("expected empty string for slice type, got %q", got)
	}
}

func TestExtractValueMap(t *testing.T) {
	typ := reflect.TypeOf(map[string]TestRequest{})
	got := ExtractValue(typ, 0)
	if got != "" {
		t.Fatalf("expected empty string for map type, got %q", got)
	}
}

func TestExtractSubValue1Arg(t *testing.T) {
	// func(ctx context.Context) — 1 arg
	fn := func(ctx context.Context) {}
	typ := reflect.TypeOf(fn)
	// context.Context is an interface — Name() returns "Context", but it's exported
	_ = ExtractSubValue(typ)
}

func TestExtractSubValue2Args(t *testing.T) {
	// func(ctx context.Context, req *TestRequest) — 2 args
	fn := func(ctx context.Context, req *TestRequest) {}
	typ := reflect.TypeOf(fn)
	got := ExtractSubValue(typ)
	if got != "TestRequest" {
		t.Fatalf("expected 'TestRequest', got %q", got)
	}
}

func TestExtractSubValue3Args(t *testing.T) {
	// func(ctx context.Context, req *TestRequest, rsp *TestResponse) — 3 args
	fn := func(ctx context.Context, req *TestRequest, rsp *TestResponse) {}
	typ := reflect.TypeOf(fn)
	got := ExtractSubValue(typ)
	if got != "TestResponse" {
		t.Fatalf("expected 'TestResponse', got %q", got)
	}
}

func TestExtractSubValueDefault(t *testing.T) {
	// func with 0 args — returns ""
	fn := func() {}
	typ := reflect.TypeOf(fn)
	got := ExtractSubValue(typ)
	if got != "" {
		t.Fatalf("expected empty string for 0-arg func, got %q", got)
	}
}
