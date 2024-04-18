package internal

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// clusterConfig prepares the config to connect to the confluent cluster
func clusterConfig() *sarama.Config {
	config := sarama.NewConfig()

	config.Net.SASL.Enable = true
	config.Net.SASL.User = viper.GetString("confluent_api_key")
	config.Net.SASL.Password = viper.GetString("confluent_api_secret")
	config.Net.SASL.Handshake = true
	config.Net.SASL.Mechanism = "PLAIN"

	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
	}

	config.MetricRegistry.UnregisterAll()

	return config
}

func collect(clusterAdmin sarama.ClusterAdmin) error {

	start := time.Now()
	logrus.Debug("starting metrics collection")

	clusterLabel := viper.GetString("confluent_cluster_label")
	topicPartitions := map[string][]int32{}

	// get broker info
	brokers, _, err := clusterAdmin.DescribeCluster()
	if err != nil {
		return err
	}

	brokerMetrics.Reset()
	for _, broker := range brokers {
		brokerMetrics.With(prometheus.Labels{
			"confluent_cluster": clusterLabel,
			"id":                fmt.Sprintf("%d", broker.ID()),
			"addr":              broker.Addr(),
			"rack":              broker.Rack(),
		}).Set(1)
	}

	// collect a list of all the topic names
	topics, err := clusterAdmin.ListTopics()
	if err != nil {
		return err
	}

	var allTopics []string
	for topicName := range topics {
		allTopics = append(allTopics, topicName)
	}

	// describe all of the topics
	topicsMetadata, err := clusterAdmin.DescribeTopics(allTopics)
	if err != nil {
		return err
	}

	topicMetrics.Reset()
	topicPartitionMetrics.Reset()
	topicPartitionDetailMetrics.Reset()
	for _, topic := range topicsMetadata {
		topicMetrics.With(prometheus.Labels{
			"confluent_cluster": clusterLabel,
			"name":              topic.Name,
			"partitions":        strconv.Itoa(len(topic.Partitions)),
			"internal":          strconv.FormatBool(topic.IsInternal),
		}).Set(1)

		topicPartitionMetrics.With(prometheus.Labels{
			"confluent_cluster": clusterLabel,
			"topic":             topic.Name,
		}).Set(float64(len(topic.Partitions)))

		for _, partition := range topic.Partitions {
			topicPartitions[topic.Name] = append(topicPartitions[topic.Name], partition.ID)
			topicPartitionDetailMetrics.With(prometheus.Labels{
				"confluent_cluster": clusterLabel,
				"topic":             topic.Name,
				"partition":         fmt.Sprintf("%d", partition.ID),
				"leader":            fmt.Sprintf("%d", partition.Leader),
				"replicas":          strconv.Itoa(len(partition.Replicas)),
				"offline_replicas":  strconv.Itoa(len(partition.OfflineReplicas)),
			}).Set(1)
		}
	}

	// collect a list of all the consumer groups
	groups, err := clusterAdmin.ListConsumerGroups()
	if err != nil {
		return err
	}
	var allGroups []string
	for group := range groups {
		allGroups = append(allGroups, group)
	}

	// describe all of the consumer groups
	groupDescriptions, err := clusterAdmin.DescribeConsumerGroups(allGroups)
	if err != nil {
		return err
	}
	consumerGroupMembersMetrics.Reset()
	consumerGroupOffsetMetrics.Reset()
	for _, consumerGroupDescription := range groupDescriptions {
		consumerGroupMembersMetrics.With(prometheus.Labels{
			"confluent_cluster": clusterLabel,
			"name":              consumerGroupDescription.GroupId,
			"state":             consumerGroupDescription.State,
			"protocol":          consumerGroupDescription.Protocol,
			"protocol_type":     consumerGroupDescription.ProtocolType,
		}).Set(float64(len(consumerGroupDescription.Members)))

		o, err := clusterAdmin.ListConsumerGroupOffsets(consumerGroupDescription.GroupId, topicPartitions)
		if err != nil {
			return err
		}

		for topic, blocks := range o.Blocks {
			for i, block := range blocks {
				// only track metrics for block offsets that aren't `-1` - those are not subscribed to the topic
				if block.Offset >= 0 {
					consumerGroupOffsetMetrics.With(prometheus.Labels{
						"confluent_cluster": clusterLabel,
						"name":              consumerGroupDescription.GroupId,
						"topic":             topic,
						"partition":         fmt.Sprintf("%d", i),
					}).Set(float64(block.Offset))
				}
			}
		}
	}

	logrus.WithField("duration", fmt.Sprintf("%dms", time.Since(start).Milliseconds())).Debug("completed metrics collection")
	return nil
}
