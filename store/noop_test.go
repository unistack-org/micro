package store

import (
	"context"
	"errors"
	"testing"
)

type testHook struct {
	f bool
}

func (t *testHook) Exists(fn FuncExists) FuncExists {
	return func(ctx context.Context, key string, opts ...ExistsOption) error {
		t.f = true
		return fn(ctx, key, opts...)
	}
}

func TestHook(t *testing.T) {
	h := &testHook{}

	s := NewStore(Hooks(HookExists(h.Exists)))

	if err := s.Init(); err != nil {
		t.Fatal(err)
	}

	err := s.Exists(context.TODO(), "test")
	if !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}

	if !h.f {
		t.Fatal("hook not works")
	}
}
