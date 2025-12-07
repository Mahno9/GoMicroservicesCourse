package consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Mahno9/GoMicroservicesCourse/order/internal/converter/kafka"
	services "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type service struct {
	shipAssembledConsumer kafkaWrapped.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
}

func NewService(shipAssebledConsumer kafkaWrapped.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder) services.ConsumerService {
	return &service{
		shipAssembledConsumer: shipAssebledConsumer,
		shipAssembledDecoder:  shipAssembledDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ShipAssembledConsumer service")

	err := s.shipAssembledConsumer.Consume(ctx, s.handleShipAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to start ShipAssembledConsumer service", zap.Error(err))
		return err
	}

	return nil
}
