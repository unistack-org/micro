package random

import (
	"testing"

	"go.unistack.org/micro/v5/selector"
)

func TestRandom(t *testing.T) {
	selector.Tests(t, NewSelector())
}

func TestReset(t *testing.T) {
	s := NewSelector()
	if err := s.Reset(); err != nil {
		t.Fatalf("Reset returned error: %v", err)
	}
}
