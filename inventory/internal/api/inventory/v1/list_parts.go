package v1

import (
	"context"

	conv "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/converter"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *genInventoryV1.ListPartsRequest) (*genInventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.ListParts(ctx, conv.APIPartToModelFilter(req.Filter))
	if err != nil {
		return nil, err
	}

	return &genInventoryV1.ListPartsResponse{
		Parts: conv.ModelToApiParts(parts),
	}, nil
}
