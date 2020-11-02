package client

import (
	"context"
	"time"

	"github.com/unistack-org/micro/v3/util/backoff"
)

// BackoffFunc is the backoff call func
type BackoffFunc func(ctx context.Context, req Request, attempts int) (time.Duration, error)

func exponentialBackoff(ctx context.Context, req Request, attempts int) (time.Duration, error) {
	return backoff.Do(attempts), nil
}
