package order

import (
	"context"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) CreateOrder(c context.Context, data model.CreateOrderData) (*model.Order, error) {
	orderParts, err := s.inventory.ListParts(c, &model.PartsFilter{
		Uuids: data.PartUuids,
	})
	if err != nil {
		return nil, err
	}

	if len(orderParts) != len(data.PartUuids) {
		return nil, model.PartsNotAvailableErr
	}

	totalPrice := float64(0.0)
	for _, part := range orderParts {
		totalPrice += part.Price
	}

	orderUuid := uuid.New().String()
	createdOrder, err := s.ordersRepo.Create(&model.Order{
		OrderUuid:  orderUuid,
		UserUuid:   data.UserUuid,
		PartUuids:  data.PartUuids,
		TotalPrice: totalPrice,
	})
	if err != nil {
		return nil, err
		// TODO: converter to v1 error response type
	}

	return createdOrder, nil
}
