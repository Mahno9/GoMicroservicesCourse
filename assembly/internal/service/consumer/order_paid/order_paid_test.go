package orderpaid

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/converter/kafka/decoder"
	"github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service/mocks"
	"github.com/Mahno9/GoMicroservicesCourse/assembly/model"
	kafkaConsumerWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
)

type OrderPaidConsumerServiceSuite struct {
	suite.Suite

	ctx context.Context

	kafkaConsumer                *kafkaConsumerMock
	orderPaidDecoder             *decoder.OrderPaidDecoder
	shipAssembledProducerService *mocks.ShipAssembledProducerService
}

type kafkaConsumerMock struct {
	mock.Mock
	capturedHandler   kafkaConsumerWrapped.MessageHandler
	testMessage       *kafkaConsumerWrapped.Message
	shouldCallHandler bool
}

func (m *kafkaConsumerMock) Consume(ctx context.Context, handler kafkaConsumerWrapped.MessageHandler) error {
	m.capturedHandler = handler
	args := m.Called(ctx, handler)

	// If configured to call handler and test message is provided, invoke it
	if m.shouldCallHandler && m.testMessage != nil {
		return handler(ctx, *m.testMessage)
	}

	return args.Error(0)
}

// Helper method to configure the mock to call handler with test message
func (m *kafkaConsumerMock) SetupHandlerCall(message kafkaConsumerWrapped.Message) {
	m.testMessage = &message
	m.shouldCallHandler = true
}

func (s *OrderPaidConsumerServiceSuite) SetupTest() {
	err := logger.Init("info", false)
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.kafkaConsumer = &kafkaConsumerMock{}
	s.orderPaidDecoder = &decoder.OrderPaidDecoder{}
	s.shipAssembledProducerService = mocks.NewShipAssembledProducerService(s.T())
}

func (s *OrderPaidConsumerServiceSuite) TestRunConsumerSuccess() {
	s.kafkaConsumer.On("Consume", s.ctx, mock.Anything).Return(nil)

	service := NewService(s.kafkaConsumer, s.orderPaidDecoder, s.shipAssembledProducerService)
	err := service.RunConsumer(s.ctx)

	s.NoError(err)
	s.kafkaConsumer.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerServiceSuite) TestRunConsumerError() {
	expectedError := errors.New("kafka consumer error")
	s.kafkaConsumer.On("Consume", s.ctx, mock.Anything).Return(expectedError)

	service := NewService(s.kafkaConsumer, s.orderPaidDecoder, s.shipAssembledProducerService)
	err := service.RunConsumer(s.ctx)

	s.ErrorIs(err, expectedError)
	s.kafkaConsumer.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerServiceSuite) TestHandleOrderPaidSuccess() {
	// Note: This test takes ~10 seconds due to time.Sleep in handleOrderPaid method
	// This simulates the ship assembly process and tests the complete message flow
	var (
		eventUuid, _       = uuid.Parse(gofakeit.UUID())
		orderUuid, _       = uuid.Parse(gofakeit.UUID())
		userUuid, _        = uuid.Parse(gofakeit.UUID())
		paymentMethod      = gofakeit.CreditCardType()
		transactionUuid, _ = uuid.Parse(gofakeit.UUID())
	)

	// Create OrderPaid protobuf message
	orderPaidMsg := &eventsV1.OrderPaid{
		Uuid:            eventUuid.String(),
		OrderUuid:       orderUuid.String(),
		UserUuid:        userUuid.String(),
		PaymentMethod:   paymentMethod,
		TransactionUuid: transactionUuid.String(),
	}

	// Serialize the protobuf message
	payload, _ := proto.Marshal(orderPaidMsg)

	// Create Kafka message
	kafkaMsg := kafkaConsumerWrapped.Message{
		Value: payload,
	}

	// Setup mock to call handler with test message
	s.kafkaConsumer.SetupHandlerCall(kafkaMsg)
	s.kafkaConsumer.On("Consume", s.ctx, mock.Anything).Return(nil)

	// Setup producer mock to verify ShipAssembled event is produced
	s.shipAssembledProducerService.On("ProduceShipAssembled", s.ctx, mock.MatchedBy(func(shipAssembled model.ShipAssembled) bool {
		s.Equal(orderUuid.String(), shipAssembled.OrderUuid, "OrderUuid should match input")
		s.Equal(userUuid.String(), shipAssembled.UserUuid, "UserUuid should match input")
		s.NotEmpty(shipAssembled.EventUuid, "EventUuid should be generated")
		s.GreaterOrEqual(shipAssembled.BuildTimeSec, int64(0), "BuildTimeSec should be non-negative")
		return true
	})).Return(nil)

	service := NewService(s.kafkaConsumer, s.orderPaidDecoder, s.shipAssembledProducerService)
	err := service.RunConsumer(s.ctx)

	s.NoError(err)
	s.shipAssembledProducerService.AssertExpectations(s.T())
}

func (s *OrderPaidConsumerServiceSuite) TestHandleOrderPaidDecoderError() {
	// Create invalid protobuf data
	invalidPayload := []byte("invalid protobuf data")

	kafkaMsg := kafkaConsumerWrapped.Message{
		Value: invalidPayload,
	}

	// Setup mock to call handler with invalid message
	s.kafkaConsumer.SetupHandlerCall(kafkaMsg)
	s.kafkaConsumer.On("Consume", s.ctx, mock.Anything).Return(nil)

	service := NewService(s.kafkaConsumer, s.orderPaidDecoder, s.shipAssembledProducerService)
	err := service.RunConsumer(s.ctx)

	// Should return error due to decoder failure
	s.Error(err)
	s.Contains(err.Error(), "filaed to unmarshal event") // From decoder error message
}

func (s *OrderPaidConsumerServiceSuite) TestHandleOrderPaidProducerError() {
	var (
		eventUuid, _       = uuid.Parse(gofakeit.UUID())
		orderUuid, _       = uuid.Parse(gofakeit.UUID())
		userUuid, _        = uuid.Parse(gofakeit.UUID())
		paymentMethod      = gofakeit.CreditCardType()
		transactionUuid, _ = uuid.Parse(gofakeit.UUID())

		orderPaidMsg = &eventsV1.OrderPaid{
			Uuid:            eventUuid.String(),
			OrderUuid:       orderUuid.String(),
			UserUuid:        userUuid.String(),
			PaymentMethod:   paymentMethod,
			TransactionUuid: transactionUuid.String(),
		}

		payload, _ = proto.Marshal(orderPaidMsg)
		kafkaMsg   = kafkaConsumerWrapped.Message{
			Value: payload,
		}

		expectedError = errors.New("kafka producer error")
	)

	// Setup mock to call handler with test message
	s.kafkaConsumer.SetupHandlerCall(kafkaMsg)
	s.kafkaConsumer.On("Consume", s.ctx, mock.Anything).Return(nil)

	// Setup producer mock to return error
	s.shipAssembledProducerService.On("ProduceShipAssembled", s.ctx, mock.Anything).Return(expectedError)

	service := NewService(s.kafkaConsumer, s.orderPaidDecoder, s.shipAssembledProducerService)
	err := service.RunConsumer(s.ctx)

	s.ErrorIs(err, expectedError)
	s.shipAssembledProducerService.AssertExpectations(s.T())
}

func TestOrderPaidConsumerServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderPaidConsumerServiceSuite))
}
