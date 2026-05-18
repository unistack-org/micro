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

func TestDefaultProfiler(t *testing.T) {
	if DefaultProfiler == nil {
		t.Fatal("expected DefaultProfiler to be non-nil")
	}
}

func TestNameOption(t *testing.T) {
	opts := Options{}
	Name("myprofile")(&opts)
	if opts.Name != "myprofile" {
		t.Fatalf("expected Name='myprofile', got %q", opts.Name)
	}
}

func TestNewProfilerWithOption(t *testing.T) {
	p := NewProfiler(Name("test"))
	if p == nil {
		t.Fatal("expected non-nil profiler")
	}
	if p.String() != "noop" {
		t.Fatalf("expected String()='noop', got %q", p.String())
	}
}
