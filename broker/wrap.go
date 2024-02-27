package broker

import (
	"reflect"
)

func As(b Broker, target any) bool {
	if b == nil {
		return false
	}
	if target == nil {
		return false
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}
	targetType := typ.Elem()
	if targetType.Kind() != reflect.Interface && !targetType.Implements(brokerType) {
		return false
	}
	return as(b, val, targetType)
}

func as(b Broker, targetVal reflect.Value, targetType reflect.Type) bool {
	if reflect.TypeOf(b).AssignableTo(targetType) {
		targetVal.Elem().Set(reflect.ValueOf(b))
		return true
	}
	return false
}

var brokerType = reflect.TypeOf((*Broker)(nil)).Elem()
