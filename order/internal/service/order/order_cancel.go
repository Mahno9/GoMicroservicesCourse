package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) OrderCancel(c context.Context, orderUuid string) error {
	order, err := s.ordersRepo.Get(orderUuid)
	if err != nil {
		return err
	}

	if order.Status != model.StatusPENDINGPAYMENT {
		return model.ErrOrderCancelConflict
	}

	order.Status = model.StatusCANCELLED

	err = s.ordersRepo.Update(order)
	if err != nil {
		return err
	}

	return nil
}
