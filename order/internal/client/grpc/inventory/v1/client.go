package v1

import (
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"

	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	connection *grpc.ClientConn
	service    inventoryV1.InventoryServiceClient
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
		service:    inventoryV1.NewInventoryServiceClient(connection),
	}, nil
}

// Call on defer
func (c *client) Close() error {
	return c.connection.Close()
}
