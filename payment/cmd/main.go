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
	"github.com/Mahno9/GoMicroservicesCourse/payment/internal/config"
	paymentService "github.com/Mahno9/GoMicroservicesCourse/payment/internal/service/payment"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

const (
	configPath = "deploy/compose/payment/.env"
)

func main() {
	// Load .env variables
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("❗ Failed to load env file: %v\n", err)
		return
	} else {
		log.Printf("✅ Env file loaded successfully: %+v\n", cfg)
	}

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

	grpcServer := grpc.NewServer()

	paymentService := paymentService.NewService()
	apiService := paymentV1API.NewAPI(paymentService)
	genPaymentV1.RegisterPaymentServiceServer(grpcServer, apiService)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("👂 gRPC server listening on port %s\n", cfg.GrpcConfig.Port())
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
