package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) OrderCancel(ctx context.Context, orderUuid uuid.UUID) error {
	reqGetCtx, cancelGet := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancelGet()

	order, err := s.orderRepository.Get(reqGetCtx, orderUuid)
	if err != nil {
		return err
	}

	if order.Status != model.StatusPENDINGPAYMENT {
		return model.ErrOrderCancelConflict
	}

	order.Status = model.StatusCANCELLED

	reqUpdateCtx, cancelUpdate := context.WithTimeout(ctx, model.RequestTimeoutUpdate)
	defer cancelUpdate()

	err = s.orderRepository.Update(reqUpdateCtx, order)
	if err != nil {
		return err
	}

	return nil
}
