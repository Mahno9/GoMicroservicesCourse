package v1

import (
	"context"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) OrderCancel(ctx context.Context, params orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	timedCtx, cancel := context.WithTimeout(ctx, commonRequestTimeout)
	defer cancel()

	err := h.orderService.OrderCancel(timedCtx, params.OrderUUID)
	return nil, err
}
