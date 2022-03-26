package client

import (
	"context"
	"math"
	"time"

	"go.unistack.org/micro/v3/util/backoff"
)

// BackoffFunc is the backoff call func
type BackoffFunc func(ctx context.Context, req Request, attempts int) (time.Duration, error)

// BackoffExp using exponential backoff func
func BackoffExp(_ context.Context, _ Request, attempts int) (time.Duration, error) {
	return backoff.Do(attempts), nil
}

// BackoffInterval specifies randomization interval for backoff func
func BackoffInterval(min time.Duration, max time.Duration) BackoffFunc {
	return func(_ context.Context, _ Request, attempts int) (time.Duration, error) {
		td := time.Duration(time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100)
		if td < min {
			return min, nil
		} else if td > max {
			return max, nil
		}
		return td, nil
	}
}
