package v1

import (
	"context"

	conv "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/converter"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/service"
	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewAPI(inventoryService service.PartService) *api {
	return &api{
		partService: inventoryService,
	}
}

func (a *api) GetPart(c context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(c, req.Uuid)
	if err != nil {
		return nil, err
	}

	return &inventoryV1.GetPartResponse{
		Part: conv.ModelPartToAPIPart(part),
	}, nil
}

func (a *api) ListParts(c context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filters := conv.APIPartToModelFilter(req.Filter)

	parts, err := a.partService.ListParts(c, filters)
	if err != nil {
		return nil, err
	}

	return conv.ModelToApiParts(parts), nil
}
