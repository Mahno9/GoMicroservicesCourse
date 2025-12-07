package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc/mocks"
	repoMocks "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/mocks"
	services "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
	serviceMocks "github.com/Mahno9/GoMicroservicesCourse/order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	inventory *clientMocks.InventoryClient
	payment   *clientMocks.PaymentClient

	repository *repoMocks.OrderRepository

	producerService *serviceMocks.ProducerService

	service services.OrderService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.inventory = clientMocks.NewInventoryClient(s.T())
	s.payment = clientMocks.NewPaymentClient(s.T())
	s.repository = repoMocks.NewOrderRepository(s.T())
	s.producerService = serviceMocks.NewProducerService(s.T())

	s.service = NewService(s.inventory, s.payment, s.repository, s.producerService)
}

func (s *ServiceSuite) SetupTest() {
	s.repository.ExpectedCalls = nil
	s.repository.Calls = nil
	s.inventory.ExpectedCalls = nil
	s.inventory.Calls = nil
	s.payment.ExpectedCalls = nil
	s.payment.Calls = nil
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
