package dns

import "testing"

func TestResolver(t *testing.T) {
	r := &Resolver{}
	recs, err := r.Resolve("unistack.org")
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) < 1 {
		t.Fatalf("records not resolved: %v", recs)
	}
}
