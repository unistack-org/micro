package reflect

import (
	"testing"
)

func TestLookup(t *testing.T) {
	type Nested2 struct {
		Name string
	}
	type Nested1 struct {
		Nested2 Nested2
	}
	type Config struct {
		Nested1 Nested1
	}

	cfg := &Config{
		Nested1: Nested1{
			Nested2: Nested2{
				Name: "NAME",
			},
		},
	}

	v, err := Lookup(cfg, "$.Nested1.Nested2.Name")
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "NAME" {
		t.Fatalf("lookup returns invalid value: %v", v)
	}
}

func TestLookupRoot(t *testing.T) {
	cfg := map[string]string{"key": "val"}
	v, err := Lookup(cfg, "$")
	if err != nil {
		t.Fatal(err)
	}
	if !v.IsValid() {
		t.Fatal("expected valid value for root lookup")
	}
}

func TestLookupBadPath(t *testing.T) {
	cfg := map[string]string{"key": "val"}
	_, err := Lookup(cfg, "bad_path")
	if err == nil {
		t.Fatal("expected error for bad path")
	}
}

func TestLookupMap(t *testing.T) {
	cfg := map[string]interface{}{
		"name": "hello",
	}
	v, err := Lookup(cfg, "$.name")
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "hello" {
		t.Fatalf("expected hello got %v", v)
	}
}

func TestLookupSliceAggregate(t *testing.T) {
	type Item struct {
		Name string
	}
	items := []Item{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	v, err := Lookup(items, "$.Name")
	if err != nil {
		t.Fatal(err)
	}
	if v.Len() != 3 {
		t.Fatalf("expected 3 items got %d", v.Len())
	}
}

func TestLookupSliceIndex(t *testing.T) {
	type Config struct {
		Items []string
	}
	// must use a value (not pointer) so parseIndex path is taken without pointer aggregate
	cfg := Config{Items: []string{"a", "b", "c"}}
	v, err := Lookup(cfg, "$.Items[1]")
	if err != nil {
		t.Fatal(err)
	}
	if v.String() != "b" {
		t.Fatalf("expected b got %v", v)
	}
}

func TestLookupMalformedIndex(t *testing.T) {
	type Config struct {
		Items []string
	}
	cfg := &Config{Items: []string{"a", "b"}}
	_, err := Lookup(cfg, "$.Items[")
	if err == nil {
		t.Fatal("expected error for malformed index")
	}
}

func TestParseIndexNoIndex(t *testing.T) {
	key, idx, err := parseIndex("name")
	if err != nil {
		t.Fatal(err)
	}
	if key != "name" || idx != -1 {
		t.Fatalf("expected name,-1 got %s,%d", key, idx)
	}
}

func TestParseIndexValid(t *testing.T) {
	key, idx, err := parseIndex("items[2]")
	if err != nil {
		t.Fatal(err)
	}
	if key != "items" || idx != 2 {
		t.Fatalf("expected items,2 got %s,%d", key, idx)
	}
}

func TestHasIndex(t *testing.T) {
	if !hasIndex("items[0]") {
		t.Fatal("expected hasIndex true")
	}
	if hasIndex("items") {
		t.Fatal("expected hasIndex false")
	}
}

func TestIsAggregable(t *testing.T) {
	// Exercise isAggregable indirectly via Lookup on a slice
	type Item struct{ Val int }
	items := []Item{{Val: 1}, {Val: 2}}
	_, err := Lookup(items, "$.Val")
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsMergeable(t *testing.T) {
	// isMergeable is exercised internally by mergeValue via aggreateAggregableValue;
	// trigger it with a slice-of-slices lookup
	type Inner struct{ Tags []string }
	items := []Inner{{Tags: []string{"a", "b"}}, {Tags: []string{"c"}}}
	v, err := Lookup(items, "$.Tags")
	if err != nil {
		t.Fatal(err)
	}
	if v.Len() != 3 {
		t.Fatalf("expected 3 merged tags got %d", v.Len())
	}
}

func TestLookupEmptySliceAggregate(t *testing.T) {
	type Item struct {
		Name string
	}
	items := []Item{}
	v, err := Lookup(items, "$.Name")
	if err != nil {
		t.Fatal(err)
	}
	if v.Len() != 0 {
		t.Fatalf("expected 0 items got %d", v.Len())
	}
}
