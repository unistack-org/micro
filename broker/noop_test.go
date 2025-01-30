package broker

import (
	"context"
	"testing"
)

type testHook struct {
	f bool
}

func (t *testHook) Publish1(fn FuncPublish) FuncPublish {
	return func(ctx context.Context, topic string, messages ...Message) error {
		t.f = true
		return fn(ctx, topic, messages...)
	}
}

func TestNoopHook(t *testing.T) {
	h := &testHook{}

	b := NewBroker(Hooks(HookPublish(h.Publish1)))

	if err := b.Init(); err != nil {
		t.Fatal(err)
	}

	if err := b.Publish(context.TODO(), "", nil); err != nil {
		t.Fatal(err)
	}

	if !h.f {
		t.Fatal("hook not works")
	}
}
