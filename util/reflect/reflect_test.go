package reflect

import (
	"fmt"
	"testing"
)

func TestMergeMapStringInterface(t *testing.T) {
	var dst interface{} //nolint:gosimple
	dst = map[string]interface{}{
		"xx": 11,
	}

	src := map[string]interface{}{
		"zz": "aa",
	}

	if err := Merge(dst, src); err != nil {
		t.Fatal(err)
	}

	mp, ok := dst.(map[string]interface{})
	if !ok || mp == nil {
		t.Fatalf("xxx %#+v\n", dst)
	}

	if fmt.Sprintf("%v", mp["xx"]) != "11" {
		t.Fatalf("xxx zzzz %#+v", mp)
	}

	if fmt.Sprintf("%v", mp["zz"]) != "aa" {
		t.Fatalf("xxx zzzz %#+v", mp)
	}
}

func TestMergeMap(t *testing.T) {
	src := map[string]interface{}{
		"skey1": "sval1",
		"skey2": map[string]interface{}{
			"skey3": "sval3",
		},
	}
	dst := map[string]interface{}{
		"skey1": "dval1",
		"skey2": map[string]interface{}{
			"skey3": "dval3",
		},
	}

	if err := Merge(src, dst); err != nil {
		t.Fatal(err)
	}

	t.Logf("%#+v", src)
}

func TestFieldName(t *testing.T) {
	src := "SomeVar"
	chk := "some_var"
	dst := FieldName(src)
	if dst != chk {
		t.Fatalf("FieldName error %s != %s", src, chk)
	}
}

func TestMergeBool(t *testing.T) {
	type str struct {
		Bool bool `json:"bool"`
	}

	mp := make(map[string]interface{})
	mp["bool"] = "true"
	s := &str{}

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if !s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = "false"

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = 1

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if !s.Bool {
		t.Fatalf("merge bool error: %#+v\n", s)
	}
}

func TestMergeString(t *testing.T) {
	type str struct {
		Bool string `json:"bool"`
	}

	mp := make(map[string]interface{})
	mp["bool"] = true
	s := &str{}

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatalf("merge with true err: %v", err)
	}

	if s.Bool != "true" {
		t.Fatalf("merge bool error: %#+v\n", s)
	}

	mp["bool"] = false
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatalf("merge with falst err: %v", err)
	}

	if s.Bool != "false" {
		t.Fatalf("merge bool error: %#+v\n", s)
	}
}

func TestMergeNested(t *testing.T) {
	type CallReqNested struct {
		Nested     *CallReqNested `json:"nested2"`
		StringArgs []string       `json:"string_args"`
		Uint64Args []uint64       `json:"uint64_args"`
	}

	type CallReq struct {
		Nested *CallReqNested `json:"nested"`
		Name   string         `json:"name"`
		Req    string         `json:"req"`
		Arg2   int            `json:"arg2"`
	}

	dst := &CallReq{
		Name: "name_old",
		Req:  "req_old",
	}

	mp := make(map[string]interface{})
	mp["name"] = "name_new"
	mp["req"] = "req_new"
	mp["arg2"] = 1
	mp["nested.string_args"] = []string{"args1", "args2"}
	mp["nested.uint64_args"] = []uint64{1, 2, 3}
	mp["nested.nested2.uint64_args"] = []uint64{1, 2, 3}

	mp = FlattenMap(mp)

	if err := Merge(dst, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if dst.Name != "name_new" || dst.Req != "req_new" || dst.Arg2 != 1 {
		t.Fatalf("merge error: %#+v", dst)
	}

	if dst.Nested == nil || len(dst.Nested.Uint64Args) != 3 ||
		len(dst.Nested.StringArgs) != 2 || dst.Nested.StringArgs[0] != "args1" ||
		len(dst.Nested.Uint64Args) != 3 || dst.Nested.Uint64Args[2] != 3 {
		t.Fatalf("merge error: %#+v", dst.Nested)
	}

	nmp := make(map[string]interface{})
	nmp["nested.uint64_args"] = []uint64{4}
	nmp = FlattenMap(nmp)

	if err := Merge(dst, nmp, SliceAppend(true), Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if dst.Nested == nil || len(dst.Nested.Uint64Args) != 4 || dst.Nested.Uint64Args[3] != 4 {
		t.Fatalf("merge error: %#+v", dst.Nested)
	}
}

func TestEqual(t *testing.T) {
	type tstr struct {
		Key1 string
		Key2 string
	}

	src := &tstr{Key1: "val1", Key2: "micro:generate"}
	dst := &tstr{Key1: "val1", Key2: "val2"}
	if !Equal(src, dst, "Key2") {
		t.Fatal("invalid Equal test")
	}
}
