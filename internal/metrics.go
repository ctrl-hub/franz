package internal

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	topicMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka",
			Subsystem: "topic",
			Name:      "info",
			Help:      "Topic information - the value is arbitrary (labels hold info).",
		},
		[]string{"cluster", "name", "partitions", "internal"},
	)
	topicPartitionMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka",
			Subsystem: "topic",
			Name:      "parititions",
			Help:      "Number of partitions in a topic",
		},
		[]string{"cluster", "topic"},
	)
	topicPartitionDetailMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka",
			Subsystem: "topic",
			Name:      "partition_info",
			Help:      "Topic partition information - the value is arbitrary (labels hold info).",
		},
		[]string{"cluster", "topic", "partition", "leader", "replicas", "offline_replicas"},
	)
	consumerGroupMembersMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka",
			Subsystem: "consumer_group",
			Name:      "members",
			Help:      "The number of members in a consumer group.",
		},
		[]string{"cluster", "name", "state", "protocol", "protocol_type"},
	)
	consumerGroupOffsetMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka",
			Subsystem: "consumer_group",
			Name:      "offsets",
			Help:      "The offset of a consumer group in a topic.",
		},
		[]string{"cluster", "name", "topic", "partition"},
	)
	brokerMetrics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kafka",
			Subsystem: "broker",
			Name:      "info",
			Help:      "Broker information - the value is arbitrary (labels hold info).",
		},
		[]string{"cluster", "id", "addr", "rack"},
	)
)

// metricsHandler creates the prom registry and returns the handler for the metrics server
func metricsHandler() http.Handler {
	r := prometheus.NewRegistry()
	r.MustRegister(
		brokerMetrics,
		consumerGroupMembersMetrics,
		consumerGroupOffsetMetrics,
		topicMetrics,
		topicPartitionMetrics,
		topicPartitionDetailMetrics,
	)
	return promhttp.HandlerFor(r, promhttp.HandlerOpts{})
}
