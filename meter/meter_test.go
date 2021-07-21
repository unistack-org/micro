package meter

import (
	"testing"
)

func TestNoopMeter(t *testing.T) {
	m := NewMeter(Path("/noop"))
	if "/noop" != m.Options().Path {
		t.Fatalf("invalid options parsing: %v", m.Options())
	}

	cnt := m.Counter("counter", "server", "noop")
	cnt.Inc()
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestBuildLabels(t *testing.T) {
	type testData struct {
		src []string
		dst []string
	}

	data := []testData{
		testData{
			src: []string{"zerolabel", "value3", "firstlabel", "value2"},
			dst: []string{"firstlabel", "value2", "zerolabel", "value3"},
		},
	}

	for _, d := range data {
		if !testEq(d.dst, BuildLabels(d.src...)) {
			t.Fatalf("slices not properly sorted: %v %v", d.dst, d.src)
		}
	}
}

func TestBuildName(t *testing.T) {
	data := map[string][]string{
		`my_metric{firstlabel="value2",zerolabel="value3"}`: []string{
			"my_metric",
			"zerolabel", "value3", "firstlabel", "value2",
		},
		`my_metric{broker="broker2",register="mdns",server="tcp"}`: []string{
			"my_metric",
			"broker", "broker1", "broker", "broker2", "server", "http", "server", "tcp", "register", "mdns",
		},
		`my_metric{aaa="aaa"}`: []string{
			"my_metric",
			"aaa", "aaa",
		},
	}

	for e, d := range data {
		if x := BuildName(d[0], d[1:]...); x != e {
			t.Fatalf("expect: %s, result: %s", e, x)
		}
	}
}
