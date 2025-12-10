package config

import (
	"github.com/joho/godotenv"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/config/env"
)

type Config struct {
	LoggerConfig LoggerConfig
	MongoConfig  MongoConfig
	GrpcConfig   GrpcConfig
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

	mongoConfig, err := env.NewMongoConfig()
	if err != nil {
		return nil, err
	}

	grpcConfig, err := env.NewGrpcConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		LoggerConfig: loggerConfig,
		MongoConfig:  mongoConfig,
		GrpcConfig:   grpcConfig,
	}, nil
}
