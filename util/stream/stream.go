// Package stream encapsulates streams within streams
package stream

import (
	"context"
	"sync"

	"go.unistack.org/micro/v4/client"
	"go.unistack.org/micro/v4/codec"
	"go.unistack.org/micro/v4/metadata"
	"go.unistack.org/micro/v4/server"
)

// Stream interface
type Stream interface {
	Context() context.Context
	SendMsg(msg interface{}) error
	RecvMsg(msg interface{}) error
	Close() error
}

type stream struct {
	Stream
	err     error
	request *request

	mu sync.RWMutex
}

type request struct {
	client.Request
	context context.Context
}

// Codec returns codec.Codec
func (r *request) Codec() codec.Codec {
	return r.Request.Codec()
}

// Header returns metadata header
func (r *request) Header() metadata.Metadata {
	md, _ := metadata.FromIncomingContext(r.context)
	return md
}

// Read returns stream data
func (r *request) Read() ([]byte, error) {
	return nil, nil
}

// Request returns server.Request
func (s *stream) Request() server.Request {
	return s.request
}

// Send sends message
func (s *stream) Send(v interface{}) error {
	err := s.Stream.SendMsg(v)
	if err != nil {
		s.mu.Lock()
		s.err = err
		s.mu.Unlock()
	}
	return err
}

// Recv receives data
func (s *stream) Recv(v interface{}) error {
	err := s.Stream.RecvMsg(v)
	if err != nil {
		s.mu.Lock()
		s.err = err
		s.mu.Unlock()
	}
	return err
}

// Error returns error that stream holds
func (s *stream) Error() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.err
}

// New returns a new encapsulated stream
// Proto stream within a server.Stream
func New(service, endpoint string, req interface{}, c client.Client, s Stream) server.Stream {
	return &stream{
		Stream: s,
		request: &request{
			context: s.Context(),
			Request: c.NewRequest(service, endpoint, req),
		},
	}
}
