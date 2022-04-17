package client

import (
	"context"
	"fmt"
	"testing"

	"go.unistack.org/micro/v3/errors"
)

func TestRetryAlways(t *testing.T) {
	tests := []error{
		nil,
		errors.InternalServerError("test", "%s", "test"),
		fmt.Errorf("test"),
	}

	for _, e := range tests {
		ok, er := RetryAlways(context.TODO(), nil, 1, e)
		if !ok || er != nil {
			t.Fatal("RetryAlways not works properly")
		}
	}
}

func TestRetryNever(t *testing.T) {
	tests := []error{
		nil,
		errors.InternalServerError("test", "%s", "test"),
		fmt.Errorf("test"),
	}

	for _, e := range tests {
		ok, er := RetryNever(context.TODO(), nil, 1, e)
		if ok || er != nil {
			t.Fatal("RetryNever not works properly")
		}
	}
}

func TestRetryOnError(t *testing.T) {
	tests := []error{
		fmt.Errorf("test"),
		errors.NotFound("test", "%s", "test"),
		errors.Timeout("test", "%s", "test"),
	}

	for i, e := range tests {
		ok, er := RetryOnError(context.TODO(), nil, 1, e)
		if i == 2 && (!ok || er != nil) {
			t.Fatal("RetryOnError not works properly")
		}
	}
}

func TestRetryOnErrors(t *testing.T) {
	tests := []error{
		fmt.Errorf("test"),
		errors.NotFound("test", "%s", "test"),
		errors.Timeout("test", "%s", "test"),
	}

	fn := RetryOnErrors(404)
	for i, e := range tests {
		ok, er := fn(context.TODO(), nil, 1, e)
		if i == 1 && (!ok || er != nil) {
			t.Fatal("RetryOnErrors not works properly")
		}
	}
}
