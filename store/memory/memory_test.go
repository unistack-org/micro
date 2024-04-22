package memory

import (
	"context"
	"testing"
	"time"

	"go.unistack.org/micro/v3/store"
)

type testHook struct {
	f bool
}

func (t *testHook) Exists(fn store.FuncExists) store.FuncExists {
	return func(ctx context.Context, key string, opts ...store.ExistsOption) error {
		t.f = true
		return fn(ctx, key, opts...)
	}
}

func TestHook(t *testing.T) {
	h := &testHook{}

	s := NewStore(store.Hooks(store.HookExists(h.Exists)))

	if err := s.Init(); err != nil {
		t.Fatal(err)
	}

	if err := s.Write(context.TODO(), "test", nil); err != nil {
		t.Error(err)
	}

	if err := s.Exists(context.TODO(), "test"); err != nil {
		t.Fatal(err)
	}

	if !h.f {
		t.Fatal("hook not works")
	}
}

func TestMemoryReInit(t *testing.T) {
	s := NewStore(store.Namespace("aaa"))
	if err := s.Init(store.Namespace("")); err != nil {
		t.Fatal(err)
	}
	if len(s.Options().Namespace) > 0 {
		t.Error("Init didn't reinitialise the store")
	}
}

func TestMemoryBasic(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func TestMemoryPrefix(t *testing.T) {
	s := NewStore()
	if err := s.Init(store.Namespace("some-prefix")); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func TestMemoryNamespace(t *testing.T) {
	s := NewStore()
	if err := s.Init(store.Namespace("some-namespace")); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func TestMemoryNamespacePrefix(t *testing.T) {
	s := NewStore()
	if err := s.Init(store.Namespace("some-namespace")); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func basictest(s store.Store, t *testing.T) {
	ctx := context.Background()
	// Read and Write an expiring Record
	if err := s.Write(ctx, "Hello", "World", store.WriteTTL(time.Millisecond*100)); err != nil {
		t.Error(err)
	}
	var val []byte
	if err := s.Read(ctx, "Hello", &val); err != nil {
		t.Error(err)
	} else if string(val) != "World" {
		t.Errorf("Expected %s, got %s", "World", val)
	}
	time.Sleep(time.Millisecond * 200)
	if err := s.Read(ctx, "Hello", &val); err != store.ErrNotFound {
		t.Errorf("Expected %# v, got %# v", store.ErrNotFound, err)
	}

	if err := s.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
}
