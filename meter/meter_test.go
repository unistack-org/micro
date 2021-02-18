package meter

import (
	"testing"
)

func TestNoopMeter(t *testing.T) {
	m := NewMeter(Path("/noop"))
	if "/noop" != m.Options().Path {
		t.Fatalf("invalid options parsing: %v", m.Options())
	}

	cnt := m.Counter("counter", Label("server", "noop"))
	cnt.Inc()
}

func TestLabelsAppend(t *testing.T) {
	var ls Labels
	ls.keys = []string{"type", "server"}
	ls.vals = []string{"noop", "http"}

	var nls Labels
	nls.keys = []string{"register"}
	nls.vals = []string{"gossip"}
	ls = ls.Append(nls)

	//ls.Sort()

	if ls.keys[0] != "type" || ls.vals[0] != "noop" {
		t.Fatalf("append error: %v", ls)
	}
}

func TestIterator(t *testing.T) {
	options := NewOptions(
		Label("name", "svc1"),
		Label("version", "0.0.1"),
		Label("id", "12345"),
		Label("type", "noop"),
		Label("server", "http"),
		Label("register", "gossip"),
		Label("aa", "kk"),
		Label("zz", "kk"),
	)

	iter := options.Labels.Iter()
	var k, v string
	cnt := 0
	for iter.Next(&k, &v) {
		if cnt == 4 && (k != "server" || v != "http") {
			t.Fatalf("iter error: %s != %s || %s != %s", k, "server", v, "http")
		}
		cnt++
	}
}
