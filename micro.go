package micro

import (
	"reflect"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/flow"
	"go.unistack.org/micro/v4/fsm"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/meter"
	"go.unistack.org/micro/v4/register"
	"go.unistack.org/micro/v4/resolver"
	"go.unistack.org/micro/v4/router"
	"go.unistack.org/micro/v4/selector"
	"go.unistack.org/micro/v4/server"
	"go.unistack.org/micro/v4/store"
	"go.unistack.org/micro/v4/sync"
	"go.unistack.org/micro/v4/tracer"
)

func As(b any, target any) bool {
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
	if targetType.Kind() != reflect.Interface &&
		!(targetType.Implements(brokerType) ||
			targetType.Implements(loggerType) ||
			targetType.Implements(clientType) ||
			targetType.Implements(serverType) ||
			targetType.Implements(codecType) ||
			targetType.Implements(flowType) ||
			targetType.Implements(fsmType) ||
			targetType.Implements(meterType) ||
			targetType.Implements(registerType) ||
			targetType.Implements(resolverType) ||
			targetType.Implements(selectorType) ||
			targetType.Implements(storeType) ||
			targetType.Implements(syncType) ||
			targetType.Implements(tracerType) ||
			targetType.Implements(serviceType) ||
			targetType.Implements(routerType)) {
		return false
	}
	if reflect.TypeOf(b).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(b))
		return true
	}
	return false
}

var brokerType = reflect.TypeOf((*broker.Broker)(nil)).Elem()
var loggerType = reflect.TypeOf((*logger.Logger)(nil)).Elem()
var clientType = reflect.TypeOf((*client.Client)(nil)).Elem()
var serverType = reflect.TypeOf((*server.Server)(nil)).Elem()
var codecType = reflect.TypeOf((*codec.Codec)(nil)).Elem()
var flowType = reflect.TypeOf((*flow.Flow)(nil)).Elem()
var fsmType = reflect.TypeOf((*fsm.FSM)(nil)).Elem()
var meterType = reflect.TypeOf((*meter.Meter)(nil)).Elem()
var registerType = reflect.TypeOf((*register.Register)(nil)).Elem()
var resolverType = reflect.TypeOf((*resolver.Resolver)(nil)).Elem()
var routerType = reflect.TypeOf((*router.Router)(nil)).Elem()
var selectorType = reflect.TypeOf((*selector.Selector)(nil)).Elem()
var storeType = reflect.TypeOf((*store.Store)(nil)).Elem()
var syncType = reflect.TypeOf((*sync.Sync)(nil)).Elem()
var tracerType = reflect.TypeOf((*tracer.Tracer)(nil)).Elem()
var serviceType = reflect.TypeOf((*Service)(nil)).Elem()
