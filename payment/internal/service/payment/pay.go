package payment

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

func (s *service) Pay(c context.Context, userUuid, orderUuid, paymentMethod string) (string, error) {
	// TODO: Validate?

	if dl, ok := c.Deadline(); ok {
		log.Printf("⌛ Context with deadline: %v\n", time.Until(dl))
	} else {
		log.Printf("⌛ Context with no timeout\n")
	}

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-c.Done():
		return "", c.Err()
	}

	paymentUuid := uuid.New().String()
	log.Printf("🆗 Оплата прошла успешно:\nuser uuid: %s, order uuid: %s, newly created payment uuid: %s\n", userUuid, orderUuid, paymentUuid)

	return paymentUuid, nil
}
