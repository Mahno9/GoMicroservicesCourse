package ship_assembled

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/notification/internal/converter/kafka"
	services "github.com/Mahno9/GoMicroservicesCourse/notification/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	"go.uber.org/zap"
)

type service struct {
	shipAssembledConsumer kafkaWrapped.Consumer
	shipAssembledDecoder  kafka.ShipAssembledDecoder
	telegramService       services.TelegramService
}

func NewService(consumer kafkaWrapped.Consumer, decoder kafka.ShipAssembledDecoder, telegramService services.TelegramService) *service {
	return &service{
		shipAssembledConsumer: consumer,
		shipAssembledDecoder:  decoder,
		telegramService:       telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ShipAssembled consumer")

	err := s.shipAssembledConsumer.Consume(ctx, s.handleShipAssembledMessage)
	if err != nil {
		logger.Error(ctx, "Failed to start ShipAssembled consumer", zap.Error(err))
		return err
	}

	return nil
}
