package meter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopMeter(t *testing.T) {
	meter := NewMeter(Path("/noop"))
	assert.NotNil(t, meter)
	assert.Equal(t, "/noop", meter.Options().Path)
	assert.Implements(t, new(Meter), meter)

	cnt := meter.Counter("counter", Label("server", "noop"))
	cnt.Inc()

}

func TestLabels(t *testing.T) {
	var ls Labels
	ls.keys = []string{"type", "server"}
	ls.vals = []string{"noop", "http"}

	ls.Sort()

	if ls.keys[0] != "server" || ls.vals[0] != "http" {
		t.Fatalf("sort error: %v", ls)
	}
}

func TestLabelsAppend(t *testing.T) {
	var ls Labels
	ls.keys = []string{"type", "server"}
	ls.vals = []string{"noop", "http"}

	var nls Labels
	nls.keys = []string{"register"}
	nls.vals = []string{"gossip"}
	ls = ls.Append(nls)

	ls.Sort()

	if ls.keys[0] != "register" || ls.vals[0] != "gossip" {
		t.Fatalf("append error: %v", ls)
	}
}

func TestIterator(t *testing.T) {
	var ls Labels
	ls.keys = []string{"type", "server", "register"}
	ls.vals = []string{"noop", "http", "gossip"}

	iter := ls.Iter()
	var k, v string
	cnt := 0
	for iter.Next(&k, &v) {
		if cnt == 1 && (k != "server" || v != "http") {
			t.Fatalf("iter error: %s != %s || %s != %s", k, "server", v, "http")
		}
		cnt++
	}
}
