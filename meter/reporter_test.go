package meter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopReporter(t *testing.T) {
	// Make a Reporter:
	reporter := NewReporter(Path("/noop"))
	assert.NotNil(t, reporter)
	assert.Equal(t, "/noop", reporter.Options().Path)

	// Check that our implementation is valid:
	assert.Implements(t, new(Reporter), reporter)
}
