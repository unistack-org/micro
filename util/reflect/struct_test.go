package reflect

import (
	"net/url"
	"testing"
)

func TestStructByTag(t *testing.T) {
	type Str struct {
		Name []string `json:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	iface, err := StructFieldByTag(val, "codec", "flatten")
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := iface.([]string); !ok {
		t.Fatalf("not []string %v", iface)
	} else if len(v) != 2 {
		t.Fatalf("invalid number %v", iface)
	}
}

func TestStructByName(t *testing.T) {
	type Str struct {
		Name []string `json:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	iface, err := StructFieldByName(val, "Name")
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := iface.([]string); !ok {
		t.Fatalf("not []string %v", iface)
	} else if len(v) != 2 {
		t.Fatalf("invalid number %v", iface)
	}
}

func TestStructURLValues(t *testing.T) {
	type Str struct {
		Str  *Str   `json:"str"`
		Name string `json:"name"`
		Args []int  `json:"args"`
	}

	val := &Str{Name: "test_name", Args: []int{1, 2, 3}, Str: &Str{Name: "nested_name"}}
	data, err := StructURLValues(val, "", []string{"json"})
	if err != nil {
		t.Fatal(err)
	}

	if data.Get("name") != "test_name" {
		t.Fatalf("invalid data: %v", data)
	}
}

func TestURLSliceVars(t *testing.T) {
	u, err := url.Parse("http://localhost/v1/test/call/my_name?key=arg1&key=arg2&key=arg3")
	if err != nil {
		t.Fatal(err)
	}

	mp, err := URLMap(u.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := mp["key"]
	if !ok {
		t.Fatalf("key not exists: %#+v", mp)
	}

	vm, ok := v.([]interface{})
	if !ok {
		t.Fatalf("invalid key value")
	}

	if len(vm) != 3 {
		t.Fatalf("missing key value: %#+v", mp)
	}
}

func TestURLVars(t *testing.T) {
	u, err := url.Parse("http://localhost/v1/test/call/my_name?req=key&arg1=arg1&arg2=12345&nested.string_args=str1&nested.string_args=str2&arg2=54321")
	if err != nil {
		t.Fatal(err)
	}

	mp, err := URLMap(u.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	_ = mp
}

func TestIsZero(t *testing.T) {
	testStr1 := struct {
		Name   string
		Value  string
		Nested struct {
			NestedName string
		}
	}{
		Name:  "test_name",
		Value: "test_value",
	}
	testStr1.Nested.NestedName = "nested_name"

	if ok := IsZero(testStr1); ok {
		t.Fatalf("zero ret on non zero struct: %#+v", testStr1)
	}

	testStr1.Name = ""
	testStr1.Value = ""
	testStr1.Nested.NestedName = ""
	if ok := IsZero(testStr1); !ok {
		t.Fatalf("non zero ret on zero struct: %#+v", testStr1)
	}

	type testStr3 struct {
		Nested string
	}
	type testStr2 struct {
		Nested *testStr3
		Name   string
	}
	vtest := &testStr2{
		Name:   "test_name",
		Nested: &testStr3{Nested: "nested_name"},
	}
	if ok := IsZero(vtest); ok {
		t.Fatalf("zero ret on non zero struct: %#+v", vtest)
	}
	vtest.Nested = nil
	vtest.Name = ""
	if ok := IsZero(vtest); !ok {
		t.Fatalf("non zero ret on zero struct: %#+v", vtest)
	}

	// t.Logf("XX %#+v\n", ok)
}
