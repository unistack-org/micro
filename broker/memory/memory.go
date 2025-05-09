package broker

import (
	"context"
	"strings"
	"sync"

	"go.unistack.org/micro/v4/broker"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/logger"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/options"
	maddr "go.unistack.org/micro/v4/util/addr"
	"go.unistack.org/micro/v4/util/id"
	mnet "go.unistack.org/micro/v4/util/net"
	"go.unistack.org/micro/v4/util/rand"
)

type Broker struct {
	funcPublish   broker.FuncPublish
	funcSubscribe broker.FuncSubscribe
	subscribers   map[string][]*Subscriber
	addr          string
	opts          broker.Options
	sync.RWMutex
	connected bool
}

type memoryMessage struct {
	c     codec.Codec
	topic string
	ctx   context.Context
	body  []byte
	hdr   metadata.Metadata
	opts  broker.PublishOptions
}

func (m *memoryMessage) Ack() error {
	return nil
}

func (m *memoryMessage) Body() []byte {
	return m.body
}

func (m *memoryMessage) Header() metadata.Metadata {
	return m.hdr
}

func (m *memoryMessage) Context() context.Context {
	return m.ctx
}

func (m *memoryMessage) Topic() string {
	return ""
}

func (m *memoryMessage) Unmarshal(dst interface{}, opts ...codec.Option) error {
	return m.c.Unmarshal(m.body, dst)
}

type Subscriber struct {
	ctx     context.Context
	exit    chan bool
	handler interface{}
	id      string
	topic   string
	opts    broker.SubscribeOptions
}

func (b *Broker) newCodec(ct string) (codec.Codec, error) {
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

func (b *Broker) Options() broker.Options {
	return b.opts
}

func (b *Broker) Address() string {
	return b.addr
}

func (b *Broker) Connect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.Lock()
	defer b.Unlock()

	if b.connected {
		return nil
	}

	// use 127.0.0.1 to avoid scan of all network interfaces
	addr, err := maddr.Extract("127.0.0.1")
	if err != nil {
		return err
	}
	var rng rand.Rand
	i := rng.Intn(20000)
	// set addr with port
	addr = mnet.HostPort(addr, 10000+i)

	b.addr = addr
	b.connected = true

	return nil
}

func (b *Broker) Disconnect(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.Lock()
	defer b.Unlock()

	if !b.connected {
		return nil
	}

	b.connected = false
	return nil
}

func (b *Broker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&b.opts)
	}

	b.funcPublish = b.fnPublish
	b.funcSubscribe = b.fnSubscribe

	b.opts.Hooks.EachPrev(func(hook options.Hook) {
		switch h := hook.(type) {
		case broker.HookPublish:
			b.funcPublish = h(b.funcPublish)
		case broker.HookSubscribe:
			b.funcSubscribe = h(b.funcSubscribe)
		}
	})

	return nil
}

func (b *Broker) NewMessage(ctx context.Context, hdr metadata.Metadata, body interface{}, opts ...broker.PublishOption) (broker.Message, error) {
	options := broker.NewPublishOptions(opts...)
	if options.ContentType == "" {
		options.ContentType = b.opts.ContentType
	}
	m := &memoryMessage{ctx: ctx, hdr: hdr, opts: options}
	c, err := b.newCodec(m.opts.ContentType)
	if err == nil {
		m.body, err = c.Marshal(body)
	}
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (b *Broker) Publish(ctx context.Context, topic string, messages ...broker.Message) error {
	return b.funcPublish(ctx, topic, messages...)
}

func (b *Broker) fnPublish(ctx context.Context, topic string, messages ...broker.Message) error {
	return b.publish(ctx, topic, messages...)
}

func (b *Broker) publish(ctx context.Context, topic string, messages ...broker.Message) error {
	b.RLock()
	if !b.connected {
		b.RUnlock()
		return broker.ErrNotConnected
	}
	b.RUnlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.RLock()
	subs, ok := b.subscribers[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	var err error

	for _, sub := range subs {
		switch s := sub.handler.(type) {
		default:
			if b.opts.Logger.V(logger.ErrorLevel) {
				b.opts.Logger.Error(ctx, "broker  handler error", broker.ErrInvalidHandler)
			}
		case func(broker.Message) error:
			for _, message := range messages {
				msg, ok := message.(*memoryMessage)
				if !ok {
					if b.opts.Logger.V(logger.ErrorLevel) {
						b.opts.Logger.Error(ctx, "broker handler error", broker.ErrInvalidMessage)
					}
				}
				msg.topic = topic
				if err = s(msg); err == nil && sub.opts.AutoAck {
					err = msg.Ack()
				}
				if err != nil {
					if b.opts.Logger.V(logger.ErrorLevel) {
						b.opts.Logger.Error(ctx, "broker handler error", err)
					}
				}
			}
		case func([]broker.Message) error:
			if err = s(messages); err == nil && sub.opts.AutoAck {
				for _, message := range messages {
					err = message.Ack()
					if err != nil {
						if b.opts.Logger.V(logger.ErrorLevel) {
							b.opts.Logger.Error(ctx, "broker handler error", err)
						}
					}
				}
			}
		}
	}

	return nil
}

func (b *Broker) Subscribe(ctx context.Context, topic string, handler interface{}, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return b.funcSubscribe(ctx, topic, handler, opts...)
}

func (b *Broker) fnSubscribe(ctx context.Context, topic string, handler interface{}, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	if err := broker.IsValidHandler(handler); err != nil {
		return nil, err
	}

	b.RLock()
	if !b.connected {
		b.RUnlock()
		return nil, broker.ErrNotConnected
	}
	b.RUnlock()

	sid, err := id.New()
	if err != nil {
		return nil, err
	}

	options := broker.NewSubscribeOptions(opts...)

	sub := &Subscriber{
		exit:    make(chan bool, 1),
		id:      sid,
		topic:   topic,
		handler: handler,
		opts:    options,
		ctx:     ctx,
	}

	b.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], sub)
	b.Unlock()

	go func() {
		<-sub.exit
		b.Lock()
		newSubscribers := make([]*Subscriber, 0, len(b.subscribers)-1)
		for _, sb := range b.subscribers[topic] {
			if sb.id == sub.id {
				continue
			}
			newSubscribers = append(newSubscribers, sb)
		}
		b.subscribers[topic] = newSubscribers
		b.Unlock()
	}()

	return sub, nil
}

func (b *Broker) String() string {
	return "memory"
}

func (b *Broker) Name() string {
	return b.opts.Name
}

func (b *Broker) Live() bool {
	return true
}

func (b *Broker) Ready() bool {
	return true
}

func (b *Broker) Health() bool {
	return true
}

func (m *Subscriber) Options() broker.SubscribeOptions {
	return m.opts
}

func (m *Subscriber) Topic() string {
	return m.topic
}

func (m *Subscriber) Unsubscribe(ctx context.Context) error {
	m.exit <- true
	return nil
}

// NewBroker return new memory broker
func NewBroker(opts ...broker.Option) broker.Broker {
	return &Broker{
		opts:        broker.NewOptions(opts...),
		subscribers: make(map[string][]*Subscriber),
	}
}
