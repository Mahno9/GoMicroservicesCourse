package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Update(_ context.Context, order *model.Order) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.orders[order.OrderUuid] = converter.ModelToRepositoryOrder(order)
	return nil
}
