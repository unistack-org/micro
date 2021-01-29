package register

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/unistack-org/micro/v3/metadata"
)

// Extract *Value from reflect.Type
func ExtractValue(v reflect.Type, d int) *Value {
	if d == 3 {
		return nil
	}
	if v == nil {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if len(v.Name()) == 0 {
		return nil
	}

	// get the rune character
	a, _ := utf8.DecodeRuneInString(string(v.Name()[0]))

	// crude check for is unexported field
	if unicode.IsLower(a) {
		return nil
	}

	arg := &Value{
		Name: v.Name(),
		Type: v.Name(),
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			val := ExtractValue(f.Type, d+1)
			if val == nil {
				continue
			}

			// if we can find a json tag use it
			if tags := f.Tag.Get("json"); len(tags) > 0 {
				parts := strings.Split(tags, ",")
				if parts[0] == "-" || parts[0] == "omitempty" {
					continue
				}
				val.Name = parts[0]
			}

			// if there's no name default it
			if len(val.Name) == 0 {
				val.Name = v.Field(i).Name
			}

			arg.Values = append(arg.Values, val)
		}
	case reflect.Slice:
		p := v.Elem()
		if p.Kind() == reflect.Ptr {
			p = p.Elem()
		}
		arg.Type = "[]" + p.Name()
	}

	return arg
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
	if request == nil || response == nil {
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
func ExtractSubValue(typ reflect.Type) *Value {
	var reqType reflect.Type
	switch typ.NumIn() {
	case 1:
		reqType = typ.In(0)
	case 2:
		reqType = typ.In(1)
	case 3:
		reqType = typ.In(2)
	default:
		return nil
	}
	return ExtractValue(reqType, 0)
}
