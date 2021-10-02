// Package vpath resolves using http path and recognised versioned urls
package vpath // import "go.unistack.org/micro/v3/api/resolver/vpath"

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"go.unistack.org/micro/v3/api/resolver"
)

// NewResolver creates new vpath api resolver
func NewResolver(opts ...resolver.Option) resolver.Resolver {
	return &vpathResolver{opts: resolver.NewOptions(opts...)}
}

type vpathResolver struct {
	opts resolver.Options
}

var re = regexp.MustCompile("^v[0-9]+$")

// Resolve endpoint
func (r *vpathResolver) Resolve(req *http.Request, opts ...resolver.ResolveOption) (*resolver.Endpoint, error) {
	if req.URL.Path == "/" {
		return nil, errors.New("unknown name")
	}

	options := resolver.NewResolveOptions(opts...)

	parts := strings.Split(req.URL.Path[1:], "/")
	if len(parts) == 1 {
		return &resolver.Endpoint{
			Name:   r.withPrefix(parts...),
			Host:   req.Host,
			Method: req.Method,
			Path:   req.URL.Path,
			Domain: options.Domain,
		}, nil
	}

	// /v1/foo
	if re.MatchString(parts[0]) {
		return &resolver.Endpoint{
			Name:   r.withPrefix(parts[0:2]...),
			Host:   req.Host,
			Method: req.Method,
			Path:   req.URL.Path,
			Domain: options.Domain,
		}, nil
	}

	return &resolver.Endpoint{
		Name:   r.withPrefix(parts[0]),
		Host:   req.Host,
		Method: req.Method,
		Path:   req.URL.Path,
		Domain: options.Domain,
	}, nil
}

func (r *vpathResolver) String() string {
	return "vpath"
}

// withPrefix transforms "foo" into "go.micro.api.foo"
func (r *vpathResolver) withPrefix(parts ...string) string {
	p := r.opts.ServicePrefix
	if len(p) > 0 {
		parts = append([]string{p}, parts...)
	}

	return strings.Join(parts, ".")
}
