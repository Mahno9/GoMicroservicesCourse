package service

import (
	"context"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	BroadcastMessage(ctx context.Context, message string) error
}
