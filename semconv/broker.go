package semconv

var (
	// PublishMessageDurationSeconds specifies meter metric name
	PublishMessageDurationSeconds = "micro_publish_message_duration_seconds"
	// PublishMessageLatencyMicroseconds specifies meter metric name
	PublishMessageLatencyMicroseconds = "micro_publish_message_latency_microseconds"
	// PublishMessageTotal specifies meter metric name
	PublishMessageTotal = "micro_publish_message_total"
	// PublishMessageInflight specifies meter metric name
	PublishMessageInflight = "micro_publish_message_inflight"
	// SubscribeMessageDurationSeconds specifies meter metric name
	SubscribeMessageDurationSeconds = "micro_subscribe_message_duration_seconds"
	// SubscribeMessageLatencyMicroseconds specifies meter metric name
	SubscribeMessageLatencyMicroseconds = "micro_subscribe_message_latency_microseconds"
	// SubscribeMessageTotal specifies meter metric name
	SubscribeMessageTotal = "micro_subscribe_message_total"
	// SubscribeMessageInflight specifies meter metric name
	SubscribeMessageInflight = "micro_subscribe_message_inflight"
	// BrokerGroupLag specifies broker lag
	BrokerGroupLag = "micro_broker_group_lag"
)
