package socket

import (
	"sync"
)

// Pool holds the socket pool
type Pool struct {
	sync.RWMutex
	pool map[string]*Socket
}

// Get socket from pool
func (p *Pool) Get(id string) (*Socket, bool) {
	// attempt to get existing socket
	p.RLock()
	socket, ok := p.pool[id]
	if ok {
		p.RUnlock()
		return socket, ok
	}
	p.RUnlock()

	// save socket
	p.Lock()
	defer p.Unlock()
	// double checked locking
	socket, ok = p.pool[id]
	if ok {
		return socket, ok
	}
	// create new socket
	socket = New(id)
	p.pool[id] = socket

	// return socket
	return socket, false
}

// Release close the socket and delete from pool
func (p *Pool) Release(s *Socket) {
	p.Lock()
	defer p.Unlock()

	// close the socket
	s.Close()
	delete(p.pool, s.id)
}

// Close the pool and delete all the sockets
func (p *Pool) Close() {
	p.Lock()
	defer p.Unlock()
	for id, sock := range p.pool {
		sock.Close()
		delete(p.pool, id)
	}
}

// NewPool returns a new socket pool
func NewPool() *Pool {
	return &Pool{
		pool: make(map[string]*Socket),
	}
}
