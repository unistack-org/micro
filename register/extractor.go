package register

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

// ExtractValue from reflect.Type from specified depth
func ExtractValue(v reflect.Type, d int) string {
	if d == 3 {
		return ""
	}
	if v == nil {
		return ""
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// slices and maps don't have a defined name
	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Map) || len(v.Name()) == 0 {
		return ""
	}

	// get the rune character
	a, _ := utf8.DecodeRuneInString(string(v.Name()[0]))

	// crude check for is unexported field
	if unicode.IsLower(a) {
		return ""
	}

	return v.Name()
}

// ExtractSubValue exctact *Value from reflect.Type
func ExtractSubValue(typ reflect.Type) string {
	var reqType reflect.Type
	switch typ.NumIn() {
	case 1:
		reqType = typ.In(0)
	case 2:
		reqType = typ.In(1)
	case 3:
		reqType = typ.In(2)
	default:
		return ""
	}
	return ExtractValue(reqType, 0)
}
