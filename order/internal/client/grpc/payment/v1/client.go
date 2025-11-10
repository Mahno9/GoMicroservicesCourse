package v1

import (
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	service genPaymentV1.PaymentServiceClient
}

func NewClient(service genPaymentV1.PaymentServiceClient) (*client, error) {
	return &client{
		service: service,
	}, nil
}
