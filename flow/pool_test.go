package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoolSizeOption(t *testing.T) {
	opts := NewOptions(PoolSize(8))
	assert.Equal(t, 8, opts.PoolSize)
}
