package converter

import (
	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	repoModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/model"
)

// nolint:dupl
// DomainToRepoPart конвертирует доменную модель в модель репозитория
func DomainToRepoPart(domainPart *domainModel.Part) *repoModel.Part {
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
// RepoToDomainPart конвертирует модель репозитория в доменную модель
func RepoToDomainPart(repoPart *repoModel.Part) *domainModel.Part {
	if repoPart == nil {
		return nil
	}

	domainPart := &domainModel.Part{
		Uuid:          repoPart.Uuid,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      domainModel.Category(repoPart.Category),
		Tags:          repoPart.Tags,
		Metadata:      repoPart.Metadata,
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}

	// Конвертация размеров
	if repoDimensions := repoPart.Dimensions; repoDimensions != nil {
		domainPart.Dimensions = &domainModel.Dimensions{
			Length: repoDimensions.Length,
			Width:  repoDimensions.Width,
			Height: repoDimensions.Height,
			Weight: repoDimensions.Weight,
		}
	}

	// Конвертация производителя
	if repoManufacturer := repoPart.Manufacturer; repoManufacturer != nil {
		domainPart.Manufacturer = &domainModel.Manufacturer{
			Name:    repoManufacturer.Name,
			Country: repoManufacturer.Country,
			Website: repoManufacturer.Website,
		}
	}

	return domainPart
}

// nolint:dupl
// DomainToRepoFilter конвертирует доменный фильтр в фильтр репозитория
func DomainToRepoFilter(domainFilter *domainModel.PartsFilter) *repoModel.PartsFilter {
	if domainFilter == nil {
		return nil
	}

	repoFilter := &repoModel.PartsFilter{}

	// Конвертация uuids из map в slice
	if len(domainFilter.Uuids) > 0 {
		repoFilter.Uuids = make([]string, len(domainFilter.Uuids))
		copy(repoFilter.Uuids, domainFilter.Uuids)
	}

	// Конвертация names из map в slice
	if len(domainFilter.Names) > 0 {
		repoFilter.Names = make([]string, len(domainFilter.Names))
		copy(repoFilter.Names, domainFilter.Names)
	}

	// Конвертация categories из map в slice
	if len(domainFilter.Categories) > 0 {
		repoFilter.Categories = make([]repoModel.Category, len(domainFilter.Categories))
		for i, category := range domainFilter.Categories {
			repoFilter.Categories[i] = repoModel.Category(category)
		}
	}

	// Конвертация manufacturerCountries из map в slice
	if len(domainFilter.ManufacturerCountries) > 0 {
		repoFilter.ManufacturerCountries = make([]string, len(domainFilter.ManufacturerCountries))
		copy(repoFilter.ManufacturerCountries, domainFilter.ManufacturerCountries)
	}

	// Конвертация tags из map в slice
	if len(domainFilter.Tags) > 0 {
		repoFilter.Tags = make([]string, len(domainFilter.Tags))
		copy(repoFilter.Tags, domainFilter.Tags)
	}

	return repoFilter
}

// nolint:dupl
// RepoToDomainFilter конвертирует фильтр репозитория в доменный фильтр
func RepoToDomainFilter(repoFilter *repoModel.PartsFilter) *domainModel.PartsFilter {
	if repoFilter == nil {
		return nil
	}

	domainFilter := &domainModel.PartsFilter{}

	// Конвертация uuids из slice в map
	if len(repoFilter.Uuids) > 0 {
		domainFilter.Uuids = make([]string, len(repoFilter.Uuids))
		copy(domainFilter.Uuids, repoFilter.Uuids)
	}

	// Конвертация names из slice в map
	if len(repoFilter.Names) > 0 {
		domainFilter.Names = make([]string, len(repoFilter.Names))
		copy(domainFilter.Names, repoFilter.Names)
	}

	// Конвертация categories из slice в map
	if len(repoFilter.Categories) > 0 {
		domainFilter.Categories = make([]domainModel.Category, len(repoFilter.Categories))
		for i, category := range repoFilter.Categories {
			domainFilter.Categories[i] = domainModel.Category(category)
		}
	}

	// Конвертация manufacturerCountries из slice в map
	if len(repoFilter.ManufacturerCountries) > 0 {
		domainFilter.ManufacturerCountries = make([]string, len(repoFilter.ManufacturerCountries))
		copy(domainFilter.ManufacturerCountries, repoFilter.ManufacturerCountries)
	}

	// Конвертация tags из slice в map
	if len(repoFilter.Tags) > 0 {
		domainFilter.Tags = make([]string, len(repoFilter.Tags))
		copy(domainFilter.Tags, repoFilter.Tags)
	}

	return domainFilter
}
