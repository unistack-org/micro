// Package grpc resolves a grpc service like /greeter.Say/Hello to greeter service
package grpc // import "go.unistack.org/micro/v3/api/resolver/grpc"

import (
	"errors"
	"net/http"
	"strings"

	"go.unistack.org/micro/v3/api/resolver"
)

// Resolver struct
type Resolver struct {
	opts resolver.Options
}

// Resolve func to resolve enndpoint
func (r *Resolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	// parse options
	options := resolver.NewResolveOptions(opts...)

	// /foo.Bar/Service
	if req.URL.Path == "/" {
		return nil, errors.New("unknown name")
	}
	// [foo.Bar, Service]
	parts := strings.Split(req.URL.Path[1:], "/")
	// [foo, Bar]
	name := strings.Split(parts[0], ".")
	// foo
	return &resolver.Endpoint{
		Name:   strings.Join(name[:len(name)-1], "."),
		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Domain: options.Domain,
	}, nil
}

func (r *Resolver) String() string {
	return "grpc"
}

// NewResolver is used to create new Resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &Resolver{opts: resolver.NewOptions(opts...)}
}
