package env

import (
	"fmt"
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
		return nil, fmt.Errorf("❗ Failed to parse Kafka config: %w", err)
	}

	return &kafkaConfig{
		brokerAddresses: strings.Split(raw.BrokersAddresses, ","),
	}, nil
}

func (c *kafkaConfig) Addresses() []string {
	return c.brokerAddresses
}
