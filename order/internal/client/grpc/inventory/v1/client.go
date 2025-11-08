package v1

import (
	def "github.com/Mahno9/GoMicroservicesCourse/order/internal/client/grpc"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	service genInventoryV1.InventoryServiceClient
}

func NewClient(service genInventoryV1.InventoryServiceClient) (*client, error) {
	return &client{
		service: service,
	}, nil
}
