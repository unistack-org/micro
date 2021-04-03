package reflect

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrInvalidStruct specifies invalid struct error
	ErrInvalidStruct = errors.New("invalid struct specified")
	// ErrInvalidParam specifies invalid url query params
	ErrInvalidParam = errors.New("invalid url query param provided")
)

var (
	bracketSplitter = regexp.MustCompile(`\[|\]`)
)

func fieldName(name string) string {
	newstr := make([]rune, 0)
	upper := false
	for idx, chr := range name {
		if idx == 0 {
			upper = true
		} else if chr == '_' {
			upper = true
			continue
		}
		if upper {
			newstr = append(newstr, unicode.ToUpper(chr))
		} else {
			newstr = append(newstr, chr)
		}
		upper = false
	}

	return string(newstr)
}

// IsEmpty returns true if value empty
func IsEmpty(v reflect.Value) bool {
	switch getKind(v) {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int:
		return v.Int() == 0
	case reflect.Uint:
		return v.Uint() == 0
	case reflect.Float32:
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
	case reflect.Struct:
		var ok bool
		for i := 0; i < v.NumField(); i++ {
			ok = IsEmpty(v.FieldByIndex([]int{i}))
			if !ok {
				return false
			}
		}
	default:
		return false
	}
	return true
}

// IsZero returns true if struct is zero (not have any defined values)
func IsZero(src interface{}) bool {
	v := reflect.ValueOf(src)
	return IsEmpty(v)
}

// Zero creates new zero interface
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

// StructFields returns slice of struct fields
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

func StructURLValues(src interface{}, pref string, tags []string) (url.Values, error) {
	data := url.Values{}

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
		if !val.CanSet() || len(fld.PkgPath) != 0 || !val.IsValid() {
			continue
		}

		var t *tag
		for _, tn := range tags {
			ts, ok := fld.Tag.Lookup(tn)
			if !ok {
				continue
			}

			tp := strings.Split(ts, ",")
			// special
			switch tn {
			case "protobuf": // special
				t = &tag{key: tn, name: tp[3][5:], opts: append(tp[:3], tp[4:]...)}
			default:
				t = &tag{key: tn, name: tp[0], opts: tp[1:]}
			}
			if t.name != "" {
				break
			}
		}

		if t.name == "" {
			// fallback to lowercase
			t.name = strings.ToLower(fld.Name)
		}
		if pref != "" {
			t.name = pref + "." + t.name
		}

		if !val.IsValid() || val.IsZero() {
			continue
		}

		switch val.Kind() {
		case reflect.Struct, reflect.Ptr:
			if val.IsNil() {
				continue
			}
			ndata, err := StructURLValues(val.Interface(), t.name, tags)
			if err != nil {
				return ndata, err
			}
			for k, v := range ndata {
				data[k] = v
			}
		default:
			switch val.Kind() {
			case reflect.Slice:
				for i := 0; i < val.Len(); i++ {
					va := val.Index(i)
					//if va.Type().Elem().Kind() != reflect.Ptr {
					if va.Kind() != reflect.Ptr {
						data.Set(t.name, fmt.Sprintf("%v", va.Interface()))
						continue
					}
					switch va.Type().Elem().String() {
					case "wrapperspb.BoolValue", "wrapperspb.BytesValue", "wrapperspb.StringValue":
						if eva := reflect.Indirect(va).FieldByName("Value"); eva.IsValid() {
							data.Add(t.name, fmt.Sprintf("%v", eva.Interface()))
						}
					case "wrapperspb.DoubleValue", "wrapperspb.FloatValue":
						if eva := reflect.Indirect(va).FieldByName("Value"); eva.IsValid() {
							data.Add(t.name, fmt.Sprintf("%v", eva.Interface()))
						}
					case "wrapperspb.Int32Value", "wrapperspb.Int64Value":
						if eva := reflect.Indirect(va).FieldByName("Value"); eva.IsValid() {
							data.Add(t.name, fmt.Sprintf("%v", eva.Interface()))
						}
					case "wrapperspb.UInt32Value", "wrapperspb.UInt64Value":
						if eva := reflect.Indirect(va).FieldByName("Value"); eva.IsValid() {
							data.Add(t.name, fmt.Sprintf("%v", eva.Interface()))
						}
					default:
						data.Add(t.name, fmt.Sprintf("%v", val.Index(i).Interface()))
					}
				}
			default:
				data.Set(t.name, fmt.Sprintf("%v", val.Interface()))
			}
		}
	}
	return data, nil
}

// URLMap returns map of url query params
func URLMap(query string) (map[string]interface{}, error) {
	var (
		mp interface{} = make(map[string]interface{})
	)

	params := strings.Split(query, "&")

	for _, part := range params {
		tm, err := queryToMap(part)
		if err != nil {
			return nil, err
		}
		mp = merge(mp, tm)
	}

	return mp.(map[string]interface{}), nil
}

// FlattenMap expand key.subkey to nested map
func FlattenMap(a map[string]interface{}) map[string]interface{} {
	// preprocess map
	nb := make(map[string]interface{}, len(a))
	for k, v := range a {
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
	return nb
}

// MergeMap merges maps
//nolint:gocyclo
func MergeMap(a interface{}, b map[string]interface{}) error {
	var err error

	ta := reflect.TypeOf(a)
	if ta.Kind() == reflect.Ptr {
		ta = ta.Elem()
	}
	va := reflect.ValueOf(a)
	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}

	for mk, mv := range b {
		vmv := reflect.ValueOf(mv)
		name := fieldName(mk)
		fva := va.FieldByName(name)
		fta, found := ta.FieldByName(name)
		if !found || !fva.IsValid() || !fva.CanSet() || fta.PkgPath != "" {
			continue
		}
		// fast path via direct assign
		if vmv.Type().AssignableTo(fta.Type) && !IsEmpty(vmv) {
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
		case reflect.Float32:
			err = mergeFloat(fva, vmv)
		case reflect.Array:
			//fmt.Printf("Array %#+v %#+v\n", fva, vmv)
		case reflect.Slice:
			err = mergeSlice(fva, vmv)
		case reflect.Ptr:
			if fva.IsNil() {
				fva.Set(reflect.New(fva.Type().Elem()))
				if fva.Elem().Type().Kind() == reflect.Struct {
					for i := 0; i < fva.Elem().NumField(); i++ {
						field := fva.Elem().Field(i)
						if field.Type().Kind() == reflect.Ptr && field.IsNil() && fva.Elem().Type().Field(i).Anonymous {
							field.Set(reflect.New(field.Type().Elem()))
						}
					}
				}
			}
			if nmp, ok := vmv.Interface().(map[string]interface{}); ok {
				err = MergeMap(fva.Interface(), nmp)
			} else {
				err = fmt.Errorf("cant fill")
			}
		case reflect.Struct:
			if nmp, ok := vmv.Interface().(map[string]interface{}); ok {
				err = MergeMap(fva.Interface(), nmp)
			} else {
				err = fmt.Errorf("cant fill")
			}
		case reflect.Map:
			//fmt.Printf("Map %#+v %#+v\n", fva, vmv)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

//nolint:gocyclo
func mergeSlice(va, vb reflect.Value) error {
	switch getKind(vb) {
	/*
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
	*/
	case reflect.String:
		var err error
		fn := func(c rune) bool { return c == ',' || c == ';' || c == ' ' }
		slice := strings.FieldsFunc(vb.String(), fn)
		va.Set(reflect.MakeSlice(va.Type(), len(slice), len(slice)))
		for idx, sl := range slice {
			vl := reflect.ValueOf(sl)
			switch va.Type().Elem().Kind() {
			case reflect.Bool:
				err = mergeBool(va.Index(idx), vl)
			case reflect.String:
				err = mergeString(va.Index(idx), vl)
			case reflect.Ptr:
				if va.Index(idx).IsNil() {
					va.Index(idx).Set(reflect.New(va.Index(idx).Type().Elem()))
				}
				switch va.Type().Elem().String() {
				case "*wrapperspb.BoolValue":
					if eva := reflect.Indirect(va.Index(idx)).FieldByName("Value"); eva.IsValid() {
						err = mergeBool(eva, vl)
					}
				case "*wrapperspb.BytesValue":
					if eva := va.Index(idx).FieldByName("Value"); eva.IsValid() {
						err = mergeUint(eva, vl)
					}
				case "*wrapperspb.DoubleValue", "*wrapperspb.FloatValue":
					if eva := reflect.Indirect(va.Index(idx)).FieldByName("Value"); eva.IsValid() {
						err = mergeFloat(eva, vl)
					}
				case "*wrapperspb.Int32Value", "*wrapperspb.Int64Value":
					if eva := reflect.Indirect(va.Index(idx)).FieldByName("Value"); eva.IsValid() {
						err = mergeInt(eva, vl)
					}
				case "*wrapperspb.StringValue":
					if eva := reflect.Indirect(va.Index(idx)).FieldByName("Value"); eva.IsValid() {
						err = mergeString(eva, vl)
					}
				case "*wrapperspb.UInt32Value", "*wrapperspb.UInt64Value":
					if eva := reflect.Indirect(va.Index(idx)).FieldByName("Value"); eva.IsValid() {
						err = mergeUint(eva, vl)
					}
				}
			}
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("cant merge %v %s with %v %s", va.Type(), va.Kind(), vb.Type(), vb.Kind())
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
	case reflect.Float32:
		if vb.Float() == 1 {
			va.SetBool(true)
		}
	case reflect.String:
		b, err := strconv.ParseBool(vb.String())
		if err != nil {
			return err
		}
		va.SetBool(b)
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
	case reflect.Float32:
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
	case reflect.Float32:
		va.SetInt(int64(vb.Float()))
	case reflect.String:
		f, err := strconv.ParseInt(vb.String(), 10, va.Type().Bits())
		if err != nil {
			return err
		}
		va.SetInt(f)
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
	case reflect.Float32:
		va.SetUint(uint64(vb.Float()))
	case reflect.String:
		f, err := strconv.ParseUint(vb.String(), 10, va.Type().Bits())
		if err != nil {
			return err
		}
		va.SetUint(f)
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
	case reflect.Float32:
		va.Set(vb)
	case reflect.String:
		f, err := strconv.ParseFloat(vb.String(), va.Type().Bits())
		if err != nil {
			return err
		}
		va.SetFloat(f)
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
		return reflect.Float32
	}
	return kind
}

func btSplitter(str string) []string {
	r := bracketSplitter.Split(str, -1)
	for idx, s := range r {
		if len(s) == 0 {
			if len(r) > idx+1 {
				copy(r[idx:], r[idx+1:])
				r = r[:len(r)-1]
			}
		}
	}
	return r
}

// queryToMap turns something like a[b][c]=4 into
//   map[string]interface{}{
//     "a": map[string]interface{}{
// 		  "b": map[string]interface{}{
// 			  "c": 4,
// 		  },
// 	  },
//   }
func queryToMap(param string) (map[string]interface{}, error) {
	rawKey, rawValue, err := splitKeyAndValue(param)
	if err != nil {
		return nil, err
	}
	rawValue, err = url.QueryUnescape(rawValue)
	if err != nil {
		return nil, err
	}
	rawKey, err = url.QueryUnescape(rawKey)
	if err != nil {
		return nil, err
	}

	pieces := btSplitter(rawKey)
	key := pieces[0]

	// If len==1 then rawKey has no [] chars and we can just
	// decode this as key=value into {key: value}
	if len(pieces) == 1 {
		return map[string]interface{}{
			key: rawValue,
		}, nil
	}

	// If len > 1 then we have something like a[b][c]=2
	// so we need to turn this into {"a": {"b": {"c": 2}}}
	// To do this we break our key into two pieces:
	//   a and b[c]
	// and then we set {"a": queryToMap("b[c]", value)}
	ret := make(map[string]interface{})
	ret[key], err = queryToMap(buildNewKey(rawKey) + "=" + rawValue)
	if err != nil {
		return nil, err
	}

	// When URL params have a set of empty brackets (eg a[]=1)
	// it is assumed to be an array. This will get us the
	// correct value for the array item and return it as an
	// []interface{} so that it can be merged properly.
	if pieces[1] == "" {
		temp := ret[key].(map[string]interface{})
		ret[key] = []interface{}{temp[""]}
	}
	return ret, nil
}

// buildNewKey will take something like:
// origKey = "bar[one][two]"
// pieces = [bar one two ]
// and return "one[two]"
func buildNewKey(origKey string) string {
	pieces := btSplitter(origKey)

	ret := origKey[len(pieces[0])+1:]
	ret = ret[:len(pieces[1])] + ret[len(pieces[1])+1:]
	return ret
}

// splitKeyAndValue splits a URL param at the last equal
// sign and returns the two strings. If no equal sign is
// found, the ErrInvalidParam error is returned.
func splitKeyAndValue(param string) (string, string, error) {
	li := strings.LastIndex(param, "=")
	if li == -1 {
		return "", "", ErrInvalidParam
	}
	return param[:li], param[li+1:], nil
}

// merge merges a with b if they are either both slices
// or map[string]interface{} types. Otherwise it returns b.
func merge(a interface{}, b interface{}) interface{} {
	if av, aok := a.(map[string]interface{}); aok {
		if bv, bok := b.(map[string]interface{}); bok {
			return mergeMapIface(av, bv)
		}
	}
	if av, aok := a.([]interface{}); aok {
		if bv, bok := b.([]interface{}); bok {
			return mergeSliceIface(av, bv)
		}
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if (va.Type().Kind() == reflect.Slice) && (va.Type().Elem().Kind() == vb.Type().Kind() || vb.Type().ConvertibleTo(va.Type().Elem())) {
		va = reflect.Append(va, vb.Convert(va.Type().Elem()))
		return va.Interface()
	}

	return b
}

// mergeMap merges a with b, attempting to merge any nested
// values in nested maps but eventually overwriting anything
// in a that can't be merged with whatever is in b.
func mergeMapIface(a map[string]interface{}, b map[string]interface{}) map[string]interface{} {
	for bK, bV := range b {
		if aV, ok := a[bK]; ok {
			if (reflect.ValueOf(aV).Type().Kind() == reflect.ValueOf(bV).Type().Kind()) ||
				((reflect.ValueOf(aV).Type().Kind() == reflect.Slice) && reflect.ValueOf(aV).Type().Elem().Kind() == reflect.ValueOf(bV).Type().Kind()) {
				nV := []interface{}{aV, bV}
				a[bK] = nV
			} else {
				a[bK] = merge(a[bK], bV)
			}
		} else {
			a[bK] = bV
		}
	}
	return a
}

// mergeSlice merges a with b and returns the result.
func mergeSliceIface(a []interface{}, b []interface{}) []interface{} {
	a = append(a, b...)
	return a
}

type tag struct {
	key  string
	name string
	opts []string
}
