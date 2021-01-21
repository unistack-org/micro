package reflect

import (
	"errors"
	"reflect"
)

var (
	ErrInvalidStruct = errors.New("invalid struct specified")
)

func IsEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem())
	case reflect.Func:
		return v.IsNil()
	case reflect.Invalid:
		return true
	}
	return false
}

func Zero(src interface{}) (interface{}, error) {
	sv := reflect.ValueOf(src)

	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}

	if sv.Kind() == reflect.Invalid {
		return nil, ErrInvalidStruct
	}

	dst := reflect.New(sv.Type())

	return dst.Interface(), nil
}

func StructFields(src interface{}) ([]reflect.StructField, error) {
	var fields []reflect.StructField

	sv := reflect.ValueOf(src)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		return nil, ErrInvalidStruct
	}

	typ := sv.Type()
	for idx := 0; idx < typ.NumField(); idx++ {
		fld := typ.Field(idx)
		val := sv.Field(idx)
		if !val.CanSet() || len(fld.PkgPath) != 0 {
			continue
		}
		if val.Kind() == reflect.Struct {
			infields, err := StructFields(val.Interface())
			if err != nil {
				return nil, err
			}
			fields = append(fields, infields...)
		} else {
			fields = append(fields, fld)
		}
	}

	return fields, nil
}

// CopyDefaults for a from b
// a and b should be pointers to the same kind of struct
func CopyDefaults(a, b interface{}) {
	pt := reflect.TypeOf(a)
	t := pt.Elem()
	va := reflect.ValueOf(a).Elem()
	vb := reflect.ValueOf(b).Elem()
	for i := 0; i < t.NumField(); i++ {
		aField := va.Field(i)
		if aField.CanSet() {
			bField := vb.Field(i)
			aField.Set(bField)
		}
	}
}

// CopyFrom sets the public members of a from b
// a and b should be pointers to structs
// a can be a different type from b
// Only the Fields which have the same name and assignable type on a
// and b will be set.
func CopyFrom(a, b interface{}) {
	ta := reflect.TypeOf(a).Elem()
	tb := reflect.TypeOf(b).Elem()
	va := reflect.ValueOf(a).Elem()
	vb := reflect.ValueOf(b).Elem()
	for i := 0; i < tb.NumField(); i++ {
		bField := vb.Field(i)
		tbField := tb.Field(i)
		name := tbField.Name
		aField := va.FieldByName(name)
		taField, found := ta.FieldByName(name)
		if found && aField.IsValid() && bField.IsValid() && aField.CanSet() && tbField.Type.AssignableTo(taField.Type) {
			aField.Set(bField)
		}
	}
}
