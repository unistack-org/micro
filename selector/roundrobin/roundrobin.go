package roundrobin // import "go.unistack.org/micro/v4/selector/roundrobin"

import (
	"go.unistack.org/micro/v4/selector"
	"go.unistack.org/micro/v4/util/rand"
)

// NewSelector returns an initialised round robin selector
func NewSelector(opts ...selector.Option) selector.Selector {
	return new(roundrobin)
}

type roundrobin struct{}

// Select return routes based on algo
func (r *roundrobin) Select(routes []string, opts ...selector.SelectOption) (selector.Next, error) {
	if len(routes) == 0 {
		return nil, selector.ErrNoneAvailable
	}
	var rng rand.Rand
	i := rng.Intn(len(routes))

	return func() string {
		route := routes[i%len(routes)]
		// increment
		i++
		return route
	}, nil
}

func (r *roundrobin) Record(addr string, err error) error { return nil }

func (r *roundrobin) Reset() error { return nil }

func (r *roundrobin) String() string {
	return "roundrobin"
}
