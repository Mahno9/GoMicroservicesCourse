package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type telegramConfigEnv struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

type telegramConfig struct {
	botToken string
}

func NewTelegramConfig() (*telegramConfig, error) {
	var raw telegramConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("❗ Failed to parse Telegram settings: %w", err)
	}

	return &telegramConfig{
		botToken: raw.BotToken,
	}, nil
}

func (c *telegramConfig) BotToken() string {
	return c.botToken
}
