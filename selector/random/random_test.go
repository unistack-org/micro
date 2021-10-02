package random

import (
	"testing"

	"go.unistack.org/micro/v3/selector"
)

func TestRandom(t *testing.T) {
	selector.Tests(t, NewSelector())
}
