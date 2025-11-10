package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderUuid string) (*model.Order, error) {
	ctxReq, cancel := context.WithTimeout(ctx, model.RequestTimeoutRead)
	defer cancel()

	return s.orderRepository.Get(ctxReq, orderUuid)
}
