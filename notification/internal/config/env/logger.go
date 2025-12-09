package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type loggerConfigEnv struct {
	Level  string `env:"LOGGER_LEVEL,required"`
	AsJson bool   `env:"LOGGER_AS_JSON,required"`
}

type loggerConfig struct {
	level  string
	asJson bool
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("❗ Failed to parse logger config: %w", err)
	}

	return &loggerConfig{
		level:  raw.Level,
		asJson: raw.AsJson,
	}, nil
}

func (c *loggerConfig) Level() string {
	return c.level
}

func (c *loggerConfig) AsJson() bool {
	return c.asJson
}
