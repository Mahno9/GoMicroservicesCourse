package v1

import (
	"context"
	"log"
	"runtime"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	model "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderReq, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {

	contextWithTracing, cancel := WithTracingCancel(ctx, "orderAPI")
	defer cancel()

	timedCtx, cancelTimed := context.WithTimeout(contextWithTracing, commonRequestTimeout)
	defer cancelTimed()

	err := h.orderService.PayOrder(timedCtx, model.PayOrderData{
		OrderUuid:     params.OrderUUID,
		PaymentMethod: converter.ApiToModelPaymentMethod(req.PaymentMethod),
		// UserUuid: is taken from order
	})

	// TODO: convert to V1 error format
	return nil, err
}

func WithTracingCancel(parent context.Context, name string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	return ctx, func() {
		buf := make([]byte, 2048)
		n := runtime.Stack(buf, false)
		log.Printf("context %s canceled from:\n%s", name, buf[:n])
		cancel()
	}
}
