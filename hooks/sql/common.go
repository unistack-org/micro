package sql

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"runtime"
)

//go:generate sh -c "go run gen.go > wrap_gen.go"

// namedValueToValue converts driver arguments of NamedValue format to Value format. Implemented in the same way as in
// database/sql ctxutil.go.
func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}

// namedValueToLabels convert driver arguments to interface{} slice
func namedValueToLabels(named []driver.NamedValue) []interface{} {
	largs := make([]interface{}, 0, len(named)*2)
	var name string
	for _, param := range named {
		if param.Name != "" {
			name = param.Name
		} else {
			name = fmt.Sprintf("$%d", param.Ordinal)
		}
		largs = append(largs, fmt.Sprintf("%s=%v", name, param.Value))
	}
	return largs
}

// getCallerName get the name of the function A where A() -> B() -> GetFunctionCallerName()
func getCallerName() string {
	pc, _, _, ok := runtime.Caller(3)
	details := runtime.FuncForPC(pc)
	var callerName string
	if ok && details != nil {
		callerName = details.Name()
	} else {
		callerName = labelUnknown
	}
	return callerName
}
