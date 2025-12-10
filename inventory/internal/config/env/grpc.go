package env

import (
	"github.com/caarlos0/env/v11"
)

type grpcEnvConfig struct {
	Host string `env:"GRPC_HOST" envDefault:""`
	Port string `env:"GRPC_PORT,required"`
}

type grpcConfig struct {
	host string
	port string
}

func NewGrpcConfig() (*grpcConfig, error) {
	var raw grpcEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &grpcConfig{
		host: raw.Host,
		port: raw.Port,
	}, nil
}

func (c *grpcConfig) Host() string {
	return c.host
}

func (c *grpcConfig) Port() string {
	return c.port
}
