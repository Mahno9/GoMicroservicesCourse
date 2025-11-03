package v1

import (
	"context"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
	"log"
)

func (h *apiHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderReq) (orderV1.CreateOrderRes, error) {
	log.Printf("Creating order with details: %v\n", req)

	timedCtx, cancel := context.WithTimeout(ctx, createOrderTimeout)
	defer cancel()

	res, err := h.orderService.CreateOrder(timedCtx, converter.ApiToModelOrderInfo(req))
	if err != nil {
		log.Printf("❗ Failed to create order: %v\nNo order is created.", err)
		return nil, err
	}

	log.Printf("❕ New order created:\n%+v\n", res)

	return &orderV1.CreateOrderCreated{
		OrderUUID:  orderV1.OrderUUID(res.OrderUuid),
		TotalPrice: orderV1.TotalPrice(res.TotalPrice),
	}, nil
}
