package v1

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/converter"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, paymentData model.PayOrderData) (string, error) {
	payOrderResponse, err := c.service.PayOrder(ctx, &genPaymentV1.PayOrderRequest{
		OrderUuid:     paymentData.OrderUuid,
		UserUuid:      paymentData.UserUuid,
		PaymentMethod: converter.ModelToPaymentPaymentMethod(paymentData.PaymentMethod),
	})
	if err != nil {
		return "", err
	}

	transactionUuid := payOrderResponse.TransactionUuid
	return transactionUuid, nil
}
