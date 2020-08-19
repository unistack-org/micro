// Package noop is a noop resolver
package noop

import (
	"github.com/unistack-org/micro/v3/resolver"
)

type Resolver struct{}

// Resolve returns the list of nodes
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	return []*resolver.Record{}, nil
}
