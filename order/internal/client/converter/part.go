package converter

import (
	"github.com/Mahno9/GoMicroservicesCourse/order/internal/model"
	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/payment/v1"
)

func InventoryToModelPart(inventoryPart *inventoryV1.Part) (*model.Part, error) {
	result := &model.Part{
		Uuid:          inventoryPart.Uuid,
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
		case *inventoryV1.Value_StringValue:
			value = val.StringValue
		case *inventoryV1.Value_Int64Value:
			value = val.Int64Value
		case *inventoryV1.Value_DoubleValue:
			value = val.DoubleValue
		case *inventoryV1.Value_BoolValue:
			value = val.BoolValue
		}
		result.Metadata[k] = &value
	}

	if inventoryPart.CreatedAt != nil {
		createdTime := inventoryPart.CreatedAt.AsTime()
		result.CreatedAt = &createdTime
	}

	if inventoryPart.UpdatedAt != nil {
		updatedTime := inventoryPart.UpdatedAt.AsTime()
		result.UpdatedAt = &updatedTime
	}

	return result, nil
}

func ModelToInventoryPartsFilter(modelFilter *model.PartsFilter) *inventoryV1.PartsFilter {
	if modelFilter == nil {
		return nil
	}

	result := &inventoryV1.PartsFilter{}

	result.Uuids = append(result.Uuids, modelFilter.Uuids...)
	result.Names = append(result.Names, modelFilter.Names...)
	for _, category := range modelFilter.Categories {
		result.Categories = append(result.Categories, inventoryV1.Category(category))
	}
	result.ManufacturerCountries = append(result.ManufacturerCountries, modelFilter.ManufacturerCountries...)
	result.Tags = append(result.Tags, modelFilter.Tags...)

	return result
}

func ModelToPaymentPaymentMethod(paymentMethod int32) paymentV1.PaymentMethod {
	switch paymentMethod {
	case 1:
		return paymentV1.PaymentMethod_CARD
	case 2:
		return paymentV1.PaymentMethod_SBP
	case 3:
		return paymentV1.PaymentMethod_CREDIT_CARD
	case 4:
		return paymentV1.PaymentMethod_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_UNKNOWN
	}
}
