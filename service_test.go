package micro

import (
	"testing"
)

type testItem struct {
	name string
}

func (ti *testItem) Name() string {
	return ti.name
}

func TestGetNameIndex(t *testing.T) {
	item1 := &testItem{name: "first"}
	item2 := &testItem{name: "second"}
	items := []interface{}{item1, item2}
	if idx := getNameIndex("second", items); idx != 1 {
		t.Fatalf("getNameIndex func error, item not found")
	}
}
