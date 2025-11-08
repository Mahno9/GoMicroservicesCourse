package converter

import (
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	genInventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
)

// APIPartToModelPart конвертирует API-модель (protobuf) в модель данных
func APIPartToModelPart(apiPart *genInventoryV1.Part) *model.Part {
	if apiPart == nil {
		return &model.Part{}
	}

	part := model.Part{
		Uuid:          apiPart.GetUuid(),
		Name:          apiPart.GetName(),
		Description:   apiPart.GetDescription(),
		Price:         apiPart.GetPrice(),
		StockQuantity: apiPart.GetStockQuantity(),
		Tags:          apiPart.GetTags(),
	}

	// Конвертация категории
	part.Category = apiToModelCategory(apiPart.GetCategory())

	// Конвертация размеров
	if apiDimensions := apiPart.GetDimensions(); apiDimensions != nil {
		part.Dimensions = &model.Dimensions{
			Length: apiDimensions.GetLength(),
			Width:  apiDimensions.GetWidth(),
			Height: apiDimensions.GetHeight(),
			Weight: apiDimensions.GetWeight(),
		}
	}

	// Конвертация производителя
	if apiManufacturer := apiPart.GetManufacturer(); apiManufacturer != nil {
		part.Manufacturer = &model.Manufacturer{
			Name:    apiManufacturer.GetName(),
			Country: apiManufacturer.GetCountry(),
			Website: apiManufacturer.GetWebsite(),
		}
	}

	// Конвертация метаданных
	if apiMetadata := apiPart.GetMetadata(); apiMetadata != nil {
		part.Metadata = apiToModelMetadata(apiMetadata)
	}

	// Конвертация времени
	if apiCreatedAt := apiPart.GetCreatedAt(); apiCreatedAt != nil {
		part.CreatedAt = apiCreatedAt.AsTime()
	}

	if apiUpdatedAt := apiPart.GetUpdatedAt(); apiUpdatedAt != nil {
		part.UpdatedAt = lo.ToPtr(apiUpdatedAt.AsTime())
	}

	return &part
}

// ModelToAPIPart конвертирует модель данных в API-модель (protobuf)
func ModelToAPIPart(modelPart *model.Part) *genInventoryV1.Part {
	apiPart := &genInventoryV1.Part{
		Uuid:          modelPart.Uuid,
		Name:          modelPart.Name,
		Description:   modelPart.Description,
		Price:         modelPart.Price,
		StockQuantity: modelPart.StockQuantity,
		Tags:          modelPart.Tags,
	}

	// Конвертация категории
	apiPart.Category = modelToAPICategory(modelPart.Category)

	// Конвертация размеров
	if modelDimensions := modelPart.Dimensions; modelDimensions != nil {
		apiPart.Dimensions = &genInventoryV1.Dimensions{
			Length: modelDimensions.Length,
			Width:  modelDimensions.Width,
			Height: modelDimensions.Height,
			Weight: modelDimensions.Weight,
		}
	}

	// Конвертация производителя
	if modelManufacturer := modelPart.Manufacturer; modelManufacturer != nil {
		apiPart.Manufacturer = &genInventoryV1.Manufacturer{
			Name:    modelManufacturer.Name,
			Country: modelManufacturer.Country,
			Website: modelManufacturer.Website,
		}
	}

	// Конвертация метаданных
	if modelMetadata := modelPart.Metadata; modelMetadata != nil {
		apiPart.Metadata = modelMetadataToAPIMetadata(modelMetadata)
	}

	// Конвертация времени
	apiPart.CreatedAt = timestamppb.New(modelPart.CreatedAt)

	if modelUpdatedAt := modelPart.UpdatedAt; modelUpdatedAt != nil {
		apiPart.UpdatedAt = timestamppb.New(*modelUpdatedAt)
	}

	return apiPart
}

// apiToModelCategory конвертирует категорию из API в модель
func apiToModelCategory(apiCategory genInventoryV1.Category) model.Category {
	switch apiCategory {
	case genInventoryV1.Category_ENGINE:
		return model.CategoryEngine
	case genInventoryV1.Category_FUEL:
		return model.CategoryFuel
	case genInventoryV1.Category_PORTHOLE:
		return model.CategoryPorthole
	case genInventoryV1.Category_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

// modelToAPICategory конвертирует категорию из модели в API
func modelToAPICategory(modelCategory model.Category) genInventoryV1.Category {
	switch modelCategory {
	case model.CategoryEngine:
		return genInventoryV1.Category_ENGINE
	case model.CategoryFuel:
		return genInventoryV1.Category_FUEL
	case model.CategoryPorthole:
		return genInventoryV1.Category_PORTHOLE
	case model.CategoryWing:
		return genInventoryV1.Category_WING
	default:
		return genInventoryV1.Category_UNKNOWN
	}
}

// apiToModelMetadata конвертирует метаданные из API в модель
func apiToModelMetadata(apiMetadata map[string]*genInventoryV1.Value) map[string]any {
	if apiMetadata == nil {
		return nil
	}

	modelMetadata := make(map[string]any)
	for key, value := range apiMetadata {
		var val any
		switch v := value.GetKind().(type) {
		case *genInventoryV1.Value_StringValue:
			val = v.StringValue
		case *genInventoryV1.Value_Int64Value:
			val = v.Int64Value
		case *genInventoryV1.Value_DoubleValue:
			val = v.DoubleValue
		case *genInventoryV1.Value_BoolValue:
			val = v.BoolValue
		}
		modelMetadata[key] = val
	}
	return modelMetadata
}

// modelMetadataToAPIMetadata конвертирует метаданные из модели в API
func modelMetadataToAPIMetadata(modelMetadata map[string]any) map[string]*genInventoryV1.Value {
	if modelMetadata == nil {
		return nil
	}

	apiMetadata := make(map[string]*genInventoryV1.Value)
	for key, value := range modelMetadata {
		if value == nil {
			continue
		}

		var apiValue *genInventoryV1.Value
		switch v := value.(type) {
		case string:
			apiValue = &genInventoryV1.Value{Kind: &genInventoryV1.Value_StringValue{StringValue: v}}
		case int64:
			apiValue = &genInventoryV1.Value{Kind: &genInventoryV1.Value_Int64Value{Int64Value: v}}
		case float64:
			apiValue = &genInventoryV1.Value{Kind: &genInventoryV1.Value_DoubleValue{DoubleValue: v}}
		case bool:
			apiValue = &genInventoryV1.Value{Kind: &genInventoryV1.Value_BoolValue{BoolValue: v}}
		default:
			// Для других типов используем строковое представление
			apiValue = &genInventoryV1.Value{Kind: &genInventoryV1.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}}
		}
		apiMetadata[key] = apiValue
	}
	return apiMetadata
}

// APIPartToModelFilter конвертирует API-фильтр в модель фильтра
func APIPartToModelFilter(apiFilter *genInventoryV1.PartsFilter) *model.PartsFilter {
	if apiFilter == nil {
		return &model.PartsFilter{}
	}

	filter := model.PartsFilter{
		Uuids:                 apiFilter.GetUuids(),
		Names:                 apiFilter.GetNames(),
		ManufacturerCountries: apiFilter.GetManufacturerCountries(),
		Tags:                  apiFilter.GetTags(),
	}

	filter.Categories = make([]model.Category, len(apiFilter.Categories))
	for i, apiCategory := range apiFilter.GetCategories() {
		filter.Categories[i] = apiToModelCategory(apiCategory)
	}

	return &filter
}

// ModelToApiParts конвертирует срез моделей частей в срез API-частей
func ModelToApiParts(modelParts []*model.Part) []*genInventoryV1.Part {
	if modelParts == nil {
		return nil
	}

	apiParts := make([]*genInventoryV1.Part, len(modelParts))
	for i, modelPart := range modelParts {
		apiParts[i] = ModelToAPIPart(modelPart)
	}

	return apiParts
}
