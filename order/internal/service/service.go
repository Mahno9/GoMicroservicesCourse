package service

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type OrderService interface {
	CreateOrder(c context.Context, data model.CreateOrderData) (*model.Order, error)
	GetOrder(c context.Context, orderUuid string) (*model.Order, error)
	OrderCancel(c context.Context, orderUuid string) error
	PayOrder(c context.Context, data model.PayOrderData) (string, error)
}
