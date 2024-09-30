package semconv

var (
	// PoolGetTotal specifies meter metric name for total number of pool get ops
	PoolGetTotal = "micro_pool_get_total"
	// PoolPutTotal specifies meter metric name for total number of pool put ops
	PoolPutTotal = "micro_pool_put_total"
	// PoolMisTotal specifies meter metric name for total number of pool misses
	PoolMisTotal = "micro_pool_mis_total"
	// PoolRetTotal specifies meter metric name for total number of pool returned to gc
	PoolRetTotal = "micro_pool_ret_total"
)
