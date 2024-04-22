package client

import (
	"context"
	"testing"
)

type testHook struct {
	f bool
}

func (t *testHook) Publish(fn FuncPublish) FuncPublish {
	return func(ctx context.Context, msg Message, opts ...PublishOption) error {
		t.f = true
		return fn(ctx, msg, opts...)
	}
}

func TestNoopHook(t *testing.T) {
	h := &testHook{}

	c := NewClient(Hooks(HookPublish(h.Publish)))

	if err := c.Init(); err != nil {
		t.Fatal(err)
	}

	if err := c.Publish(context.TODO(), c.NewMessage("", nil, MessageContentType("application/octet-stream"))); err != nil {
		t.Fatal(err)
	}

	if !h.f {
		t.Fatal("hook not works")
	}
}
