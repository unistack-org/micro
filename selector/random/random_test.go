package random

import (
	"testing"

	"go.unistack.org/micro/v4/selector"
)

func TestRandom(t *testing.T) {
	selector.Tests(t, NewSelector())
}
