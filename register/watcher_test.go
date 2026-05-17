package register

import (
	"context"
	"testing"
)

func TestWatcherStop(t *testing.T) {
	r := NewRegister()
	w, err := r.Watch(context.Background())
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	// Stop should not panic
	w.Stop()
}

func TestWatcherNext(t *testing.T) {
	r := NewRegister()
	w, err := r.Watch(context.Background())
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	// noop watcher Next returns nil, nil
	result, err := w.Next()
	if err != nil {
		t.Fatalf("Next returned unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result from noop watcher, got %v", result)
	}
}

func TestWatcherStopThenNext(t *testing.T) {
	r := NewRegister()
	w, err := r.Watch(context.Background())
	if err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}
	w.Stop()
	// After Stop, Next may return nil or an error — both are acceptable
	_, _ = w.Next()
}

func TestEventTypeString(t *testing.T) {
	cases := []struct {
		ev   EventType
		want string
	}{
		{EventCreate, "create"},
		{EventDelete, "delete"},
		{EventUpdate, "update"},
		{EventType(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.ev.String(); got != c.want {
			t.Fatalf("EventType(%d).String() = %q, want %q", c.ev, got, c.want)
		}
	}
}
