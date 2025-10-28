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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	httpPort          = "8082"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrdersStorage struct {
}

type OrderHandler struct{}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderReq) (*orderV1.CreateOrderCreated, error) {
	return nil, nil // TODO:
}

func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	return nil, nil // TODO:
}

func (h *OrderHandler) OrderCancel(ctx context.Context, params orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	return nil, nil // TODO:
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderReq, params orderV1.PayOrderParams) (*orderV1.PayOrderOK, error) {
	return nil, nil // TODO:
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	// TODO: create store
	orderHandler := &OrderHandler{}

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
