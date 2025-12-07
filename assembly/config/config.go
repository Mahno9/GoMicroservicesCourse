package config

import (
	"github.com/joho/godotenv"

	env "github.com/Mahno9/GoMicroservicesCourse/assembly/config/env"
)

type Config struct {
	Logger            LoggerConfig
	Kafka             KafkaConfig
	OrderPaidConsumer OrderPaidConsumerConfig
	ShipAssembled     ShipAssembledConfig
}

func Load(path ...string) (*Config, error) {
	err := godotenv.Load(path...)
	if err != nil {
		return nil, err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	kafka, err := env.NewKafkaConfig()
	if err != nil {
		return nil, err
	}

	orderPaidConsumer, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return nil, err
	}

	shipAssembled, err := env.NewShipAssembledConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Logger:            logger,
		Kafka:             kafka,
		OrderPaidConsumer: orderPaidConsumer,
		ShipAssembled:     shipAssembled,
	}, nil
}
