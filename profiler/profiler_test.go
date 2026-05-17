package profiler

import "testing"

func TestNoopProfiler(t *testing.T) {
	p := NewProfiler()

	if err := p.Start(); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
	if p.String() != "noop" {
		t.Fatalf("expected String()='noop', got %q", p.String())
	}
}
