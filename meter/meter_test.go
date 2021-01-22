package meter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopMeter(t *testing.T) {
	meter := NewMeter(Path("/noop"))
	assert.NotNil(t, meter)
	assert.Equal(t, "/noop", meter.Options().Path)
	assert.Implements(t, new(Meter), meter)
}
