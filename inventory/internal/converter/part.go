package converter

import (
	"fmt"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	inventoryV1 "github.com/Mahno9/GoMicroservicesCourse/shared/pkg/proto/inventory/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// APIPartToModelPart конвертирует API-модель (protobuf) в модель данных
func APIPartToModelPart(apiPart *inventoryV1.Part) *model.Part {
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
	part.Category = apiCategoryToModelCategory(apiPart.GetCategory())

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
		part.Metadata = apiMetadataToModelMetadata(&apiMetadata)
	}

	// Конвертация времени
	if apiCreatedAt := apiPart.GetCreatedAt(); apiCreatedAt != nil {
		part.CreatedAt = lo.ToPtr(apiCreatedAt.AsTime())
	}

	if apiUpdatedAt := apiPart.GetUpdatedAt(); apiUpdatedAt != nil {
		part.UpdatedAt = lo.ToPtr(apiUpdatedAt.AsTime())
	}

	return &part
}

// ModelPartToAPIPart конвертирует модель данных в API-модель (protobuf)
func ModelPartToAPIPart(modelPart *model.Part) *inventoryV1.Part {
	apiPart := &inventoryV1.Part{
		Uuid:          modelPart.Uuid,
		Name:          modelPart.Name,
		Description:   modelPart.Description,
		Price:         modelPart.Price,
		StockQuantity: modelPart.StockQuantity,
		Tags:          modelPart.Tags,
	}

	// Конвертация категории
	apiPart.Category = modelCategoryToAPICategory(modelPart.Category)

	// Конвертация размеров
	if modelDimensions := modelPart.Dimensions; modelDimensions != nil {
		apiPart.Dimensions = &inventoryV1.Dimensions{
			Length: modelDimensions.Length,
			Width:  modelDimensions.Width,
			Height: modelDimensions.Height,
			Weight: modelDimensions.Weight,
		}
	}

	// Конвертация производителя
	if modelManufacturer := modelPart.Manufacturer; modelManufacturer != nil {
		apiPart.Manufacturer = &inventoryV1.Manufacturer{
			Name:    modelManufacturer.Name,
			Country: modelManufacturer.Country,
			Website: modelManufacturer.Website,
		}
	}

	// Конвертация метаданных
	if modelMetadata := modelPart.Metadata; modelMetadata != nil {
		apiPart.Metadata = modelMetadataToAPIMetadata(&modelMetadata)
	}

	// Конвертация времени
	if modelCreatedAt := modelPart.CreatedAt; modelCreatedAt != nil {
		apiPart.CreatedAt = timestamppb.New(*modelCreatedAt)
	}

	if modelUpdatedAt := modelPart.UpdatedAt; modelUpdatedAt != nil {
		apiPart.UpdatedAt = timestamppb.New(*modelUpdatedAt)
	}

	return apiPart
}

// apiCategoryToModelCategory конвертирует категорию из API в модель
func apiCategoryToModelCategory(apiCategory inventoryV1.Category) model.Category {
	switch apiCategory {
	case inventoryV1.Category_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

// modelCategoryToAPICategory конвертирует категорию из модели в API
func modelCategoryToAPICategory(modelCategory model.Category) inventoryV1.Category {
	switch modelCategory {
	case model.CategoryEngine:
		return inventoryV1.Category_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_WING
	default:
		return inventoryV1.Category_UNKNOWN
	}
}

// apiMetadataToModelMetadata конвертирует метаданные из API в модель
func apiMetadataToModelMetadata(apiMetadata *map[string]*inventoryV1.Value) map[string]*any {
	if apiMetadata == nil {
		return nil
	}

	modelMetadata := make(map[string]*any)
	for key, value := range *apiMetadata {
		var val any
		switch v := value.GetKind().(type) {
		case *inventoryV1.Value_StringValue:
			val = v.StringValue
		case *inventoryV1.Value_Int64Value:
			val = v.Int64Value
		case *inventoryV1.Value_DoubleValue:
			val = v.DoubleValue
		case *inventoryV1.Value_BoolValue:
			val = v.BoolValue
		}
		modelMetadata[key] = &val
	}
	return modelMetadata
}

// modelMetadataToAPIMetadata конвертирует метаданные из модели в API
func modelMetadataToAPIMetadata(modelMetadata *map[string]*any) map[string]*inventoryV1.Value {
	if modelMetadata == nil {
		return nil
	}

	apiMetadata := make(map[string]*inventoryV1.Value)
	for key, value := range *modelMetadata {
		if value == nil {
			continue
		}

		var apiValue *inventoryV1.Value
		switch v := (*value).(type) {
		case string:
			apiValue = &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: v}}
		case int64:
			apiValue = &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: v}}
		case float64:
			apiValue = &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: v}}
		case bool:
			apiValue = &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: v}}
		default:
			// Для других типов используем строковое представление
			apiValue = &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: fmt.Sprintf("%v", v)}}
		}
		apiMetadata[key] = apiValue
	}
	return apiMetadata
}

// APIPartToModelFilter конвертирует API-фильтр в модель фильтра
func APIPartToModelFilter(apiFilter *inventoryV1.PartsFilter) *model.PartsFilter {
	if apiFilter == nil {
		return &model.PartsFilter{}
	}

	filter := model.PartsFilter{
		Uuids:                 apiFilter.GetUuids(),
		Names:                 apiFilter.GetNames(),
		ManufacturerCountries: apiFilter.GetManufacturerCountries(),
		Tags:                  apiFilter.GetTags(),
	}

	// Конвертация категорий
	apiCategories := apiFilter.GetCategories()
	modelCategories := make([]model.Category, len(apiCategories))
	for i, apiCategory := range apiCategories {
		modelCategories[i] = apiCategoryToModelCategory(apiCategory)
	}
	filter.Categories = modelCategories

	return &filter
}

// ModelToApiParts конвертирует срез моделей частей в срез API-частей
func ModelToApiParts(modelParts []*model.Part) *inventoryV1.ListPartsResponse {
	if modelParts == nil {
		return &inventoryV1.ListPartsResponse{
			Parts: []*inventoryV1.Part{},
		}
	}

	apiParts := make([]*inventoryV1.Part, len(modelParts))
	for i, modelPart := range modelParts {
		apiParts[i] = ModelPartToAPIPart(modelPart)
	}

	return &inventoryV1.ListPartsResponse{
		Parts: apiParts,
	}
}
