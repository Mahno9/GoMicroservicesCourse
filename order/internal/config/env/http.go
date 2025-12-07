package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type httpEnvConfig struct {
	Host string `env:"SERVICE_URL,required"`
	Port string `env:"SERVICE_PORT,required"`
}

type httpConfig struct {
	host string
	port string
}

func NewHttpConfig() (*httpConfig, error) {
	var raw httpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &httpConfig{
		host: raw.Host,
		port: raw.Port,
	}, nil
}

func (c *httpConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

func (c *httpConfig) Host() string {
	return c.host
}

func (c *httpConfig) Port() string {
	return c.port
}
