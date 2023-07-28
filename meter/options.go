package meter

import (
	"context"
	"reflect"

	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/options"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

// Options for metrics implementations
type Options struct {
	// Logger used for logging
	Logger logger.Logger
	// Context holds external options
	Context context.Context
	// Name holds the meter name
	Name string
	// Address holds the address that serves metrics
	Address string
	// Path holds the path for metrics
	Path string
	// MetricPrefix holds the prefix for all metrics
	MetricPrefix string
	// LabelPrefix holds the prefix for all labels
	LabelPrefix string
	// Labels holds the default labels
	Labels []string
	// WriteProcessMetrics flag to write process metrics
	WriteProcessMetrics bool
	// WriteFDMetrics flag to write fd metrics
	WriteFDMetrics bool
}

// NewOptions prepares a set of options:
func NewOptions(opt ...options.Option) Options {
	opts := Options{
		Address:      DefaultAddress,
		Path:         DefaultPath,
		Context:      context.Background(),
		Logger:       logger.DefaultLogger,
		MetricPrefix: DefaultMetricPrefix,
		LabelPrefix:  DefaultLabelPrefix,
	}

	for _, o := range opt {
		o(&opts)
	}

	return opts
}

// LabelPrefix sets the labels prefix
func LabelPrefix(pref string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, pref, ".LabelPrefix")
	}
}

// MetricPrefix sets the metric prefix
func MetricPrefix(pref string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, pref, ".MetricPrefix")
	}
}

// Path used to serve metrics over HTTP
func Path(path string) options.Option {
	return func(src interface{}) error {
		return options.Set(src, path, ".Path")
	}
}

/*
// TimingObjectives defines the desired spread of statistics for histogram / timing metrics:
func TimingObjectives(value map[float64]float64) Option {
	return func(o *Options) {
		o.TimingObjectives = value
	}
}
*/

// Labels sets the meter labels
func Labels(ls ...string) options.Option {
	return func(src interface{}) error {
		v, err := options.Get(src, ".Labels")
		if err != nil {
			return err
		} else if rutil.IsZero(v) {
			v = reflect.MakeSlice(reflect.TypeOf(v), 0, len(ls)).Interface()
		}
		cv := reflect.ValueOf(v)
		for _, l := range ls {
			reflect.Append(cv, reflect.ValueOf(l))
		}
		err = options.Set(src, cv, ".Labels")
		return err
	}
}

// WriteProcessMetrics enable process metrics output for write
func WriteProcessMetrics(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".WriteProcessMetrics")
	}
}

// WriteFDMetrics enable fd metrics output for write
func WriteFDMetrics(b bool) options.Option {
	return func(src interface{}) error {
		return options.Set(src, b, ".WriteFDMetrics")
	}
}
