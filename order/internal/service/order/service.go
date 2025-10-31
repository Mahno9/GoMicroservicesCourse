package order

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
)

type service struct {
}

var _ def.OrderService = (*service)(nil)

func (s *service) CreateOrder(c context.Context, data model.CreateOrderData) (model.CreateOrderData, error) {
	// TODO:
}

func (s *service) GetOrder(c context.Context, orderUuid string) (model.OrderInfo, error) {
	// TODO:
}

func (s *service) OrderCnacel(c context.Context, orderUuid string) error {
	// TODO: converter to v1 error response type
}

func (s *service) PayOrder(c context.Context, data model.PayOrderData) error {
	// TODO: converter to v1 error response type
}
