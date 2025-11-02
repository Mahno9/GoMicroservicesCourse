package v1

import (
	"context"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/converter"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	timedContext, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	response, err := c.service.ListParts(timedContext, &inventoryV1.ListPartsRequest{
		Filter: converter.ModelToInventoryPartsFilter(filter),
	})
	if err != nil {
		return nil, err
	}

	parts := make([]*model.Part, 0)
	for _, part := range response.Parts {
		modelPart, err := converter.InventoryToModelPart(part)
		if err != nil {
			return nil, err
		}
		parts = append(parts, modelPart)
	}

	return parts, nil
}
