package client

import (
	"context"
)

type clientCallOptions struct {
	Client
	opts []CallOption
}

func (s *clientCallOptions) Call(ctx context.Context, req Request, rsp interface{}, opts ...CallOption) error {
	return s.Client.Call(ctx, req, rsp, append(s.opts, opts...)...)
}

func (s *clientCallOptions) Stream(ctx context.Context, req Request, opts ...CallOption) (Stream, error) {
	return s.Client.Stream(ctx, req, append(s.opts, opts...)...)
}

// NewClientCallOptions add CallOption to every call
func NewClientCallOptions(c Client, opts ...CallOption) Client {
	return &clientCallOptions{c, opts}
}
