package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderMainFlow() {
	var (
		userUuid  = gofakeit.UUID()
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()
		partUuid3 = gofakeit.UUID()

		partPrice1 = gofakeit.Float64()
		partPrice2 = gofakeit.Float64()
		partPrice3 = gofakeit.Float64()

		createOrderData = model.CreateOrderData{
			UserUuid:  userUuid,
			PartUuids: []string{partUuid1, partUuid2, partUuid3},
		}

		orderParts = []*model.Part{
			{
				Uuid:  partUuid1,
				Price: partPrice1,
			},
			{
				Uuid:  partUuid2,
				Price: partPrice2,
			},
			{
				Uuid:  partUuid3,
				Price: partPrice3,
			},
		}

		totalPrice = partPrice1 + partPrice2 + partPrice3

		expectedOrder = &model.Order{
			OrderUuid:  gofakeit.UUID(), // Будет сгенерирован в функции
			UserUuid:   userUuid,
			PartUuids:  []string{partUuid1, partUuid2, partUuid3},
			TotalPrice: totalPrice,
			Status:     model.StatusPENDINGPAYMENT, // По умолчанию при создании
		}
	)

	s.inventory.On("ListParts", mock.Anything, &model.PartsFilter{
		Uuids: createOrderData.PartUuids,
	}).Return(orderParts, nil)

	s.repository.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.UserUuid == userUuid &&
			len(order.PartUuids) == 3 &&
			order.PartUuids[0] == partUuid1 &&
			order.PartUuids[1] == partUuid2 &&
			order.PartUuids[2] == partUuid3 &&
			order.TotalPrice == totalPrice
	})).Return(expectedOrder, nil)

	result, err := s.service.CreateOrder(s.ctx, createOrderData)

	s.NoError(err)
	s.Equal(expectedOrder, result)
}

func (s *ServiceSuite) TestCreateOrderInventoryError() {
	var (
		userUuid  = gofakeit.UUID()
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		createOrderData = model.CreateOrderData{
			UserUuid:  userUuid,
			PartUuids: []string{partUuid1, partUuid2},
		}

		expectedError = model.ErrPartsNotAvailable
	)

	s.inventory.On("ListParts", mock.Anything, &model.PartsFilter{
		Uuids: createOrderData.PartUuids,
	}).Return(nil, expectedError)

	_, err := s.service.CreateOrder(s.ctx, createOrderData)

	s.ErrorIs(err, expectedError)
	s.repository.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *ServiceSuite) TestCreateOrderPartsNotAvailable() {
	var (
		userUuid  = gofakeit.UUID()
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()
		partUuid3 = gofakeit.UUID()

		createOrderData = model.CreateOrderData{
			UserUuid:  userUuid,
			PartUuids: []string{partUuid1, partUuid2, partUuid3},
		}

		orderParts = []*model.Part{
			{
				Uuid: partUuid1,
			},
			{
				Uuid: partUuid2,
			},
			// Только 2 части из 3 запрошенных
		}

		expectedError = model.ErrPartsNotAvailable
	)

	s.inventory.On("ListParts", mock.Anything, &model.PartsFilter{
		Uuids: createOrderData.PartUuids,
	}).Return(orderParts, nil)

	_, err := s.service.CreateOrder(s.ctx, createOrderData)

	s.ErrorIs(err, expectedError)
	s.repository.AssertNumberOfCalls(s.T(), "Create", 0)
}

func (s *ServiceSuite) TestCreateOrderRepositoryCreateError() {
	var (
		userUuid  = gofakeit.UUID()
		partUuid1 = gofakeit.UUID()
		partUuid2 = gofakeit.UUID()

		partPrice1 = gofakeit.Float64()
		partPrice2 = gofakeit.Float64()

		createOrderData = model.CreateOrderData{
			UserUuid:  userUuid,
			PartUuids: []string{partUuid1, partUuid2},
		}

		orderParts = []*model.Part{
			{
				Uuid:  partUuid1,
				Price: partPrice1,
			},
			{
				Uuid:  partUuid2,
				Price: partPrice2,
			},
		}

		expectedError = model.ErrOrderDoesNotExist
	)

	s.inventory.On("ListParts", mock.Anything, &model.PartsFilter{
		Uuids: createOrderData.PartUuids,
	}).Return(orderParts, nil)

	s.repository.On("Create", mock.Anything, mock.MatchedBy(func(order *model.Order) bool {
		return order.UserUuid == userUuid &&
			len(order.PartUuids) == 2 &&
			order.PartUuids[0] == partUuid1 &&
			order.PartUuids[1] == partUuid2 &&
			order.TotalPrice == partPrice1+partPrice2
	})).Return(nil, expectedError)

	_, err := s.service.CreateOrder(s.ctx, createOrderData)

	s.ErrorIs(err, expectedError)
}
