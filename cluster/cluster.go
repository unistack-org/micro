package cluster

import (
	"context"

	"go.unistack.org/micro/v3/metadata"
)

// Message sent to member in cluster
type Message interface {
	// Header returns message headers
	Header() metadata.Metadata
	// Body returns broker message may be []byte slice or some go struct or interface
	Body() interface{}
}

type Node interface {
	// Name returns node name
	Name() string
	// Address returns node address
	Address() string
	// Metadata returns node metadata
	Metadata() metadata.Metadata
}

// Cluster interface used for cluster communication across nodes
type Cluster interface {
	// Join is used to take an existing members and performing state sync
	Join(ctx context.Context, addr ...string) error
	// Leave broadcast a leave message and stop listeners
	Leave(ctx context.Context) error
	// Ping is used to probe live status of the node
	Ping(ctx context.Context, node Node, payload []byte) error
	// Members returns the cluster members
	Members() ([]Node, error)
	// Broadcast send message for all members in cluster, if filter is not nil, nodes may be filtered
	// by key/value pairs
	Broadcast(ctx context.Context, msg Message, filter ...string) error
	// Unicast send message to single member in cluster
	Unicast(ctx context.Context, node Node, msg Message) error
	// Live returns cluster liveness
	Live() bool
	// Ready returns cluster readiness
	Ready() bool
	// Health returns cluster health
	Health() bool
}
