package v1

import (
	"context"
	"log"

	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *genPaymentV1.PayOrderRequest) (*genPaymentV1.PayOrderResponse, error) {
	paymentUuid, err := a.paymentService.Pay(ctx, req.UserUuid, req.OrderUuid, req.PaymentMethod.String())
	if err != nil {
		log.Printf("❗ Payment API error: %v", err)
		return nil, err
	}

	return &genPaymentV1.PayOrderResponse{
		TransactionUuid: paymentUuid,
	}, nil
}
