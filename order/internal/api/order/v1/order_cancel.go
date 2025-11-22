package v1

import (
	"context"

	"github.com/google/uuid"

	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) OrderCancel(ctx context.Context, params genOrderV1.OrderCancelParams) (genOrderV1.OrderCancelRes, error) {
	orderUuid := uuid.UUID(params.OrderUUID)

	err := h.orderService.OrderCancel(ctx, orderUuid)
	if err != nil {
		return nil, err
	}

	return &genOrderV1.OrderCancelNoContent{}, nil
}
