package unwrap

import (
	"testing"

	"go.unistack.org/micro/v3/codec"
)

func TestUnwrapOmit(t *testing.T) {
	type val struct {
		MP  map[string]string `json:"mp" logger:"omit"`
		STR string            `json:"str"`
		AR  []string          `json:"ar"`
	}

	v1 := &val{AR: []string{"string1", "string2"}, STR: "string", MP: map[string]string{"key": "val"}}

	t.Logf("output: %#v", v1)
	t.Logf("output: %#v", Unwrap(v1))
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

	t.Logf("output: %#v", Unwrap(v1))

	type val2 struct {
		mp  map[string]string
		val *val2
		str string
		ar  []string
	}

	v2 := &val2{ar: []string{string1, string2}, str: string1, val: &val2{str: string2}, mp: map[string]string{"key": "val"}}

	t.Logf("output: %#v", v2)
}

func TestUnwrapCodec(t *testing.T) {
	type val struct {
		MP  map[string]string `json:"mp"`
		STR string            `json:"str"`
		AR  []string          `json:"ar"`
	}

	v1 := &val{AR: []string{"string1", "string2"}, STR: "string", MP: map[string]string{"key": "val"}}

	t.Logf("output: %#v", Unwrap(v1, UnwrapCodec(codec.NewCodec())))
}
