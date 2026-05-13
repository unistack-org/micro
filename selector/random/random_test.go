package random

import (
	"testing"

	"go.unistack.org/micro/v5/selector"
)

func TestRandom(t *testing.T) {
	selector.Tests(t, NewSelector())
}
