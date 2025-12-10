package service

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/model"
)

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipAssembledProducerService interface {
	ProduceShipAssembled(ctx context.Context, shipAssembled model.ShipAssembled) error
}
