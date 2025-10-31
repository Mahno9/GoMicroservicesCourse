package service

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

type OrderService interface {
	CreateOrder(c context.Context, data model.CreateOrderData) (model.CreateOrderData, error)
	// GetOrder
	// OrderCancel
	// PayOrder
}
