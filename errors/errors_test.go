package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestIsRetrayable(t *testing.T) {
	err := fmt.Errorf("ORA-")
	if !IsRetryable(err, RetrayableOracleErrors...) {
		t.Fatalf("IsRetrayable not works")
	}
}

func TestMarshalJSON(t *testing.T) {
	e := InternalServerError("id", "err: %v", fmt.Errorf("err: %v", `xxx: "UNIX_TIMESTAMP": invalid identifier`))
	_, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmpty(t *testing.T) {
	msg := "test"
	var err *Error
	err = FromError(errors.New(msg))
	if err.Detail != msg {
		t.Fatalf("invalid error %v", err)
	}
	err = FromError(fmt.Errorf(`{"id":"","detail":"%s","status":"%s","code":0}`, msg, msg))
	if err.Detail != msg || err.Status != msg {
		t.Fatalf("invalid error %#+v", err)
	}
}

func TestFromError(t *testing.T) {
	err := NotFound("go.micro.test", "%s", "example")
	merr := FromError(err)
	if merr.ID != "go.micro.test" || merr.Code != 404 {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
	err = errors.New(err.Error())
	merr = FromError(err)
	if merr.ID != "go.micro.test" || merr.Code != 404 {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
}

func TestEqual(t *testing.T) {
	err1 := NotFound("myid1", "msg1")
	err2 := NotFound("myid2", "msg2")

	if !Equal(err1, err2) {
		t.Fatal("errors must be equal")
	}

	err3 := errors.New("my test err")
	if Equal(err1, err3) {
		t.Fatal("errors must be not equal")
	}
}

func TestErrors(t *testing.T) {
	testData := []*Error{
		{
			ID:     "test",
			Code:   500,
			Detail: "Internal server error",
			Status: http.StatusText(500),
		},
	}

	for _, e := range testData {
		ne := New(e.ID, e.Detail, e.Code)

		if e.Error() != ne.Error() {
			t.Fatalf("Expected %s got %s", e.Error(), ne.Error())
		}

		pe := Parse(ne.Error())

		if pe == nil {
			t.Fatalf("Expected error got nil %v", pe)
		}

		if pe.ID != e.ID {
			t.Fatalf("Expected %s got %s", e.ID, pe.ID)
		}

		if pe.Detail != e.Detail {
			t.Fatalf("Expected %s got %s", e.Detail, pe.Detail)
		}

		if pe.Code != e.Code {
			t.Fatalf("Expected %d got %d", e.Code, pe.Code)
		}

		if pe.Status != e.Status {
			t.Fatalf("Expected %s got %s", e.Status, pe.Status)
		}
	}
}

func TestCodeIn(t *testing.T) {
	err := InternalServerError("id", "%s", "msg")

	if ok := CodeIn(err, 400, 500); !ok {
		t.Fatalf("CodeIn not works: %v", err)
	}

	if ok := CodeIn(err.(*Error).Code, 500); !ok {
		t.Fatalf("CodeIn not works: %v", err)
	}

	if ok := CodeIn(err, 100); ok {
		t.Fatalf("CodeIn not works: %v", err)
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name string
		fn   func() error
		code int
	}{
		{"BadRequest", func() error { return BadRequest("id", "msg") }, 400},
		{"Unauthorized", func() error { return Unauthorized("id", "msg") }, 401},
		{"Forbidden", func() error { return Forbidden("id", "msg") }, 403},
		{"MethodNotAllowed", func() error { return MethodNotAllowed("id", "msg") }, 405},
		{"Timeout", func() error { return Timeout("id", "msg") }, 408},
		{"Conflict", func() error { return Conflict("id", "msg") }, 409},
		{"NotImplemented", func() error { return NotImplemented("id", "msg") }, 501},
		{"BadGateway", func() error { return BadGateway("id", "msg") }, 502},
		{"ServiceUnavailable", func() error { return ServiceUnavailable("id", "msg") }, 503},
		{"GatewayTimeout", func() error { return GatewayTimeout("id", "msg") }, 504},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			e, ok := err.(*Error)
			if !ok {
				t.Fatalf("expected *Error, got %T", err)
			}
			if e.Code != int32(tt.code) {
				t.Errorf("expected code %d, got %d", tt.code, e.Code)
			}
		})
	}
}

func TestReset(t *testing.T) {
	e := &Error{ID: "test", Code: 500, Detail: "detail", Status: "Internal Server Error"}
	e.Reset()
	if e.ID != "" || e.Code != 0 || e.Detail != "" || e.Status != "" {
		t.Fatalf("Reset did not clear Error: %+v", e)
	}
}

func TestRetryable(t *testing.T) {
	inner := fmt.Errorf("inner error")
	re := Retryable(inner)
	if re.Error() != inner.Error() {
		t.Fatalf("expected %q, got %q", inner.Error(), re.Error())
	}
	// Unwrap via errors.Is on wrapped error
	if !errors.Is(re, inner) {
		t.Fatal("errors.Is should find wrapped inner error")
	}
}

func TestRetryableNilInner(t *testing.T) {
	re := Retryable(nil)
	if re.Error() != "" {
		t.Fatalf("expected empty string for nil inner error, got %q", re.Error())
	}
}

func TestIsRetryableMicroErrors(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{Unauthorized("id", "msg"), true},
		{Forbidden("id", "msg"), true},
		{Timeout("id", "msg"), true},
		{InternalServerError("id", "msg"), true},
		{NotImplemented("id", "msg"), true},
		{BadGateway("id", "msg"), true},
		{ServiceUnavailable("id", "msg"), true},
		{GatewayTimeout("id", "msg"), true},
		{BadRequest("id", "msg"), false},
		{NotFound("id", "msg"), false},
		{Retryable(fmt.Errorf("some error")), true},
	}

	for _, tt := range tests {
		got := IsRetryable(tt.err, RetryableMicroErrors...)
		if got != tt.expected {
			t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.expected)
		}
	}
}

func TestIsRetryableNoFuncs(t *testing.T) {
	err := fmt.Errorf("some error")
	if IsRetryable(err) {
		t.Fatal("IsRetryable with no funcs should return false")
	}
}

func TestEqualBothNonError(t *testing.T) {
	err1 := fmt.Errorf("same error")
	err2 := fmt.Errorf("different error")
	// Two plain errors are compared by identity, not message
	if Equal(err1, err1) != (err1 == err1) {
		t.Fatal("Equal should compare plain errors by identity")
	}
	if Equal(err1, err2) {
		t.Fatal("two different plain errors should not be equal")
	}
}

func TestCodeInDefault(t *testing.T) {
	// Pass a type that is neither *Error nor int32
	if CodeIn("not an error", 500) {
		t.Fatal("CodeIn with unsupported type should return false")
	}
}

func TestFromErrorNil(t *testing.T) {
	if FromError(nil) != nil {
		t.Fatal("FromError(nil) should return nil")
	}
}

func TestIsRetryablePostgresErrors(t *testing.T) {
	tests := []string{
		"number of field descriptions must equal number of",
		"not a pointer",
		"values, but dst struct has only",
		"struct doesn't have corresponding row field",
		"cannot find field",
		"cannot scan",
		"cannot convert",
		"failed to connect to",
	}
	for _, msg := range tests {
		err := fmt.Errorf("%s", msg)
		if !IsRetryable(err, RetrayablePostgresErrors...) {
			t.Errorf("IsRetryable postgres: expected true for %q", msg)
		}
	}
	// Non-retryable postgres error
	err := fmt.Errorf("some unrelated error")
	if IsRetryable(err, RetrayablePostgresErrors...) {
		t.Error("IsRetryable postgres: expected false for unrelated error")
	}
}

func TestUnmarshalInvalidData(t *testing.T) {
	e := &Error{}
	err := e.Unmarshal([]byte("short"))
	if err == nil {
		t.Fatal("expected error for short data")
	}
}
