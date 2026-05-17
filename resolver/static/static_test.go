package static

import "testing"

func TestResolveEmpty(t *testing.T) {
	r := &Resolver{}
	records, err := r.Resolve("my.service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 || records[0].Address != "my.service" {
		t.Fatalf("expected fallback to name, got %v", records)
	}
}

func TestResolveWithNodes(t *testing.T) {
	r := &Resolver{
		Nodes: []string{"192.168.1.1:8080", "192.168.1.2:8080"},
	}
	records, err := r.Resolve("ignored")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Address != "192.168.1.1:8080" {
		t.Fatalf("unexpected address: %s", records[0].Address)
	}
}
