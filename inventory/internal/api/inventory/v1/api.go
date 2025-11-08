package v1

import (
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

type api struct {
	genInventoryV1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewAPI(inventoryService service.PartService) *api {
	return &api{
		partService: inventoryService,
	}
}
