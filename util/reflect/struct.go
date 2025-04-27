package reflect

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// ErrInvalidParam specifies invalid url query params
var ErrInvalidParam = errors.New("invalid url query param provided")

// var timeKind = reflect.ValueOf(time.Time{}).Kind()

// bracketSplitter
var bracketSplitter = regexp.MustCompile(`\[|\]`)

// StructField contains struct field path its value and field
type StructField struct {
	Value reflect.Value
	Path  string
	Field reflect.StructField
}

// StructFieldNameByTag get struct field name by tag key and its value
func StructFieldNameByTag(src interface{}, tkey string, tval string) (string, interface{}, error) {
	sv := reflect.ValueOf(src)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		return "", nil, ErrInvalidStruct
	}

	typ := sv.Type()
	for idx := 0; idx < typ.NumField(); idx++ {
		fld := typ.Field(idx)
		val := sv.Field(idx)
		if len(fld.PkgPath) != 0 {
			continue
		}

		if ts, ok := fld.Tag.Lookup(tkey); ok {
			for _, p := range strings.Split(ts, ",") {
				if p == tval {
					return fld.Name, val.Interface(), nil
				}
			}
		}

		switch val.Kind() {
		case reflect.Ptr:
			if val = val.Elem(); val.Kind() == reflect.Struct {
				if name, fld, err := StructFieldNameByTag(val.Interface(), tkey, tval); err == nil {
					return name, fld, nil
				}
			}
		case reflect.Struct:
			if name, fld, err := StructFieldNameByTag(val.Interface(), tkey, tval); err == nil {
				return name, fld, nil
			}
		}
	}
	return "", nil, ErrNotFound
}

// StructFieldByTag get struct field by tag key and its value
func StructFieldByTag(src interface{}, tkey string, tval string) (interface{}, error) {
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
		if len(fld.PkgPath) != 0 {
			continue
		}

		if ts, ok := fld.Tag.Lookup(tkey); ok {
			for _, p := range strings.Split(ts, ",") {
				if p == tval {
					return val.Interface(), nil
				}
			}
		}

		switch val.Kind() {
		case reflect.Ptr:
			if val = val.Elem(); val.Kind() == reflect.Struct {
				if iface, err := StructFieldByTag(val.Interface(), tkey, tval); err == nil {
					return iface, nil
				}
			}
		case reflect.Struct:
			if iface, err := StructFieldByTag(val.Interface(), tkey, tval); err == nil {
				return iface, nil
			}
		}
	}
	return nil, ErrNotFound
}

// ZeroFieldByPath clean struct field by its path
func ZeroFieldByPath(src interface{}, path string) error {
	if src == nil {
		return nil
	}
	var err error
	val := reflect.ValueOf(src)

	if IsEmpty(val) {
		return nil
	}

	for _, p := range strings.Split(path, ".") {
		if IsEmpty(val) {
			return nil
		}

		val, err = structValueByName(val, p)
		if err != nil {
			return err
		}
	}

	if IsEmpty(val) {
		return nil
	}

	if !val.CanSet() {
		return ErrInvalidStruct
	}

	val.Set(reflect.Zero(val.Type()))

	return nil
}

// SetFieldByPath set struct field by its path
func SetFieldByPath(src interface{}, dst interface{}, path string) error {
	var err error
	val := reflect.ValueOf(src)

	for _, p := range strings.Split(path, ".") {
		val, err = structValueByName(val, p)
		if err != nil {
			return err
		}
	}

	if !val.CanSet() {
		return ErrInvalidStruct
	}

	val.Set(reflect.ValueOf(dst))

	return nil
}

// structValueByName get struct field by its name
func structValueByName(sv reflect.Value, tkey string) (reflect.Value, error) {
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		return reflect.Zero(reflect.TypeOf(sv)), ErrInvalidStruct
	}

	typ := sv.Type()
	for idx := 0; idx < typ.NumField(); idx++ {
		fld := typ.Field(idx)
		val := sv.Field(idx)
		if len(fld.PkgPath) != 0 {
			continue
		}

		if fld.Name == tkey || strings.EqualFold(strings.ToLower(fld.Name), strings.ToLower(tkey)) {
			return val, nil
		}

		switch val.Kind() {
		case reflect.Ptr:
			if val = val.Elem(); val.Kind() == reflect.Struct {
				if iface, err := structValueByName(val, tkey); err == nil {
					return iface, nil
				}
			}
		case reflect.Struct:
			if iface, err := structValueByName(val, tkey); err == nil {
				return iface, nil
			}
		}
	}
	return reflect.Zero(reflect.TypeOf(sv)), ErrNotFound
}

// StructFieldByPath get struct field by its path
func StructFieldByPath(src interface{}, path string) (interface{}, error) {
	var err error
	for _, p := range strings.Split(path, ".") {
		src, err = StructFieldByName(src, p)
		if err != nil {
			return nil, err
		}
	}
	return src, err
}

// StructFieldByName get struct field by its name
func StructFieldByName(src interface{}, tkey string) (interface{}, error) {
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
		if len(fld.PkgPath) != 0 {
			continue
		}
		if fld.Name == tkey || strings.EqualFold(strings.ToLower(fld.Name), strings.ToLower(tkey)) {
			return val.Interface(), nil
		}

		switch val.Kind() {
		case reflect.Ptr:
			if val = val.Elem(); val.Kind() == reflect.Struct {
				if iface, err := StructFieldByName(val.Interface(), tkey); err == nil {
					return iface, nil
				}
			}
		case reflect.Struct:
			if iface, err := StructFieldByName(val.Interface(), tkey); err == nil {
				return iface, nil
			}
		}
	}
	return nil, ErrNotFound
}

// StructFieldsMap returns struct map[string]interface{} or error
func StructFieldsMap(src interface{}) (map[string]interface{}, error) {
	fields, err := StructFields(src)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		mp[field.Path] = field.Value.Interface()
	}
	return mp, nil
}

// StructFields returns slice of struct fields
func StructFields(src interface{}) ([]StructField, error) {
	var fields []StructField

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
		if !val.IsValid() || len(fld.PkgPath) != 0 {
			continue
		}

		switch val.Interface().(type) {
		case time.Time, *time.Time:
			fields = append(fields, StructField{Field: fld, Value: val, Path: fld.Name})
			continue
		case time.Duration, *time.Duration:
			fields = append(fields, StructField{Field: fld, Value: val, Path: fld.Name})
			continue
		}

		switch val.Kind() {
		case reflect.Ptr:
			if val.CanSet() && fld.Type.Elem().Kind() == reflect.Struct {
				if val.IsNil() {
					val.Set(reflect.New(fld.Type.Elem()))
				}
			}
			switch reflect.Indirect(val).Kind() {
			case reflect.Struct:
				infields, err := StructFields(val.Interface())
				if err != nil {
					return nil, err
				}
				for _, infield := range infields {
					infield.Path = fmt.Sprintf("%s.%s", fld.Name, infield.Path)
					fields = append(fields, infield)
				}
			default:
				fields = append(fields, StructField{Field: fld, Value: val, Path: fld.Name})
			}
		case reflect.Struct:
			infields, err := StructFields(val.Interface())
			if err != nil {
				return nil, err
			}
			for _, infield := range infields {
				infield.Path = fmt.Sprintf("%s.%s", fld.Name, infield.Path)
				fields = append(fields, infield)
			}
		default:

			fields = append(fields, StructField{Field: fld, Value: val, Path: fld.Name})
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

// StructURLValues get struct fields via url.Values
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
		if len(fld.PkgPath) != 0 || !val.IsValid() {
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
					// if va.Type().Elem().Kind() != reflect.Ptr {
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
	var mp interface{} = make(map[string]interface{})

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

// FlattenMapFixed flattens a nested map into a single-level map using dot notation for nested keys.
// In case of key conflicts, all nested levels will be discarded in favor of the first-level key.
//
// Example #1:
//
//	Input:
//	  {
//	    "user.name": "alex",
//	    "user.document.id": "document_id"
//	    "user.document.number": "document_number"
//	  }
//	Output:
//	  {
//	    "user": {
//	      "name": "alex",
//	      "document": {
//	        "id": "document_id"
//	        "number": "document_number"
//	      }
//	    }
//	  }
//
// Example #2 (with conflicts):
//
//	Input:
//	  {
//	    "user": "alex",
//	    "user.document.id": "document_id"
//	    "user.document.number": "document_number"
//	  }
//	Output:
//	  {
//	    "user": "alex"
//	  }
func FlattenMapFixed(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range input {
		parts := strings.Split(k, ".")

		if len(parts) == 1 {
			result[k] = v
			continue
		}

		current := result

		for i, part := range parts {
			// last element in the path
			if i == len(parts)-1 {
				current[part] = v
				break
			}

			// initialize map for current level if not exist
			if _, ok := current[part]; !ok {
				current[part] = make(map[string]interface{})
			}

			if nested, ok := current[part].(map[string]interface{}); ok {
				current = nested // continue to the nested map
			} else {
				break // if current element is not a map, ignore it
			}
		}
	}

	return result
}

/*
	case reflect.String:
		fn := func(c rune) bool { return c == ',' || c == ';' || c == ' ' }
		slice := strings.FieldsFunc(vb.String(), fn)
		if va.IsNil() {
			va.Set(reflect.MakeSlice(va.Type(), len(slice), len(slice)))
		}
*/

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
//
//	  map[string]interface{}{
//	    "a": map[string]interface{}{
//			  "b": map[string]interface{}{
//				  "c": 4,
//			  },
//		  },
//	  }
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
