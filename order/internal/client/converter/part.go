package converter

import (
	"log"

	"github.com/google/uuid"

	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	genPaymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func InventoryToModelParts(inventoryParts []*genInventoryV1.Part) ([]*model.Part, error) {
	parts := make([]*model.Part, len(inventoryParts))
	for i, part := range inventoryParts {
		parts[i] = inventoryToModelPart(part)
	}

	return parts, nil
}

func inventoryToModelPart(inventoryPart *genInventoryV1.Part) *model.Part {
	partUuid, err := uuid.Parse(inventoryPart.Uuid)
	if err != nil {
		log.Printf("Failed to parse UUID %s: %v", inventoryPart.Uuid, err)
	}

	result := &model.Part{
		Uuid:          partUuid,
		Name:          inventoryPart.Name,
		Description:   inventoryPart.Description,
		Price:         inventoryPart.Price,
		StockQuantity: inventoryPart.StockQuantity,
		Category:      int32(inventoryPart.Category),
		Tags:          inventoryPart.Tags,
		Metadata:      make(map[string]*any),
	}

	if inventoryPart.Dimensions != nil {
		result.Dimensions = &model.Dimensions{
			Length: inventoryPart.Dimensions.Length,
			Width:  inventoryPart.Dimensions.Width,
			Height: inventoryPart.Dimensions.Height,
			Weight: inventoryPart.Dimensions.Weight,
		}
	}

	if inventoryPart.Manufacturer != nil {
		result.Manufacturer = &model.Manufacturer{
			Name:    inventoryPart.Manufacturer.Name,
			Country: inventoryPart.Manufacturer.Country,
			Website: inventoryPart.Manufacturer.Website,
		}
	}

	for k, v := range inventoryPart.Metadata {
		var value any
		switch val := v.Kind.(type) {
		case *genInventoryV1.Value_StringValue:
			value = val.StringValue
		case *genInventoryV1.Value_Int64Value:
			value = val.Int64Value
		case *genInventoryV1.Value_DoubleValue:
			value = val.DoubleValue
		case *genInventoryV1.Value_BoolValue:
			value = val.BoolValue
		}
		result.Metadata[k] = &value
	}

	if inventoryPart.CreatedAt != nil {
		createdTime := inventoryPart.CreatedAt.AsTime()
		result.CreatedAt = createdTime
	} else {
		log.Printf("InventoryPart.CreatedAt is nil")
	}

	if inventoryPart.UpdatedAt != nil {
		updatedTime := inventoryPart.UpdatedAt.AsTime()
		result.UpdatedAt = &updatedTime
	}

	return result
}

func ModelToInventoryPartsFilter(modelFilter *model.PartsFilter) *genInventoryV1.PartsFilter {
	if modelFilter == nil {
		return nil
	}

	result := &genInventoryV1.PartsFilter{}

	for _, partUuid := range modelFilter.Uuids {
		result.Uuids = append(result.Uuids, partUuid.String())
	}
	result.Names = append(result.Names, modelFilter.Names...)
	for _, category := range modelFilter.Categories {
		result.Categories = append(result.Categories, genInventoryV1.Category(category))
	}
	result.ManufacturerCountries = append(result.ManufacturerCountries, modelFilter.ManufacturerCountries...)
	result.Tags = append(result.Tags, modelFilter.Tags...)

	return result
}

func ModelToPaymentPaymentMethod(paymentMethod int32) genPaymentV1.PaymentMethod {
	switch paymentMethod {
	case 1:
		return genPaymentV1.PaymentMethod_CARD
	case 2:
		return genPaymentV1.PaymentMethod_SBP
	case 3:
		return genPaymentV1.PaymentMethod_CREDIT_CARD
	case 4:
		return genPaymentV1.PaymentMethod_INVESTOR_MONEY
	default:
		return genPaymentV1.PaymentMethod_UNKNOWN
	}
}
