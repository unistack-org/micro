package broker

import (
	"context"
	"strings"
)

type NoopBroker struct {
	opts Options
}

func NewBroker(opts ...Option) *NoopBroker {
	b := &NoopBroker{opts: NewOptions(opts...)}
	return b
}

func (b *NoopBroker) Name() string {
	return b.opts.Name
}

func (b *NoopBroker) String() string {
	return "noop"
}

func (b *NoopBroker) Options() Options {
	return b.opts
}

func (b *NoopBroker) Init(opts ...Option) error {
	for _, opt := range opts {
		opt(&b.opts)
	}
	return nil
}

func (b *NoopBroker) Connect(_ context.Context) error {
	return nil
}

func (b *NoopBroker) Disconnect(_ context.Context) error {
	return nil
}

func (b *NoopBroker) Address() string {
	return strings.Join(b.opts.Addrs, ",")
}

func (b *NoopBroker) BatchPublish(_ context.Context, _ []*Message, _ ...PublishOption) error {
	return nil
}

func (b *NoopBroker) Publish(_ context.Context, _ string, _ *Message, _ ...PublishOption) error {
	return nil
}

type NoopSubscriber struct {
	ctx          context.Context
	topic        string
	handler      Handler
	batchHandler BatchHandler
	opts         SubscribeOptions
}

func (b *NoopBroker) BatchSubscribe(ctx context.Context, topic string, handler BatchHandler, opts ...SubscribeOption) (Subscriber, error) {
	return &NoopSubscriber{ctx: ctx, topic: topic, opts: NewSubscribeOptions(opts...), batchHandler: handler}, nil
}

func (b *NoopBroker) Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	return &NoopSubscriber{ctx: ctx, topic: topic, opts: NewSubscribeOptions(opts...), handler: handler}, nil
}

func (s *NoopSubscriber) Options() SubscribeOptions {
	return s.opts
}

func (s *NoopSubscriber) Topic() string {
	return s.topic
}

func (s *NoopSubscriber) Unsubscribe(ctx context.Context) error {
	return nil
}
