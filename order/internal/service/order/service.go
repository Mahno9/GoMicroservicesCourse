package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/order"
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
	"github.com/google/uuid"
)

type service struct {
	inventory  grpc.InventoryClient
	payment    grpc.PaymentClient
	ordersRepo repository.OrderRepository
}

var _ def.OrderService = (*service)(nil)

func NewService(inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		inventory:  inventoryClient,
		payment:    paymentClient,
		ordersRepo: order.NewRepository(),
	}
}

func (s *service) CreateOrder(c context.Context, data model.CreateOrderData) (model.CreateOrderData, error) {
	orderParts, err := s.inventory.ListParts(c, &model.PartsFilter{
		Uuids: data.PartUuids,
	})
	if err != nil {
		return model.CreateOrderData{}, err
	}

	if len(orderParts) != len(data.PartUuids) {
		return model.CreateOrderData{}, model.PartsNotAvailableErr
	}

	totalPrice := float64(0.0)
	for _, part := range orderParts {
		totalPrice += part.Price
	}

	orderUuid := uuid.New().String()
	err = s.ordersRepo.Create(&model.Order{
		OrderUuid:  orderUuid,
		UserUuid:   data.UserUuid,
		PartUuids:  data.PartUuids,
		TotalPrice: totalPrice,
	})
	if err != nil {
		return model.CreateOrderData{}, err
	}

	return data, nil
}

func (s *service) GetOrder(c context.Context, orderUuid string) (model.Order, error) {
	// TODO:
}

func (s *service) OrderCnacel(c context.Context, orderUuid string) error {
	// TODO: converter to v1 error response type
}

func (s *service) PayOrder(c context.Context, data model.PayOrderData) error {
	// TODO: converter to v1 error response type
}
