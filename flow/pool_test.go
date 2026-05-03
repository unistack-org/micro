package flow

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPoolSizeOption(t *testing.T) {
	opts := NewOptions(PoolSize(8))
	assert.Equal(t, 8, opts.PoolSize)
}

func TestNewFlowCreatesPool(t *testing.T) {
	f := NewFlow(PoolSize(4))
	mf, ok := f.(*microFlow)
	require.True(t, ok)
	assert.NotNil(t, mf.pool)
	assert.Equal(t, 4, mf.pool.Cap())
	require.NoError(t, f.Close())
}

func TestFlowClose(t *testing.T) {
	f := NewFlow(PoolSize(2))
	require.NoError(t, f.Close())
	// Second close must not panic
	require.NoError(t, f.Close())
}
