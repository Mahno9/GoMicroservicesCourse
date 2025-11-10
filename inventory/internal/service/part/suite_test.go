package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	repoMocks "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	repository *repoMocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupSuite() {
	s.ctx = context.Background()

	s.repository = repoMocks.NewPartRepository(s.T())

	s.service = NewService(s.repository)
}

func (s *ServiceSuite) SetupTest() {
	s.repository.ExpectedCalls = nil
	s.repository.Calls = nil
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
