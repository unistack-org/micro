// Package jitter provides a random jitter
package jitter

import (
	"time"

	"github.com/unistack-org/micro/v3/util/rand"
)

// Do returns a random time to jitter with max cap specified
func Do(d time.Duration) time.Duration {
	var rng rand.Rand
	v := rng.Float64() * float64(d.Nanoseconds())
	return time.Duration(v)
}
