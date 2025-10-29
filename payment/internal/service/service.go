package service

import (
	"context"
)

type PaymentService interface {
	Pay(ctx context.Context, userUuid string, orderUuid string, paymentMethod string) (string, error)
}
