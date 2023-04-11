package register

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"go.unistack.org/micro/v4/metadata"
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

// ExtractEndpoint extract *Endpoint from reflect.Method
func ExtractEndpoint(method reflect.Method) *Endpoint {
	if method.PkgPath != "" {
		return nil
	}

	var rspType, reqType reflect.Type
	var stream bool
	mt := method.Type

	switch mt.NumIn() {
	case 3:
		reqType = mt.In(1)
		rspType = mt.In(2)
	case 4:
		reqType = mt.In(2)
		rspType = mt.In(3)
	default:
		return nil
	}

	// are we dealing with a stream?
	switch rspType.Kind() {
	case reflect.Func, reflect.Interface:
		stream = true
	}

	request := ExtractValue(reqType, 0)
	response := ExtractValue(rspType, 0)
	if request == "" || response == "" {
		return nil
	}

	ep := &Endpoint{
		Name:     method.Name,
		Request:  request,
		Response: response,
		Metadata: metadata.New(0),
	}

	if stream {
		ep.Metadata.Set("stream", fmt.Sprintf("%v", stream))
	}

	return ep
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
