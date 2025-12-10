package v1

import (
	"context"

	"go.uber.org/zap"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/converter"
	"github.com/Mahno9/GoMicroservicesCourse/platform/pkg/logger"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *genInventoryV1.GetPartRequest) (*genInventoryV1.GetPartResponse, error) {
	if err := checkUUID(req.Uuid); err != nil {
		logger.Error(ctx, "invalid UUID", zap.String("uuid", req.Uuid), zap.Error(err))
		return nil, err
	}

	part, err := a.partService.GetPart(ctx, req.Uuid)
	if err != nil {
		logger.Error(ctx, "failed to get part", zap.String("uuid", req.Uuid), zap.Error(err))
		return nil, converter.ModelToAPIError(err)
	}

	logger.Info(ctx, "part retrieved successfully", zap.String("uuid", req.Uuid))
	return &genInventoryV1.GetPartResponse{
		Part: converter.ModelToAPIPart(part),
	}, nil
}
