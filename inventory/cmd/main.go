package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"syscall"

	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50052
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("❗ failed to listen: %v\n", err)
		return
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Fatalf("❗ failed to close listener: %v\n", err)
		}
	}()

	grpcServer := grpc.NewServer()

	service := &inventoryService{
		parts: map[string]*inventoryV1.Part{},
	}
	service.initParts() // init with dummy data
	inventoryV1.RegisterInventoryServiceServer(grpcServer, service)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("👂 gRPC server listening on port %d\n", grpcPort)
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("❗ failed to serve: %v\n", err)
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
