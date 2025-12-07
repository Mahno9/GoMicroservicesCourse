package env

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type kafkaConfigEnv struct {
	BrokersAddresses string `env:"KAFKA_BROKERS,required"`
}

type kafkaConfig struct {
	brokerAddresses []string
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var raw kafkaConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{
		brokerAddresses: strings.Split(raw.BrokersAddresses, ","),
	}, nil
}

func (c *kafkaConfig) Addresses() []string {
	return c.brokerAddresses
}
