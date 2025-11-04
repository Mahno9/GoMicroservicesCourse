package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderMainFlow() {
	var (
		orderUuid   = gofakeit.UUID()
		userUuid    = gofakeit.UUID()
		partUuid1   = gofakeit.UUID()
		partUuid2   = gofakeit.UUID()
		totalPrice  = gofakeit.Float64()
		paymentData = gofakeit.UUID()

		expectedOrder = &model.Order{
			OrderUuid:       orderUuid,
			UserUuid:        userUuid,
			PartUuids:       []string{partUuid1, partUuid2},
			TotalPrice:      totalPrice,
			Status:          model.StatusPAID,
			TransactionUuid: paymentData,
			PaymentMethod:   int32(gofakeit.Number(0, 4)),
		}
	)

	s.repository.On("Get", orderUuid).Return(expectedOrder, nil)

	result, err := s.service.GetOrder(s.ctx, orderUuid)

	s.NoError(err)
	s.Equal(expectedOrder, result)
}

func (s *ServiceSuite) TestGetOrderRepositoryError() {
	var (
		orderUuid     = gofakeit.UUID()
		expectedError = model.ErrOrderDoesNotExist
	)

	s.repository.On("Get", orderUuid).Return(nil, expectedError)

	_, err := s.service.GetOrder(s.ctx, orderUuid)

	s.ErrorIs(err, expectedError)
}
