package api

import (
	"errors"
	"regexp"
	"strings"

	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/server"
)

type Api interface {
	// Initialise options
	Init(...Option) error
	// Get the options
	Options() Options
	// Register a http handler
	Register(*Endpoint) error
	// Register a route
	Deregister(*Endpoint) error
	// Implementation of api
	String() string
}

type Options struct{}

type Option func(*Options) error

// Endpoint is a mapping between an RPC method and HTTP endpoint
type Endpoint struct {
	// RPC Method e.g. Greeter.Hello
	Name string
	// Description e.g what's this endpoint for
	Description string
	// API Handler e.g rpc, proxy
	Handler string
	// HTTP Host e.g example.com
	Host []string
	// HTTP Methods e.g GET, POST
	Method []string
	// HTTP Path e.g /greeter. Expect POSIX regex
	Path []string
	// Body destination
	// "*" or "" - top level message value
	// "string" - inner message value
	Body string
	// Stream flag
	Stream bool
}

// Service represents an API service
type Service struct {
	// Name of service
	Name string
	// The endpoint for this service
	Endpoint *Endpoint
	// Versions of this service
	Services []*register.Service
}

// Encode encodes an endpoint to endpoint metadata
func Encode(e *Endpoint) map[string]string {
	if e == nil {
		return nil
	}

	// endpoint map
	ep := make(map[string]string)

	// set vals only if they exist
	set := func(k, v string) {
		if len(v) == 0 {
			return
		}
		ep[k] = v
	}

	set("endpoint", e.Name)
	set("description", e.Description)
	set("handler", e.Handler)
	set("method", strings.Join(e.Method, ","))
	set("path", strings.Join(e.Path, ","))
	set("host", strings.Join(e.Host, ","))
	set("body", e.Body)

	return ep
}

// Decode decodes endpoint metadata into an endpoint
func Decode(e metadata.Metadata) *Endpoint {
	if e == nil {
		return nil
	}

	ep := &Endpoint{}
	ep.Name, _ = e.Get("endpoint")
	ep.Description, _ = e.Get("description")
	epmethod, _ := e.Get("method")
	ep.Method = []string{epmethod}
	eppath, _ := e.Get("path")
	ep.Path = []string{eppath}
	ephost, _ := e.Get("host")
	ep.Host = []string{ephost}
	ep.Handler, _ = e.Get("handler")
	ep.Body, _ = e.Get("body")

	return ep
}

// Validate validates an endpoint to guarantee it won't blow up when being served
func Validate(e *Endpoint) error {
	if e == nil {
		return errors.New("endpoint is nil")
	}

	if len(e.Name) == 0 {
		return errors.New("name required")
	}

	for _, p := range e.Path {
		ps := p[0]
		pe := p[len(p)-1]

		if ps == '^' && pe == '$' {
			_, err := regexp.CompilePOSIX(p)
			if err != nil {
				return err
			}
		} else if ps == '^' && pe != '$' {
			return errors.New("invalid path")
		} else if ps != '^' && pe == '$' {
			return errors.New("invalid path")
		}
	}

	if len(e.Handler) == 0 {
		return errors.New("invalid handler")
	}

	return nil
}

/*
Design ideas

// Gateway is an api gateway interface
type Gateway interface {
	// Register a http handler
	Handle(pattern string, http.Handler)
	// Register a route
	RegisterRoute(r Route)
	// Init initialises the command line.
	// It also parses further options.
	Init(...Option) error
	// Run the gateway
	Run() error
}

// NewGateway returns a new api gateway
func NewGateway() Gateway {
	return newGateway()
}
*/

// WithEndpoint returns a server.HandlerOption with endpoint metadata set
//
// Usage:
//
// 	proto.RegisterHandler(service.Server(), new(Handler), api.WithEndpoint(
//		&api.Endpoint{
//			Name: "Greeter.Hello",
//			Path: []string{"/greeter"},
//		},
//	))
func WithEndpoint(e *Endpoint) server.HandlerOption {
	return server.EndpointMetadata(e.Name, Encode(e))
}
