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

// DomainToRepoFilter конвертирует доменный фильтр в фильтр репозитория
func DomainToRepoFilter(domainFilter *domainModel.PartsFilter) *repoModel.PartsFilter {
	if domainFilter == nil {
		return nil
	}

	repoFilter := &repoModel.PartsFilter{}

	// Конвертация uuids из map в slice
	if len(domainFilter.Uuids) > 0 {
		repoFilter.Uuids = make([]string, 0, len(domainFilter.Uuids))
		for uuid := range domainFilter.Uuids {
			repoFilter.Uuids = append(repoFilter.Uuids, uuid)
		}
	}

	// Конвертация names из map в slice
	if len(domainFilter.Names) > 0 {
		repoFilter.Names = make([]string, 0, len(domainFilter.Names))
		for name := range domainFilter.Names {
			repoFilter.Names = append(repoFilter.Names, name)
		}
	}

	// Конвертация categories из map в slice
	if len(domainFilter.Categories) > 0 {
		repoFilter.Categories = make([]repoModel.Category, 0, len(domainFilter.Categories))
		for category := range domainFilter.Categories {
			repoFilter.Categories = append(repoFilter.Categories, repoModel.Category(category))
		}
	}

	// Конвертация manufacturerCountries из map в slice
	if len(domainFilter.ManufacturerCountries) > 0 {
		repoFilter.ManufacturerCountries = make([]string, 0, len(domainFilter.ManufacturerCountries))
		for country := range domainFilter.ManufacturerCountries {
			repoFilter.ManufacturerCountries = append(repoFilter.ManufacturerCountries, country)
		}
	}

	// Конвертация tags из map в slice
	if len(domainFilter.Tags) > 0 {
		repoFilter.Tags = make([]string, 0, len(domainFilter.Tags))
		for tag := range domainFilter.Tags {
			repoFilter.Tags = append(repoFilter.Tags, tag)
		}
	}

	return repoFilter
}

// RepoToDomainFilter конвертирует фильтр репозитория в доменный фильтр
func RepoToDomainFilter(repoFilter *repoModel.PartsFilter) *domainModel.PartsFilter {
	if repoFilter == nil {
		return nil
	}

	domainFilter := &domainModel.PartsFilter{}

	// Конвертация uuids из slice в map
	if len(repoFilter.Uuids) > 0 {
		domainFilter.Uuids = make(map[string]any, len(repoFilter.Uuids))
		for _, uuid := range repoFilter.Uuids {
			domainFilter.Uuids[uuid] = true
		}
	}

	// Конвертация names из slice в map
	if len(repoFilter.Names) > 0 {
		domainFilter.Names = make(map[string]any, len(repoFilter.Names))
		for _, name := range repoFilter.Names {
			domainFilter.Names[name] = true
		}
	}

	// Конвертация categories из slice в map
	if len(repoFilter.Categories) > 0 {
		domainFilter.Categories = make(map[domainModel.Category]any, len(repoFilter.Categories))
		for _, category := range repoFilter.Categories {
			domainFilter.Categories[domainModel.Category(category)] = true
		}
	}

	// Конвертация manufacturerCountries из slice в map
	if len(repoFilter.ManufacturerCountries) > 0 {
		domainFilter.ManufacturerCountries = make(map[string]any, len(repoFilter.ManufacturerCountries))
		for _, country := range repoFilter.ManufacturerCountries {
			domainFilter.ManufacturerCountries[country] = true
		}
	}

	// Конвертация tags из slice в map
	if len(repoFilter.Tags) > 0 {
		domainFilter.Tags = make(map[string]any, len(repoFilter.Tags))
		for _, tag := range repoFilter.Tags {
			domainFilter.Tags[tag] = true
		}
	}

	return domainFilter
}
