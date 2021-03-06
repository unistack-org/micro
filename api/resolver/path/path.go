// Package path resolves using http path
package path

import (
	"net/http"
	"strings"

	"github.com/unistack-org/micro/v3/api/resolver"
)

// Resolver the path resolver
type Resolver struct {
	opts resolver.Options
}

// Resolve resolves endpoint
func (r *Resolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	// parse options
	options := resolver.NewResolveOptions(opts...)

	if req.URL.Path == "/" {
		return nil, resolver.ErrNotFound
	}

	parts := strings.Split(req.URL.Path[1:], "/")

	return &resolver.Endpoint{
		Name:   r.opts.ServicePrefix + "." + parts[0],
		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Domain: options.Domain,
	}, nil
}

// String retruns the string representation
func (r *Resolver) String() string {
	return "path"
}

// NewResolver returns new path resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{opts: resolver.NewOptions(opts...)}
}
