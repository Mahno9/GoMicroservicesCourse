package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApiHandlerV1 "github.com/Mahno9/GoMicroservicesCourse/order/internal/api/order/v1"
	inventoryClientV1 "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc/payment/v1"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/config"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/migrator"
	orderRepo "github.com/Mahno9/GoMicroservicesCourse/order/internal/repository/order"
	orderModel "github.com/Mahno9/GoMicroservicesCourse/order/internal/service/order"
	genOrderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	configPath = "deploy/compose/order/.env"
)

func main() {
	ctx := context.Background()

	// Load .env variables
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("❗ Failed to load env file: %v\n", err)
		return
	}

	// Inventory
	inventoryConn, err := grpc.NewClient(cfg.ClientsConfig.InventoryAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❗ Failed to create inventory connection: %v\n", err)
		return
	}
	defer func() {
		err := inventoryConn.Close()
		if err != nil {
			log.Printf("❗ Failed to close inventory connection: %v\n", err)
		}
	}()

	inventoryClient := genInventoryV1.NewInventoryServiceClient(inventoryConn)

	inventory, err := inventoryClientV1.NewClient(inventoryClient)
	if err != nil {
		log.Printf("❗ Failed to create inventory client: %v\n", err)
	}

	// Payment
	paymentConn, err := grpc.NewClient(cfg.ClientsConfig.PaymentAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❗ Failed to create payment connection: %v\n", err)
		return
	}
	defer func() {
		err := paymentConn.Close()
		if err != nil {
			log.Printf("❗ Failed to close payment connection: %v\n", err)
		}
	}()

	paymenService := genPaymentV1.NewPaymentServiceClient(paymentConn)

	payment, err := paymentClientV1.NewClient(paymenService)
	if err != nil {
		log.Printf("❗ Failed to create payment client: %v\n", err)
	}

	// Order Repository: DB connection
	dbConnPool, err := pgxpool.New(ctx, cfg.PostgresConfig.URI())
	if err != nil {
		log.Printf("❗ Failed to create database connection pool: %v\n", err)
		return
	}
	defer dbConnPool.Close()
	err = dbConnPool.Ping(ctx)
	if err != nil {
		log.Printf("❗ Failed to ping database: %v\n", err)
		return
	}

	// Order Repository: migration
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*dbConnPool.Config().Copy().ConnConfig), cfg.PostgresConfig.MigrationsDir())
	err = migratorRunner.Up()
	if err != nil {
		log.Printf("❗ Failed to run migrations: %v\n", err)
		return
	}

	// Order Repository: instance
	repository := orderRepo.NewRepository(dbConnPool)

	// Order Service
	orderService := orderModel.NewService(inventory, payment, repository)

	// HTTP requests handler
	apiHandler := orderApiHandlerV1.NewAPIHandler(orderService)

	orderServer, err := genOrderV1.NewServer(apiHandler)
	if err != nil {
		panic(err)
	}

	// HTTP Serving
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Mount("/", orderServer)

	httpServer := &http.Server{
		Addr:              cfg.HttpConfig.Address(),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("👂 HTTP-сервер запущен на порту %s\n", cfg.HttpConfig.Port())
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❗ Ошибка при запуске HTTP-сервера: %v\n", err)
		}
	}()

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🟧 Shutting down gRPC server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("❗ Ошибка при остановке HTTP-сервера: %v\n", err)
	}

	log.Println("✅ gRPC server gracefully stopped")
}
