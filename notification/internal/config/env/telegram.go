package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type telegramConfigEnv struct {
	BotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	ChatID   int64  `env:"TELEGRAM_CHAT_ID,required"`
}

type telegramConfig struct {
	botToken string
	chatID   int64
}

func NewTelegramConfig() (*telegramConfig, error) {
	var raw telegramConfigEnv
	if err := env.Parse(&raw); err != nil {
		return nil, fmt.Errorf("❗ Failed to parse Telegram settings: %w", err)
	}

	return &telegramConfig{
		botToken: raw.BotToken,
		chatID:   raw.ChatID,
	}, nil
}

func (c *telegramConfig) BotToken() string {
	return c.botToken
}

func (c *telegramConfig) ChatID() int64 {
	return c.chatID
}
