// Package mock provides a mock implementation of broker.Broker for testing.
package mock

import (
	"context"

	"go.unistack.org/micro/v5/broker"
	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/metadata"
)

type mockMessage struct {
	ctx   context.Context
	topic string
	hdr   metadata.Metadata
	body  []byte
	c     codec.Codec
	err   error
}

func (m *mockMessage) Context() context.Context  { return m.ctx }
func (m *mockMessage) Topic() string             { return m.topic }
func (m *mockMessage) Header() metadata.Metadata { return m.hdr }
func (m *mockMessage) Body() []byte              { return m.body }
func (m *mockMessage) Ack() error                { return nil }
func (m *mockMessage) Error() error              { return m.err }

func (m *mockMessage) Unmarshal(dst any, opts ...codec.Option) error {
	if m.c == nil {
		return codec.ErrUnknownContentType
	}
	return m.c.Unmarshal(m.body, dst, opts...)
}

// NewMockMessage creates a broker.Message with pre-marshaled body for use with InjectMessage.
// Optionally pass a codec to enable Unmarshal in the handler under test.
func NewMockMessage(ctx context.Context, topic string, hdr metadata.Metadata, body []byte, c ...codec.Codec) broker.Message {
	msg := &mockMessage{ctx: ctx, topic: topic, hdr: hdr, body: body}
	if len(c) > 0 {
		msg.c = c[0]
	}
	return msg
}
