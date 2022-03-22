package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	// SplitToken used to detect path components
	SplitToken = "."
	// IndexCloseChar used to detect index end
	IndexCloseChar = "]"
	// IndexOpenChar used to detect index start
	IndexOpenChar = "["
)

var (
	// ErrMalformedIndex returns when index key have invalid format
	ErrMalformedIndex = errors.New("malformed index key")
	// ErrInvalidIndexUsage returns when index key usage error
	ErrInvalidIndexUsage = errors.New("invalid index key usage")
	// ErrKeyNotFound returns when key not found
	ErrKeyNotFound = errors.New("unable to find the key")
	// ErrBadJSONPath returns when path have invalid syntax
	ErrBadJSONPath = errors.New("bad path: must start with $ and have more then 2 chars")
)

// Lookup performs a lookup into a value, using a path of keys. The key should
// match with a Field or a MapIndex. For slice you can use the syntax key[index]
// to access a specific index. If one key owns to a slice and an index is not
// specificied the rest of the path will be apllied to evaley value of the
// slice, and the value will be merged into a slice.
func Lookup(i interface{}, path string) (reflect.Value, error) {
	if path == "" || path[0:1] != "$" {
		return reflect.Value{}, ErrBadJSONPath
	}

	if path == "$" {
		return reflect.ValueOf(i), nil
	}

	if len(path) < 2 {
		return reflect.Value{}, ErrBadJSONPath
	}

	return lookup(i, strings.Split(path[2:], SplitToken)...)
}

func lookup(i interface{}, path ...string) (reflect.Value, error) {
	value := reflect.ValueOf(i)
	var parent reflect.Value
	var err error

	for i, part := range path {
		parent = value

		value, err = getValueByName(value, part)
		if err == nil {
			continue
		}

		if !isAggregable(parent) {
			break
		}

		value, err = aggreateAggregableValue(parent, path[i:])

		break
	}

	return value, err
}

func getValueByName(v reflect.Value, key string) (reflect.Value, error) {
	var value reflect.Value
	var index int
	var err error

	key, index, err = parseIndex(key)
	if err != nil {
		return value, err
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return getValueByName(v.Elem(), key)
	case reflect.Struct:
		value = v.FieldByName(key)
	case reflect.Map:
		kvalue := reflect.Indirect(reflect.New(v.Type().Key()))
		kvalue.SetString(key)
		value = v.MapIndex(kvalue)
	}

	if !value.IsValid() {
		return reflect.Value{}, ErrKeyNotFound
	}

	if index != -1 {
		if value.Type().Kind() != reflect.Slice {
			return reflect.Value{}, ErrInvalidIndexUsage
		}

		value = value.Index(index)
	}

	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	return value, nil
}

func aggreateAggregableValue(v reflect.Value, path []string) (reflect.Value, error) {
	values := make([]reflect.Value, 0)

	l := v.Len()
	if l == 0 {
		ty, ok := lookupType(v.Type(), path...)
		if !ok {
			return reflect.Value{}, ErrKeyNotFound
		}
		return reflect.MakeSlice(reflect.SliceOf(ty), 0, 0), nil
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		break
	default:
		return reflect.Value{}, fmt.Errorf("unsuported kind for index")
	}

	index := indexFunction(v)
	for i := 0; i < l; i++ {
		value, err := lookup(index(i).Interface(), path...)
		if err != nil {
			return reflect.Value{}, err
		}

		values = append(values, value)
	}

	return mergeValue(values), nil
}

func indexFunction(v reflect.Value) func(i int) reflect.Value {
	switch v.Kind() {
	case reflect.Slice:
		return v.Index
	case reflect.Map:
		keys := v.MapKeys()
		return func(i int) reflect.Value {
			return v.MapIndex(keys[i])
		}
	}
	return func(i int) reflect.Value { return reflect.Value{} }
}

func mergeValue(values []reflect.Value) reflect.Value {
	values = removeZeroValues(values)
	l := len(values)
	if l == 0 {
		return reflect.Value{}
	}

	sample := values[0]
	mergeable := isMergeable(sample)

	t := sample.Type()
	if mergeable {
		t = t.Elem()
	}

	value := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	for i := 0; i < l; i++ {
		if !values[i].IsValid() {
			continue
		}

		if mergeable {
			value = reflect.AppendSlice(value, values[i])
		} else {
			value = reflect.Append(value, values[i])
		}
	}

	return value
}

func removeZeroValues(values []reflect.Value) []reflect.Value {
	l := len(values)

	var v []reflect.Value
	for i := 0; i < l; i++ {
		if values[i].IsValid() {
			v = append(v, values[i])
		}
	}

	return v
}

func isAggregable(v reflect.Value) bool {
	k := v.Kind()

	return k == reflect.Map || k == reflect.Slice
}

func isMergeable(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Map || k == reflect.Slice
}

func hasIndex(s string) bool {
	return strings.Contains(s, IndexOpenChar)
}

func parseIndex(s string) (string, int, error) {
	start := strings.Index(s, IndexOpenChar)
	end := strings.Index(s, IndexCloseChar)

	if start == -1 && end == -1 {
		return s, -1, nil
	}

	if (start != -1 && end == -1) || (start == -1 && end != -1) {
		return "", -1, ErrMalformedIndex
	}

	index, err := strconv.Atoi(s[start+1 : end])
	if err != nil {
		return "", -1, ErrMalformedIndex
	}

	return s[:start], index, nil
}

func lookupType(ty reflect.Type, path ...string) (reflect.Type, bool) {
	if len(path) == 0 {
		return ty, true
	}

	switch ty.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		if hasIndex(path[0]) {
			return lookupType(ty.Elem(), path[1:]...)
		}
		// Aggregate.
		return lookupType(ty.Elem(), path...)
	case reflect.Ptr:
		return lookupType(ty.Elem(), path...)
	case reflect.Interface:
		// We can't know from here without a value. Let's just return this type.
		return ty, true
	case reflect.Struct:
		f, ok := ty.FieldByName(path[0])
		if ok {
			return lookupType(f.Type, path[1:]...)
		}
	}
	return nil, false
}
