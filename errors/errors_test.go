package errors

import (
	"encoding/json"
	er "errors"
	"fmt"
	"net/http"
	"testing"
)

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
	err = FromError(fmt.Errorf(msg))
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
	err = er.New(err.Error())
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

	err3 := er.New("my test err")
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
