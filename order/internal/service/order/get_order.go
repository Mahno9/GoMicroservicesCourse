package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) GetOrder(c context.Context, orderUuid string) (*model.Order, error) {
	order, err := s.ordersRepo.Get(orderUuid)
	if err != nil {
		return nil, err
	}

	return order, nil
}
