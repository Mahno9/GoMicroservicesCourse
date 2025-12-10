package config

import (
	"github.com/joho/godotenv"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/config/env"
)

type Config struct {
	Logger LoggerConfig

	Postgres PostgresConfig

	Http HttpConfig

	Clients ClientsConfig

	Kafka                 KafkaConfig
	OrderPaidProducer     OrderPaidProducerConfig
	ShipAssembledConsumer ShipAssembledConsumerConfig
}

func Load(path ...string) (*Config, error) {
	err := godotenv.Load(path...)
	if err != nil {
		return nil, err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return nil, err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	httpConfig, err := env.NewHttpConfig()
	if err != nil {
		return nil, err
	}

	clientsConfig, err := env.NewClientsConfig()
	if err != nil {
		return nil, err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return nil, err
	}

	orderPaidProducer, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return nil, err
	}

	shipAssembledConsumer, err := env.NewShipAssembledConsumerConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Logger:                loggerConfig,
		Postgres:              postgresConfig,
		Http:                  httpConfig,
		Clients:               clientsConfig,
		Kafka:                 kafkaConfig,
		OrderPaidProducer:     orderPaidProducer,
		ShipAssembledConsumer: shipAssembledConsumer,
	}, nil
}
