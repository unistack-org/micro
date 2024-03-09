package broker

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	subSig      = "func(context.Context, interface{}) error"
	batchSubSig = "func([]context.Context, []interface{}) error"
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

// ValidateSubscriber func signature
func ValidateSubscriber(sub interface{}) error {
	typ := reflect.TypeOf(sub)
	var argType reflect.Type
	switch typ.Kind() {
	case reflect.Func:
		name := "Func"
		switch typ.NumIn() {
		case 1: // func(Message) error

		case 2: // func(context.Context, Message) error or func(context.Context, []Message) error
			argType = typ.In(2)
			// if sub.Options().Batch {
			if argType.Kind() != reflect.Slice {
				return fmt.Errorf("subscriber %v dont have required signature %s", name, batchSubSig)
			}
			if strings.Compare(fmt.Sprintf("%v", argType), "[]interface{}") == 0 {
				return fmt.Errorf("subscriber %v dont have required signaure %s", name, batchSubSig)
			}
		//	}
		default:
			return fmt.Errorf("subscriber %v takes wrong number of args: %v required signature %s or %s", name, typ.NumIn(), subSig, batchSubSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("subscriber %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("subscriber %v has wrong number of return values: %v require signature %s or %s",
				name, typ.NumOut(), subSig, batchSubSig)
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
				return fmt.Errorf("subscriber %v.%v takes wrong number of args: %v required signature %s or %s",
					name, method.Name, method.Type.NumIn(), subSig, batchSubSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"subscriber %v.%v has wrong number of return values: %v require signature %s or %s",
					name, method.Name, method.Type.NumOut(), subSig, batchSubSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("subscriber %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}
