package roundrobin

import (
	"testing"

	"go.unistack.org/micro/v4/selector"
)

func TestRoundRobin(t *testing.T) {
	selector.Tests(t, NewSelector())

	r1 := "127.0.0.1:8000"
	r2 := "127.0.0.1:8001"
	r3 := "127.0.0.1:8002"

	sel := NewSelector()

	// By passing r1 and r2 first, it forces a set sequence of (r1 => r2 => r3 => r1)

	next, err := sel.Select([]string{r1})
	if err != nil {
		t.Fatal(err)
	}
	r := next()

	if r1 != r {
		t.Fatal("Expected route to be r == r1")
	}

	next, err = sel.Select([]string{r2})
	if err != nil {
		t.Fatal(err)
	}
	r = next()
	if r2 != r {
		t.Fatal("Expected route to be r2")
	}

	routes := []string{r1, r2, r3}
	next, err = sel.Select(routes)
	if err != nil {
		t.Fatal(err)
	}
	n1, n2, n3, n4 := next(), next(), next(), next()

	// start element is random but then it should loop through in order
	start := -1
	for i := 0; i < 3; i++ {
		if n1 == routes[i] {
			start = i
			break
		}
	}
	if start == -1 {
		t.Fatalf("start == -1 %v %v", start, -1)
	}
	if routes[start] != n1 {
		t.Fatal("Unexpected route")
	}
	if routes[(start+1)%3] != n2 {
		t.Fatal("Unexpected route")
	}
	if routes[(start+2)%3] != n3 {
		t.Fatal("Unexpected route")
	}
	if routes[(start+3)%3] != n4 {
		t.Fatal("Unexpected route")
	}
}
