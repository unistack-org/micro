package logger // import "go.unistack.org/micro/v3/logger/wrapper"

import (
	"context"
	"reflect"

	rutil "go.unistack.org/micro/v3/util/reflect"
)

// LogFunc function used for Log method
type LogFunc func(ctx context.Context, level Level, args ...interface{})

// LogfFunc function used for Logf method
type LogfFunc func(ctx context.Context, level Level, msg string, args ...interface{})

type Wrapper interface {
	// Log logs message with needed level
	Log(LogFunc) LogFunc
	// Logf logs message with needed level
	Logf(LogfFunc) LogfFunc
}

var _ Logger = &OmitLogger{}

type OmitLogger struct {
	l Logger
}

func NewOmitLogger(l Logger) Logger {
	return &OmitLogger{l: l}
}

func (w *OmitLogger) Init(opts ...Option) error {
	return w.l.Init(append(opts, WrapLogger(NewOmitWrapper()))...)
}

func (w *OmitLogger) V(level Level) bool {
	return w.l.V(level)
}

func (w *OmitLogger) Level(level Level) {
	w.l.Level(level)
}

func (w *OmitLogger) Clone(opts ...Option) Logger {
	return w.l.Clone(opts...)
}

func (w *OmitLogger) Options() Options {
	return w.l.Options()
}

func (w *OmitLogger) Fields(fields ...interface{}) Logger {
	return w.l.Fields(fields...)
}

func (w *OmitLogger) Info(ctx context.Context, args ...interface{}) {
	w.l.Info(ctx, args...)
}

func (w *OmitLogger) Trace(ctx context.Context, args ...interface{}) {
	w.l.Trace(ctx, args...)
}

func (w *OmitLogger) Debug(ctx context.Context, args ...interface{}) {
	w.l.Debug(ctx, args...)
}

func (w *OmitLogger) Warn(ctx context.Context, args ...interface{}) {
	w.l.Warn(ctx, args...)
}

func (w *OmitLogger) Error(ctx context.Context, args ...interface{}) {
	w.l.Error(ctx, args...)
}

func (w *OmitLogger) Fatal(ctx context.Context, args ...interface{}) {
	w.l.Fatal(ctx, args...)
}

func (w *OmitLogger) Infof(ctx context.Context, msg string, args ...interface{}) {
	w.l.Infof(ctx, msg, args...)
}

func (w *OmitLogger) Tracef(ctx context.Context, msg string, args ...interface{}) {
	w.l.Tracef(ctx, msg, args...)
}

func (w *OmitLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {
	w.l.Debugf(ctx, msg, args...)
}

func (w *OmitLogger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	w.l.Warnf(ctx, msg, args...)
}

func (w *OmitLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	w.l.Errorf(ctx, msg, args...)
}

func (w *OmitLogger) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	w.l.Fatalf(ctx, msg, args...)
}

func (w *OmitLogger) Log(ctx context.Context, level Level, args ...interface{}) {
	w.l.Log(ctx, level, args...)
}

func (w *OmitLogger) Logf(ctx context.Context, level Level, msg string, args ...interface{}) {
	w.l.Logf(ctx, level, msg, args...)
}

func (w *OmitLogger) String() string {
	return w.l.String()
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
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		narg := arg
		if val.Kind() != reflect.Struct {
			nargs = append(nargs, narg)
			continue
		}

		if narg, err = rutil.Zero(arg); err != nil {
			nargs = append(nargs, narg)
			continue
		}

		rutil.CopyDefaults(narg, arg)
		if flds, ferr := rutil.StructFields(narg); ferr == nil {
			for _, fld := range flds {
				if tv, ok := fld.Field.Tag.Lookup("logger"); ok && tv == "omit" {
					fld.Value.Set(reflect.Zero(fld.Value.Type()))
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
