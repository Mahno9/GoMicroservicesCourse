package producer

import (
	services "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
)

type service struct {
	orderPaidProducer kafkaWrapped.Producer
}

func NewService(orderPaidProducer kafkaWrapped.Producer) services.ProducerService {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}
