package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, orderUuid string) (*model.Order, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	repoOrder, exists := r.orders[orderUuid]
	if !exists {
		return nil, model.ErrOrderDoesNotExist
	}

	return converter.RepositoryOrderToModel(repoOrder), nil
}
