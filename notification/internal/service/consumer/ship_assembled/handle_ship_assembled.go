package ship_assembled

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/notification/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func (s *service) handleShipAssembledMessage(ctx context.Context, message consumer.Message) error {
	logger.Info(ctx, "Processing ShipAssembled message")

	event, err := s.shipAssembledDecoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled message", zap.Error(err))
		return err
	}

	return s.handleShipAssembled(ctx, event)
}

func (s *service) handleShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	logger.Info(ctx, "Handling ShipAssembled event", zap.Any("event", event))

	err := s.telegramService.SendShipAssembledMessage(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send telegram notification", zap.Error(err))
		return model.ErrFailedToSendNotification
	}

	logger.Info(ctx, "Successfully sent ShipAssembled notification",
		zap.String("orderID", event.OrderID),
		zap.String("trackingID", event.TrackingID),
		zap.String("userID", event.UserID))

	return nil
}
