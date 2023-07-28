package client

import (
	"context"

	"go.unistack.org/micro/v4/options"
)

type clientCallOptions struct {
	Client
	opts []options.Option
}

func (s *clientCallOptions) Call(ctx context.Context, req Request, rsp interface{}, opts ...options.Option) error {
	return s.Client.Call(ctx, req, rsp, append(s.opts, opts...)...)
}

func (s *clientCallOptions) Stream(ctx context.Context, req Request, opts ...options.Option) (Stream, error) {
	return s.Client.Stream(ctx, req, append(s.opts, opts...)...)
}

// NewClientCallOptions add CallOption to every call
func NewClientCallOptions(c Client, opts ...options.Option) Client {
	return &clientCallOptions{c, opts}
}
