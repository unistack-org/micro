// Package stream encapsulates streams within streams
package stream

import (
	"context"
	"sync"

	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/codec"
	"github.com/unistack-org/micro/v3/metadata"
	"github.com/unistack-org/micro/v3/server"
)

type Stream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}

type stream struct {
	Stream

	sync.RWMutex
	err     error
	request *request
}

type request struct {
	client.Request
	context context.Context
}

func (r *request) Codec() codec.Codec {
	return r.Request.Codec()
}

func (r *request) Header() metadata.Metadata {
	md, _ := metadata.FromContext(r.context)
	return md
}

func (r *request) Read() ([]byte, error) {
	return nil, nil
}

func (s *stream) Request() server.Request {
	return s.request
}

func (s *stream) Send(v interface{}) error {
	err := s.Stream.SendMsg(v)
	if err != nil {
		s.Lock()
		s.err = err
		s.Unlock()
	}
	return err
}

func (s *stream) Recv(v interface{}) error {
	err := s.Stream.RecvMsg(v)
	if err != nil {
		s.Lock()
		s.err = err
		s.Unlock()
	}
	return err
}

func (s *stream) Error() error {
	s.RLock()
	defer s.RUnlock()
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
