package assembly

import (
	services "github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
)

type service struct {
	shipAssembledProducer kafkaWrapped.Producer
}

func NewService(shipAssembledProducer kafkaWrapped.Producer) services.ShipAssembledProducerService {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}
