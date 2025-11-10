package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	newOrder := converter.ModelToRepositoryOrder(order)
	r.orders[order.OrderUuid] = newOrder

	return converter.RepositoryOrderToModel(newOrder), nil
}
