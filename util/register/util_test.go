package register

import (
	"os"
	"testing"

	"go.unistack.org/micro/v5/register"
)

func TestAddNodes(t *testing.T) {
	old := []*register.Node{
		{ID: "old-1", Address: "localhost:1111"},
	}
	neu := []*register.Node{
		{ID: "new-1", Address: "localhost:2222"},
	}
	result := addNodes(old, neu)
	if len(result) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(result))
	}
}

func TestAddNodesDedup(t *testing.T) {
	node := &register.Node{ID: "shared-1", Address: "localhost:1111"}
	result := addNodes([]*register.Node{node}, []*register.Node{node})
	if len(result) != 1 {
		t.Fatalf("expected 1 node after dedup, got %d", len(result))
	}
}

func TestCopyService(t *testing.T) {
	svc := &register.Service{
		Name:    "mysvc",
		Version: "1.0.0",
		Nodes: []*register.Node{
			{ID: "n1", Address: "localhost:9999"},
		},
	}
	copied := CopyService(svc)
	if copied.Name != svc.Name {
		t.Fatalf("expected name %q, got %q", svc.Name, copied.Name)
	}
	if len(copied.Nodes) != len(svc.Nodes) {
		t.Fatalf("expected %d nodes, got %d", len(svc.Nodes), len(copied.Nodes))
	}
	// Verify deep copy
	copied.Nodes[0].ID = "changed"
	if svc.Nodes[0].ID == "changed" {
		t.Fatal("CopyService not performing deep copy of nodes")
	}
}

func TestCopy(t *testing.T) {
	services := []*register.Service{
		{Name: "svc1", Version: "1.0.0", Nodes: []*register.Node{{ID: "n1", Address: "a:1"}}},
		{Name: "svc2", Version: "2.0.0", Nodes: []*register.Node{{ID: "n2", Address: "b:2"}}},
	}
	copied := Copy(services)
	if len(copied) != len(services) {
		t.Fatalf("expected %d services, got %d", len(services), len(copied))
	}
}

func TestMerge(t *testing.T) {
	olist := []*register.Service{
		{Name: "svc", Version: "1.0.0", Nodes: []*register.Node{{ID: "n1", Address: "a:1"}}},
	}
	nlist := []*register.Service{
		{Name: "svc", Version: "1.0.0", Nodes: []*register.Node{{ID: "n2", Address: "b:2"}}},
	}
	merged := Merge(olist, nlist)
	if len(merged) == 0 {
		t.Fatal("expected merged services, got none")
	}
	// Should have both nodes merged
	if len(merged[0].Nodes) != 2 {
		t.Fatalf("expected 2 nodes after merge, got %d", len(merged[0].Nodes))
	}
}

func TestMergeNewVersion(t *testing.T) {
	olist := []*register.Service{
		{Name: "svc", Version: "1.0.0", Nodes: []*register.Node{{ID: "n1", Address: "a:1"}}},
	}
	nlist := []*register.Service{
		{Name: "svc", Version: "2.0.0", Nodes: []*register.Node{{ID: "n2", Address: "b:2"}}},
	}
	merged := Merge(olist, nlist)
	// Both versions should appear
	if len(merged) == 0 {
		t.Fatal("expected merged services, got none")
	}
}

func TestRemove(t *testing.T) {
	services := []*register.Service{
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*register.Node{
				{
					ID:      "foo-123",
					Address: "localhost:9999",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*register.Node{
				{
					ID:      "foo-123",
					Address: "localhost:6666",
				},
			},
		},
	}

	servs := Remove([]*register.Service{services[0]}, []*register.Service{services[1]})
	if i := len(servs); i > 0 {
		t.Errorf("Expected 0 nodes, got %d: %+v", i, servs)
	}
	if len(os.Getenv("INTEGRATION_TESTS")) == 0 {
		t.Logf("Services %+v", servs)
	}
}

func TestRemoveNodes(t *testing.T) {
	services := []*register.Service{
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*register.Node{
				{
					ID:      "foo-123",
					Address: "localhost:9999",
				},
				{
					ID:      "foo-321",
					Address: "localhost:6666",
				},
			},
		},
		{
			Name:    "foo",
			Version: "1.0.0",
			Nodes: []*register.Node{
				{
					ID:      "foo-123",
					Address: "localhost:6666",
				},
			},
		},
	}

	nodes := delNodes(services[0].Nodes, services[1].Nodes)
	if i := len(nodes); i != 1 {
		t.Errorf("Expected only 1 node, got %d: %+v", i, nodes)
	}
	if len(os.Getenv("INTEGRATION_TESTS")) == 0 {
		t.Logf("Nodes %+v", nodes)
	}
}
