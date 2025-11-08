package v1

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) GetOrder(ctx context.Context, params genOrderV1.GetOrderParams) (genOrderV1.GetOrderRes, error) {
	order, err := h.orderService.GetOrder(ctx, params.OrderUUID)

	return converter.ModelToApiGetOrder(order), err
}
