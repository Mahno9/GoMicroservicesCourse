package env

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledConsumerConfigEnv struct {
	ShipAssembledTopicName       string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
	ShipAssembledConsumerGroupId string `env:"SHIP_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type shipAssembledConsumerConfig struct {
	topicName       string
	consumerGroupId string
}

func NewShipAssembledConsumerConfig() (*shipAssembledConsumerConfig, error) {
	var raw shipAssembledConsumerConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("❗ Failed to parse ship assembled consumer config: %w", err)
	}

	return &shipAssembledConsumerConfig{
		topicName:       raw.ShipAssembledTopicName,
		consumerGroupId: raw.ShipAssembledConsumerGroupId,
	}, nil
}

func (c *shipAssembledConsumerConfig) TopicName() string {
	return c.topicName
}

func (c *shipAssembledConsumerConfig) ConsumerGroupId() string {
	return c.consumerGroupId
}

func (c *shipAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
