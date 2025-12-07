package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partV1API "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/api/inventory/v1"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/config"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository"
	partRepository "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/part"
	partService "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service/part"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

const (
	initPartsTimeout = 5 * time.Second

	configPath = "deploy/compose/inventory/.env"
)

func main() {
	ctx := context.Background()

	// Load .env variables
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("❗ Failed to load env file: %v\n", err)
		return
	} else {
		log.Printf("✅ Env file loaded successfully: %+v\n", cfg)
	}

	// Inventory Repository
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConfig.URI()))
	if err != nil {
		log.Printf("❗ failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("❗ failed to disconnect from database: %v\n", cerr)
		}
	}()

	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("❗ failed to ping database: %v\n", err)
		return
	}

	inventoryDb := client.Database(cfg.MongoConfig.DatabaseName())
	wrappedDb := &repository.MongoDatabaseAdapter{Database: inventoryDb}
	partRepo, err := partRepository.NewRepository(ctx, wrappedDb)
	if err != nil {
		log.Printf("❗ failed to create repository: %v\n", err)
		return
	}

	// Inventory Service
	service := partService.NewService(partRepo)

	timedContext, cancel := context.WithTimeout(ctx, initPartsTimeout)
	defer cancel()

	err = service.InitWithDummy(timedContext) // init with dummy data
	if err != nil {
		log.Printf("❗ failed to init with dummy data: %v\n", err)
		return
	}

	// Inventory API
	serviceAPI := partV1API.NewAPI(service)

	// gRPC service serving
	grpcServer := grpc.NewServer()
	genInventoryV1.RegisterInventoryServiceServer(grpcServer, serviceAPI)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcConfig.Port()))
	if err != nil {
		log.Printf("❗ failed to listen: %v\n", err)
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("❗ failed to close listener: %v\n", err)
		}
	}()

	go func() {
		log.Printf("👂 gRPC server listening on port %s\n", cfg.GrpcConfig.Port())
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Printf("❗ failed to serve: %v\n", err)
			return
		}
	}()

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🟧 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ gRPC server gracefully stopped")
}
