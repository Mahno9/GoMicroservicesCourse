package v1

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	model "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) PayOrder(ctx context.Context, req *genOrderV1.PayOrderReq, params genOrderV1.PayOrderParams) (genOrderV1.PayOrderRes, error) {
	transactionUuid, err := h.orderService.PayOrder(ctx, model.PayOrderData{
		OrderUuid:     params.OrderUUID,
		PaymentMethod: converter.ApiToModelPaymentMethod(req.PaymentMethod),
		// UserUuid: is taken from order
	})
	if err != nil {
		return nil, err
	}

	return &genOrderV1.PayOrderOK{
		TransactionUUID: genOrderV1.TransactionUUID(transactionUuid),
	}, nil
}
