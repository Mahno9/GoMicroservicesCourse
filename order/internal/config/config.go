package config

import (
	"github.com/joho/godotenv"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/config/env"
)

type Config struct {
	LoggerConfig   LoggerConfig
	PostgresConfig PostgresConfig
	HttpConfig     HttpConfig
	ClientsConfig  ClientsConfig
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

	return &Config{
		LoggerConfig:   loggerConfig,
		PostgresConfig: postgresConfig,
		HttpConfig:     httpConfig,
		ClientsConfig:  clientsConfig,
	}, nil
}
