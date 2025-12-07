package app

import (
	"context"
	"fmt"
	"time"

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
	repository "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/migrator"
	orderRepo "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/order"
	orderModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/service/order"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/closer"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	config *config.Config

	orderV1API *genOrderV1.Server
	inventory  clients.InventoryClient
	payment    clients.PaymentClient
	repository repository.OrderRepository
	router     *chi.Mux
}

func NewDIContainer(cfg *config.Config) *diContainer {
	return &diContainer{config: cfg}
}

func (c *diContainer) OrderHTTPServer(ctx context.Context) *genOrderV1.Server {
	if c.orderV1API != nil {
		return c.orderV1API
	}

	orderService := orderModel.NewService(c.Inventory(), c.Payment(), c.Repository(ctx))

	apiHandler := orderApiHandlerV1.NewAPIHandler(orderService)

	var err error
	c.orderV1API, err = genOrderV1.NewServer(apiHandler)
	if err != nil {
		panic(fmt.Sprintf("❗ failed to create server: %v\n", err.Error()))
	}

	return c.orderV1API
}

func (c *diContainer) Inventory() clients.InventoryClient {
	if c.inventory != nil {
		return c.inventory
	}

	inventoryConn, err := grpc.NewClient(c.config.ClientsConfig.InventoryAddress(),
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

	paymentConn, err := grpc.NewClient(c.config.ClientsConfig.PaymentAddress(),
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

	dbConnPool, err := pgxpool.New(ctx, c.config.PostgresConfig.URI())
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
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*dbConnPool.Config().Copy().ConnConfig), c.config.PostgresConfig.MigrationsDir())
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
	c.router.Mount("/", c.OrderHTTPServer(ctx))

	return c.router
}
