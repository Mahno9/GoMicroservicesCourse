package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *ServiceSuite) TestOrderCancelMainFlow() {
	var (
		orderUuid  = gofakeit.UUID()
		userUuid   = gofakeit.UUID()
		partUuid1  = gofakeit.UUID()
		partUuid2  = gofakeit.UUID()
		totalPrice = gofakeit.Float64()

		order = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			PartUuids:  []string{partUuid1, partUuid2},
			TotalPrice: totalPrice,
			Status:     model.StatusPENDINGPAYMENT,
		}
	)

	s.repository.On("Get", orderUuid).Return(order, nil)
	s.repository.On("Update", mock.MatchedBy(func(order *model.Order) bool {
		return order.OrderUuid == orderUuid &&
			order.UserUuid == userUuid &&
			len(order.PartUuids) == 2 &&
			order.PartUuids[0] == partUuid1 &&
			order.PartUuids[1] == partUuid2 &&
			order.TotalPrice == totalPrice &&
			order.Status == model.StatusCANCELLED
	})).Return(nil)

	err := s.service.OrderCancel(s.ctx, orderUuid)

	s.NoError(err)
}

func (s *ServiceSuite) TestOrderCancelRepositoryGetError() {
	var (
		orderUuid     = gofakeit.UUID()
		expectedError = model.ErrOrderDoesNotExist
	)

	s.repository.On("Get", orderUuid).Return(nil, expectedError)

	err := s.service.OrderCancel(s.ctx, orderUuid)

	s.ErrorIs(err, expectedError)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *ServiceSuite) TestOrderCancelInvalidStatus() {
	var (
		orderUuid  = gofakeit.UUID()
		userUuid   = gofakeit.UUID()
		partUuid1  = gofakeit.UUID()
		totalPrice = gofakeit.Float64()

		order = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			PartUuids:  []string{partUuid1},
			TotalPrice: totalPrice,
			Status:     model.StatusPAID, // Не pending payment
		}

		expectedError = model.ErrOrderCancelConflict
	)

	s.repository.On("Get", orderUuid).Return(order, nil)

	err := s.service.OrderCancel(s.ctx, orderUuid)

	s.ErrorIs(err, expectedError)
	s.repository.AssertNumberOfCalls(s.T(), "Update", 0)
}

func (s *ServiceSuite) TestOrderCancelRepositoryUpdateError() {
	var (
		orderUuid  = gofakeit.UUID()
		userUuid   = gofakeit.UUID()
		partUuid1  = gofakeit.UUID()
		totalPrice = gofakeit.Float64()

		order = &model.Order{
			OrderUuid:  orderUuid,
			UserUuid:   userUuid,
			PartUuids:  []string{partUuid1},
			TotalPrice: totalPrice,
			Status:     model.StatusPENDINGPAYMENT,
		}

		expectedError = model.ErrOrderDoesNotExist
	)

	s.repository.On("Get", orderUuid).Return(order, nil)
	s.repository.On("Update", mock.MatchedBy(func(order *model.Order) bool {
		return order.OrderUuid == orderUuid &&
			order.UserUuid == userUuid &&
			len(order.PartUuids) == 1 &&
			order.PartUuids[0] == partUuid1 &&
			order.TotalPrice == totalPrice &&
			order.Status == model.StatusCANCELLED
	})).Return(expectedError)

	err := s.service.OrderCancel(s.ctx, orderUuid)

	s.ErrorIs(err, expectedError)
}
