package converter

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

// nolint:dupl
// ModelToRepositoryPart конвертирует доменную модель в модель репозитория
func ModelToRepositoryPart(domainPart *model.Part) *repoModel.Part {
	if domainPart == nil {
		return nil
	}

	repoPart := &repoModel.Part{
		Uuid:          domainPart.Uuid,
		Name:          domainPart.Name,
		Description:   domainPart.Description,
		Price:         domainPart.Price,
		StockQuantity: domainPart.StockQuantity,
		Category:      repoModel.Category(domainPart.Category),
		Tags:          domainPart.Tags,
		Metadata:      domainPart.Metadata,
		CreatedAt:     domainPart.CreatedAt,
		UpdatedAt:     domainPart.UpdatedAt,
	}

	// Конвертация размеров
	if domainDimensions := domainPart.Dimensions; domainDimensions != nil {
		repoPart.Dimensions = &repoModel.Dimensions{
			Length: domainDimensions.Length,
			Width:  domainDimensions.Width,
			Height: domainDimensions.Height,
			Weight: domainDimensions.Weight,
		}
	}

	// Конвертация производителя
	if domainManufacturer := domainPart.Manufacturer; domainManufacturer != nil {
		repoPart.Manufacturer = &repoModel.Manufacturer{
			Name:    domainManufacturer.Name,
			Country: domainManufacturer.Country,
			Website: domainManufacturer.Website,
		}
	}

	return repoPart
}

// nolint:dupl
// RepositoryToModelPart конвертирует модель репозитория в доменную модель
func RepositoryToModelPart(repoPart *repoModel.Part) *model.Part {
	if repoPart == nil {
		return nil
	}

	domainPart := &model.Part{
		Uuid:          repoPart.Uuid,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      model.Category(repoPart.Category),
		Tags:          repoPart.Tags,
		Metadata:      repoPart.Metadata,
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}

	// Конвертация размеров
	if repoDimensions := repoPart.Dimensions; repoDimensions != nil {
		domainPart.Dimensions = &model.Dimensions{
			Length: repoDimensions.Length,
			Width:  repoDimensions.Width,
			Height: repoDimensions.Height,
			Weight: repoDimensions.Weight,
		}
	}

	// Конвертация производителя
	if repoManufacturer := repoPart.Manufacturer; repoManufacturer != nil {
		domainPart.Manufacturer = &model.Manufacturer{
			Name:    repoManufacturer.Name,
			Country: repoManufacturer.Country,
			Website: repoManufacturer.Website,
		}
	}

	return domainPart
}

// nolint:dupl
// ModelToRepositoryFilter конвертирует доменный фильтр в фильтр репозитория
func ModelToRepositoryFilter(filter *model.PartsFilter) bson.M {
	if filter == nil {
		return nil
	}

	repoFilter := bson.M{}

	if len(filter.Uuids) > 0 {
		repoFilter["uuid"] = bson.M{"$in": filter.Uuids}
	}

	if len(filter.Names) > 0 {
		repoFilter["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		repoFilter["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		repoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		repoFilter["tags"] = bson.M{"$all": filter.Tags}
	}

	return repoFilter
}
