package roundrobin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unistack-org/micro/v3/selector"
)

func TestRoundRobin(t *testing.T) {
	selector.Tests(t, NewSelector())

	r1 := "127.0.0.1:8000"
	r2 := "127.0.0.1:8001"
	r3 := "127.0.0.1:8002"

	sel := NewSelector()

	// By passing r1 and r2 first, it forces a set sequence of (r1 => r2 => r3 => r1)

	next, err := sel.Select([]string{r1})
	r := next()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, r1, r, "Expected route to be r1")

	next, err = sel.Select([]string{r2})
	r = next()
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, r2, r, "Expected route to be r2")

	routes := []string{r1, r2, r3}
	next, err = sel.Select(routes)
	assert.Nil(t, err, "Error should be nil")
	n1, n2, n3, n4 := next(), next(), next(), next()

	// start element is random but then it should loop through in order
	start := -1
	for i := 0; i < 3; i++ {
		if n1 == routes[i] {
			start = i
			break
		}
	}
	assert.NotEqual(t, start, -1)
	assert.Equal(t, routes[start], n1, "Unexpected route")
	assert.Equal(t, routes[(start+1)%3], n2, "Unexpected route")
	assert.Equal(t, routes[(start+2)%3], n3, "Unexpected route")
	assert.Equal(t, routes[(start+3)%3], n4, "Unexpected route")
}
