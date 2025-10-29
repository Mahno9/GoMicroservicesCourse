package v1

import (
	"context"
	"log"

	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(c context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	paymentUuid, err := a.paymentService.Pay(c, req.UserUuid, req.OrderUuid, req.PaymentMethod.String())
	if err != nil {
		log.Printf("❗ Payment API error: %v", err)
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: paymentUuid,
	}, nil
}
