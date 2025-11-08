package v1

import (
	"context"
	"log"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/converter"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
)

func (h *apiHandler) CreateOrder(ctx context.Context, req *genOrderV1.CreateOrderReq) (genOrderV1.CreateOrderRes, error) {
	log.Printf("Creating order with details: %v\n", req)

	res, err := h.orderService.CreateOrder(ctx, converter.ApiToModelOrderInfo(req))
	if err != nil {
		log.Printf("❗ Failed to create order: %v\nNo order is created.", err)
		return nil, err
	}

	log.Printf("❕ New order created:\n%+v\n", res)

	return &genOrderV1.CreateOrderCreated{
		OrderUUID:  genOrderV1.OrderUUID(res.OrderUuid),
		TotalPrice: genOrderV1.TotalPrice(res.TotalPrice),
	}, nil
}
