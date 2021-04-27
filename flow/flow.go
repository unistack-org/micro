// Package flow is an interface used for saga pattern microservice workflow
package flow

type Step interface {
	// Endpoint returns rpc endpoint service_name.service_method or broker topic
	Endpoint() string
}

type Workflow interface {
	Steps() [][]Step
	Stop() error
}

type Flow interface {
	Start(Workflow) error
	Stop(Workflow)
}
