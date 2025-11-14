package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *ServiceSuite) TestGetOrderMainFlow() {
	var (
		orderUuid, _       = uuid.Parse(gofakeit.UUID())
		userUuid, _        = uuid.Parse(gofakeit.UUID())
		partUuid1, _       = uuid.Parse(gofakeit.UUID())
		partUuid2, _       = uuid.Parse(gofakeit.UUID())
		totalPrice         = gofakeit.Float64()
		transactionUuid, _ = uuid.Parse(gofakeit.UUID())

		expectedOrder = &model.Order{
			OrderUuid:       orderUuid,
			UserUuid:        userUuid,
			PartUuids:       []uuid.UUID{partUuid1, partUuid2},
			TotalPrice:      totalPrice,
			Status:          model.StatusPAID,
			TransactionUuid: &transactionUuid,
			PaymentMethod:   int32(gofakeit.Number(0, 4)),
		}
	)

	s.repository.On("Get", mock.Anything, orderUuid).Return(expectedOrder, nil)

	result, err := s.service.GetOrder(s.ctx, orderUuid)

	s.NoError(err)
	s.Equal(expectedOrder, result)
}

func (s *ServiceSuite) TestGetOrderRepositoryError() {
	var (
		orderUuid, _  = uuid.Parse(gofakeit.UUID())
		expectedError = model.ErrOrderDoesNotExist
	)

	s.repository.On("Get", mock.Anything, orderUuid).Return(nil, expectedError)

	_, err := s.service.GetOrder(s.ctx, orderUuid)

	s.ErrorIs(err, expectedError)
}
