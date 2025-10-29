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

	paymentV1API "github.com/Mahno9/GoMicroservicesCourse/payment/internal/api/payment/v1"
	paymentService "github.com/Mahno9/GoMicroservicesCourse/payment/internal/service/payment"
	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
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

	paymentService := paymentService.NewService()
	apiService := paymentV1API.NewAPI(paymentService)
	paymentV1.RegisterPaymentServiceServer(grpcServer, apiService)

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
