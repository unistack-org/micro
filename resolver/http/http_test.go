package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.unistack.org/micro/v5/resolver"
)

func TestResolveHTTPSuccess(t *testing.T) {
	nodes := []*resolver.Record{
		{Address: "10.0.0.1:8080"},
		{Address: "10.0.0.2:8080"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Response{Nodes: nodes})
	}))
	defer srv.Close()

	r := &HTTPResolver{
		Proto: "http",
		Host:  srv.Listener.Addr().String(),
		Path:  "/network/nodes",
	}
	records, err := r.Resolve("any")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestResolveHTTPServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	r := &HTTPResolver{
		Proto: "http",
		Host:  srv.Listener.Addr().String(),
	}
	_, err := r.Resolve("any")
	if err == nil {
		t.Fatal("expected error from server 500")
	}
}

func TestResolveHTTPDefaults(t *testing.T) {
	// Test with default Proto/Path/Host (will fail to connect, just checking defaults)
	r := &HTTPResolver{}
	_, err := r.Resolve("any")
	// Connection refused expected — just verify it doesn't panic
	_ = err
}
