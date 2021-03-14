package meter

import (
	"testing"
)

func TestNoopMeter(t *testing.T) {
	m := NewMeter(Path("/noop"))
	if "/noop" != m.Options().Path {
		t.Fatalf("invalid options parsing: %v", m.Options())
	}

	cnt := m.Counter("counter", Labels("server", "noop"))
	cnt.Inc()
}

func TestLabelsSort(t *testing.T) {
	ls := []string{"server", "http", "register", "mdns", "broker", "broker1", "broker", "broker2", "server", "tcp"}
	Sort(&ls)

	if ls[0] != "broker" || ls[1] != "broker2" {
		t.Fatalf("sort error: %v", ls)
	}
}
