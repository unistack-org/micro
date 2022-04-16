package client

import (
	"context"
	"testing"
	"time"
)

func TestNewClientCallOptions(t *testing.T) {
	var flag bool
	w := func(fn CallFunc) CallFunc {
		flag = true
		return fn
	}
	c := NewClientCallOptions(NewClient(),
		WithAddress("127.0.0.1"),
		WithCallWrapper(w),
		WithRequestTimeout(1*time.Millisecond),
		WithRetries(0),
		WithBackoff(BackoffInterval(10*time.Millisecond, 100*time.Millisecond)),
	)
	_ = c.Call(context.TODO(), c.NewRequest("service", "endpoint", nil), nil)
	if !flag {
		t.Fatalf("NewClientCallOptions not works")
	}
}
