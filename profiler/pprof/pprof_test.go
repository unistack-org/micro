package pprof

import (
	"testing"

	profile "go.unistack.org/micro/v5/profiler"
)

func TestNewProfile(t *testing.T) {
	p := NewProfile()
	if p == nil {
		t.Fatal("expected non-nil profiler")
	}
	if p.String() != "pprof" {
		t.Fatalf("expected String()='pprof', got %q", p.String())
	}
}

func TestNewProfileWithName(t *testing.T) {
	p := NewProfile(profile.Name("testprofile"))
	if p == nil {
		t.Fatal("expected non-nil profiler")
	}
}

func TestStartStop(t *testing.T) {
	p := NewProfile(profile.Name("test_coverage"))
	if err := p.Start(); err != nil {
		t.Fatalf("Start returned error: %v", err)
	}
	// Starting again should be a no-op (already running)
	if err := p.Start(); err != nil {
		t.Fatalf("Second Start returned error: %v", err)
	}
	if err := p.Stop(); err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
	// Stopping again should be a no-op
	if err := p.Stop(); err != nil {
		t.Fatalf("Second Stop returned error: %v", err)
	}
}

func TestString(t *testing.T) {
	p := NewProfile()
	if p.String() != "pprof" {
		t.Fatalf("expected 'pprof', got %q", p.String())
	}
}
