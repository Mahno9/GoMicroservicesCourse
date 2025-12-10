package env

import (
	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Port string `env:"GRPC_PORT,required"`
}

type grpcConfig struct {
	port string
}

func NewGrpcConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcConfig{
		port: raw.Port,
	}, nil
}

func (c *grpcConfig) Port() string {
	return c.port
}
