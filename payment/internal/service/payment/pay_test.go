package payment

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

var paymentMethods = []string{"1", "2", "3"}

func (s *ServiceSuite) TestPay() {
	var (
		userUuid      = gofakeit.UUID()
		orderUuid     = gofakeit.UUID()
		paymentMethod = gofakeit.RandomString(paymentMethods)
	)

	transactionUuid, err := s.service.Pay(s.ctx, userUuid, orderUuid, paymentMethod)

	s.NoError(err)
	s.NotEmpty(transactionUuid)
}

func (s *ServiceSuite) TestPayTimeout() {
	// Создаем контекст с тайм-аутом 500мс, что меньше, чем payDelay (1с)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	var (
		userUuid      = gofakeit.UUID()
		orderUuid     = gofakeit.UUID()
		paymentMethod = gofakeit.RandomString(paymentMethods)
	)

	// Вызываем Pay с контекстом, который истечет раньше, чем завершится операция
	transactionUuid, err := s.service.Pay(ctx, userUuid, orderUuid, paymentMethod)

	// Проверяем, что произошла ошибка тайм-аута
	s.Error(err)
	s.Empty(transactionUuid)
	s.Equal(context.DeadlineExceeded, err)
}

func (s *ServiceSuite) TestPayCancelledContext() {
	ctx, cancel := context.WithCancel(context.Background())

	var (
		userUuid      = gofakeit.UUID()
		orderUuid     = gofakeit.UUID()
		paymentMethod = gofakeit.RandomString(paymentMethods)
	)

	type PayResult struct {
		uuid string
		err  error
	}

	resultChan := make(chan PayResult, 1)

	go func() {
		uuid, err := s.service.Pay(ctx, userUuid, orderUuid, paymentMethod)
		resultChan <- PayResult{uuid: uuid, err: err}
	}()

	time.Sleep(100 * time.Millisecond) // nolint: forbidigo

	cancel()

	result := <-resultChan

	s.Error(result.err)
	s.Empty(result.uuid)
	s.Equal(context.Canceled, result.err)
}
