package service

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, data model.CreateOrderData) (*model.Order, error)
	GetOrder(ctx context.Context, orderUuid string) (*model.Order, error)
	OrderCancel(ctx context.Context, orderUuid string) error
	PayOrder(ctx context.Context, data model.PayOrderData) (string, error)
}
