// Package host resolves using http host
package host

import (
	"net/http"

	"github.com/unistack-org/micro/v3/api/resolver"
)

type hostResolver struct {
	opts resolver.Options
}

// Resolve endpoint
func (r *hostResolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	// parse options
	options := resolver.NewResolveOptions(opts...)

	return &resolver.Endpoint{
		Name:   req.Host,
		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Domain: options.Domain,
	}, nil
}

func (r *hostResolver) String() string {
	return "host"
}

// NewResolver creates new host api resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &hostResolver{opts: resolver.NewOptions(opts...)}
}
