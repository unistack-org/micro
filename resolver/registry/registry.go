// Package registry resolves names using the micro registry
package registry

import (
	"context"

	"github.com/unistack-org/micro/v3/registry"
	"github.com/unistack-org/micro/v3/resolver"
)

// Resolver is a registry network resolver
type Resolver struct {
	// Registry is the registry to use otherwise we use the defaul
	Registry registry.Registry
}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	services, err := r.Registry.GetService(context.TODO(), name)
	if err != nil {
		return nil, err
	}

	records := make([]*resolver.Record, 0, len(services))
	for _, service := range services {
		for _, node := range service.Nodes {
			records = append(records, &resolver.Record{
				Address: node.Address,
			})
		}
	}

	return records, nil
}
