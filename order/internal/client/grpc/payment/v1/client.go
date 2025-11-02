package v1

import (
	"time"

	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"

	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	connectionTimeout = 1 * time.Second
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	connection *grpc.ClientConn
	service    paymentV1.PaymentServiceClient
}

// Call Close() for the new client on defer
func NewClient(address string) (*client, error) {
	connection, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		connection: connection,
		service:    paymentV1.NewPaymentServiceClient(connection),
	}, nil
}

// Call on defer
func (c *client) Close() error {
	return c.connection.Close()
}
