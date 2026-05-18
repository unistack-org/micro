package noop

import "testing"

func TestResolve(t *testing.T) {
	r := &Resolver{}
	records, err := r.Resolve("test.service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty records, got %d", len(records))
	}
}
