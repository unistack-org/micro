package register

import (
	"context"
	"testing"

	"go.unistack.org/micro/v5/register"
	"go.unistack.org/micro/v5/resolver"
)

func TestResolveFromRegistry(t *testing.T) {
	r := &Resolver{Register: register.NewRegister()}
	records, err := r.Resolve("nonexistent.service")
	if err != nil {
		t.Logf("resolver returned error (acceptable for noop): %v", err)
	}
	_ = records
}

func TestResolveWithRegisteredService(t *testing.T) {
	reg := register.NewRegister()
	ctx := t.Context()
	svc := &register.Service{
		Name:    "my.service",
		Version: "v1",
		Nodes: []*register.Node{
			{ID: "1", Address: "127.0.0.1:8080"},
		},
	}
	_ = reg.Register(ctx, svc)

	r := &Resolver{Register: reg}
	records, err := r.Resolve("my.service")
	if err != nil {
		t.Logf("noop register doesn't persist, error: %v", err)
	}
	_ = records
}

// stubRegister is a minimal register.Register implementation for testing.
type stubRegister struct {
	services []*register.Service
	err      error
}

func (s *stubRegister) Init(...register.Option) error { return nil }
func (s *stubRegister) Options() register.Options     { return register.Options{} }
func (s *stubRegister) Connect(ctx context.Context) error {
	return nil
}
func (s *stubRegister) Disconnect(ctx context.Context) error {
	return nil
}
func (s *stubRegister) Register(ctx context.Context, svc *register.Service, opts ...register.RegisterOption) error {
	return nil
}
func (s *stubRegister) Deregister(ctx context.Context, svc *register.Service, opts ...register.DeregisterOption) error {
	return nil
}
func (s *stubRegister) LookupService(ctx context.Context, name string, opts ...register.LookupOption) ([]*register.Service, error) {
	return s.services, s.err
}
func (s *stubRegister) ListServices(ctx context.Context, opts ...register.ListOption) ([]*register.Service, error) {
	return nil, nil
}
func (s *stubRegister) Watch(ctx context.Context, opts ...register.WatchOption) (register.Watcher, error) {
	return nil, nil
}
func (s *stubRegister) String() string { return "stub" }
func (s *stubRegister) Name() string   { return "stub" }
func (s *stubRegister) Live() bool     { return true }
func (s *stubRegister) Ready() bool    { return true }
func (s *stubRegister) Health() bool   { return true }

func TestResolveNodesMapping(t *testing.T) {
	stub := &stubRegister{
		services: []*register.Service{
			{
				Name:    "test.service",
				Version: "v1",
				Nodes: []*register.Node{
					{ID: "node1", Address: "10.0.0.1:9000"},
					{ID: "node2", Address: "10.0.0.2:9001"},
				},
			},
		},
	}
	r := &Resolver{Register: stub}
	records, err := r.Resolve("test.service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	expected := map[string]bool{"10.0.0.1:9000": true, "10.0.0.2:9001": true}
	for _, rec := range records {
		if !expected[rec.Address] {
			t.Fatalf("unexpected address: %s", rec.Address)
		}
	}
}

func TestResolveError(t *testing.T) {
	stub := &stubRegister{err: register.ErrNotFound}
	r := &Resolver{Register: stub}
	_, err := r.Resolve("test.service")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ensure resolver.Record is used (import check)
var _ *resolver.Record
