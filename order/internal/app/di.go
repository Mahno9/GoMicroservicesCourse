package app

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApiHandlerV1 "github.com/Mahno9/GoMicroservicesCourse/order/internal/api/order/v1"
	clients "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	inventoryClientV1 "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc/payment/v1"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/config"
	kafkaConverter "github.com/Mahno9/GoMicroservicesCourse/order/internal/converter/kafka"
	kafkaConverterDecoder "github.com/Mahno9/GoMicroservicesCourse/order/internal/converter/kafka/decoder"
	repository "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
	orderRepo "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/order"
	services "github.com/Mahno9/GoMicroservicesCourse/order/internal/service"
	consumerService "github.com/Mahno9/GoMicroservicesCourse/order/internal/service/consumer"
	orderService "github.com/Mahno9/GoMicroservicesCourse/order/internal/service/order"
	producerService "github.com/Mahno9/GoMicroservicesCourse/order/internal/service/producer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	wrappedKafka "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Mahno9/GoMicroservicesCourse/platform/pkg/kafka/producer"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/migrator"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	config *config.Config

	orderService    services.OrderService
	producerService services.ProducerService
	consumerService services.ConsumerService

	orderV1API *genOrderV1.Server

	inventory clients.InventoryClient
	payment   clients.PaymentClient

	repository repository.OrderRepository
	router     *chi.Mux

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.Producer

	consumerGroup         sarama.ConsumerGroup
	shipAssembledConsumer wrappedKafka.Consumer
	shipAssembledDecoder  kafkaConverter.ShipAssembledDecoder
}

func NewDIContainer(cfg *config.Config) *diContainer {
	return &diContainer{config: cfg}
}

func (c *diContainer) OrderV1API(ctx context.Context) *genOrderV1.Server {
	if c.orderV1API != nil {
		return c.orderV1API
	}

	apiHandler := orderApiHandlerV1.NewAPIHandler(c.OrderService(ctx))

	var err error
	c.orderV1API, err = genOrderV1.NewServer(apiHandler)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create server: %v\n", err.Error()))
	}

	return c.orderV1API
}

func (c *diContainer) OrderService(ctx context.Context) services.OrderService {
	if c.orderService == nil {
		c.orderService = orderService.NewService(c.Inventory(), c.Payment(), c.Repository(ctx), c.ProducerService(ctx))
	}

	return c.orderService
}

func (c *diContainer) ProducerService(ctx context.Context) services.ProducerService {
	if c.producerService == nil {
		c.producerService = producerService.NewService(c.OrderPaidProducer(ctx))
	}

	return c.producerService
}

func (c *diContainer) ConsumerService(ctx context.Context) services.ConsumerService {
	if c.consumerService == nil {
		c.consumerService = consumerService.NewService(c.ShipAssembledConsumer(ctx), c.ShipAssembledDecoder())
	}

	return c.consumerService
}

func (c *diContainer) Inventory() clients.InventoryClient {
	if c.inventory != nil {
		return c.inventory
	}

	inventoryConn, err := grpc.NewClient(c.config.Clients.InventoryAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create inventory connection: %v\n", err.Error()))
	}

	closer.AddNamed("Inventory client", func(ctx context.Context) error {
		return inventoryConn.Close()
	})

	inventoryClient := genInventoryV1.NewInventoryServiceClient(inventoryConn)

	c.inventory, err = inventoryClientV1.NewClient(inventoryClient)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create inventory client: %v\n", err.Error()))
	}

	return c.inventory
}

func (c *diContainer) Payment() clients.PaymentClient {
	if c.payment != nil {
		return c.payment
	}

	paymentConn, err := grpc.NewClient(c.config.Clients.PaymentAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create payment connection: %v\n", err.Error()))
	}
	closer.AddNamed("Payment client", func(ctx context.Context) error {
		return paymentConn.Close()
	})

	paymenService := genPaymentV1.NewPaymentServiceClient(paymentConn)

	c.payment, err = paymentClientV1.NewClient(paymenService)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create payment client: %v\n", err.Error()))
	}

	return c.payment
}

func (c *diContainer) Repository(ctx context.Context) repository.OrderRepository {
	if c.repository != nil {
		return c.repository
	}

	dbConnPool, err := pgxpool.New(ctx, c.config.Postgres.URI())
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create database connection pool: %v\n", err.Error()))
	}

	closer.AddNamed("Database connection pool", func(ctx context.Context) error {
		dbConnPool.Close()
		return nil
	})

	err = dbConnPool.Ping(ctx)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to ping database: %v\n", err.Error()))
	}

	migratePGDatabase(dbConnPool, c)

	c.repository = orderRepo.NewRepository(dbConnPool)

	return c.repository
}

func migratePGDatabase(dbConnPool *pgxpool.Pool, c *diContainer) {
	migratorRunner := migrator.New(stdlib.OpenDB(*dbConnPool.Config().Copy().ConnConfig), c.config.Postgres.MigrationsDir())
	err := migratorRunner.Up()
	if err != nil {
		panic(fmt.Sprintf("❗ failed to run migrations: %v\n", err.Error()))
	}
}

func (c *diContainer) Router(ctx context.Context) *chi.Mux {
	if c.router != nil {
		return c.router
	}

	c.router = chi.NewRouter()

	c.router.Use(middleware.Logger)
	c.router.Use(middleware.Recoverer)
	c.router.Use(middleware.Timeout(10 * time.Second))
	c.router.Mount("/", c.OrderV1API(ctx))

	return c.router
}

func (c *diContainer) SyncProducer() sarama.SyncProducer {
	if c.syncProducer != nil {
		return c.syncProducer
	}

	var err error
	c.syncProducer, err = sarama.NewSyncProducer(
		c.config.Kafka.BrokersAddresses(),
		c.config.OrderPaidProducer.Config(),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create sync producer: %v\n", err.Error()))
	}

	closer.AddNamed("Sync producer", func(ctx context.Context) error {
		return c.syncProducer.Close()
	})

	return c.syncProducer
}

func (c *diContainer) OrderPaidProducer(ctx context.Context) wrappedKafka.Producer {
	if c.orderPaidProducer != nil {
		return c.orderPaidProducer
	}

	c.orderPaidProducer = wrappedKafkaProducer.NewProducer(
		c.SyncProducer(),
		c.config.OrderPaidProducer.TopicName(),
		logger.Logger(),
	)

	return c.orderPaidProducer
}

func (c *diContainer) ShipAssembledConsumer(ctx context.Context) wrappedKafka.Consumer {
	if c.shipAssembledConsumer != nil {
		return c.shipAssembledConsumer
	}

	c.shipAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
		c.ConsumerGroup(ctx),
		[]string{c.config.ShipAssembledConsumer.TopicName()},
		logger.Logger(),
	)

	return c.shipAssembledConsumer
}

func (c *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if c.consumerGroup != nil {
		return c.consumerGroup
	}

	var err error
	c.consumerGroup, err = sarama.NewConsumerGroup(
		c.config.Kafka.BrokersAddresses(),
		c.config.ShipAssembledConsumer.ConsumerGroupID(),
		c.config.ShipAssembledConsumer.Config(),
	)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create consumer group: %v", err))
	}

	return c.consumerGroup
}

func (c *diContainer) ShipAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if c.shipAssembledDecoder == nil {
		c.shipAssembledDecoder = kafkaConverterDecoder.NewShipAssembledDecoder()
	}

	return c.shipAssembledDecoder
}
