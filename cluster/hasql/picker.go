package sql

import (
	"fmt"
	"math"
	"time"

	"golang.yandex/hasql/v2"
)

// compile time guard
var _ hasql.NodePicker[Querier] = (*CustomPicker[Querier])(nil)

// CustomPickerOptions holds options to pick nodes
type CustomPickerOptions struct {
	MaxLag   int
	Priority map[string]int32
	Retries  int
}

// CustomPickerOption func apply option to CustomPickerOptions
type CustomPickerOption func(*CustomPickerOptions)

// CustomPickerMaxLag specifies max lag for which node can be used
func CustomPickerMaxLag(n int) CustomPickerOption {
	return func(o *CustomPickerOptions) {
		o.MaxLag = n
	}
}

// NewCustomPicker creates new node picker
func NewCustomPicker[T Querier](opts ...CustomPickerOption) *CustomPicker[Querier] {
	options := CustomPickerOptions{}
	for _, o := range opts {
		o(&options)
	}
	return &CustomPicker[Querier]{opts: options}
}

// CustomPicker holds node picker options
type CustomPicker[T Querier] struct {
	opts CustomPickerOptions
}

// PickNode used to return specific node
func (p *CustomPicker[T]) PickNode(cnodes []hasql.CheckedNode[T]) hasql.CheckedNode[T] {
	for _, n := range cnodes {
		fmt.Printf("node %s\n", n.Node.String())
	}
	return cnodes[0]
}

func (p *CustomPicker[T]) getPriority(nodeName string) int32 {
	if prio, ok := p.opts.Priority[nodeName]; ok {
		return prio
	}
	return math.MaxInt32 // Default to lowest priority
}

// CompareNodes used to sort nodes
func (p *CustomPicker[T]) CompareNodes(a, b hasql.CheckedNode[T]) int {
	// Get replication lag values
	aLag := a.Info.(interface{ ReplicationLag() int }).ReplicationLag()
	bLag := b.Info.(interface{ ReplicationLag() int }).ReplicationLag()

	// First check that lag lower then MaxLag
	if aLag > p.opts.MaxLag && bLag > p.opts.MaxLag {
		return 0 // both are equal
	}

	// If one node exceeds MaxLag and the other doesn't, prefer the one that doesn't
	if aLag > p.opts.MaxLag {
		return 1 // b is better
	}
	if bLag > p.opts.MaxLag {
		return -1 // a is better
	}

	// Get node priorities
	aPrio := p.getPriority(a.Node.String())
	bPrio := p.getPriority(b.Node.String())

	// if both priority equals
	if aPrio == bPrio {
		// First compare by replication lag
		if aLag < bLag {
			return -1
		}
		if aLag > bLag {
			return 1
		}
		// If replication lag is equal, compare by latency
		aLatency := a.Info.(interface{ Latency() time.Duration }).Latency()
		bLatency := b.Info.(interface{ Latency() time.Duration }).Latency()

		if aLatency < bLatency {
			return -1
		}
		if aLatency > bLatency {
			return 1
		}

		// If lag and latency is equal
		return 0
	}

	// If priorities are different, prefer the node with lower priority value
	if aPrio < bPrio {
		return -1
	}

	return 1
}
