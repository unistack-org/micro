package random

import (
	"testing"

	"github.com/unistack-org/micro/v3/selector"
)

func TestRandom(t *testing.T) {
	selector.Tests(t, NewSelector())
}
