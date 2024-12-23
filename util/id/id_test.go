package id

import "testing"

func TestUUIDv8(t *testing.T) {
	id, err := New()
	if err != nil {
		t.Fatal(err)
	}
	_ = id
}

func TestToUUID(t *testing.T) {
	id, err := New()
	if err != nil {
		t.Fatal(err)
	}
	u := ToUUID(id)
	_ = u
}
