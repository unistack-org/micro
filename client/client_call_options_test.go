package client

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v4/options"
)

func TestNewClientCallOptions(t *testing.T) {
	var flag bool
	w := func(fn CallFunc) CallFunc {
		flag = true
		return fn
	}
	c := NewClientCallOptions(NewClient(),
		options.Address("127.0.0.1"),
		WithCallWrapper(w),
		RequestTimeout(1*time.Millisecond),
		Retries(0),
		Backoff(BackoffInterval(10*time.Millisecond, 100*time.Millisecond)),
	)
	_ = c.Call(context.TODO(), c.NewRequest("service", "endpoint", nil), nil)
	if !flag {
		t.Fatalf("NewClientCallOptions not works")
	}
}
