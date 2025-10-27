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

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = 50051
)

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (ps paymentService) PayOrder(c context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// TODO: Validate?

	time.Sleep(2 * time.Second)

	paymentUuid := uuid.New().String()
	log.Println("🆗 Оплата прошла успешно,", paymentUuid)

	return &paymentV1.PayOrderResponse{
		TransactionUuid: paymentUuid,
	}, nil
}

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

	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(grpcServer, service)

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
