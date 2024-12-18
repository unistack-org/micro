package broker

import (
	"context"
	"strings"

	"go.unistack.org/micro/v3/options"
)

type NoopBroker struct {
	funcPublish        FuncPublish
	funcBatchPublish   FuncBatchPublish
	funcSubscribe      FuncSubscribe
	funcBatchSubscribe FuncBatchSubscribe
	opts               Options
}

func NewBroker(opts ...Option) *NoopBroker {
	b := &NoopBroker{opts: NewOptions(opts...)}
	b.funcPublish = b.fnPublish
	b.funcBatchPublish = b.fnBatchPublish
	b.funcSubscribe = b.fnSubscribe
	b.funcBatchSubscribe = b.fnBatchSubscribe

	return b
}

func (b *NoopBroker) Health() bool {
	return true
}

func (b *NoopBroker) Live() bool {
	return true
}

func (b *NoopBroker) Ready() bool {
	return true
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

	b.funcPublish = b.fnPublish
	b.funcBatchPublish = b.fnBatchPublish
	b.funcSubscribe = b.fnSubscribe
	b.funcBatchSubscribe = b.fnBatchSubscribe

	b.opts.Hooks.EachPrev(func(hook options.Hook) {
		switch h := hook.(type) {
		case HookPublish:
			b.funcPublish = h(b.funcPublish)
		case HookBatchPublish:
			b.funcBatchPublish = h(b.funcBatchPublish)
		case HookSubscribe:
			b.funcSubscribe = h(b.funcSubscribe)
		case HookBatchSubscribe:
			b.funcBatchSubscribe = h(b.funcBatchSubscribe)
		}
	})

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

func (b *NoopBroker) fnBatchPublish(_ context.Context, _ []*Message, _ ...PublishOption) error {
	return nil
}

func (b *NoopBroker) BatchPublish(ctx context.Context, msgs []*Message, opts ...PublishOption) error {
	return b.funcBatchPublish(ctx, msgs, opts...)
}

func (b *NoopBroker) fnPublish(_ context.Context, _ string, _ *Message, _ ...PublishOption) error {
	return nil
}

func (b *NoopBroker) Publish(ctx context.Context, topic string, msg *Message, opts ...PublishOption) error {
	return b.funcPublish(ctx, topic, msg, opts...)
}

type NoopSubscriber struct {
	ctx          context.Context
	topic        string
	handler      Handler
	batchHandler BatchHandler
	opts         SubscribeOptions
}

func (b *NoopBroker) fnBatchSubscribe(ctx context.Context, topic string, handler BatchHandler, opts ...SubscribeOption) (Subscriber, error) {
	return &NoopSubscriber{ctx: ctx, topic: topic, opts: NewSubscribeOptions(opts...), batchHandler: handler}, nil
}

func (b *NoopBroker) BatchSubscribe(ctx context.Context, topic string, handler BatchHandler, opts ...SubscribeOption) (Subscriber, error) {
	return b.funcBatchSubscribe(ctx, topic, handler, opts...)
}

func (b *NoopBroker) fnSubscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	return &NoopSubscriber{ctx: ctx, topic: topic, opts: NewSubscribeOptions(opts...), handler: handler}, nil
}

func (b *NoopBroker) Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	return b.funcSubscribe(ctx, topic, handler, opts...)
}

func (s *NoopSubscriber) Options() SubscribeOptions {
	return s.opts
}

func (s *NoopSubscriber) Topic() string {
	return s.topic
}

func (s *NoopSubscriber) Unsubscribe(_ context.Context) error {
	return nil
}
