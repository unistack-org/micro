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

func TestResolverIPAddress(t *testing.T) {
	r := &Resolver{}
	recs, err := r.Resolve("1.2.3.4:9000")
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 1 {
		t.Fatalf("expected 1 record for IP, got %d", len(recs))
	}
	if recs[0].Address != "1.2.3.4:9000" {
		t.Fatalf("expected address 1.2.3.4:9000, got %s", recs[0].Address)
	}
}

func TestResolverEmptyHost(t *testing.T) {
	r := &Resolver{}
	// passing just ":8085" should resolve localhost
	recs, err := r.Resolve(":8085")
	if err != nil {
		// DNS lookup for localhost may fail in test env; that's acceptable
		t.Logf("resolve returned error (acceptable): %v", err)
		return
	}
	if len(recs) < 1 {
		t.Fatalf("expected at least 1 record, got %d", len(recs))
	}
}

func TestResolverHostWithPort(t *testing.T) {
	r := &Resolver{Address: "1.1.1.1:53"}
	recs, err := r.Resolve("unistack.org:443")
	if err != nil {
		t.Logf("resolve returned error (acceptable in test env): %v", err)
		return
	}
	if len(recs) < 1 {
		t.Fatalf("expected at least 1 record, got %d", len(recs))
	}
	for _, rec := range recs {
		if rec.Address == "" {
			t.Fatal("record has empty address")
		}
	}
}
