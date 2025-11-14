package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/converter"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, paymentData model.PayOrderData) (uuid.UUID, error) {
	payOrderResponse, err := c.service.PayOrder(ctx, &genPaymentV1.PayOrderRequest{
		OrderUuid:     paymentData.OrderUuid.String(),
		UserUuid:      paymentData.UserUuid.String(),
		PaymentMethod: converter.ModelToPaymentPaymentMethod(paymentData.PaymentMethod),
	})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(payOrderResponse.TransactionUuid)
}
