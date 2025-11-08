package payment

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

const (
	payDelay = 1 * time.Second
)

func (s *service) Pay(ctx context.Context, userUuid, orderUuid, paymentMethod string) (string, error) {
	if dl, ok := ctx.Deadline(); ok {
		log.Printf("⌛ Context with deadline: %v\n", time.Until(dl))
	} else {
		log.Printf("⌛ Context with no timeout\n")
	}

	timer := time.NewTimer(payDelay)
	defer timer.Stop()

	select {
	case <-timer.C:
		paymentUuid := uuid.New().String()

		log.Printf("🆗 Оплата прошла успешно:\nuser uuid: %s, order uuid: %s, newly created payment uuid: %s\n", userUuid, orderUuid, paymentUuid)
		return paymentUuid, nil

	case <-ctx.Done():
		return "", ctx.Err()
	}
}
