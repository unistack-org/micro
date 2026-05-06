package metadata

import (
	"testing"
)

func TestNewClientWrapper(t *testing.T) {
	w := NewClientWrapper("key1", "key2")
	if w == nil {
		t.Error("expected non-nil wrapper")
	}
_ = w
}