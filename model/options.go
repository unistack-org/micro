// Package model is an interface for data modelling
package model

import (
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/logger"
	"github.com/unistack-org/micro/v3/store"
	"github.com/unistack-org/micro/v3/sync"
)

type Options struct {
	// Database to write to
	Database string
	// for serialising
	Codec codec.Codec
	// for locking
	Sync sync.Sync
	// for storage
	Store store.Store
	// for logger
	Logger logger.Logger
}

type Option func(o *Options)

// Logger sets the logger
func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

type ReadOptions struct{}

type ReadOption func(o *ReadOptions)

type DeleteOptions struct{}

type DeleteOption func(o *DeleteOptions)
