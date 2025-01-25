package http

import (
	"errors"
	"net/http"

	"go.unistack.org/micro/v4/router"
	"go.unistack.org/micro/v4/selector"
)

type roundTripper struct {
	rt   http.RoundTripper
	st   selector.Selector
	opts Options
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	routes, err := r.opts.Router.Lookup(router.QueryService(req.URL.Host))
	if err != nil {
		return nil, err
	}

	// rudimentary retry 3 times
	for _, route := range routes {
		req.URL.Host = route.Address
		w, err := r.rt.RoundTrip(req)
		if err != nil {
			continue
		}
		return w, nil
	}

	return nil, errors.New("failed request")
}
