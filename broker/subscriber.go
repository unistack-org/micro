package broker

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

const (
	messageSig  = "func(broker.Message) error"
	messagesSig = "func([]broker.Message) error"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// Is this an exported - upper case - name?
func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// IsValidHandler func signature
func IsValidHandler(sub interface{}) error {
	typ := reflect.TypeOf(sub)
	var argType reflect.Type
	switch typ.Kind() {
	case reflect.Func:
		name := "Func"
		switch typ.NumIn() {
		case 1:
			argType = typ.In(0)
		default:
			return fmt.Errorf("subscriber %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), messageSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("subscriber %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("subscriber %v has wrong number of return values: %v require signature %s",
				name, typ.NumOut(), messageSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("subscriber %v returns %v not error", name, returnType.String())
		}
	default:
		hdlr := reflect.ValueOf(sub)
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			switch method.Type.NumIn() {
			case 3:
				argType = method.Type.In(2)
			default:
				return fmt.Errorf("subscriber %v.%v takes wrong number of args: %v required signature %s",
					name, method.Name, method.Type.NumIn(), messageSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"subscriber %v.%v has wrong number of return values: %v require signature %s",
					name, method.Name, method.Type.NumOut(), messageSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("subscriber %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}
