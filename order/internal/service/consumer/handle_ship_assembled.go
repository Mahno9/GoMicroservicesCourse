package consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConsumerWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func (s *service) handleShipAssembled(ctx context.Context, msg kafkaConsumerWrapped.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ship assembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Ship assembled event received", zap.String("order_uuid", event.OrderUuid))
	// TODO: send message to Telegram

	return nil
}
