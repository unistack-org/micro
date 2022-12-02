package unwrap

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"go.unistack.org/micro/v3/codec"
)

func TestSql(t *testing.T) {
	type val struct {
		S sql.NullString
	}
	v := &val{S: sql.NullString{Valid: true, String: "test"}}
	buf := fmt.Sprintf("%#v", Unwrap(v))
	t.Logf(buf)
}

func TestUnwrap(t *testing.T) {
	string1 := "string1"
	string2 := "string2"

	type val1 struct {
		mp  map[string]string
		val *val1
		str *string
		ar  []*string
	}

	v1 := &val1{ar: []*string{&string1, &string2}, str: &string1, val: &val1{str: &string2}, mp: map[string]string{"key": "val"}}

	buf := fmt.Sprintf("%#v", Unwrap(v1))
	chk := `&unwrap.val1{mp:map[string]string{"key":"val"}, val:(*unwrap.val1){mp:map[string]string<nil>, val:(*unwrap.val1)<nil>, str:(*string)"string2", ar:[]*string<nil>}, str:(*string)"string1", ar:[]*string{<*><shown>, <*>"string2"}}`
	if strings.Compare(buf, chk) != 0 {
		t.Fatalf("not proper written\n%s\n%s", buf, chk)
	}

	type val2 struct {
		mp  map[string]string
		val *val2
		str string
		ar  []string
	}

	v2 := &val2{ar: []string{string1, string2}, str: string1, val: &val2{str: string2}, mp: map[string]string{"key": "val"}}
	_ = v2
	// t.Logf("output: %#v", v2)
}

func TestCodec(t *testing.T) {
	type val struct {
		MP  map[string]string `json:"mp"`
		STR string            `json:"str"`
		AR  []string          `json:"ar"`
	}

	v1 := &val{AR: []string{"string1", "string2"}, STR: "string", MP: map[string]string{"key": "val"}}

	buf := fmt.Sprintf("%#v", Unwrap(v1, Codec(codec.NewCodec())))
	if strings.Compare(buf, `{"mp":{"key":"val"},"str":"string","ar":["string1","string2"]}`) != 0 {
		t.Fatalf("not proper written %s", buf)
	}
}

func TestOmit(t *testing.T) {
	type val struct {
		Key1 string `logger:"omit"`
		Key2 string `logger:"take"`
	}
	v1 := &val{Key1: "val1", Key2: "val2"}
	buf := fmt.Sprintf("%#v", Unwrap(v1))
	if strings.Compare(buf, `&unwrap.val{Key2:"val2"}`) != 0 {
		t.Fatalf("not proper written %s", buf)
	}
}

func TestTagged(t *testing.T) {
	type val struct {
		Key1 string `logger:"take"`
		Key2 string
	}

	v1 := &val{Key1: "val1", Key2: "val2"}
	buf := fmt.Sprintf("%#v", Unwrap(v1, Tagged(true)))
	if strings.Compare(buf, `&unwrap.val{Key1:"val1"}`) != 0 {
		t.Fatalf("not proper written %s", buf)
	}
}
