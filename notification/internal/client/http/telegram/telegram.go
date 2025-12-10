package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	clients "github.com/Mahno9/GoMicroservicesCourse/notification/internal/client/http"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type client struct {
	bot          *bot.Bot
	startHandler clients.HandlerCallback
}

func NewClient(tBot *bot.Bot) clients.TelegramClient {
	c := &client{
		bot: tBot,
	}

	tBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, c.baseStartHandler)

	return c
}

func (c *client) baseStartHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	logger.Info(ctx, "/start called", zap.Int64("chatId", update.Message.Chat.ID))

	if c.startHandler != nil {
		err := c.startHandler(ctx)
		if err != nil {
			logger.Warn(ctx, "❗ Failed while run external start handler", zap.Error(err))
		}
	} else {
		logger.Info(ctx, "No custom handler specified for /start")
	}
}

func (c *client) SendMessage(ctx context.Context, chatId int64, message string) error {
	logger.Info(ctx, "Sending telegram message",
		zap.Int64("chatId", chatId),
		zap.String("message", message))

	_, err := c.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      message,
		ParseMode: "Markdown",
	})

	return err
}

func (c *client) SetStartHandler(ctx context.Context, handler clients.HandlerCallback) error {
	c.startHandler = handler
	return nil
}
