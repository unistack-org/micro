package tracer

import (
	"testing"
)

func TestUniqLabels(t *testing.T) {
	labels := []interface{}{"key1", "val1", "key1", "val2"}
	labels = UniqLabels(labels)
	if labels[1] != "val2" {
		t.Fatalf("UniqLabels not works")
	}
}
