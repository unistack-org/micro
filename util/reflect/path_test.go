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
