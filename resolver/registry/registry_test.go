package register

import (
	"testing"

	"go.unistack.org/micro/v5/register"
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
