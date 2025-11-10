package v1

import (
	"context"

	conv "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/converter"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *genInventoryV1.GetPartRequest) (*genInventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}

	return &genInventoryV1.GetPartResponse{
		Part: conv.ModelToAPIPart(part),
	}, nil
}
