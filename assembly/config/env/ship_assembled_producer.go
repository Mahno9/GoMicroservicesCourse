package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type shipAssembledConfigEnv struct {
	ShipAssembledTopicName string `env:"SHIP_ASSEMBLED_TOPIC_NAME,required"`
}

type shipAssembledConfig struct {
	topicName string
}

func NewShipAssembledConfig() (*shipAssembledConfig, error) {
	var raw shipAssembledConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &shipAssembledConfig{
		topicName: raw.ShipAssembledTopicName,
	}, nil
}

func (c *shipAssembledConfig) TopicName() string {
	return c.topicName
}

func (c *shipAssembledConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true // TODO: check wether it need

	return config
}
