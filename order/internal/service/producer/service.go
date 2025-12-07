package producer

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	services "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafkaWrapped.Producer
}

func NewService(orderPaidProducer kafkaWrapped.Producer) services.ProducerService {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := &eventsV1.OrderPaid{
		Uuid:            event.Uuid,
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUuid,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%w: %w", model.ErrKafkaInvalidInputData, err)
	}

	err = s.orderPaidProducer.Send(ctx, []byte(event.OrderUuid), payload)
	if err != nil {
		logger.Error(ctx, "Failed to send OrderPaid event", zap.Error(err))
		return fmt.Errorf("%w: %w", model.ErrKafkaSend, err)
	}

	return nil
}
