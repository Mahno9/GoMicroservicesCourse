package kafka

import "github.com/Mahno9/GoMicroservicesCourse/order/internal/model"

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
