package reflect_test

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	rutil "go.unistack.org/micro/v4/util/reflect"
)

func TestStructFields(t *testing.T) {
	type NestedStr struct {
		BBB string
		CCC int
	}
	type Str struct {
		Name   []string `json:"name" codec:"flatten"`
		XXX    string   `json:"xxx"`
		Nested NestedStr
	}

	val := &Str{Name: []string{"first", "second"}, XXX: "ttt", Nested: NestedStr{BBB: "ddd", CCC: 9}}
	fields, err := rutil.StructFields(val)
	if err != nil {
		t.Fatal(err)
	}
	var ok bool
	for _, field := range fields {
		if field.Path == "Nested.CCC" {
			ok = true
		}
	}
	if !ok {
		t.Fatalf("struct fields returns invalid path: %v", fields)
	}
}

func TestStructFieldsNested(t *testing.T) {
	type NestedConfig struct {
		Value string
	}
	type Config struct {
		Time     time.Time
		Nested   *NestedConfig
		Metadata map[string]int
		Broker   string
		Addr     []string
		Wait     time.Duration
		Verbose  bool
	}
	cfg := &Config{Nested: &NestedConfig{}}
	fields, err := rutil.StructFields(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 7 {
		for _, field := range fields {
			t.Logf("field %#+v\n", field)
		}
		t.Fatalf("invalid fields number: %d != %d", 7, len(fields))
	}
}

func TestSetFieldByPath(t *testing.T) {
	type NestedStr struct {
		BBB string `json:"bbb"`
		CCC int    `json:"ccc"`
	}
	type Str1 struct {
		Name   []string  `json:"name" codec:"flatten"`
		XXX    string    `json:"xxx"`
		Nested NestedStr `json:"nested"`
	}
	type Str2 struct {
		XXX    string     `json:"xxx"`
		Nested *NestedStr `json:"nested"`
		Name   []string   `json:"name" codec:"flatten"`
	}
	var err error
	val1 := &Str1{Name: []string{"first", "second"}, XXX: "ttt", Nested: NestedStr{BBB: "ddd", CCC: 9}}
	val2 := &Str2{Name: []string{"first", "second"}, XXX: "ttt", Nested: &NestedStr{BBB: "ddd", CCC: 9}}
	err = rutil.SetFieldByPath(val1, "xxx", "Nested.BBB")
	if err != nil {
		t.Fatal(err)
	}
	if val1.Nested.BBB != "xxx" {
		t.Fatalf("SetFieldByPath not works: %#+v", val1)
	}
	err = rutil.SetFieldByPath(val2, "xxx", "Nested.BBB")
	if err != nil {
		t.Fatal(err)
	}
	if val2.Nested.BBB != "xxx" {
		t.Fatalf("SetFieldByPath not works: %#+v", val1)
	}
}

func TestZeroFieldByPath(t *testing.T) {
	type NestedStr struct {
		BBB string `json:"bbb"`
		CCC int    `json:"ccc"`
	}
	type Str1 struct {
		Name   []string  `json:"name" codec:"flatten"`
		XXX    string    `json:"xxx"`
		Nested NestedStr `json:"nested"`
	}
	type Str2 struct {
		XXX    string     `json:"xxx"`
		Nested *NestedStr `json:"nested"`
		Name   []string   `json:"name" codec:"flatten"`
	}
	var err error
	val1 := &Str1{Name: []string{"first", "second"}, XXX: "ttt", Nested: NestedStr{BBB: "ddd", CCC: 9}}

	err = rutil.ZeroFieldByPath(val1, "Nested.BBB")
	if err != nil {
		t.Fatal(err)
	}
	err = rutil.ZeroFieldByPath(val1, "Nested")
	if err != nil {
		t.Fatal(err)
	}
	if val1.Nested.BBB == "ddd" {
		t.Fatalf("zero field not works: %v", val1)
	}

	val2 := &Str2{Name: []string{"first", "second"}, XXX: "ttt", Nested: &NestedStr{BBB: "ddd", CCC: 9}}
	err = rutil.ZeroFieldByPath(val2, "Nested")
	if err != nil {
		t.Fatal(err)
	}
	if val2.Nested != nil {
		t.Fatalf("zero field not works: %v", val2)
	}
}

func TestStructFieldsMap(t *testing.T) {
	type NestedStr struct {
		BBB string
		CCC int
	}
	type Str struct {
		Name   []string `json:"name" codec:"flatten"`
		XXX    string   `json:"xxx"`
		Nested NestedStr
	}

	val := &Str{Name: []string{"first", "second"}, XXX: "ttt", Nested: NestedStr{BBB: "ddd", CCC: 9}}
	fields, err := rutil.StructFieldsMap(val)
	if err != nil {
		t.Fatal(err)
	}
	if v, ok := fields["Nested.BBB"]; !ok || v != "ddd" {
		t.Fatalf("invalid field from %v", fields)
	}
}

func TestStructFieldByPath(t *testing.T) {
	type NestedStr struct {
		BBB string
		CCC int
	}
	type Str struct {
		Name   []string `json:"name" codec:"flatten"`
		XXX    string   `json:"xxx"`
		Nested NestedStr
	}

	val := &Str{Name: []string{"first", "second"}, XXX: "ttt", Nested: NestedStr{BBB: "ddd", CCC: 9}}
	field, err := rutil.StructFieldByPath(val, "Nested.CCC")
	if err != nil {
		t.Fatal(err)
	}
	if reflect.Indirect(reflect.ValueOf(field)).Int() != 9 {
		t.Fatalf("invalid elem returned: %v", field)
	}
}

func TestStructFieldByTag(t *testing.T) {
	type Str struct {
		Name []string `json:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	iface, err := rutil.StructFieldByTag(val, "codec", "flatten")
	if err != nil {
		t.Fatal(err)
	}

	if v, ok := iface.(*[]string); !ok {
		t.Fatalf("not *[]string %v", iface)
	} else if len(*v) != 2 {
		t.Fatalf("invalid number %v", iface)
	}
}

func TestStructFieldByName(t *testing.T) {
	type Str struct {
		Name []string `json:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	iface, err := rutil.StructFieldByName(val, "Name")
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
	data, err := rutil.StructURLValues(val, "", []string{"json"})
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

	mp, err := rutil.URLMap(u.RawQuery)
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

func TestURLMap(t *testing.T) {
	u, err := url.Parse("http://localhost/v1/test/call/my_name?req=key&arg1=arg1&arg2=12345&nested.string_args=str1&nested.string_args=str2&arg2=54321")
	if err != nil {
		t.Fatal(err)
	}

	mp, err := rutil.URLMap(u.RawQuery)
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

	if ok := rutil.IsZero(testStr1); ok {
		t.Fatalf("zero ret on non zero struct: %#+v", testStr1)
	}

	testStr1.Name = ""
	testStr1.Value = ""
	testStr1.Nested.NestedName = ""
	if ok := rutil.IsZero(testStr1); !ok {
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
	if ok := rutil.IsZero(vtest); ok {
		t.Fatalf("zero ret on non zero struct: %#+v", vtest)
	}
	vtest.Nested = nil
	vtest.Name = ""
	if ok := rutil.IsZero(vtest); !ok {
		t.Fatalf("non zero ret on zero struct: %#+v", vtest)
	}

	// t.Logf("XX %#+v\n", ok)
}
