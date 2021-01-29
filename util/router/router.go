package router

import (
	"github.com/unistack-org/micro/v3/register"
	"github.com/unistack-org/micro/v3/router"
)

type apiRouter struct {
	routes []router.Route
	router.Router
}

func (r *apiRouter) Lookup(...router.QueryOption) ([]router.Route, error) {
	return r.routes, nil
}

func (r *apiRouter) String() string {
	return "api"
}

// New router is a hack for API routing
func New(srvs []*register.Service) router.Router {
	var routes []router.Route

	for _, srv := range srvs {
		for _, n := range srv.Nodes {
			routes = append(routes, router.Route{Address: n.Address, Metadata: n.Metadata})
		}
	}

	return &apiRouter{routes: routes}
}
