// +build ignore

// Package model is an interface for data modelling
package model

// Model provides an interface for data modelling
type Model interface {
	// Initialise options
	Init(...Option) error
	// NewEntity creates a new entity to store or access
	NewEntity(name string, value interface{}) Entity
	// Create a value
	Create(Entity) error
	// Read values
	Read(...ReadOption) ([]Entity, error)
	// Update the value
	Update(Entity) error
	// Delete an entity
	Delete(...DeleteOption) error
	// Implementation of the model
	String() string
}

type Entity interface {
	// Unique id of the entity
	Id() string
	// Name of the entity
	Name() string
	// The value associated with the entity
	Value() interface{}
	// Attributes of the entity
	Attributes() map[string]interface{}
	// Read a value as a concrete type
	Read(v interface{}) error
}
