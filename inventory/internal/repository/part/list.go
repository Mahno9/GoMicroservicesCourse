package part

import (
	"context"

	"github.com/samber/lo"

	domainModel "github.com/Mahno9/GoMicroservicesCourse/inventory/internal/model"
	"github.com/Mahno9/GoMicroservicesCourse/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filters *domainModel.PartsFilter) ([]*domainModel.Part, error) {
	parts := make([]*domainModel.Part, 0)
	repoFilters := converter.DomainToRepoFilter(filters)

	for _, part := range r.GetAll() {
		// Check if the part matches all the filters
		matches := true

		if len(repoFilters.Uuids) > 0 && !lo.Contains(repoFilters.Uuids, part.Uuid) {
			matches = false
		}
		if len(repoFilters.Names) > 0 && !lo.Contains(repoFilters.Names, part.Name) {
			matches = false
		}
		if len(repoFilters.Categories) > 0 && !lo.Contains(repoFilters.Categories, part.Category) {
			matches = false
		}
		if len(repoFilters.ManufacturerCountries) > 0 && !lo.Contains(repoFilters.ManufacturerCountries, part.Manufacturer.Country) {
			matches = false
		}

		// Check if any of the part's tags match the filter tags
		if len(repoFilters.Tags) > 0 {
			tagMatch := false
			for _, tag := range repoFilters.Tags {
				if lo.Contains(part.Tags, tag) {
					tagMatch = true
					break
				}
			}
			if !tagMatch {
				matches = false
			}
		}

		// If part matches all filters, add it to the result
		if matches {
			parts = append(parts, converter.RepoToDomainPart(part))
		}
	}

	return parts, nil
}
