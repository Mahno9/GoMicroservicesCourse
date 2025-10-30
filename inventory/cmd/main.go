package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partV1API "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/api/inventory/v1"
	partRepository "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/part"
	partService "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service/part"
	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("❗ failed to listen: %v\n", err)
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("❗ failed to close listener: %v\n", err)
		}
	}()

	grpcServer := grpc.NewServer()

	partRepo := partRepository.NewRepository()
	service := partService.NewService(partRepo)
	service.InitWithDummy() // init with dummy data
	serviceAPI := partV1API.NewAPI(service)
	inventoryV1.RegisterInventoryServiceServer(grpcServer, serviceAPI)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("👂 gRPC server listening on port %d\n", grpcPort)
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Printf("❗ failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🟧 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ gRPC server gracefully stopped")
}
