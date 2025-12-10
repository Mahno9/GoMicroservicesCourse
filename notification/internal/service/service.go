package service

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/notification/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendShipAssembledMessage(ctx context.Context, event model.ShipAssembledEvent) error
	SendOrderPaidMessage(ctx context.Context, event model.OrderPaidEvent) error
}
