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
	bot         *bot.Bot
	newChatSubs []clients.NewChatSubscriber
}

func NewClient(tBot *bot.Bot) clients.TelegramClient {
	c := &client{
		bot: tBot,
	}

	tBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, c.startHandler)

	return c
}

func (c *client) startHandler(ctx context.Context, bot *bot.Bot, update *models.Update) {
	logger.Info(ctx, "New chat started", zap.Int64("chatId", update.Message.Chat.ID))

	for _, sub := range c.newChatSubs {
		if err := sub.NewChatStarted(ctx, update.Message.Chat.ID); err != nil {
			logger.Error(ctx, "😡 Failed to register new chat", zap.Error(err))
		}
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

func (c *client) RegisterNewChatSubscriber(sub clients.NewChatSubscriber) {
	c.newChatSubs = append(c.newChatSubs, sub)
}
