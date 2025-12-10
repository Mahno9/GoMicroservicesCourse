package client

import "context"

type HandlerCallback func(ctx context.Context) error

type TelegramClient interface {
	SendMessage(ctx context.Context, chatId int64, message string) error
	SetStartHandler(ctx context.Context, handler HandlerCallback) error
}
