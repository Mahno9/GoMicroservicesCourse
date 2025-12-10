package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderPaidConsumerService     *mocks.OrderPaidConsumerService
	shipAssembledProducerService *mocks.ShipAssembledProducerService
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.orderPaidConsumerService = mocks.NewOrderPaidConsumerService(s.T())
	s.shipAssembledProducerService = mocks.NewShipAssembledProducerService(s.T())
}

func (s *ServiceSuite) SetupTest() {
	s.orderPaidConsumerService.ExpectedCalls = nil
	s.orderPaidConsumerService.Calls = nil
	s.shipAssembledProducerService.ExpectedCalls = nil
	s.shipAssembledProducerService.Calls = nil
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
