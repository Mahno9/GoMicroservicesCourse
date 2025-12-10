package assembly

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/model"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
)

func (s *service) ProduceShipAssembled(ctx context.Context, data model.ShipAssembled) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    data.EventUuid,
		OrderUuid:    data.OrderUuid,
		UserUuid:     data.UserUuid,
		BuildTimeSec: data.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%w: %w", model.ErrKafkaInvalidInputData, err)
	}

	err = s.shipAssembledProducer.Send(ctx, []byte(data.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "Failed to send ShipAssembled event", zap.Error(err))
		return err
	}

	return nil
}
