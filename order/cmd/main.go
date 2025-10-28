package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inventoryAddress  = "localhost:50052"
	paymentAddress    = "localhost:50051"
	httpPort          = "8082"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❗ Failed to connect to inventory service (%s): %v\n", inventoryAddress, err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("❗ Failed to close inventory connection: %v\n", err)
		}
	}()
	inventory := inventoryV1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❗ Failed to connect to payment service (%s): %v\n", paymentAddress, err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("❗ Failed to close payment connection: %v\n", err)
		}
	}()
	payment := paymentV1.NewPaymentServiceClient(paymentConn)

	orderHandler := &OrderHandler{
		store:     &OrdersStorage{},
		inventory: &inventory,
		payment:   &payment,
	}

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Mount("/", orderServer)

	httpServer := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("👂 HTTP-сервер запущен на порту %s\n", httpPort)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❗ Ошибка при запуске HTTP-сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🟧 Shutting down gRPC server...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("❗ Ошибка при остановке HTTP-сервера: %v\n", err)
	}

	log.Println("✅ gRPC server gracefully stopped")
}
