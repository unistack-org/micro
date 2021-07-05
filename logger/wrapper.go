package logger

import (
	"context"
	"reflect"

	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

// LogFunc function used for Log method
type LogFunc func(ctx context.Context, level Level, args ...interface{})

// LogfFunc function used for Logf method
type LogfFunc func(ctx context.Context, level Level, msg string, args ...interface{})

type Wrapper interface {
	// Log logs message with needed level
	Log(LogFunc) LogFunc
	//	Log(ctx context.Context, level Level, args ...interface{})
	// Logf logs message with needed level
	Logf(LogfFunc) LogfFunc
	//Logf(ctx context.Context, level Level, msg string, args ...interface{})
}

type OmitWrapper struct{}

func NewOmitWrapper() Wrapper {
	return &OmitWrapper{}
}

func getArgs(args []interface{}) []interface{} {
	nargs := make([]interface{}, 0, len(args))
	var err error
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		switch val.Kind() {
		case reflect.Ptr:
			val = val.Elem()
		}
		narg := arg
		if val.Kind() == reflect.Struct {
			if narg, err = rutil.Zero(arg); err == nil {
				rutil.CopyDefaults(narg, arg)
				if flds, ferr := rutil.StructFields(narg); ferr == nil {
					for _, fld := range flds {
						if tv, ok := fld.Field.Tag.Lookup("logger"); ok && tv == "omit" {
							fld.Value.Set(reflect.Zero(fld.Value.Type()))
						}
					}
				}
			}
		}
		nargs = append(nargs, narg)
	}
	return nargs
}

func (w *OmitWrapper) Log(fn LogFunc) LogFunc {
	return func(ctx context.Context, level Level, args ...interface{}) {
		fn(ctx, level, getArgs(args)...)
	}
}

func (w *OmitWrapper) Logf(fn LogfFunc) LogfFunc {
	return func(ctx context.Context, level Level, msg string, args ...interface{}) {
		fn(ctx, level, msg, getArgs(args)...)
	}
}
