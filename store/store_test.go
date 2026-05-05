package store

import (
	"context"
	"testing"
)

func TestNoopStore_Name(t *testing.T) {
	s := NewStore()
	if name := s.String(); name != "noop" {
		t.Errorf("expected 'noop', got %q", name)
	}
}

func TestNoopStore_Options(t *testing.T) {
	s := NewStore()
	opts := s.Options()
	if opts.Context == nil {
		t.Error("expected non-nil context in options")
	}
}

func TestNoopStore_String(t *testing.T) {
	s := NewStore()
	if str := s.String(); str != "noop" {
		t.Errorf("expected 'noop', got %q", str)
	}
}

func TestNoopStore_Init(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNoopStore_ConnectDisconnect(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()
	if err := s.Connect(ctx); err != nil {
		t.Errorf("unexpected error on Connect: %v", err)
	}
	if err := s.Disconnect(ctx); err != nil {
		t.Errorf("unexpected error on Disconnect: %v", err)
	}
}

func TestNoopStore_ReadWriteDelete(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()

	if err := s.Write(ctx, "key", "value"); err != nil {
		t.Errorf("unexpected error on Write: %v", err)
	}

	var val string
	if err := s.Read(ctx, "key", &val); err != nil {
		t.Errorf("unexpected error on Read: %v", err)
	}

	if err := s.Delete(ctx, "key"); err != nil {
		t.Errorf("unexpected error on Delete: %v", err)
	}
}

func TestNoopStore_Exists(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()
	if err := s.Exists(ctx, "nonexistent"); err == nil {
		t.Error("expected error for nonexistent key")
	}
}

func TestNoopStore_List(t *testing.T) {
	s := NewStore()
	if err := s.Init(); err != nil {
		t.Errorf("unexpected error on Init: %v", err)
	}
	ctx := context.Background()
	list, err := s.List(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if list != nil {
		t.Error("expected nil list")
	}
}

func TestNoopStore_LiveReadyHealth(t *testing.T) {
	s := NewStore()
	if !s.Live() {
		t.Error("expected store to be live")
	}
	if !s.Ready() {
		t.Error("expected store to be ready")
	}
	if !s.Health() {
		t.Error("expected store to be healthy")
	}
}
