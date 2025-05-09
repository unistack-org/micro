package broker

import (
	"context"
	"strings"
	"sync"

	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
)

type NoopBroker struct {
	funcPublish   FuncPublish
	funcSubscribe FuncSubscribe
	opts          Options
	sync.RWMutex
}

func (b *NoopBroker) newCodec(ct string) (codec.Codec, error) {
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}
	b.RLock()
	c, ok := b.opts.Codecs[ct]
	b.RUnlock()
	if ok {
		return c, nil
	}
	return nil, codec.ErrUnknownContentType
}

func NewBroker(opts ...Option) *NoopBroker {
	b := &NoopBroker{opts: NewOptions(opts...)}
	b.funcPublish = b.fnPublish
	b.funcSubscribe = b.fnSubscribe

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
	b.funcSubscribe = b.fnSubscribe

	b.opts.Hooks.EachPrev(func(hook options.Hook) {
		switch h := hook.(type) {
		case HookPublish:
			b.funcPublish = h(b.funcPublish)
		case HookSubscribe:
			b.funcSubscribe = h(b.funcSubscribe)
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

type noopMessage struct {
	c    codec.Codec
	ctx  context.Context
	body []byte
	hdr  metadata.Metadata
	opts PublishOptions
}

func (m *noopMessage) Ack() error {
	return nil
}

func (m *noopMessage) Body() []byte {
	return m.body
}

func (m *noopMessage) Header() metadata.Metadata {
	return m.hdr
}

func (m *noopMessage) Context() context.Context {
	return m.ctx
}

func (m *noopMessage) Topic() string {
	return ""
}

func (m *noopMessage) Unmarshal(dst interface{}, opts ...codec.Option) error {
	return m.c.Unmarshal(m.body, dst)
}

func (b *NoopBroker) NewMessage(ctx context.Context, hdr metadata.Metadata, body interface{}, opts ...PublishOption) (Message, error) {
	options := NewPublishOptions(opts...)
	if options.ContentType == "" {
		options.ContentType = b.opts.ContentType
	}
	m := &noopMessage{ctx: ctx, hdr: hdr, opts: options}
	c, err := b.newCodec(m.opts.ContentType)
	if err == nil {
		m.body, err = c.Marshal(body)
	}
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (b *NoopBroker) fnPublish(_ context.Context, _ string, _ ...Message) error {
	return nil
}

func (b *NoopBroker) Publish(ctx context.Context, topic string, msg ...Message) error {
	return b.funcPublish(ctx, topic, msg...)
}

type NoopSubscriber struct {
	ctx     context.Context
	topic   string
	handler interface{}
	opts    SubscribeOptions
}

func (b *NoopBroker) fnSubscribe(ctx context.Context, topic string, handler interface{}, opts ...SubscribeOption) (Subscriber, error) {
	return &NoopSubscriber{ctx: ctx, topic: topic, opts: NewSubscribeOptions(opts...), handler: handler}, nil
}

func (b *NoopBroker) Subscribe(ctx context.Context, topic string, handler interface{}, opts ...SubscribeOption) (Subscriber, error) {
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
