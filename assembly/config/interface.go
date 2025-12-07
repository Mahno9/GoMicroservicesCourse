package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Addresses() []string
}

type OrderPaidConsumerConfig interface {
	TopicName() string
	ConsumerGroupId() string
	Config() *sarama.Config
}

type ShipAssembledConfig interface {
	TopicName() string
	Config() *sarama.Config
}
