package v1

import (
	"context"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	timedCtx, cancel := context.WithTimeout(ctx, commonRequestTimeout)
	defer cancel()

	order, err := h.orderService.GetOrder(timedCtx, params.OrderUUID)

	return converter.ModelToApiGetOrder(order), err
}
