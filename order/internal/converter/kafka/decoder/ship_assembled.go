package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	model "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewShipAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var event eventsV1.ShipAssembled
	err := proto.Unmarshal(data, &event)
	if err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}, nil
}
