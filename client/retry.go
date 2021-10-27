package client

import (
	"context"

	"go.unistack.org/micro/v3/errors"
)

// RetryFunc that returning either false or a non-nil error will result in the call not being retried
type RetryFunc func(ctx context.Context, req Request, retryCount int, err error) (bool, error)

// RetryAlways always retry on error
func RetryAlways(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	return true, nil
}

// RetryNever never retry on error
func RetryNever(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	return false, nil
}

// RetryOnError retries a request on a 500 or timeout error
func RetryOnError(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	me := errors.FromError(err)
	switch me.Code {
	// retry on timeout or internal server error
	case 408, 500:
		return true, nil
	}

	return false, nil
}
