package order

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
	services "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
)

type service struct {
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
	orderRepository repository.OrderRepository
	producerService services.ProducerService
}

func NewService(inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, repository repository.OrderRepository, producerService services.ProducerService) services.OrderService {
	return &service{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: repository,
		producerService: producerService,
	}
}
