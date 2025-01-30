package options

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

func toInt8SliceE(i interface{}) ([]int8, error) {
	if i == nil {
		return []int8{}, fmt.Errorf("unable to cast %#v of type %T to []int8", i, i)
	}

	switch v := i.(type) {
	case []int8:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int8, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt8E(s.Index(j).Interface())
			if err != nil {
				return []int8{}, fmt.Errorf("unable to cast %#v of type %T to []int8", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int8{}, fmt.Errorf("unable to cast %#v of type %T to []int8", i, i)
	}
}

func toInt16SliceE(i interface{}) ([]int16, error) {
	if i == nil {
		return []int16{}, fmt.Errorf("unable to cast %#v of type %T to []int16", i, i)
	}

	switch v := i.(type) {
	case []int16:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int16, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt16E(s.Index(j).Interface())
			if err != nil {
				return []int16{}, fmt.Errorf("unable to cast %#v of type %T to []int16", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int16{}, fmt.Errorf("unable to cast %#v of type %T to []int16", i, i)
	}
}

func toInt32SliceE(i interface{}) ([]int32, error) {
	if i == nil {
		return []int32{}, fmt.Errorf("unable to cast %#v of type %T to []int32", i, i)
	}

	switch v := i.(type) {
	case []int32:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int32, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt32E(s.Index(j).Interface())
			if err != nil {
				return []int32{}, fmt.Errorf("unable to cast %#v of type %T to []int32", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int32{}, fmt.Errorf("unable to cast %#v of type %T to []int32", i, i)
	}
}

func toInt64SliceE(i interface{}) ([]int64, error) {
	if i == nil {
		return []int64{}, fmt.Errorf("unable to cast %#v of type %T to []int64", i, i)
	}

	switch v := i.(type) {
	case []int64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToInt64E(s.Index(j).Interface())
			if err != nil {
				return []int64{}, fmt.Errorf("unable to cast %#v of type %T to []int64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []int64{}, fmt.Errorf("unable to cast %#v of type %T to []int64", i, i)
	}
}

func toUintSliceE(i interface{}) ([]uint, error) {
	if i == nil {
		return []uint{}, fmt.Errorf("unable to cast %#v of type %T to []uint", i, i)
	}

	switch v := i.(type) {
	case []uint:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]uint, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToUintE(s.Index(j).Interface())
			if err != nil {
				return []uint{}, fmt.Errorf("unable to cast %#v of type %T to []uint", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []uint{}, fmt.Errorf("unable to cast %#v of type %T to []uint", i, i)
	}
}

func toUint8SliceE(i interface{}) ([]uint8, error) {
	if i == nil {
		return []uint8{}, fmt.Errorf("unable to cast %#v of type %T to []uint8", i, i)
	}

	switch v := i.(type) {
	case []uint8:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]uint8, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToUint8E(s.Index(j).Interface())
			if err != nil {
				return []uint8{}, fmt.Errorf("unable to cast %#v of type %T to []uint8", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []uint8{}, fmt.Errorf("unable to cast %#v of type %T to []uint8", i, i)
	}
}

func toUint16SliceE(i interface{}) ([]uint16, error) {
	if i == nil {
		return []uint16{}, fmt.Errorf("unable to cast %#v of type %T to []uint16", i, i)
	}

	switch v := i.(type) {
	case []uint16:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]uint16, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToUint16E(s.Index(j).Interface())
			if err != nil {
				return []uint16{}, fmt.Errorf("unable to cast %#v of type %T to []uint16", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []uint16{}, fmt.Errorf("unable to cast %#v of type %T to []uint16", i, i)
	}
}

func toUint32SliceE(i interface{}) ([]uint32, error) {
	if i == nil {
		return []uint32{}, fmt.Errorf("unable to cast %#v of type %T to []uint32", i, i)
	}

	switch v := i.(type) {
	case []uint32:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]uint32, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToUint32E(s.Index(j).Interface())
			if err != nil {
				return []uint32{}, fmt.Errorf("unable to cast %#v of type %T to []uint32", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []uint32{}, fmt.Errorf("unable to cast %#v of type %T to []uint32", i, i)
	}
}

func toUint64SliceE(i interface{}) ([]uint64, error) {
	if i == nil {
		return []uint64{}, fmt.Errorf("unable to cast %#v of type %T to []uint64", i, i)
	}

	switch v := i.(type) {
	case []uint64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]uint64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToUint64E(s.Index(j).Interface())
			if err != nil {
				return []uint64{}, fmt.Errorf("unable to cast %#v of type %T to []uint64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []uint64{}, fmt.Errorf("unable to cast %#v of type %T to []uint64", i, i)
	}
}

func toFloat32SliceE(i interface{}) ([]float32, error) {
	if i == nil {
		return []float32{}, fmt.Errorf("unable to cast %#v of type %T to []float32", i, i)
	}

	switch v := i.(type) {
	case []float32:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float32, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat32E(s.Index(j).Interface())
			if err != nil {
				return []float32{}, fmt.Errorf("unable to cast %#v of type %T to []float32", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float32{}, fmt.Errorf("unable to cast %#v of type %T to []float32", i, i)
	}
}

func toFloat64SliceE(i interface{}) ([]float64, error) {
	if i == nil {
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}

	switch v := i.(type) {
	case []float64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat64E(s.Index(j).Interface())
			if err != nil {
				return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float32", i, i)
	}
}

func setMap(src interface{}, dst interface{}) error {
	var err error

	if src == nil {
		return fmt.Errorf("unable to cast %#v of type %T", src, src)
	}
	if dst == nil {
		return fmt.Errorf("unable to cast %#v of type %T", dst, dst)
	}

	val := reflect.ValueOf(src)

	keyKind := val.Type().Key().Kind()
	valKind := val.Type().Elem().Kind()

	switch v := dst.(type) {
	case []interface{}:
		if len(v) == 1 {
			dstVal := reflect.ValueOf(v[0])
			if dstVal.Kind() != reflect.Map {
				return nil
			}
			mapIter := dstVal.MapRange()
			for mapIter.Next() {
				var (
					keyVal interface{}
					valVal interface{}
				)
				switch keyKind {
				case reflect.Bool:
					keyVal, err = cast.ToBoolE(mapIter.Key())
				case reflect.String:
					keyVal, err = cast.ToStringE(mapIter.Key())
				case reflect.Float32:
					keyVal, err = cast.ToFloat32E(mapIter.Key())
				case reflect.Float64:
					keyVal, err = cast.ToFloat64E(mapIter.Key())
				case reflect.Int8:
					keyVal, err = cast.ToInt8E(mapIter.Key())
				case reflect.Int:
					keyVal, err = cast.ToIntE(mapIter.Key())
				case reflect.Int16:
					keyVal, err = cast.ToInt16E(mapIter.Key())
				case reflect.Int32:
					keyVal, err = cast.ToInt32E(mapIter.Key())
				case reflect.Int64:
					keyVal, err = cast.ToInt64E(mapIter.Key())
				case reflect.Uint8:
					keyVal, err = cast.ToUint8E(mapIter.Key())
				case reflect.Uint:
					keyVal, err = cast.ToUintE(mapIter.Key())
				case reflect.Uint16:
					keyVal, err = cast.ToUint16E(mapIter.Key())
				case reflect.Uint32:
					keyVal, err = cast.ToUint32E(mapIter.Key())
				case reflect.Uint64:
					keyVal, err = cast.ToUint64E(mapIter.Key())
				}
				if err != nil {
					return err
				}
				switch valKind {
				case reflect.Bool:
					valVal, err = cast.ToBoolE(mapIter.Value())
				case reflect.String:
					valVal, err = cast.ToStringE(mapIter.Value())
				case reflect.Float32:
					valVal, err = cast.ToFloat32E(mapIter.Value())
				case reflect.Float64:
					valVal, err = cast.ToFloat64E(mapIter.Value())
				case reflect.Int8:
					valVal, err = cast.ToInt8E(mapIter.Value())
				case reflect.Int:
					valVal, err = cast.ToIntE(mapIter.Value())
				case reflect.Int16:
					valVal, err = cast.ToInt16E(mapIter.Value())
				case reflect.Int32:
					valVal, err = cast.ToInt32E(mapIter.Value())
				case reflect.Int64:
					valVal, err = cast.ToInt64E(mapIter.Value())
				case reflect.Uint8:
					valVal, err = cast.ToUint8E(mapIter.Value())
				case reflect.Uint:
					valVal, err = cast.ToUintE(mapIter.Value())
				case reflect.Uint16:
					valVal, err = cast.ToUint16E(mapIter.Value())
				case reflect.Uint32:
					valVal, err = cast.ToUint32E(mapIter.Value())
				case reflect.Uint64:
					valVal, err = cast.ToUint64E(mapIter.Value())
				}
				if err != nil {
					return err
				}

				val.SetMapIndex(reflect.ValueOf(keyVal), reflect.ValueOf(valVal))
			}
			return nil
		}
		if l := len(v) % 2; l == 1 {
			v = v[:len(v)-1]
		}
		var (
			keyVal interface{}
			valVal interface{}
		)
		for i := 0; i < len(v); i += 2 {
			switch keyKind {
			case reflect.Bool:
				keyVal, err = cast.ToBoolE(v[i])
			case reflect.String:
				keyVal, err = cast.ToStringE(v[i])
			case reflect.Float32:
				keyVal, err = cast.ToFloat32E(v[i])
			case reflect.Float64:
				keyVal, err = cast.ToFloat64E(v[i])
			case reflect.Int8:
				keyVal, err = cast.ToInt8E(v[i])
			case reflect.Int:
				keyVal, err = cast.ToIntE(v[i])
			case reflect.Int16:
				keyVal, err = cast.ToInt16E(v[i])
			case reflect.Int32:
				keyVal, err = cast.ToInt32E(v[i])
			case reflect.Int64:
				keyVal, err = cast.ToInt64E(v[i])
			case reflect.Uint8:
				keyVal, err = cast.ToUint8E(v[i])
			case reflect.Uint:
				keyVal, err = cast.ToUintE(v[i])
			case reflect.Uint16:
				keyVal, err = cast.ToUint16E(v[i])
			case reflect.Uint32:
				keyVal, err = cast.ToUint32E(v[i])
			case reflect.Uint64:
				keyVal, err = cast.ToUint64E(v[i])
			}
			if err != nil {
				return err
			}
			switch valKind {
			case reflect.Bool:
				valVal, err = cast.ToBoolE(v[i+1])
			case reflect.String:
				valVal, err = cast.ToStringE(v[i+1])
			case reflect.Float32:
				valVal, err = cast.ToFloat32E(v[i+1])
			case reflect.Float64:
				valVal, err = cast.ToFloat64E(v[i+1])
			case reflect.Int8:
				valVal, err = cast.ToInt8E(v[i+1])
			case reflect.Int:
				valVal, err = cast.ToIntE(v[i+1])
			case reflect.Int16:
				valVal, err = cast.ToInt16E(v[i+1])
			case reflect.Int32:
				valVal, err = cast.ToInt32E(v[i+1])
			case reflect.Int64:
				valVal, err = cast.ToInt64E(v[i+1])
			case reflect.Uint8:
				valVal, err = cast.ToUint8E(v[i+1])
			case reflect.Uint:
				valVal, err = cast.ToUintE(v[i+1])
			case reflect.Uint16:
				valVal, err = cast.ToUint16E(v[i+1])
			case reflect.Uint32:
				valVal, err = cast.ToUint32E(v[i+1])
			case reflect.Uint64:
				valVal, err = cast.ToUint64E(v[i+1])
			}
			if err != nil {
				return err
			}

			val.SetMapIndex(reflect.ValueOf(keyVal), reflect.ValueOf(valVal))
		}
	default:
		dstVal := reflect.ValueOf(dst)
		if dstVal.Kind() != reflect.Map {
			return nil
		}
		mapIter := dstVal.MapRange()
		for mapIter.Next() {
			var (
				keyVal interface{}
				valVal interface{}
			)
			switch keyKind {
			case reflect.Bool:
				keyVal, err = cast.ToBoolE(mapIter.Key())
			case reflect.String:
				keyVal, err = cast.ToStringE(mapIter.Key())
			case reflect.Float32:
				keyVal, err = cast.ToFloat32E(mapIter.Key())
			case reflect.Float64:
				keyVal, err = cast.ToFloat64E(mapIter.Key())
			case reflect.Int8:
				keyVal, err = cast.ToInt8E(mapIter.Key())
			case reflect.Int:
				keyVal, err = cast.ToIntE(mapIter.Key())
			case reflect.Int16:
				keyVal, err = cast.ToInt16E(mapIter.Key())
			case reflect.Int32:
				keyVal, err = cast.ToInt32E(mapIter.Key())
			case reflect.Int64:
				keyVal, err = cast.ToInt64E(mapIter.Key())
			case reflect.Uint8:
				keyVal, err = cast.ToUint8E(mapIter.Key())
			case reflect.Uint:
				keyVal, err = cast.ToUintE(mapIter.Key())
			case reflect.Uint16:
				keyVal, err = cast.ToUint16E(mapIter.Key())
			case reflect.Uint32:
				keyVal, err = cast.ToUint32E(mapIter.Key())
			case reflect.Uint64:
				keyVal, err = cast.ToUint64E(mapIter.Key())
			}
			if err != nil {
				return err
			}
			switch valKind {
			case reflect.Bool:
				valVal, err = cast.ToBoolE(mapIter.Value())
			case reflect.String:
				valVal, err = cast.ToStringE(mapIter.Value())
			case reflect.Float32:
				valVal, err = cast.ToFloat32E(mapIter.Value())
			case reflect.Float64:
				valVal, err = cast.ToFloat64E(mapIter.Value())
			case reflect.Int8:
				valVal, err = cast.ToInt8E(mapIter.Value())
			case reflect.Int:
				valVal, err = cast.ToIntE(mapIter.Value())
			case reflect.Int16:
				valVal, err = cast.ToInt16E(mapIter.Value())
			case reflect.Int32:
				valVal, err = cast.ToInt32E(mapIter.Value())
			case reflect.Int64:
				valVal, err = cast.ToInt64E(mapIter.Value())
			case reflect.Uint8:
				valVal, err = cast.ToUint8E(mapIter.Value())
			case reflect.Uint:
				valVal, err = cast.ToUintE(mapIter.Value())
			case reflect.Uint16:
				valVal, err = cast.ToUint16E(mapIter.Value())
			case reflect.Uint32:
				valVal, err = cast.ToUint32E(mapIter.Value())
			case reflect.Uint64:
				valVal, err = cast.ToUint64E(mapIter.Value())
			}
			if err != nil {
				return err
			}

			val.SetMapIndex(reflect.ValueOf(keyVal), reflect.ValueOf(valVal))
		}
		return nil
	}
	return nil
}
