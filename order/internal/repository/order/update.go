package order

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Update(order *model.Order) error {
	r.orders[order.OrderUuid] = converter.ModelToRepositoryOrder(order)
	return nil
}
