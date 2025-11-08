package v1

import (
	"context"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/client/converter"
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	response, err := c.service.ListParts(ctx, &genInventoryV1.ListPartsRequest{
		Filter: converter.ModelToInventoryPartsFilter(filter),
	})
	if err != nil {
		return nil, err
	}

	return converter.InventoryToModelParts(response.Parts)
}
