package env

import "github.com/caarlos0/env/v11"

type kafkaConfigEnv struct {
	BrokersAddress string `env:"KAFKA_BROKERS,required"`
}

type kafkaConfig struct {
	brokersAddress string
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var raw kafkaConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{
		brokersAddress: raw.BrokersAddress,
	}, nil
}

func (c *kafkaConfig) BrokersAddresses() []string {
	return []string{c.brokersAddress}
}
