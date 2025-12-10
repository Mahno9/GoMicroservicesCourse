package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type mongoEnvConfig struct {
	Host         string `env:"MONGO_HOST,required"`
	Port         int    `env:"MONGO_PORT,required"`
	AuthDB       string `env:"MONGO_AUTH_DB,required"`
	DatabaseName string `env:"MONGO_INITDB_DATABASE,required"`
	User         string `env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password     string `env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}

type mongoConfig struct {
	databaseName string
	uRI          string
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw mongoEnvConfig

	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &mongoConfig{
		databaseName: raw.DatabaseName,
		uRI: fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			raw.User,
			raw.Password,
			raw.Host,
			raw.Port,
			raw.DatabaseName,
			raw.AuthDB),
	}, nil
}

func (c *mongoConfig) DatabaseName() string {
	return c.databaseName
}

func (c *mongoConfig) URI() string {
	return c.uRI
}
