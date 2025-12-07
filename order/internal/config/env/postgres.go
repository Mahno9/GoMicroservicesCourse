package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	User          string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	MigrationsDir string `env:"MIGRATIONS_DIR"`
}

type postgresConfig struct {
	uri           string
	migrationsDir string
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgresConfig{
		uri:           fmt.Sprintf("postgres://%s:%s@%s:%s", raw.User, raw.Password, raw.Host, raw.Port),
		migrationsDir: raw.MigrationsDir,
	}, nil
}

func (c *postgresConfig) URI() string {
	return c.uri
}

func (c *postgresConfig) MigrationsDir() string {
	return c.migrationsDir
}
