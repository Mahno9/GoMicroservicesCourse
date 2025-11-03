package order

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/converter"
)

func (r *repository) Get(orderUuid string) (*model.Order, error) {
	repoOrder, exists := r.orders[orderUuid]
	if !exists {
		return nil, model.ErrOrderDoesNotExist
	}

	return converter.RepositoryOrderToModel(repoOrder), nil
}
