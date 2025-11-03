package order

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Create(order *model.Order) (*model.Order, error) {
	r.orders[order.OrderUuid] = converter.ModelToRepositoryOrder(order)

	storedOrder, ok := r.Get(order.OrderUuid)
	if ok != nil {
		return nil, model.ErrOrderDoesNotExist
	}
	return storedOrder, nil
}
