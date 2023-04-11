// Package proxy is a transparent proxy built on the micro/server
package proxy // import "go.unistack.org/micro/v4/proxy"

import (
	"context"

	"go.unistack.org/micro/v4/server"
)

// DefaultEndpoint holds default proxy address
var DefaultEndpoint = "localhost:9090"

// Proxy can be used as a proxy server for micro services
type Proxy interface {
	// ProcessMessage handles inbound messages
	ProcessMessage(context.Context, server.Message) error
	// ServeRequest handles inbound requests
	ServeRequest(context.Context, server.Request, server.Response) error
	// Name of the proxy protocol
	String() string
}
