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

	message := s.formatOrderPaidMessage(event)

	err := s.telegramService.BroadcastMessage(ctx, message)
	if err != nil {
		logger.Error(ctx, "Failed to send telegram notification", zap.Error(err))
		return model.ErrFailedToSendNotification
	}

	logger.Info(ctx, "Successfully sent OrderPaid notification",
		zap.String("orderID", event.OrderID),
		zap.String("userID", event.UserID))

	return nil
}

func (s *service) formatOrderPaidMessage(event model.OrderPaidEvent) string {
	return "🎉 Order Paid!\n\n" +
		"Order ID: " + event.OrderID + "\n" +
		"Payment Method: " + event.PaymentMethod + "\n" +
		"Transaction ID: " + event.TransactionUUID + "\n" +
		"Thank you for your payment!"
}
