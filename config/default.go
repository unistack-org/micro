package config

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type defaultConfig struct {
	opts Options
}

func (c *defaultConfig) Options() Options {
	return c.opts
}

func (c *defaultConfig) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	return nil
}

func (c *defaultConfig) Load(ctx context.Context, opts ...LoadOption) error {
	for _, fn := range c.opts.BeforeLoad {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	options := NewLoadOptions(opts...)
	mopts := []func(*mergo.Config){mergo.WithTypeCheck}
	if options.Override {
		mopts = append(mopts, mergo.WithOverride)
	}
	if options.Append {
		mopts = append(mopts, mergo.WithAppendSlice)
	}

	dst := c.opts.Struct
	if options.Struct != nil {
		dst = options.Struct
	}

	src, err := rutil.Zero(dst)
	if err == nil {
		if err = fillValues(reflect.ValueOf(src), c.opts.StructTag); err == nil {
			err = mergo.Merge(dst, src, mopts...)
		}
	}

	if err != nil && !c.opts.AllowFail {
		return err
	}

	for _, fn := range c.opts.AfterLoad {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	return nil
}

//nolint:gocyclo
func fillValue(value reflect.Value, val string) error {
	if !rutil.IsEmpty(value) {
		return nil
	}
	switch value.Kind() {
	case reflect.Map:
		t := value.Type()
		nvals := strings.FieldsFunc(val, func(c rune) bool { return c == ',' || c == ';' })
		if value.IsNil() {
			value.Set(reflect.MakeMapWithSize(t, len(nvals)))
		}
		kt := t.Key()
		et := t.Elem()
		for _, nval := range nvals {
			kv := strings.FieldsFunc(nval, func(c rune) bool { return c == '=' })
			mkey := reflect.Indirect(reflect.New(kt))
			mval := reflect.Indirect(reflect.New(et))
			if err := fillValue(mkey, kv[0]); err != nil {
				return err
			}
			if err := fillValue(mval, kv[1]); err != nil {
				return err
			}
			value.SetMapIndex(mkey, mval)
		}
	case reflect.Slice, reflect.Array:
		nvals := strings.FieldsFunc(val, func(c rune) bool { return c == ',' || c == ';' })
		value.Set(reflect.MakeSlice(reflect.SliceOf(value.Type().Elem()), len(nvals), len(nvals)))
		for idx, nval := range nvals {
			nvalue := reflect.Indirect(reflect.New(value.Type().Elem()))
			if err := fillValue(nvalue, nval); err != nil {
				return err
			}
			value.Index(idx).Set(nvalue)
		}
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	case reflect.String:
		value.Set(reflect.ValueOf(val))
	case reflect.Float32:
		v, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(float32(v)))
	case reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	case reflect.Int:
		v, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int(v)))
	case reflect.Int8:
		v, err := strconv.ParseInt(val, 10, 8)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	case reflect.Int16:
		v, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int16(v)))
	case reflect.Int32:
		v, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(int32(v)))
	case reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	case reflect.Uint:
		v, err := strconv.ParseUint(val, 10, 0)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint(v)))
	case reflect.Uint8:
		v, err := strconv.ParseUint(val, 10, 8)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint8(v)))
	case reflect.Uint16:
		v, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint16(v)))
	case reflect.Uint32:
		v, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(uint32(v)))
	case reflect.Uint64:
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		value.Set(reflect.ValueOf(v))
	}
	return nil
}

func fillValues(valueOf reflect.Value, tname string) error {
	var values reflect.Value

	if valueOf.Kind() == reflect.Ptr {
		values = valueOf.Elem()
	} else {
		values = valueOf
	}

	if values.Kind() == reflect.Invalid {
		return ErrInvalidStruct
	}

	fields := values.Type()

	for idx := 0; idx < fields.NumField(); idx++ {
		field := fields.Field(idx)
		value := values.Field(idx)
		if !value.CanSet() {
			continue
		}
		if len(field.PkgPath) != 0 {
			continue
		}
		switch value.Kind() {
		case reflect.Struct:
			value.Set(reflect.Indirect(reflect.New(value.Type())))
			if err := fillValues(value, tname); err != nil {
				return err
			}
			continue
		case reflect.Ptr:
			if value.IsNil() {
				if value.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone
					break
				}
				// nil pointer to struct: create a zero instance
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
			if err := fillValues(value, tname); err != nil {
				return err
			}
			continue
		}
		tag, ok := field.Tag.Lookup(tname)
		if !ok {
			continue
		}

		if err := fillValue(value, tag); err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultConfig) Save(ctx context.Context, opts ...SaveOption) error {
	for _, fn := range c.opts.BeforeSave {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	for _, fn := range c.opts.AfterSave {
		if err := fn(ctx, c); err != nil && !c.opts.AllowFail {
			return err
		}
	}

	return nil
}

func (c *defaultConfig) String() string {
	return "default"
}

func (c *defaultConfig) Name() string {
	return c.opts.Name
}

func (c *defaultConfig) Watch(ctx context.Context, opts ...WatchOption) (Watcher, error) {
	w := &defaultWatcher{
		opts:  c.opts,
		wopts: NewWatchOptions(opts...),
		done:  make(chan struct{}),
		vchan: make(chan map[string]interface{}),
		echan: make(chan error),
	}

	go w.run()

	return w, nil
}

// NewConfig returns new default config source
func NewConfig(opts ...Option) Config {
	options := NewOptions(opts...)
	if len(options.StructTag) == 0 {
		options.StructTag = "default"
	}
	return &defaultConfig{opts: options}
}
