package options

import (
	"context"
	"crypto/tls"
	"reflect"
	"time"

	"go.unistack.org/micro/v4/metadata"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

// Option func signature
type Option func(interface{}) error

// Set assign value to struct by its path
func Set(src interface{}, dst interface{}, path string) error {
	return rutil.SetFieldByPath(src, dst, path)
}

// Get returns value from struct by its path
func Get(src interface{}, path string) (interface{}, error) {
	return rutil.StructFieldByPath(src, path)
}

// Name set Name value
func Name(v ...string) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Name")
	}
}

// Address set Address value to single string or slice of strings
func Address(v ...string) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Address")
	}
}

// Broker set Broker value
func Broker(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Broker")
	}
}

// Logger set Logger value
func Logger(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Logger")
	}
}

// Meter set Meter value
func Meter(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Meter")
	}
}

// Tracer set Tracer value
func Tracer(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Tracer")
	}
}

// Store set Store value
func Store(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Store")
	}
}

// Register set Register value
func Register(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Register")
	}
}

// Router set Router value
func Router(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Router")
	}
}

// Codec set Codec value
func Codec(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Codec")
	}
}

// Client set Client value
func Client(v interface{}) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Client")
	}
}

// Codecs to be used to encode/decode requests for a given content type
func Codecs(ct string, v interface{}) Option {
	return func(src interface{}) error {
		cm, err := Get(src, ".Codecs")
		if err != nil {
			return err
		} else if rutil.IsZero(cm) {
			cm = reflect.MakeMap(reflect.TypeOf(cm)).Interface()
		}
		cv := reflect.ValueOf(cm)
		cv.SetMapIndex(reflect.ValueOf(ct), reflect.ValueOf(v))
		return Set(src, cv.Interface(), ".Codecs")
	}
}

// Context set Context value
func Context(v context.Context) Option {
	return func(src interface{}) error {
		return Set(src, v, ".Context")
	}
}

// TLSConfig set TLSConfig value
func TLSConfig(v *tls.Config) Option {
	return func(src interface{}) error {
		return Set(src, v, ".TLSConfig")
	}
}

func ContextOption(k, v interface{}) Option {
	return func(src interface{}) error {
		ctx, err := Get(src, ".Context")
		if err != nil {
			return err
		}
		if ctx == nil {
			ctx = context.Background()
		}
		err = Set(src, context.WithValue(ctx.(context.Context), k, v), ".Context")
		return err
	}
}

// ContentType pass ContentType for message data
func ContentType(ct string) Option {
	return func(src interface{}) error {
		return Set(src, ct, ".ContentType")
	}
}

// Metadata pass additional metadata
func Metadata(md ...any) Option {
	var result metadata.Metadata
	if len(md) == 1 {
		switch vt := md[0].(type) {
		case metadata.Metadata:
			result = metadata.Copy(vt)
		case map[string]string:
			result = make(metadata.Metadata, len(vt))
			for k, v := range vt {
				result.Set(k, v)
			}
		case map[string][]string:
			result = metadata.Copy(vt)
		default:
			result = metadata.New(0)
		}
	} else {
		result = metadata.New(len(md) / 2)
		for idx := 0; idx <= len(md)/2; idx += 2 {
			k, ok := md[idx].(string)
			switch vt := md[idx+1].(type) {
			case string:
				if ok {
					result.Set(k, vt)
				}
			case []string:
				if ok {
					result.Append(k, vt...)
				}
			}
		}
	}

	return func(src interface{}) error {
		return Set(src, result, ".Metadata")
	}
}

// Namespace to use
func Namespace(ns string) Option {
	return func(src interface{}) error {
		return Set(src, ns, ".Namespace")
	}
}

// Labels sets the labels
func Labels(ls ...interface{}) Option {
	return func(src interface{}) error {
		v, err := Get(src, ".Labels")
		if err != nil {
			return err
		} else if rutil.IsZero(v) {
			v = reflect.MakeSlice(reflect.TypeOf(v), 0, len(ls)).Interface()
		}
		cv := reflect.ValueOf(v)
		for _, l := range ls {
			cv = reflect.Append(cv, reflect.ValueOf(l))
		}
		return Set(src, cv.Interface(), ".Labels")
	}
}

// Timeout pass timeout time.Duration
func Timeout(td time.Duration) Option {
	return func(src interface{}) error {
		return Set(src, td, ".Timeout")
	}
}
