package reflect

import (
	"fmt"
	"testing"
)

func TestMergeMapStringInterface(t *testing.T) {
	var dst interface{} //nolint:staticcheck
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
		dst.Nested.Uint64Args[2] != 3 {
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

func TestZero(t *testing.T) {
	type Str struct {
		Name  string
		Value int
	}
	src := &Str{Name: "hello", Value: 42}
	z, err := Zero(src)
	if err != nil {
		t.Fatal(err)
	}
	if z == nil {
		t.Fatal("expected non-nil zero value")
	}
	zs, ok := z.(*Str)
	if !ok {
		t.Fatalf("expected *Str got %T", z)
	}
	if zs.Name != "" || zs.Value != 0 {
		t.Fatalf("expected zero struct, got %+v", zs)
	}
}

func TestMergeFloat(t *testing.T) {
	type str struct {
		F32 float32 `json:"f32"`
		F64 float64 `json:"f64"`
	}

	mp := make(map[string]interface{})
	mp["f32"] = "3.14"
	mp["f64"] = "2.718"
	s := &str{}

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.F32 < 3.13 || s.F32 > 3.15 {
		t.Fatalf("merge float32 error: %#+v\n", s)
	}
	if s.F64 < 2.717 || s.F64 > 2.719 {
		t.Fatalf("merge float64 error: %#+v\n", s)
	}
}

func TestMergeInt(t *testing.T) {
	type str struct {
		I   int   `json:"i"`
		I32 int32 `json:"i32"`
		I64 int64 `json:"i64"`
	}

	mp := make(map[string]interface{})
	mp["i"] = "10"
	mp["i32"] = "20"
	mp["i64"] = "30"
	s := &str{}

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.I != 10 || s.I32 != 20 || s.I64 != 30 {
		t.Fatalf("merge int error: %#+v\n", s)
	}
}

func TestMergeUint(t *testing.T) {
	type str struct {
		U   uint   `json:"u"`
		U32 uint32 `json:"u32"`
		U64 uint64 `json:"u64"`
	}

	mp := make(map[string]interface{})
	mp["u"] = "10"
	mp["u32"] = "20"
	mp["u64"] = "30"
	s := &str{}

	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}

	if s.U != 10 || s.U32 != 20 || s.U64 != 30 {
		t.Fatalf("merge uint error: %#+v\n", s)
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

func TestEqualSlice(t *testing.T) {
	src := []string{"a", "b"}
	dst := []string{"a", "b"}
	if !Equal(src, dst) {
		t.Fatal("expected equal slices")
	}
}

func TestEqualMap(t *testing.T) {
	src := map[string]string{"k": "v"}
	dst := map[string]string{"k": "v"}
	if !Equal(src, dst) {
		t.Fatal("expected equal maps")
	}
}

func TestMergeBoolFromUint(t *testing.T) {
	type str struct {
		B bool `json:"b"`
	}
	mp := map[string]interface{}{"b": uint(1)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if !s.B {
		t.Fatalf("expected true from uint(1), got %v", s.B)
	}
}

func TestMergeBoolFromFloat(t *testing.T) {
	type str struct {
		B bool `json:"b"`
	}
	mp := map[string]interface{}{"b": float64(0)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.B {
		t.Fatalf("expected false from float64(0), got %v", s.B)
	}
}

func TestMergeIntFromUint(t *testing.T) {
	type str struct {
		I int `json:"i"`
	}
	mp := map[string]interface{}{"i": uint(42)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.I != 42 {
		t.Fatalf("expected 42 got %d", s.I)
	}
}

func TestMergeIntFromFloat(t *testing.T) {
	type str struct {
		I int `json:"i"`
	}
	mp := map[string]interface{}{"i": float64(7)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.I != 7 {
		t.Fatalf("expected 7 got %d", s.I)
	}
}

func TestMergeIntFromBool(t *testing.T) {
	type str struct {
		I int `json:"i"`
	}
	mp := map[string]interface{}{"i": true}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.I != 1 {
		t.Fatalf("expected 1 from bool true, got %d", s.I)
	}
}

func TestMergeUintFromInt(t *testing.T) {
	type str struct {
		U uint `json:"u"`
	}
	mp := map[string]interface{}{"u": int(5)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.U != 5 {
		t.Fatalf("expected 5 got %d", s.U)
	}
}

func TestMergeUintFromFloat(t *testing.T) {
	type str struct {
		U uint `json:"u"`
	}
	mp := map[string]interface{}{"u": float64(9)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.U != 9 {
		t.Fatalf("expected 9 got %d", s.U)
	}
}

func TestMergeUintFromBool(t *testing.T) {
	type str struct {
		U uint `json:"u"`
	}
	mp := map[string]interface{}{"u": false}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.U != 0 {
		t.Fatalf("expected 0 from bool false, got %d", s.U)
	}
}

func TestMergeFloatFromInt(t *testing.T) {
	type str struct {
		F float64 `json:"f"`
	}
	mp := map[string]interface{}{"f": int(3)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.F != 3.0 {
		t.Fatalf("expected 3.0 got %v", s.F)
	}
}

func TestMergeFloatFromUint(t *testing.T) {
	type str struct {
		F float64 `json:"f"`
	}
	mp := map[string]interface{}{"f": uint(4)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.F != 4.0 {
		t.Fatalf("expected 4.0 got %v", s.F)
	}
}

func TestMergeFloatFromBool(t *testing.T) {
	type str struct {
		F float64 `json:"f"`
	}
	mp := map[string]interface{}{"f": true}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.F != 1.0 {
		t.Fatalf("expected 1.0 from bool true, got %v", s.F)
	}
}

func TestMergeStringFromInt(t *testing.T) {
	type str struct {
		S string `json:"s"`
	}
	// mergeString uses sval.Type().Bits() as base; int8 = 8 bits → base 8 (valid)
	mp := map[string]interface{}{"s": int8(7)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.S == "" {
		t.Fatalf("expected non-empty string from int8, got %q", s.S)
	}
}

func TestMergeStringFromUint(t *testing.T) {
	type str struct {
		S string `json:"s"`
	}
	// mergeString uses sval.Type().Bits() as base; uint8 = 8 bits → base 8 (valid)
	mp := map[string]interface{}{"s": uint8(9)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.S == "" {
		t.Fatalf("expected non-empty string from uint8, got %q", s.S)
	}
}

func TestMergeStringFromFloat(t *testing.T) {
	type str struct {
		S string `json:"s"`
	}
	mp := map[string]interface{}{"s": float64(1.5)}
	s := &str{}
	if err := Merge(s, mp, Tags([]string{"json"})); err != nil {
		t.Fatal(err)
	}
	if s.S == "" {
		t.Fatalf("expected non-empty string from float, got %q", s.S)
	}
}

func TestZeroInvalid(t *testing.T) {
	_, err := Zero(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestMergeNonStruct(t *testing.T) {
	// Merge to a non-struct, non-map should return ErrInvalidStruct
	var s string = "hello"
	err := Merge(&s, map[string]interface{}{"x": 1})
	if err == nil {
		t.Fatal("expected error merging into non-struct pointer")
	}
}

func TestMergeBoolInvalidValue(t *testing.T) {
	type str struct {
		B bool `json:"b"`
	}
	// string that doesn't parse as bool
	mp := map[string]interface{}{"b": "notabool"}
	s := &str{}
	err := Merge(s, mp, Tags([]string{"json"}))
	if err == nil {
		t.Fatal("expected error for invalid bool string")
	}
}

func TestMergeIntInvalidString(t *testing.T) {
	type str struct {
		I int `json:"i"`
	}
	mp := map[string]interface{}{"i": "notanumber"}
	s := &str{}
	err := Merge(s, mp, Tags([]string{"json"}))
	if err == nil {
		t.Fatal("expected error for invalid int string")
	}
}

func TestMergeUintInvalidString(t *testing.T) {
	type str struct {
		U uint `json:"u"`
	}
	mp := map[string]interface{}{"u": "notanumber"}
	s := &str{}
	err := Merge(s, mp, Tags([]string{"json"}))
	if err == nil {
		t.Fatal("expected error for invalid uint string")
	}
}

func TestMergeFloatInvalidString(t *testing.T) {
	type str struct {
		F float64 `json:"f"`
	}
	mp := map[string]interface{}{"f": "notafloat"}
	s := &str{}
	err := Merge(s, mp, Tags([]string{"json"}))
	if err == nil {
		t.Fatal("expected error for invalid float string")
	}
}

func TestIsEmptyFunc(t *testing.T) {
	// non-nil func: IsZero should return false
	import_reflect := func() {}
	if IsZero(import_reflect) {
		t.Fatal("non-nil func should not be empty")
	}
	// nil func: IsZero should return true
	var fn func()
	if !IsZero(fn) {
		t.Fatal("nil func should be empty")
	}
}

func TestIsEmptyStructDefault(t *testing.T) {
	// default case in IsEmpty: chan — IsZero returns false for channels
	ch := make(chan int)
	if IsZero(ch) {
		t.Fatal("channel default case should return false")
	}
}

func TestEqualPointerBothNil(t *testing.T) {
	var a *struct{ X int }
	var b *struct{ X int }
	if !Equal(a, b) {
		t.Fatal("two nil pointers should be equal")
	}
}

func TestEqualPointerNilVsNonNil(t *testing.T) {
	type S struct{ X int }
	var a *S
	b := &S{X: 1}
	if Equal(a, b) {
		t.Fatal("nil vs non-nil pointer should not be equal")
	}
}
