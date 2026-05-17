// Package mock provides a mock implementation of broker.Broker for testing.
package mock

import (
	"context"

	"go.unistack.org/micro/v5/codec"
	"go.unistack.org/micro/v5/metadata"
)

// MockMessage is a broker.Message implementation for use in tests.
// It tracks whether Ack was called so tests can assert handler behaviour.
type MockMessage struct {
	ctx   context.Context
	topic string
	hdr   metadata.Metadata
	body  []byte
	c     codec.Codec
	err   error
	acked bool
}

// Context implements broker.Message.
func (m *MockMessage) Context() context.Context { return m.ctx }

// Topic implements broker.Message.
func (m *MockMessage) Topic() string { return m.topic }

// Header implements broker.Message.
func (m *MockMessage) Header() metadata.Metadata { return m.hdr }

// Body implements broker.Message.
func (m *MockMessage) Body() []byte { return m.body }

// Ack implements broker.Message. Records that the message was acknowledged.
func (m *MockMessage) Ack() error { m.acked = true; return nil }

// Error implements broker.Message.
func (m *MockMessage) Error() error { return m.err }

// Acked reports whether Ack was called on this message.
func (m *MockMessage) Acked() bool { return m.acked }

// Unmarshal implements broker.Message.
func (m *MockMessage) Unmarshal(dst any, opts ...codec.Option) error {
	if m.c == nil {
		return codec.ErrUnknownContentType
	}
	return m.c.Unmarshal(m.body, dst, opts...)
}

// NewMockMessage creates a MockMessage with pre-marshaled body for use with InjectMessage.
// Pass a non-nil codec to enable Unmarshal in the handler under test.
func NewMockMessage(ctx context.Context, topic string, hdr metadata.Metadata, body []byte, c codec.Codec) *MockMessage {
	return &MockMessage{ctx: ctx, topic: topic, hdr: hdr, body: body, c: c}
}
