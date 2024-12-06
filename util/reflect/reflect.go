package reflect

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

var (
	// ErrInvalidStruct happens when passed not struct and not struct pointer
	ErrInvalidStruct = errors.New("invalid struct specified")
	// ErrInvalidValue happens when passed invalid value for field
	ErrInvalidValue = errors.New("invalid value specified")
	// ErrNotFound happens when struct field not found
	ErrNotFound = errors.New("struct field not found")
)

// Option func signature
type Option func(*Options)

// Options for merge
type Options struct {
	// Tags specifies tags to lookup
	Tags []string
	// SliceAppend controls slice appending
	SliceAppend bool
}

// Tags sets the merge tags for lookup
func Tags(t []string) Option {
	return func(o *Options) {
		o.Tags = t
	}
}

// SliceAppend sets the option
func SliceAppend(b bool) Option {
	return func(o *Options) {
		o.SliceAppend = b
	}
}

var maxDepth = 32

func mergeMap(dst, src map[string]interface{}, depth int) map[string]interface{} {
	if depth > maxDepth {
		return dst
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := mapify(srcVal)
			dstMap, dstMapOk := mapify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = mergeMap(dstMap, srcMap, depth+1)
			}
		}
		dst[key] = srcVal
	}
	return dst
}

func mapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)
	if value.Kind() == reflect.Map {
		m := map[string]interface{}{}
		for _, k := range value.MapKeys() {
			m[k.String()] = value.MapIndex(k).Interface()
		}
		return m, true
	}
	return map[string]interface{}{}, false
}

// Merge merges map[string]interface{} to destination struct
func Merge(dst interface{}, mp map[string]interface{}, opts ...Option) error {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	if unmarshaler, ok := dst.(json.Unmarshaler); ok {
		buf, err := json.Marshal(mp)
		if err == nil {
			err = unmarshaler.UnmarshalJSON(buf)
		}
		return err
	}

	if mapper, ok := dst.(map[string]interface{}); ok {
		mergeMap(mapper, mp, 0)
		return nil
	}

	var err error
	var sval reflect.Value
	var fname string

	dviface := reflect.ValueOf(dst)
	if dviface.Kind() == reflect.Ptr {
		dviface = dviface.Elem()
	}

	if dviface.Kind() != reflect.Struct {
		return ErrInvalidStruct
	}

	dtype := dviface.Type()
	for idx := 0; idx < dtype.NumField(); idx++ {
		dfld := dtype.Field(idx)
		dval := dviface.Field(idx)
		if !dval.CanSet() || len(dfld.PkgPath) != 0 || !dval.IsValid() {
			continue
		}

		fname = ""
		for _, tname := range options.Tags {
			tvalue, ok := dfld.Tag.Lookup(tname)
			if !ok {
				continue
			}

			tpart := strings.Split(tvalue, ",")
			switch tname {
			case "protobuf":
				for _, p := range tpart {
					if idx := strings.Index(p, "name="); idx > 0 {
						fname = p[idx:]
					}
				}
			default:
				fname = tpart[0]
			}

			if fname != "" {
				break
			}
		}

		if fname == "" {
			fname = FieldName(dfld.Name)
		}

		val, ok := mp[fname]
		if !ok {
			val, ok = mp[dfld.Name]
			if !ok {
				continue
			}
		}

		sval = reflect.ValueOf(val)

		switch getKind(dval) {
		case reflect.Bool:
			err = mergeBool(dval, sval)
		case reflect.String:
			err = mergeString(dval, sval)
		case reflect.Int:
			err = mergeInt(dval, sval)
		case reflect.Uint:
			err = mergeUint(dval, sval)
		case reflect.Float32:
			err = mergeFloat(dval, sval)
		case reflect.Struct:
			mp, ok := sval.Interface().(map[string]interface{})
			if !ok {
				return ErrInvalidValue
			}
			err = Merge(dval.Interface(), mp, opts...)
		case reflect.Ptr:
			mp, ok := sval.Interface().(map[string]interface{})
			if !ok {
				return ErrInvalidValue
			}
			if dval.IsNil() {
				dval.Set(reflect.New(dval.Type().Elem()))
				if dval.Elem().Type().Kind() == reflect.Struct {
					for i := 0; i < dval.Elem().NumField(); i++ {
						field := dval.Elem().Field(i)
						if field.Type().Kind() == reflect.Ptr && field.IsNil() && dval.Elem().Type().Field(i).Anonymous {
							field.Set(reflect.New(field.Type().Elem()))
						}
					}
				}
			}
			err = Merge(dval.Interface(), mp, opts...)
		case reflect.Slice:
			if !options.SliceAppend && dval.IsNil() {
				dval.Set(reflect.MakeSlice(dval.Type(), sval.Len(), sval.Len()))
			} else if options.SliceAppend && dval.IsNil() {
				dval.Set(reflect.MakeSlice(dval.Type(), 0, 0))
			}

			for idx := 0; idx < sval.Len(); idx++ {
				var idval reflect.Value
				if !options.SliceAppend {
					idval = dval.Index(idx)
				} else {
					idval = reflect.Indirect(reflect.New(dval.Type().Elem()))
				}
				if getKind(idval) == reflect.Ptr && idval.IsNil() {
					idval.Set(reflect.New(idval.Type().Elem()))
					if idval.Elem().Type().Kind() == reflect.Struct {
						for i := 0; i < idval.Elem().NumField(); i++ {
							field := idval.Elem().Field(i)
							if field.Type().Kind() == reflect.Ptr && field.IsNil() && idval.Elem().Type().Field(i).Anonymous {
								field.Set(reflect.New(field.Type().Elem()))
							}
						}
					}
				}
				switch getKind(idval) {
				case reflect.Bool:
					err = mergeBool(idval, sval.Index(idx))
				case reflect.String:
					err = mergeString(idval, sval.Index(idx))
				case reflect.Int:
					err = mergeInt(idval, sval.Index(idx))
				case reflect.Uint:
					err = mergeUint(idval, sval.Index(idx))
				case reflect.Float32:
					err = mergeFloat(idval, sval.Index(idx))
				case reflect.Struct:
					imp, ok := sval.Index(idx).Interface().(map[string]interface{})
					if !ok {
						return ErrInvalidValue
					}
					err = Merge(idval.Interface(), imp, opts...)
				case reflect.Ptr:
					nsval := sval.Index(idx)
					if getKind(sval.Index(idx)) == reflect.Interface {
						nsval = reflect.ValueOf(nsval.Interface())
					}
					switch reflect.Indirect(idval).Type().String() {
					case "wrapperspb.BoolValue":
						if eva := reflect.Indirect(idval).FieldByName("Value"); eva.IsValid() {
							err = mergeBool(eva, nsval)
						}
					case "wrapperspb.BytesValue":
						if eva := reflect.Indirect(idval).FieldByName("Value"); eva.IsValid() {
							err = mergeUint(eva, nsval)
						}
					case "wrapperspb.DoubleValue", "wrapperspb.FloatValue":
						if eva := reflect.Indirect(idval).FieldByName("Value"); eva.IsValid() {
							err = mergeFloat(eva, nsval)
						}
					case "wrapperspb.Int32Value", "wrapperspb.Int64Value":
						if eva := reflect.Indirect(idval).FieldByName("Value"); eva.IsValid() {
							err = mergeInt(eva, nsval)
						}
					case "wrapperspb.StringValue":
						if eva := reflect.Indirect(idval).FieldByName("Value"); eva.IsValid() {
							err = mergeString(eva, nsval)
						}
					case "wrapperspb.UInt32Value", "wrapperspb.UInt64Value":
						if eva := reflect.Indirect(idval).FieldByName("Value"); eva.IsValid() {
							err = mergeUint(eva, nsval)
						}
					default:
						imp, ok := nsval.Interface().(map[string]interface{})
						if !ok {
							return ErrInvalidValue
						}
						if idval.IsNil() {
							idval.Set(reflect.New(idval.Type().Elem()))
							if idval.Elem().Type().Kind() == reflect.Struct {
								for i := 0; i < idval.Elem().NumField(); i++ {
									field := idval.Elem().Field(i)
									if field.Type().Kind() == reflect.Ptr && field.IsNil() && idval.Elem().Type().Field(i).Anonymous {
										field.Set(reflect.New(field.Type().Elem()))
									}
								}
							}
						}
						err = Merge(idval.Interface(), imp, opts...)
					}
				}
				if options.SliceAppend {
					dval.Set(reflect.Append(dval, idval))
				}
			}
			/*
				  case reflect.Interface:
							  err = d.decodeBasic(name, input, outVal)
							case reflect.Map:
								err = merge(dval, sval)
							case reflect.Array:
								err = mergeArray(dval, sval)
			*/
		default:
			err = ErrInvalidValue
		}

		if err != nil {
			err = fmt.Errorf("%v key %v invalid val %v", err, fname, sval.Interface())
		}

	}

	return nil
}

func mergeBool(dval reflect.Value, sval reflect.Value) error {
	switch getKind(sval) {
	case reflect.Int:
		switch sval.Int() {
		case 1:
			dval.SetBool(true)
		case 0:
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Uint:
		switch sval.Uint() {
		case 1:
			dval.SetBool(true)
		case 0:
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Float32:
		switch sval.Float() {
		case 1:
			dval.SetBool(true)
		case 0:
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Bool:
		dval.SetBool(sval.Bool())
	case reflect.String:
		switch sval.String() {
		case "t", "T", "true", "TRUE", "True", "1", "yes":
			dval.SetBool(true)
		case "f", "F", "false", "FALSE", "False", "0", "no":
			dval.SetBool(false)
		default:
			return ErrInvalidValue
		}
	case reflect.Interface:
		return mergeBool(dval, reflect.ValueOf(fmt.Sprintf("%v", sval.Interface())))
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeString(dval reflect.Value, sval reflect.Value) error {
	switch getKind(sval) {
	case reflect.Int:
		dval.SetString(strconv.FormatInt(sval.Int(), sval.Type().Bits()))
	case reflect.Uint:
		dval.SetString(strconv.FormatUint(sval.Uint(), sval.Type().Bits()))
	case reflect.Float32:
		dval.SetString(strconv.FormatFloat(sval.Float(), 'f', -1, sval.Type().Bits()))
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetString(strconv.FormatBool(true))
		case false:
			dval.SetString(strconv.FormatBool(false))
		}
	case reflect.String:
		dval.SetString(sval.String())
	case reflect.Interface:
		return mergeString(dval, reflect.ValueOf(fmt.Sprintf("%v", sval.Interface())))
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeInt(dval reflect.Value, sval reflect.Value) error {
	switch getKind(sval) {
	case reflect.Int:
		dval.SetInt(sval.Int())
	case reflect.Uint:
		dval.SetInt(int64(sval.Uint()))
	case reflect.Float32:
		dval.SetInt(int64(sval.Float()))
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetInt(1)
		case false:
			dval.SetInt(0)
		}
	case reflect.String:
		l, err := strconv.ParseInt(sval.String(), 0, dval.Type().Bits())
		if err != nil {
			return err
		}
		dval.SetInt(l)
	case reflect.Interface:
		return mergeInt(dval, reflect.ValueOf(fmt.Sprintf("%v", sval.Interface())))
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeUint(dval reflect.Value, sval reflect.Value) error {
	switch getKind(sval) {
	case reflect.Int:
		dval.SetUint(uint64(sval.Int()))
	case reflect.Uint:
		dval.SetUint(sval.Uint())
	case reflect.Float32:
		dval.SetUint(uint64(sval.Float()))
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetUint(1)
		case false:
			dval.SetUint(0)
		}
	case reflect.String:
		l, err := strconv.ParseUint(sval.String(), 0, dval.Type().Bits())
		if err != nil {
			return err
		}
		dval.SetUint(l)
	case reflect.Interface:
		return mergeUint(dval, reflect.ValueOf(fmt.Sprintf("%v", sval.Interface())))
	default:
		return ErrInvalidValue
	}
	return nil
}

func mergeFloat(dval reflect.Value, sval reflect.Value) error {
	switch getKind(sval) {
	case reflect.Int:
		dval.SetFloat(float64(sval.Int()))
	case reflect.Uint:
		dval.SetFloat(float64(sval.Uint()))
	case reflect.Float32:
		dval.SetFloat(sval.Float())
	case reflect.Bool:
		switch sval.Bool() {
		case true:
			dval.SetFloat(1)
		case false:
			dval.SetFloat(0)
		}
	case reflect.String:
		l, err := strconv.ParseFloat(sval.String(), dval.Type().Bits())
		if err != nil {
			return err
		}
		dval.SetFloat(l)
	case reflect.Interface:
		return mergeFloat(dval, reflect.ValueOf(fmt.Sprintf("%v", sval.Interface())))
	default:
		return ErrInvalidValue
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

// FieldName returns map field name that can be looked up in struct field
func FieldName(name string) string {
	newstr := make([]rune, 0, len(name))
	for idx, chr := range name {
		if idx == 0 {
			newstr = append(newstr, unicode.ToLower(chr))
			continue
		} else if unicode.IsUpper(chr) {
			newstr = append(newstr, '_')
			newstr = append(newstr, unicode.ToLower(chr))
			continue
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func Equal(src interface{}, dst interface{}, excptFields ...string) bool {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	switch srcVal.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < srcVal.Len(); i++ {
			e := srcVal.Index(i).Interface()
			a := dstVal.Index(i).Interface()
			if !Equal(e, a, excptFields...) {
				return false
			}
		}
		return true
	case reflect.Map:
		for i := 0; i < len(srcVal.MapKeys()); i++ {
			key := srcVal.MapKeys()[i]
			keyStr := fmt.Sprintf("%v", key.Interface())
			if stringContains(keyStr, excptFields) {
				continue
			}
			s := srcVal.MapIndex(key)
			d := dstVal.MapIndex(key)
			if !Equal(s.Interface(), d.Interface(), excptFields...) {
				return false
			}
		}
		return true
	case reflect.Struct, reflect.Interface:
		for i := 0; i < srcVal.NumField(); i++ {
			typeField := srcVal.Type().Field(i)
			if stringContains(typeField.Name, excptFields) {
				continue
			}
			s := srcVal.Field(i)
			d := dstVal.FieldByName(typeField.Name)
			if s.CanInterface() && d.CanInterface() {
				if !Equal(s.Interface(), d.Interface(), excptFields...) {
					return false
				}
			} else {
				return false
			}
		}
		return true
	case reflect.Ptr:
		if srcVal.IsNil() {
			return dstVal.IsNil()
		}
		s := srcVal.Elem()
		d := reflect.Indirect(dstVal)
		if s.CanInterface() && d.CanInterface() {
			return Equal(s.Interface(), d.Interface(), excptFields...)
		}
		return false
	case reflect.String, reflect.Int, reflect.Int64, reflect.Float32, reflect.Float64, reflect.Bool:
		return src == dst
	default:
		return srcVal.Interface() == dstVal.Interface()
	}
}

func stringContains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
