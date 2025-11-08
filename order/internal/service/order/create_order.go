package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, data model.CreateOrderData) (*model.Order, error) {
	ctxListReq, cancelList := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancelList()

	orderParts, err := s.inventoryClient.ListParts(ctxListReq, &model.PartsFilter{
		Uuids: data.PartUuids,
	})
	if err != nil {
		return nil, err
	}

	if len(orderParts) != len(data.PartUuids) {
		return nil, model.ErrPartsNotAvailable
	}

	totalPrice := 0.0
	for _, part := range orderParts {
		totalPrice += part.Price
	}

	orderUuid := uuid.New().String()

	ctxCreateReq, cancelCreate := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancelCreate()

	createdOrder, err := s.orderRepository.Create(ctxCreateReq, &model.Order{
		OrderUuid:  orderUuid,
		UserUuid:   data.UserUuid,
		PartUuids:  data.PartUuids,
		TotalPrice: totalPrice,
	})
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}
