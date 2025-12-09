package env

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerConfigEnv struct {
	OrderPaidTopicName       string `env:"ORDER_PAID_TOPIC_NAME,required"`
	OrderPaidConsumerGroupId string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

type orderPaidConsumerConfig struct {
	topicName       string
	consumerGroupId string
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	var raw orderPaidConsumerConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("❗ Failed to parse order paid consumer config: %w", err)
	}

	return &orderPaidConsumerConfig{
		topicName:       raw.OrderPaidTopicName,
		consumerGroupId: raw.OrderPaidConsumerGroupId,
	}, nil
}

func (c *orderPaidConsumerConfig) TopicName() string {
	return c.topicName
}

func (c *orderPaidConsumerConfig) ConsumerGroupId() string {
	return c.consumerGroupId
}

func (c *orderPaidConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
