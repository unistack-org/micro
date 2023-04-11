// Package noop is a noop resolver
package noop // import "go.unistack.org/micro/v4/resolver/noop"

import (
	"go.unistack.org/micro/v4/resolver"
)

// Resolver contains noop resolver
type Resolver struct{}

// Resolve returns the list of nodes
func (r *Resolver) Resolve(name string) ([]*resolver.Record, error) {
	return []*resolver.Record{}, nil
}
