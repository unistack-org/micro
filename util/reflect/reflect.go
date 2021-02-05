package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func MergeMap(a interface{}, b map[string]interface{}) error {
	var err error

	// preprocess map
	nb := make(map[string]interface{}, len(b))
	for k, v := range b {
		ps := strings.Split(k, ".")
		if len(ps) == 1 {
			nb[k] = v
			continue
		}
		em := make(map[string]interface{})
		em[ps[len(ps)-1]] = v
		for i := len(ps) - 2; i > 0; i-- {
			nm := make(map[string]interface{})
			nm[ps[i]] = em
			em = nm
		}
		if vm, ok := nb[ps[0]]; ok {
			// nested map
			nm := vm.(map[string]interface{})
			for vk, vv := range em {
				nm[vk] = vv
			}
			nb[ps[0]] = nm
		} else {
			nb[ps[0]] = em
		}
	}

	ta := reflect.TypeOf(a)
	if ta.Kind() == reflect.Ptr {
		ta = ta.Elem()
	}
	va := reflect.ValueOf(a)
	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}

	for mk, mv := range nb {
		vmv := reflect.ValueOf(mv)
		//		tmv := reflect.TypeOf(mv)
		name := strings.Title(mk)
		fva := va.FieldByName(name)
		fta, found := ta.FieldByName(name)
		if !found || !fva.IsValid() || !fva.CanSet() || fta.PkgPath != "" {
			continue
		}
		// fast path via direct assign
		if vmv.Type().AssignableTo(fta.Type) {
			fva.Set(vmv)
			continue
		}
		switch getKind(fva) {
		case reflect.Bool:
			err = mergeBool(fva, vmv)
		case reflect.String:
			err = mergeString(fva, vmv)
		case reflect.Int:
			err = mergeInt(fva, vmv)
		case reflect.Uint:
			err = mergeUint(fva, vmv)
		case reflect.Float64:
			err = mergeFloat(fva, vmv)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeBool(va, vb reflect.Value) error {
	switch getKind(vb) {
	case reflect.Int:
		if vb.Int() == 1 {
			va.SetBool(true)
		}
	case reflect.Uint:
		if vb.Uint() == 1 {
			va.SetBool(true)
		}
	case reflect.Float64:
		if vb.Float() == 1 {
			va.SetBool(true)
		}
	case reflect.String:
		if vb.String() == "1" || vb.String() == "true" {
			vb.SetBool(true)
		}
	default:
		return fmt.Errorf("cant merge %v %s with %v %s", va, va.Kind(), vb, vb.Kind())
	}
	return nil
}

func mergeString(va, vb reflect.Value) error {
	switch getKind(vb) {
	case reflect.Int:
		va.SetString(fmt.Sprintf("%d", vb.Int()))
	case reflect.Uint:
		va.SetString(fmt.Sprintf("%d", vb.Uint()))
	case reflect.Float64:
		va.SetString(fmt.Sprintf("%f", vb.Float()))
	case reflect.String:
		va.Set(vb)
	default:
		return fmt.Errorf("cant merge %v %s with %v %s", va, va.Kind(), vb, vb.Kind())
	}
	return nil
}

func mergeInt(va, vb reflect.Value) error {
	switch getKind(vb) {
	case reflect.Int:
		va.Set(vb)
	case reflect.Uint:
		va.SetInt(int64(vb.Uint()))
	case reflect.Float64:
		va.SetInt(int64(vb.Float()))
	case reflect.String:
		if f, err := strconv.ParseInt(vb.String(), 10, va.Type().Bits()); err != nil {
			return err
		} else {
			va.SetInt(f)
		}
	default:
		return fmt.Errorf("cant merge %v %s with %v %s", va, va.Kind(), vb, vb.Kind())
	}
	return nil
}

func mergeUint(va, vb reflect.Value) error {
	switch getKind(vb) {
	case reflect.Int:
		va.SetUint(uint64(vb.Int()))
	case reflect.Uint:
		va.Set(vb)
	case reflect.Float64:
		va.SetUint(uint64(vb.Float()))
	case reflect.String:
		if f, err := strconv.ParseUint(vb.String(), 10, va.Type().Bits()); err != nil {
			return err
		} else {
			va.SetUint(f)
		}
	default:
		return fmt.Errorf("cant merge %v %s with %v %s", va, va.Kind(), vb, vb.Kind())
	}
	return nil
}

func mergeFloat(va, vb reflect.Value) error {
	switch getKind(vb) {
	case reflect.Int:
		va.SetFloat(float64(vb.Int()))
	case reflect.Uint:
		va.SetFloat(float64(vb.Uint()))
	case reflect.Float64:
		va.Set(vb)
	case reflect.String:
		if f, err := strconv.ParseFloat(vb.String(), va.Type().Bits()); err != nil {
			return err
		} else {
			va.SetFloat(f)
		}
	default:
		return fmt.Errorf("cant merge %v %s with %v %s", va, va.Kind(), vb, vb.Kind())
	}

	return nil
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()
	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float64
	}
	return kind
}
