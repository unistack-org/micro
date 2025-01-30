package options

import (
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
	mreflect "go.unistack.org/micro/v3/util/reflect"
)

// Options interface must be used by all options
type Validator interface {
	// Validate returns nil, if all options are correct,
	// otherwise returns an error explaining the mistake
	Validate() error
}

// Option func signature
type Option func(interface{}) error

// Apply assign options to struct src
func Apply(src interface{}, opts ...Option) error {
	for _, opt := range opts {
		if err := opt(src); err != nil {
			return err
		}
	}
	return nil
}

// SetValueByPath set src struct field to val dst via path
func SetValueByPath(src interface{}, dst interface{}, path string) error {
	var err error

	switch v := dst.(type) {
	case []interface{}:
		if len(v) == 1 {
			dst = v[0]
		}
	}

	var sv reflect.Value
	switch t := src.(type) {
	case reflect.Value:
		sv = t
	default:
		sv = reflect.ValueOf(src)
	}

	parts := strings.Split(path, ".")

	for _, p := range parts {

		if sv.Kind() == reflect.Ptr {
			sv = sv.Elem()
		}
		if sv.Kind() != reflect.Struct {
			return mreflect.ErrInvalidStruct
		}

		typ := sv.Type()
		for idx := 0; idx < typ.NumField(); idx++ {
			fld := typ.Field(idx)
			val := sv.Field(idx)

			/*
				if len(fld.PkgPath) != 0 {
					continue
				}
			*/

			if fld.Anonymous {
				if len(parts) == 1 && val.Kind() == reflect.Struct {
					if err = SetValueByPath(val, dst, p); err != nil {
						return err
					}
				}
			}

			if fld.Name != p && !strings.EqualFold(strings.ToLower(fld.Name), strings.ToLower(p)) {
				continue
			}

			switch val.Interface().(type) {
			case []time.Duration:
				dst, err = cast.ToDurationSliceE(dst)
				if err != nil {
					return err
				}
				reflect.Copy(val, reflect.ValueOf(dst))
				return nil
			case time.Duration:
				dst, err = cast.ToDurationE(dst)
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(dst))
				return nil
			case time.Time:
				dst, err = cast.ToTimeE(dst)
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(dst))
				return nil
			}

			switch val.Kind() {
			case reflect.Map:
				if val.IsZero() {
					val.Set(reflect.MakeMap(val.Type()))
				}

				return setMap(val.Interface(), dst)
			case reflect.Array, reflect.Slice:
				switch val.Type().Elem().Kind() {
				case reflect.Bool:
					dst, err = cast.ToBoolSliceE(dst)
				case reflect.String:
					dst, err = cast.ToStringSliceE(dst)
				case reflect.Float32:
					dst, err = toFloat32SliceE(dst)
				case reflect.Float64:
					dst, err = toFloat64SliceE(dst)
				case reflect.Int8:
					dst, err = toInt8SliceE(dst)
				case reflect.Int:
					dst, err = cast.ToIntSliceE(dst)
				case reflect.Int16:
					dst, err = toInt16SliceE(dst)
				case reflect.Int32:
					dst, err = toInt32SliceE(dst)
				case reflect.Int64:
					dst, err = toInt64SliceE(dst)
				case reflect.Uint8:
					dst, err = toUint8SliceE(dst)
				case reflect.Uint:
					dst, err = toUintSliceE(dst)
				case reflect.Uint16:
					dst, err = toUint16SliceE(dst)
				case reflect.Uint32:
					dst, err = toUint32SliceE(dst)
				case reflect.Uint64:
					dst, err = toUint64SliceE(dst)
				}
				if err != nil {
					return err
				}
				if val.Kind() == reflect.Slice {
					val.Set(reflect.ValueOf(dst))
				} else {
					reflect.Copy(val, reflect.ValueOf(dst))
				}
				return nil
			case reflect.Float32:
				dst, err = toFloat32SliceE(dst)
			case reflect.Float64:
				dst, err = toFloat64SliceE(dst)
			case reflect.Bool:
				dst, err = cast.ToBoolE(dst)
			case reflect.String:
				dst, err = cast.ToStringE(dst)
			case reflect.Int8:
				dst, err = cast.ToInt8E(dst)
			case reflect.Int:
				dst, err = cast.ToIntE(dst)
			case reflect.Int16:
				dst, err = cast.ToInt16E(dst)
			case reflect.Int32:
				dst, err = cast.ToInt32E(dst)
			case reflect.Int64:
				dst, err = cast.ToInt64E(dst)
			case reflect.Uint8:
				dst, err = cast.ToUint8E(dst)
			case reflect.Uint:
				dst, err = cast.ToUintE(dst)
			case reflect.Uint16:
				dst, err = cast.ToUint16E(dst)
			case reflect.Uint32:
				dst, err = cast.ToUint32E(dst)
			case reflect.Uint64:
				dst, err = cast.ToUint64E(dst)
			default:
			}
			if err != nil {
				return err
			}
			val.Set(reflect.ValueOf(dst))
		}
	}

	return nil
}

// NewOption create new option with name
func NewOption(name string) func(...interface{}) Option {
	return func(dst ...interface{}) Option {
		return func(src interface{}) error {
			return SetValueByPath(src, dst, name)
		}
	}
}

var (
	Address   = NewOption("Address")
	Name      = NewOption("Name")
	Broker    = NewOption("Broker")
	Logger    = NewOption("Logger")
	Meter     = NewOption("Meter")
	Tracer    = NewOption("Tracer")
	Store     = NewOption("Store")
	Register  = NewOption("Register")
	Router    = NewOption("Router")
	Codec     = NewOption("Codec")
	Codecs    = NewOption("Codecs")
	Client    = NewOption("Client")
	Context   = NewOption("Context")
	TLSConfig = NewOption("TLSConfig")
	Metadata  = NewOption("Metadata")
	Timeout   = NewOption("Timeout")
)
