package app

import (
	"context"

	paymentV1API "github.com/Mahno9/GoMicroservicesCourse/payment/internal/api/payment/v1"
	"github.com/Mahno9/GoMicroservicesCourse/payment/internal/config"
	service "github.com/Mahno9/GoMicroservicesCourse/payment/internal/service"
	paymentService "github.com/Mahno9/GoMicroservicesCourse/payment/internal/service/payment"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	config *config.Config

	paymentV1API   genPaymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDIContainer(cfg *config.Config) *diContainer {
	return &diContainer{config: cfg}
}

func (c *diContainer) PaymentV1API(ctx context.Context) genPaymentV1.PaymentServiceServer {
	if c.paymentV1API == nil {
		c.paymentV1API = paymentV1API.NewAPI(c.PaymentService(ctx))
	}

	return c.paymentV1API
}

func (c *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if c.paymentService == nil {
		c.paymentService = paymentService.NewService()
	}

	return c.paymentService
}
