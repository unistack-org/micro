package mock

import (
	"testing"

	"go.unistack.org/micro/v5/meter"
)

func TestMockMeter_Implements(t *testing.T) {
	var _ meter.Meter = NewMockMeter()
}
