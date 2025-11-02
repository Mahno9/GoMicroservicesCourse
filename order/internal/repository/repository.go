package repository

import "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"

type OrderRepository interface {
	Create(order *model.Order) (*model.Order, error)
	Get(orderUuid string) (*model.Order, error)
	Update(order *model.Order) error
}
