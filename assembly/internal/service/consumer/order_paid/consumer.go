package orderpaid

import (
	"context"

	"go.uber.org/zap"

	kafkaDecoder "github.com/Mahno9/GoMicroservicesCourse/assembly/converter/kafka"
	services "github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type service struct {
	shipAssembledProducer services.ShipAssembledProducerService

	orderPaidConsumer kafkaWrapped.Consumer
	orderPaidDecoder  kafkaDecoder.OrderPaidDecoder
}

func NewService(orderPaidConsumer kafkaWrapped.Consumer, orderPaidDecoder kafkaDecoder.OrderPaidDecoder, producerService services.ShipAssembledProducerService) services.OrderPaidConsumerService {
	return &service{
		shipAssembledProducer: producerService,
		orderPaidConsumer:     orderPaidConsumer,
		orderPaidDecoder:      orderPaidDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderPaid consumer")

	err := s.orderPaidConsumer.Consume(ctx, s.handleOrderPaid)
	if err != nil {
		logger.Error(ctx, "Failed to start OrderPaid consumer", zap.Error(err))
		return err
	}

	return nil
}
