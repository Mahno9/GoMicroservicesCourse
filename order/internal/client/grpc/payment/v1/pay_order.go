package v1

import (
	"context"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, paymentData model.PayOrderData) (string, error) {
	timedContext, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	payOrderResponse, err := c.service.PayOrder(timedContext, &paymentV1.PayOrderRequest{
		OrderUuid:     paymentData.OrderUuid,
		UserUuid:      paymentData.UserUuid,
		PaymentMethod: paymentV1.PaymentMethod(paymentData.PaymentMethod),
	})
	if err != nil {
		return "", err
	}

	transactionUuid := payOrderResponse.TransactionUuid
	return transactionUuid, nil
}
