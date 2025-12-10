package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/Mahno9/GoMicroservicesCourse/assembly/config"
	kafkaDecoders "github.com/Mahno9/GoMicroservicesCourse/assembly/converter/kafka"
	kafkaDecoder "github.com/Mahno9/GoMicroservicesCourse/assembly/converter/kafka/decoder"
	services "github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service"
	orderPaidConsumerService "github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service/consumer/order_paid"
	shipAssembledProducerService "github.com/Mahno9/GoMicroservicesCourse/assembly/internal/service/producer/ship_assembled"
	kafkaWrapped "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	kafkaWrappedConsumer "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	kafkaWrappedProducer "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/producer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
)

type diContainer struct {
	config *config.Config

	orderPaidConsumerService services.OrderPaidConsumerService
	orderPaidConsumer        kafkaWrapped.Consumer
	orderPaidConsumerGroup   sarama.ConsumerGroup
	orderPaidDecoder         kafkaDecoders.OrderPaidDecoder

	shipAssembledProducerService services.ShipAssembledProducerService
	shipAssembledProducer        kafkaWrapped.Producer
	shipAssembledSyncPruducer    sarama.SyncProducer
}

func NewDIContainer(cfg *config.Config) *diContainer {
	return &diContainer{config: cfg}
}

func (c *diContainer) ConsumerService(ctx context.Context) services.OrderPaidConsumerService {
	if c.orderPaidConsumerService == nil {
		c.orderPaidConsumerService = orderPaidConsumerService.NewService(c.OrderPaidConsumer(ctx), c.OrderPaidDecoder(), c.ProducerService(ctx))
	}
	return c.orderPaidConsumerService
}

func (c *diContainer) OrderPaidConsumer(ctx context.Context) kafkaWrapped.Consumer {
	if c.orderPaidConsumer == nil {
		c.orderPaidConsumer = kafkaWrappedConsumer.NewConsumer(
			c.OrderPaidConsumerGroup(),
			[]string{c.config.OrderPaidConsumer.TopicName()},
			logger.Logger(),
		)
	}
	return c.orderPaidConsumer
}

func (c *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if c.orderPaidConsumerGroup != nil {
		return c.orderPaidConsumerGroup
	}

	var err error
	c.orderPaidConsumerGroup, err = sarama.NewConsumerGroup(
		c.config.Kafka.Addresses(),
		c.config.OrderPaidConsumer.ConsumerGroupId(),
		c.config.OrderPaidConsumer.Config(),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create consumer group: %v", err))
	}

	return c.orderPaidConsumerGroup
}

func (c *diContainer) OrderPaidDecoder() kafkaDecoders.OrderPaidDecoder {
	if c.orderPaidDecoder == nil {
		c.orderPaidDecoder = &kafkaDecoder.OrderPaidDecoder{}
	}

	return c.orderPaidDecoder
}

func (c *diContainer) ProducerService(ctx context.Context) services.ShipAssembledProducerService {
	if c.shipAssembledProducerService == nil {
		c.shipAssembledProducerService = shipAssembledProducerService.NewService(
			c.ShipAssembledProducer(),
		)
	}

	return c.shipAssembledProducerService
}

func (c *diContainer) ShipAssembledProducer() kafkaWrapped.Producer {
	if c.shipAssembledProducer == nil {
		c.shipAssembledProducer = kafkaWrappedProducer.NewProducer(
			c.ShipAssembledSyncPruducer(),
			c.config.ShipAssembled.TopicName(),
			logger.Logger(),
		)
	}
	return c.shipAssembledProducer
}

func (c *diContainer) ShipAssembledSyncPruducer() sarama.SyncProducer {
	if c.shipAssembledSyncPruducer != nil {
		return c.shipAssembledSyncPruducer
	}

	var err error
	c.shipAssembledSyncPruducer, err = sarama.NewSyncProducer(
		c.config.Kafka.Addresses(),
		c.config.ShipAssembled.Config(),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create sync producer: %v", err))
	}

	return c.shipAssembledSyncPruducer
}
