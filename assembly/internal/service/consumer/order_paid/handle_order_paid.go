package orderpaid

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/model"
	kafkaConsumerWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

func (s *service) handleOrderPaid(ctx context.Context, msg kafkaConsumerWrapped.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "OrderPaid event received", zap.String("uuid", event.Uuid))

	startAssemblyTime := time.Now()

	logger.Info(ctx, "Assembling...")
	time.Sleep(10 * time.Second) //nolint:forbidigo

	assemblyRemainTimeSec := time.Since(startAssemblyTime).Seconds()
	logger.Info(ctx, fmt.Sprintf("Assembling was accomplished in %.2f seconds", assemblyRemainTimeSec))

	err = s.shipAssembledProducer.ProduceShipAssembled(ctx, model.ShipAssembled{
		EventUuid:    uuid.New().String(),
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: int64(assemblyRemainTimeSec),
	})
	if err != nil {
		logger.Error(ctx, "Failed to produce ShipAssembled event", zap.Error(err))
		return err
	}

	return nil
}
