package order

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/order"
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
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
