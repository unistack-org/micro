// Package flow is an interface used for saga pattern messaging
package flow

type Step interface {
	// Endpoint returns service_name.service_method
	Endpoint() string
}
