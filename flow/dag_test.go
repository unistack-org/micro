package flow

import (
	"fmt"
	"testing"

	"github.com/silas/dag"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestDag(t *testing.T) {
	d1 := &dag.AcyclicGraph{}
	d2 := &dag.AcyclicGraph{}
	d2v1 := d2.Add(&node{"Substep.Create"})
	v1 := d1.Add(&node{"AccountService.Create"})
	v2 := d1.Add(&node{"AuthzService.Create"})
	v3 := d1.Add(&node{"AuthnService.Create"})
	v4 := d1.Add(&node{"ProjectService.Create"})
	v5 := d1.Add(&node{"ContactService.Create"})
	v6 := d1.Add(&node{"NetworkService.Create"})
	v7 := d1.Add(&node{"MailerService.Create"})
	v8 := d1.Add(&node{"NestedService.Create"})
	v9 := d1.Add(d2v1)
	d1.Connect(dag.BasicEdge(v1, v2))
	d1.Connect(dag.BasicEdge(v1, v3))
	d1.Connect(dag.BasicEdge(v1, v4))
	d1.Connect(dag.BasicEdge(v1, v5))
	d1.Connect(dag.BasicEdge(v1, v6))
	d1.Connect(dag.BasicEdge(v1, v7))
	d1.Connect(dag.BasicEdge(v7, v8))
	d1.Connect(dag.BasicEdge(v8, v9))

	if err := d1.Validate(); err != nil {
		t.Fatal(err)
	}

	d1.TransitiveReduction()

	var steps [][]string
	fn := func(n dag.Vertex, idx int) error {
		if idx == 0 {
			steps = make([][]string, 1, 1)
			steps[0] = make([]string, 0, 1)
		} else if idx >= len(steps) {
			tsteps := make([][]string, idx+1, idx+1)
			copy(tsteps, steps)
			steps = tsteps
			steps[idx] = make([]string, 0, 1)
		}
		steps[idx] = append(steps[idx], fmt.Sprintf("%s", n))
		return nil
	}

	start := &node{"AccountService.Create"}
	err := d1.SortedDepthFirstWalk([]dag.Vertex{start}, fn)
	checkErr(t, err)
	if len(steps) != 4 {
		t.Fatalf("invalid steps: %#+v", steps)
	}
	if steps[3][0] != "Substep.Create" {
		t.Fatalf("invalid last step: %#+v", steps)
	}
}
