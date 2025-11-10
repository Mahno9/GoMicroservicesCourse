package service

import (
	"context"
)

type PaymentService interface {
	Pay(ctx context.Context, userUuid, orderUuid, paymentMethod string) (string, error)
}
