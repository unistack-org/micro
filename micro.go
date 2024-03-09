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
	if targetType.Kind() != reflect.Interface {
		switch {
		case targetType.Implements(brokerType):
			break
		case targetType.Implements(loggerType):
			break
		case targetType.Implements(clientType):
			break
		case targetType.Implements(serverType):
			break
		case targetType.Implements(codecType):
			break
		case targetType.Implements(flowType):
			break
		case targetType.Implements(fsmType):
			break
		case targetType.Implements(meterType):
			break
		case targetType.Implements(registerType):
			break
		case targetType.Implements(resolverType):
			break
		case targetType.Implements(selectorType):
			break
		case targetType.Implements(storeType):
			break
		case targetType.Implements(syncType):
			break
		case targetType.Implements(serviceType):
			break
		case targetType.Implements(routerType):
			break
		case targetType.Implements(tracerType):
			break
		default:
			return false
		}
	}
	if reflect.TypeOf(b).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(b))
		return true
	}
	return false
}

var (
	brokerType   = reflect.TypeOf((*broker.Broker)(nil)).Elem()
	loggerType   = reflect.TypeOf((*logger.Logger)(nil)).Elem()
	clientType   = reflect.TypeOf((*client.Client)(nil)).Elem()
	serverType   = reflect.TypeOf((*server.Server)(nil)).Elem()
	codecType    = reflect.TypeOf((*codec.Codec)(nil)).Elem()
	flowType     = reflect.TypeOf((*flow.Flow)(nil)).Elem()
	fsmType      = reflect.TypeOf((*fsm.FSM)(nil)).Elem()
	meterType    = reflect.TypeOf((*meter.Meter)(nil)).Elem()
	registerType = reflect.TypeOf((*register.Register)(nil)).Elem()
	resolverType = reflect.TypeOf((*resolver.Resolver)(nil)).Elem()
	routerType   = reflect.TypeOf((*router.Router)(nil)).Elem()
	selectorType = reflect.TypeOf((*selector.Selector)(nil)).Elem()
	storeType    = reflect.TypeOf((*store.Store)(nil)).Elem()
	syncType     = reflect.TypeOf((*sync.Sync)(nil)).Elem()
	tracerType   = reflect.TypeOf((*tracer.Tracer)(nil)).Elem()
	serviceType  = reflect.TypeOf((*Service)(nil)).Elem()
)
