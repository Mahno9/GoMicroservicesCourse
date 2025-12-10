package order_paid

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/notification/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func (s *service) handleOrderPaidMessage(ctx context.Context, message consumer.Message) error {
	logger.Info(ctx, "Processing OrderPaid message")

	event, err := s.orderPaidDecoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid message", zap.Error(err))
		return err
	}

	return s.handleOrderPaid(ctx, event)
}

func (s *service) handleOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	logger.Info(ctx, "Handling OrderPaid event", zap.Any("event", event))

	err := s.telegramService.SendOrderPaidMessage(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send telegram notification", zap.Error(err))
		return model.ErrFailedToSendNotification
	}

	logger.Info(ctx, "Successfully sent OrderPaid notification",
		zap.String("orderID", event.OrderID),
		zap.String("userID", event.UserID))

	return nil
}
