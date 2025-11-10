package order

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
)

type service struct {
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
	orderRepository repository.OrderRepository
}

var _ def.OrderService = (*service)(nil)

func NewService(inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, repository repository.OrderRepository) *service {
	return &service{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: repository,
	}
}
