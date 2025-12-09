package order_paid

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/converter/kafka"
	services "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	orderPaidConsumer kafkaWrapped.Consumer
	orderPaidDecoder  kafka.OrderPaidDecoder
	telegramService   services.TelegramService
}

func NewService(consumer kafkaWrapped.Consumer, decoder kafka.OrderPaidDecoder, telegramService services.TelegramService) *service {
	return &service{
		orderPaidConsumer: consumer,
		orderPaidDecoder:  decoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting OrderPaid consumer")

	err := s.orderPaidConsumer.Consume(ctx, s.handleOrderPaidMessage)
	if err != nil {
		logger.Error(ctx, "Failed to start OrderPaid consumer", zap.Error(err))
		return err
	}

	return nil
}
