package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderMainFlow() {
	var (
		orderUuid, _  = uuid.Parse(gofakeit.UUID())
		userUuid, _   = uuid.Parse(gofakeit.UUID())
		paymentMethod = int32(gofakeit.Number(0, 4))
		totalPrice    = gofakeit.Float64()

		transactionUuid, _ = uuid.Parse(gofakeit.UUID())

		payOrderData = model.PayOrderData{
			OrderUuid:     orderUuid,
			UserUuid:      uuid.Nil, // Will be set from order
			PaymentMethod: paymentMethod,
		}

		repoEntity = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			Status:     model.StatusPENDINGPAYMENT,
			TotalPrice: totalPrice,
		}
	)

	s.repository.On("Get", mock.Anything, payOrderData.OrderUuid).Return(repoEntity, nil)

	// В методе uuid должен подставиться из заказа.
	// Делаем здесь это руками для проверки структуры на входе в PayOrder.
	payOrderData.UserUuid = userUuid
	s.payment.On("PayOrder", mock.Anything, payOrderData).Return(transactionUuid, nil)

	s.repository.On("Update", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.OrderUuid == orderUuid &&
			order.UserUuid == userUuid &&
			len(order.PartUuids) == 0 && // В оригинальном заказе PartUuids не установлен
			order.TotalPrice == totalPrice &&
			order.Status == model.StatusPAID &&
			*order.TransactionUuid == transactionUuid &&
			order.PaymentMethod == paymentMethod
	})).Return(nil)

	resultTransactionUuid, err := s.service.PayOrder(s.ctx, payOrderData)
	s.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 1)
	s.Equal(transactionUuid, resultTransactionUuid)
}

func (s *ServiceSuite) TestPayOrderRepositoryGetError() {
	var (
		orderUuid, _  = uuid.Parse(gofakeit.UUID())
		paymentMethod = int32(gofakeit.Number(0, 4))

		payOrderData = model.PayOrderData{
			OrderUuid:     orderUuid,
			UserUuid:      uuid.Nil, // Will be set from order
			PaymentMethod: paymentMethod,
		}

		repoError = model.ErrOrderDoesNotExist
	)

	s.repository.On("Get", mock.Anything, payOrderData.OrderUuid).Return(nil, repoError)

	_, err := s.service.PayOrder(s.ctx, payOrderData)

	s.ErrorIs(err, repoError)
	s.payment.AssertNotCalled(s.T(), "PayOrder", s.ctx, payOrderData)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *ServiceSuite) TestPayOrderInvalidStatus() {
	var (
		orderUuid, _  = uuid.Parse(gofakeit.UUID())
		userUuid, _   = uuid.Parse(gofakeit.UUID())
		paymentMethod = int32(gofakeit.Number(0, 4))
		totalPrice    = gofakeit.Float64()

		payOrderData = model.PayOrderData{
			OrderUuid:     orderUuid,
			UserUuid:      uuid.Nil, // Will be set from order
			PaymentMethod: paymentMethod,
		}

		repoEntity = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			Status:     model.StatusCANCELLED, // Не pending payment
			TotalPrice: totalPrice,
		}
	)

	s.repository.On("Get", mock.Anything, payOrderData.OrderUuid).Return(repoEntity, nil)
	_, err := s.service.PayOrder(s.ctx, payOrderData)

	s.ErrorIs(err, model.ErrOrderCancelConflict)
	s.payment.AssertNotCalled(s.T(), "PayOrder", s.ctx, payOrderData)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *ServiceSuite) TestPayOrderPaymentServiceError() {
	var (
		orderUuid, _  = uuid.Parse(gofakeit.UUID())
		userUuid, _   = uuid.Parse(gofakeit.UUID())
		paymentMethod = int32(gofakeit.Number(0, 4))
		totalPrice    = gofakeit.Float64()

		payOrderData = model.PayOrderData{
			OrderUuid:     orderUuid,
			UserUuid:      uuid.Nil, // Will be set from order
			PaymentMethod: paymentMethod,
		}

		repoEntity = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			Status:     model.StatusPENDINGPAYMENT,
			TotalPrice: totalPrice,
		}

		paymentError = model.ErrUnknownPaymentMethod
	)

	s.repository.On("Get", mock.Anything, payOrderData.OrderUuid).Return(repoEntity, nil)

	payOrderData.UserUuid = userUuid
	s.payment.On("PayOrder", mock.Anything, payOrderData).Return(uuid.Nil, paymentError)

	_, err := s.service.PayOrder(s.ctx, payOrderData)

	s.ErrorIs(err, paymentError)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *ServiceSuite) TestPayOrderRepositoryUpdateError() {
	var (
		orderUuid, _  = uuid.Parse(gofakeit.UUID())
		userUuid, _   = uuid.Parse(gofakeit.UUID())
		paymentMethod = int32(gofakeit.Number(0, 4))
		totalPrice    = gofakeit.Float64()

		transactionUuid, _ = uuid.Parse(gofakeit.UUID())

		payOrderData = model.PayOrderData{
			OrderUuid:     orderUuid,
			UserUuid:      uuid.Nil, // Will be set from order
			PaymentMethod: paymentMethod,
		}

		repoEntity = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			Status:     model.StatusPENDINGPAYMENT,
			TotalPrice: totalPrice,
		}

		updateError = model.ErrOrderDoesNotExist
	)

	s.repository.On("Get", mock.Anything, payOrderData.OrderUuid).Return(repoEntity, nil)

	payOrderData.UserUuid = userUuid
	s.payment.On("PayOrder", mock.Anything, payOrderData).Return(transactionUuid, nil)

	s.repository.On("Update", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.OrderUuid == orderUuid &&
			order.UserUuid == userUuid &&
			order.TotalPrice == totalPrice &&
			order.Status == model.StatusPAID &&
			*order.TransactionUuid == transactionUuid &&
			order.PaymentMethod == paymentMethod
	})).Return(updateError)

	_, err := s.service.PayOrder(s.ctx, payOrderData)

	s.ErrorIs(err, updateError)
}
