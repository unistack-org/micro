package noop

import (
	"testing"

	"github.com/unistack-org/micro/v3/metrics"

	"github.com/stretchr/testify/assert"
)

func TestNoopReporter(t *testing.T) {

	// Make a Reporter:
	reporter := New(metrics.Path("/noop"))
	assert.NotNil(t, reporter)
	assert.Equal(t, "/noop", reporter.options.Path)

	// Check that our implementation is valid:
	assert.Implements(t, new(metrics.Reporter), reporter)
}
