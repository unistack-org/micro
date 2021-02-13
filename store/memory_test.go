package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/unistack-org/micro/v3/store"
)

func TestMemoryReInit(t *testing.T) {
	s := store.NewStore(store.Table("aaa"))
	if err := s.Init(store.Table("")); err != nil {
		t.Fatal(err)
	}
	if len(s.Options().Table) > 0 {
		t.Error("Init didn't reinitialise the store")
	}
}

func TestMemoryBasic(t *testing.T) {
	s := store.NewStore()
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func TestMemoryPrefix(t *testing.T) {
	s := store.NewStore()
	if err := s.Init(store.Table("some-prefix")); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func TestMemoryNamespace(t *testing.T) {
	s := store.NewStore()
	if err := s.Init(store.Database("some-namespace")); err != nil {
		t.Fatal(err)
	}
	basictest(s, t)
}

func TestMemoryNamespacePrefix(t *testing.T) {
	s := store.NewStore()
	if err := s.Init(store.Table("some-prefix"), store.Database("some-namespace")); err != nil {
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
