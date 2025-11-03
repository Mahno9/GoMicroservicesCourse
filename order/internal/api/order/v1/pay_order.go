package v1

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	model "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderReq, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	timedCtx, cancelTimed := context.WithTimeout(ctx, commonRequestTimeout)
	defer cancelTimed()

	transactionUuid, err := h.orderService.PayOrder(timedCtx, model.PayOrderData{
		OrderUuid:     params.OrderUUID,
		PaymentMethod: converter.ApiToModelPaymentMethod(req.PaymentMethod),
		// UserUuid: is taken from order
	})
	if err != nil {
		return nil, err
	}

	return &orderV1.PayOrderOK{
		TransactionUUID: orderV1.TransactionUUID(transactionUuid),
	}, nil
}
