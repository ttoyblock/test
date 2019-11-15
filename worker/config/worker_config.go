package config

import (
	"github.com/Shopify/sarama"
)

const (
	WorkerTopic   = "async_job"
	WorkerChannel = "excute_job"
)

var KafkaConfig *sarama.Config

func init() {
	KafkaConfig = sarama.NewConfig()
	KafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	KafkaConfig.Producer.Retry.Max = 5
	KafkaConfig.Producer.Return.Successes = true
	KafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
}
