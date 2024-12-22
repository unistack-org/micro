package id

import "testing"

func TestUUIDv8(t *testing.T) {
	id, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("xxx %s\n", id)
}
