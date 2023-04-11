// Package register resolves names using the micro register
package register // import "go.unistack.org/micro/v4/resolver/registry"

import (
	"context"

	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/resolver"
)

// Resolver is a register network resolver
type Resolver struct {
	// Register is the register to use otherwise we use the defaul
	Register register.Register
}

// Resolve assumes ID is a domain name e.g micro.mu
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	services, err := r.Register.LookupService(context.TODO(), name)
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
