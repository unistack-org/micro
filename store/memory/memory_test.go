package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.unistack.org/micro/v5/metadata"
	"go.unistack.org/micro/v5/store"
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

func TestMemoryStore_ReadWrite(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()

	if err := s.Write(ctx, "key1", "value1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var val string
	if err := s.Read(ctx, "key1", &val); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got %q", val)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()

	_ = s.Write(ctx, "key1", "value1")
	if err := s.Delete(ctx, "key1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var val string
	if err := s.Read(ctx, "key1", &val); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_List(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()

	_ = s.Write(ctx, "key1", "value1")
	_ = s.Write(ctx, "key2", "value2")

	list, err := s.List(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
}

func TestMemoryStore_LazyConnect(t *testing.T) {
	s := NewStore(store.LazyConnect(true))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// Read should trigger connect
	if err := s.Read(ctx, "key", nil); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	// Write should trigger connect
	if err := s.Write(ctx, "key", "val"); err != nil {
		t.Fatal(err)
	}
	// Read should work
	var val string
	if err := s.Read(ctx, "key", &val); err != nil {
		t.Fatal(err)
	}
	if val != "val" {
		t.Errorf("expected 'val', got %q", val)
	}
}

func TestMemoryStore_ListSuffix(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := s.Write(ctx, "keysuffix", "val1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, "other", "val2"); err != nil {
		t.Fatal(err)
	}
	list, err := s.List(ctx, store.ListSuffix("suffix"))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0] != "keysuffix" {
		t.Errorf("expected [keysuffix], got %v", list)
	}
}

func TestMemoryStore_ListLimitOffset(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if err := s.Write(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i)); err != nil {
			t.Fatal(err)
		}
	}
	// Limit 2, Offset 1 → keys 1,2
	list, err := s.List(ctx, store.ListLimit(2), store.ListOffset(1))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 items, got %d", len(list))
	}
	if list[0] != "key1" || list[1] != "key2" {
		t.Errorf("unexpected list: %v", list)
	}
}

func TestMemoryWatcherStop(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	w, err := s.Watch(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Stop should not panic
	w.Stop()
}

func TestMemoryStore_Connect(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := s.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := s.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
	// Write a key, disconnect, check if it's gone
	if err := s.Write(ctx, "key", "val"); err != nil {
		t.Fatal(err)
	}
	if err := s.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
	var val string
	if err := s.Read(ctx, "key", &val); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound after disconnect, got %v", err)
	}
}

func TestMemoryStore_String(t *testing.T) {
	s := NewStore()
	if s.String() != "memory" {
		t.Errorf("expected 'memory', got %q", s.String())
	}
}

func TestMemoryStore_Name(t *testing.T) {
	s := NewStore()
	if s.Name() != "" {
		t.Errorf("expected empty name, got %q", s.Name())
	}
	s2 := NewStore(store.Name("mem"))
	if s2.Name() != "mem" {
		t.Errorf("expected 'mem', got %q", s2.Name())
	}
}

func TestMemoryStore_LiveReadyHealth(t *testing.T) {
	s := NewStore()
	if !s.Live() {
		t.Error("expected Live() to return true")
	}
	if !s.Ready() {
		t.Error("expected Ready() to return true")
	}
	if !s.Health() {
		t.Error("expected Health() to return true")
	}
}

func TestMemoryStore_Watch(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	w, err := s.Watch(ctx)
	if err != nil {
		t.Fatalf("unexpected error on Watch: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
	evt, err := w.Next()
	if err != nil {
		t.Errorf("unexpected error on Next: %v", err)
	}
	if evt != nil {
		t.Errorf("expected nil event, got %v", evt)
	}
	w.Stop()
}

func TestMemoryStore_ReadWriteWithOptions(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// Write with TTL
	if err := s.Write(ctx, "key1", "val1", store.WriteTTL(time.Millisecond*100)); err != nil {
		t.Fatal(err)
	}
	var val string
	if err := s.Read(ctx, "key1", &val); err != nil {
		t.Fatal(err)
	}
	if val != "val1" {
		t.Errorf("expected 'val1', got %q", val)
	}
	time.Sleep(time.Millisecond * 200)
	if err := s.Read(ctx, "key1", &val); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound after TTL, got %v", err)
	}
	// Write with Metadata
	md := metadata.Metadata{"foo": []string{"bar"}}
	if err := s.Write(ctx, "key2", "val2", store.WriteMetadata(md)); err != nil {
		t.Fatal(err)
	}
}

func TestMemoryStore_ListPrefix(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err := s.Write(ctx, "pre1", "val1"); err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, "pre2", "val2"); err != nil {
		t.Fatal(err)
	}
	if err := s.Write(ctx, "other", "val3"); err != nil {
		t.Fatal(err)
	}
	list, err := s.List(ctx, store.ListPrefix("pre"))
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}
	hasPre1 := false
	hasPre2 := false
	for _, k := range list {
		if k == "pre1" {
			hasPre1 = true
		}
		if k == "pre2" {
			hasPre2 = true
		}
	}
	if !hasPre1 || !hasPre2 {
		t.Error("expected pre1 and pre2 in list")
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			if err := s.Write(ctx, key, fmt.Sprintf("val%d", i)); err != nil {
				t.Errorf("write error: %v", err)
			}
			var val string
			if err := s.Read(ctx, key, &val); err != nil {
				t.Errorf("read error: %v", err)
			}
		}(i)
	}
	wg.Wait()
}

func TestMemoryStore_Exists(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	// Non-existing key
	if err := s.Exists(ctx, "nonexistent"); err != store.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	// Existing key
	if err := s.Write(ctx, "exists", "val"); err != nil {
		t.Fatal(err)
	}
	if err := s.Exists(ctx, "exists"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// With ExistsNamespace
	s2 := NewStore(store.Namespace("ns"))
	if err := s2.Init(); err != nil {
		t.Fatal(err)
	}
	if err := s2.Write(ctx, "key", "val"); err != nil {
		t.Fatal(err)
	}
	if err := s2.Exists(ctx, "key", store.ExistsNamespace("ns")); err != nil {
		t.Errorf("unexpected error with namespace: %v", err)
	}
}
