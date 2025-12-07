package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, data model.CreateOrderData) (*model.Order, error)
	GetOrder(ctx context.Context, orderUuid uuid.UUID) (*model.Order, error)
	OrderCancel(ctx context.Context, orderUuid uuid.UUID) error
	PayOrder(ctx context.Context, data model.PayOrderData) (uuid.UUID, error)
}

type ProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
