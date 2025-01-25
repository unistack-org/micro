package server

import (
	"testing"

	"go.unistack.org/micro/v4/errors"
)

func TestError(t *testing.T) {
	e := NewError("svc1")
	err := e.BadRequest("%s", "test")
	merr, ok := err.(*errors.Error)
	if !ok {
		t.Fatal("error not *errors.Error")
	}
	if merr.ID != "svc1" {
		t.Fatal("id != svc1")
	}
}
