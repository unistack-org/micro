package sql

import (
	"context"
	"math"

	"golang.yandex/hasql/v2"
)

// ClusterOptions contains cluster specific options
type ClusterOptions struct {
	NodeChecker        hasql.NodeChecker
	NodePicker         hasql.NodePicker[Querier]
	NodeDiscoverer     hasql.NodeDiscoverer[Querier]
	Options            []hasql.ClusterOpt[Querier]
	Context            context.Context
	Retries            int
	NodePriority       map[string]int32
	NodeStateCriterion hasql.NodeStateCriterion
}

// ClusterOption apply cluster options to ClusterOptions
type ClusterOption func(*ClusterOptions)

// WithClusterNodeChecker pass hasql.NodeChecker to cluster options
func WithClusterNodeChecker(c hasql.NodeChecker) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodeChecker = c
	}
}

// WithClusterNodePicker pass hasql.NodePicker to cluster options
func WithClusterNodePicker(p hasql.NodePicker[Querier]) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodePicker = p
	}
}

// WithClusterNodeDiscoverer pass hasql.NodeDiscoverer to cluster options
func WithClusterNodeDiscoverer(d hasql.NodeDiscoverer[Querier]) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodeDiscoverer = d
	}
}

// WithRetries retry count on other nodes in case of error
func WithRetries(n int) ClusterOption {
	return func(o *ClusterOptions) {
		o.Retries = n
	}
}

// WithClusterContext pass context.Context to cluster options and used for checks
func WithClusterContext(ctx context.Context) ClusterOption {
	return func(o *ClusterOptions) {
		o.Context = ctx
	}
}

// WithClusterOptions pass hasql.ClusterOpt
func WithClusterOptions(opts ...hasql.ClusterOpt[Querier]) ClusterOption {
	return func(o *ClusterOptions) {
		o.Options = append(o.Options, opts...)
	}
}

// WithClusterNodeStateCriterion pass default hasql.NodeStateCriterion
func WithClusterNodeStateCriterion(c hasql.NodeStateCriterion) ClusterOption {
	return func(o *ClusterOptions) {
		o.NodeStateCriterion = c
	}
}

type ClusterNode struct {
	Name     string
	DB       Querier
	Priority int32
}

// WithClusterNodes create cluster with static NodeDiscoverer
func WithClusterNodes(cns ...ClusterNode) ClusterOption {
	return func(o *ClusterOptions) {
		nodes := make([]*hasql.Node[Querier], 0, len(cns))
		if o.NodePriority == nil {
			o.NodePriority = make(map[string]int32, len(cns))
		}
		for _, cn := range cns {
			nodes = append(nodes, hasql.NewNode(cn.Name, cn.DB))
			if cn.Priority == 0 {
				cn.Priority = math.MaxInt32
			}
			o.NodePriority[cn.Name] = cn.Priority
		}
		o.NodeDiscoverer = hasql.NewStaticNodeDiscoverer(nodes...)
	}
}

type nodeStateCriterionKey struct{}

// NodeStateCriterion inject hasql.NodeStateCriterion to context
func NodeStateCriterion(ctx context.Context, c hasql.NodeStateCriterion) context.Context {
	return context.WithValue(ctx, nodeStateCriterionKey{}, c)
}

func (c *Cluster) getNodeStateCriterion(ctx context.Context) hasql.NodeStateCriterion {
	if v, ok := ctx.Value(nodeStateCriterionKey{}).(hasql.NodeStateCriterion); ok {
		return v
	}
	return c.options.NodeStateCriterion
}
