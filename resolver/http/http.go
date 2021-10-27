// Package http resolves names to network addresses using a http request
package http // import "go.unistack.org/micro/v3/resolver/http"

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"go.unistack.org/micro/v3/resolver"
)

// nolint: golint,revive
// HTTPResolver is a HTTP network resolver
type HTTPResolver struct {
	// Proto if not set, defaults to http
	Proto string
	// Path sets the path to lookup. Defaults to /network
	Path string
	// Host url to use for the query
	Host string
}

// Response contains resolver.Record
type Response struct {
	Nodes []*resolver.Record `json:"nodes,omitempty"`
}

// Resolve assumes ID is a domain which can be converted to a http://name/network request
func (r *HTTPResolver) Resolve(name string) ([]*resolver.Record, error) {
	proto := "http"
	host := "localhost:8080"
	path := "/network/nodes"

	if len(r.Proto) > 0 {
		proto = r.Proto
	}

	if len(r.Path) > 0 {
		path = r.Path
	}

	if len(r.Host) > 0 {
		host = r.Host
	}

	uri := &url.URL{
		Scheme: proto,
		Path:   path,
		Host:   host,
	}
	q := uri.Query()
	q.Set("name", name)
	uri.RawQuery = q.Encode()

	// nolint: noctx
	rsp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return nil, errors.New("non 200 response")
	}
	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	// encoding format is assumed to be json
	var response *Response

	if err := json.Unmarshal(b, &response); err != nil {
		return nil, err
	}

	return response.Nodes, nil
}
