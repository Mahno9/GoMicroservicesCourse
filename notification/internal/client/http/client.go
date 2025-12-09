package client

import "context"

type TelegramClient interface {
	SendMessage(ctx context.Context, chatId int64, message string) error
	RegisterNewChatSubscriber(sub NewChatSubscriber)
}

type NewChatSubscriber interface {
	NewChatStarted(ctx context.Context, chatId int64) error
}
