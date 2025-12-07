package env

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type shipAssembledConsumerConfigEnv struct {
	ShipAssembledTopicName       string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
	ShipAssembledConsumerGroupID string `env:"SHIP_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

type shipAssembledConsumerConfig struct {
	shipAssembledTopicName       string
	shipAssembledConsumerGroupID string
}

func NewShipAssembledConsumerConfig() (*shipAssembledConsumerConfig, error) {
	var raw shipAssembledConsumerConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	err := raw.checkValues()
	if err != nil {
		return nil, err
	}

	return &shipAssembledConsumerConfig{
		shipAssembledTopicName:       raw.ShipAssembledTopicName,
		shipAssembledConsumerGroupID: raw.ShipAssembledConsumerGroupID,
	}, nil
}

func (c *shipAssembledConsumerConfigEnv) checkValues() error {
	if c.ShipAssembledTopicName == "" {
		return fmt.Errorf("%w: %s", model.ErrInvalidConfigValue, "SHIP_ASSEMBLED_TOPIC_NAME")
	}

	if c.ShipAssembledConsumerGroupID == "" {
		return fmt.Errorf("%w: %s", model.ErrInvalidConfigValue, "SHIP_ASSEMBLED_CONSUMER_GROUP_ID")
	}

	return nil
}

func (c *shipAssembledConsumerConfig) TopicName() string {
	return c.shipAssembledTopicName
}

func (c *shipAssembledConsumerConfig) ConsumerGroupID() string {
	return c.shipAssembledConsumerGroupID
}

func (c *shipAssembledConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
