package decoder

import (
	"github.com/Mahno9/GoMicroservicesCourse/notification/model"
	eventsV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type ShipAssembledDecoder struct{}

func NewShipAssembledDecoder() *ShipAssembledDecoder {
	return &ShipAssembledDecoder{}
}

func (d *ShipAssembledDecoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var protoEvent eventsV1.ShipAssembled
	err := proto.Unmarshal(data, &protoEvent)
	if err != nil {
		return model.ShipAssembledEvent{}, err
	}

	return model.ShipAssembledEvent{
		OrderID:    protoEvent.GetOrderUuid(),
		UserID:     protoEvent.GetUserUuid(),
		TrackingID: protoEvent.GetEventUuid(),
	}, nil
}
