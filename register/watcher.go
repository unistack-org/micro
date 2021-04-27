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
	Service *Service
	// Action holds the action
	Action string
}

// EventType defines register event type
type EventType int

const (
	// Create is emitted when a new service is registered
	Create EventType = iota
	// Delete is emitted when an existing service is deregistered
	Delete
	// Update is emitted when an existing service is updated
	Update
)

// String returns human readable event type
func (t EventType) String() string {
	switch t {
	case Create:
		return "create"
	case Delete:
		return "delete"
	case Update:
		return "update"
	default:
		return "unknown"
	}
}

// Event is register event
type Event struct {
	// Timestamp is event timestamp
	Timestamp time.Time
	// Service is register service
	Service *Service
	// ID is register id
	ID string
	// Type defines type of event
	Type EventType
}
