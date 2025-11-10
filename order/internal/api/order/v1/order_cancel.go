package v1

import (
	"context"

	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) OrderCancel(ctx context.Context, params genOrderV1.OrderCancelParams) (genOrderV1.OrderCancelRes, error) {
	err := h.orderService.OrderCancel(ctx, params.OrderUUID)
	return nil, err
}
