package assembly

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
)

type ShipAssembledProducerServiceSuite struct {
	suite.Suite

	ctx context.Context

	kafkaProducer *kafkaProducerMock
}

type kafkaProducerMock struct {
	mock.Mock
}

func (m *kafkaProducerMock) Send(ctx context.Context, key []byte, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (s *ShipAssembledProducerServiceSuite) SetupTest() {
	err := logger.Init("info", false)
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.kafkaProducer = &kafkaProducerMock{}
}

func (s *ShipAssembledProducerServiceSuite) TestProduceShipAssembledSuccess() {
	var (
		eventUuid, _ = uuid.Parse(gofakeit.UUID())
		orderUuid, _ = uuid.Parse(gofakeit.UUID())
		userUuid, _  = uuid.Parse(gofakeit.UUID())
		buildTimeSec = int64(gofakeit.IntRange(1, 100))

		shipAssembled = model.ShipAssembled{
			EventUuid:    eventUuid.String(),
			OrderUuid:    orderUuid.String(),
			UserUuid:     userUuid.String(),
			BuildTimeSec: buildTimeSec,
		}
	)

	s.kafkaProducer.On("Send", s.ctx, []byte(eventUuid.String()),
		mock.MatchedBy(func(payload []byte) bool {
			// Unmarshal the protobuf message
			var msg eventsV1.ShipAssembled
			err := proto.Unmarshal(payload, &msg)
			s.Require().NoError(err, "Failed to unmarshal ShipAssembled protobuf message")

			// Verify all fields match the input data
			s.Equal(eventUuid.String(), msg.EventUuid, "EventUuid should match input")
			s.Equal(orderUuid.String(), msg.OrderUuid, "OrderUuid should match input")
			s.Equal(userUuid.String(), msg.UserUuid, "UserUuid should match input")
			s.Equal(buildTimeSec, msg.BuildTimeSec, "BuildTimeSec should match input")

			return true
		})).Return(nil)

	service := NewService(s.kafkaProducer)
	err := service.ProduceShipAssembled(s.ctx, shipAssembled)

	s.NoError(err)
	s.kafkaProducer.AssertExpectations(s.T())
}

func (s *ShipAssembledProducerServiceSuite) TestProduceShipAssembledProducerError() {
	var (
		eventUuid, _ = uuid.Parse(gofakeit.UUID())
		orderUuid, _ = uuid.Parse(gofakeit.UUID())
		userUuid, _  = uuid.Parse(gofakeit.UUID())
		buildTimeSec = int64(gofakeit.IntRange(1, 100))

		shipAssembled = model.ShipAssembled{
			EventUuid:    eventUuid.String(),
			OrderUuid:    orderUuid.String(),
			UserUuid:     userUuid.String(),
			BuildTimeSec: buildTimeSec,
		}

		expectedError = errors.New("kafka producer error")
	)

	s.kafkaProducer.On("Send", s.ctx, []byte(eventUuid.String()), mock.AnythingOfType("[]uint8")).Return(expectedError)

	service := NewService(s.kafkaProducer)
	err := service.ProduceShipAssembled(s.ctx, shipAssembled)

	s.ErrorIs(err, expectedError)
	s.kafkaProducer.AssertExpectations(s.T())
}

func TestShipAssembledProducerServiceSuite(t *testing.T) {
	suite.Run(t, new(ShipAssembledProducerServiceSuite))
}
