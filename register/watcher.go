package register

import "time"

// Watcher is an interface that returns updates
// about services within the register.
type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	// Stop stops the watcher
	Stop()
}

// Result is returned by a call to Next on
// the watcher. Actions can be create, update, delete
type Result struct {
	// Service holds register service
	Service *Service `json:"service,omitempty"`
	// Action holds the action
	Action EventType `json:"action,omitempty"`
}

// EventType defines register event type
type EventType int

const (
	// EventCreate is emitted when a new service is registered
	EventCreate EventType = iota
	// EventDelete is emitted when an existing service is deregistered
	EventDelete
	// EventUpdate is emitted when an existing service is updated
	EventUpdate
)

// String returns human readable event type
func (t EventType) String() string {
	switch t {
	case EventCreate:
		return "create"
	case EventDelete:
		return "delete"
	case EventUpdate:
		return "update"
	default:
		return "unknown"
	}
}

// Event is register event
type Event struct {
	// Timestamp is event timestamp
	Timestamp time.Time `json:"timestamp,omitempty"`
	// Service is register service
	Service *Service `json:"service,omitempty"`
	// ID is register id
	ID string `json:"id,omitempty"`
	// Type defines type of event
	Type EventType `json:"type,omitempty"`
}
